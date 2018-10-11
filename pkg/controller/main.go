package controller

import(
    "bytes"
    "fmt"
    "image"
    "image/jpeg"
    _ "image/png"
    _ "image/gif"
    "log"
    "os"
    "strings"

    "github.com/Japodrilo/MyP-Proyecto2/pkg/model"
	"github.com/Japodrilo/MyP-Proyecto2/pkg/view"

    "github.com/bogem/id3v2"
    "github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"github.com/gotk3/gotk3/glib"
)
/**
 * Controlador de la ventana principal.
 * Contiene como campos un cliente del paquete modelo, la ventana,
 * un objeto Cuaderno, la ListBox de la ventana, un objeto menu
 * del paquete vista, un diccionario de renglones de la List box,
 * un canal para coordinar la información de invitaciones de salas y
 * una rebanada de salas.
 */
type Principal struct {
    database   *model.Database
    mainWindow *view.MainWindow
	treeview   *TreeView
    treeSel    *gtk.TreeSelection
}

/**
 * Constructor de la ventana principal.
 * Le asigna accines al menú, y activa/desactiva las opciones
 * pertinenetes.
 */
func NewPrincipal() *Principal {
    database, exists := model.NewDatabase()
	mainWindow := view.SetupMainWindow()
    treeview := NewTreeView(mainWindow.TreeView)
    sel, _ := treeview.TreeView.TreeView.GetSelection()

    principal := &Principal{
        database:   database,
		mainWindow: mainWindow,
		treeview:   treeview,
        treeSel:    sel,
	}

    if !exists {
        database.CreateDB()
        principal.Populate()
    }

    database.LoadDB()

    mainWindow.Buttons["search"].Connect("clicked", func() {
        mainWindow.SearchBar.SetSearchMode(true)
        mainWindow.SearchBar.GrabFocus()
    })

    sel.Connect("changed", func() {principal.SelectionChanged(sel)})

    mainWindow.Buttons["populate"].Connect("clicked", func() {
        principal.treeSel.UnselectAll()
        principal.treeSel.SetMode(gtk.SELECTION_NONE)
        //glib.IdleAdd(principal.treeview.ListStore.Clear, )
        principal.database.Database.Close()
        principal.database, _ = model.NewDatabase()
        principal.database.CreateDB()
        mainWindow.Buttons["populate"].SetSensitive(false)
        glib.IdleAdd(principal.treeview.ListStore.Clear, )
        principal.Populate()
    })

    mainWindow.Buttons["edit"].Connect("clicked", func() {
        sel.UnselectAll()
    })

    mainWindow.SearchEntry.Connect("activate", func() {
        text := view.GetTextEntryClean(mainWindow.SearchEntry)
        principal.SimpleSearch(text)
    })

    principal.mainWindow.Win.ShowAll()
	return principal
}

func (principal *Principal) Populate() {
    miner := model.NewMiner()
    miner.Traverse()
    go miner.Extract()

    go miner.Populate(principal.database)
    go principal.PopulateOnTheFly(miner)
}

func (principal *Principal) Repopulate() {
    principal.database.LoadDB()
    principal.PopulateFromExistingDB(principal.database)
}

type SongInfo struct {
    image *gtk.Image
    title string
    artist string
    album string
}

// Handler of "activate" signal of TreeView's selection
func (principal *Principal) SelectionChanged(s *gtk.TreeSelection) {
	rows := s.GetSelectedRows(principal.treeview.Filter)
	items := make([]string, 0, rows.Length())

    if rows == nil {
        principal.defaultImage("Title", "Artist", "Album")
        return
    }

	for l := rows; l != nil; l = l.Next() {
		path := l.Data().(*gtk.TreePath)
		iter, _ := principal.treeview.ListStore.GetIter(path)
        cell, _ := principal.treeview.ListStore.GetValue(iter, 5)
        vis, _ := cell.GoValue()
        fmt.Println(vis)
        for i := 0; i < 5; i++ {
            cell, _ := principal.treeview.ListStore.GetValue(iter, i)
		    str, _ := cell.GetString()
		    items = append(items, str)
        }
	}
    for i := 0; i < 4; i++ {
        if i < 3 {
            fmt.Print(items[i] + " - ")
        } else {
            fmt.Println(items[i])
        }
    }
    tag, err := id3v2.Open(items[4], id3v2.Options{Parse: true})
    if err != nil {
        log.Fatal("error while opening mp3 file: ", items[4] + " ", err)
    }
    defer tag.Close()
    pictures := tag.GetFrames(tag.CommonID("Attached picture"))
    if len(pictures) == 0 {
        principal.defaultImage(items[0], items[1], items[2])
    } else {
        pic, ok := pictures[0].(id3v2.PictureFrame)
        if !ok {
                log.Fatal("could not assert picture frame")
            }
        file, _ := os.Create("./image.jpg")
        imageReader := bytes.NewReader(pic.Picture)
        loadedImage, _, _ := image.Decode(imageReader)
        if loadedImage != nil {
            err = jpeg.Encode(file, loadedImage, nil)
            pix, _ := gdk.PixbufNewFromFileAtScale("./image.jpg", 250, 250, false)
            image, _ := gtk.ImageNewFromPixbuf(pix)
            glib.IdleAdd(principal.AttachInfo, &SongInfo{image, items[0], items[1], items[2]})
        } else {
            principal.defaultImage(items[0], items[1], items[2])
        }
    }
}

func (principal *Principal) defaultImage(title, artist, album string) {
    pix, _ := gdk.PixbufNewFromFileAtScale("./noimage.png", 250, 250, false)
    image, _ := gtk.ImageNewFromPixbuf(pix)
    glib.IdleAdd(principal.AttachInfo, &SongInfo{image, title, artist, album})
    glib.IdleAdd(principal.mainWindow.Win.ShowAll, )
}

func (principal *Principal) AttachInfo(songInfo *SongInfo) {
    previous, err := principal.mainWindow.Grid.GetChildAt(0,0)
    if err != nil {
        log.Fatal("unable to get child from grid:", err)
    }
    previous.Destroy()
    principal.mainWindow.Grid.Attach(songInfo.image, 0, 0, 1, 1)
    principal.mainWindow.SongInfo[0].SetText("\n\n\n\t" + songInfo.title + "\n\n\n")
    principal.mainWindow.SongInfo[1].SetText("\t" + songInfo.artist + "\n\n\n")
    principal.mainWindow.SongInfo[2].SetText("\t" + songInfo.album)
    principal.mainWindow.Win.ShowAll()
}

func (principal *Principal) PopulateOnTheFly(miner *model.Miner) {
    for rola := range miner.TrackList {
	    glib.IdleAdd(principal.treeview.addRowFromRola, rola)
    }
    principal.mainWindow.Buttons["populate"].SetSensitive(true)
    principal.treeSel.SetMode(gtk.SELECTION_SINGLE)
}

func (principal *Principal) SimpleSearch(wildcard string) {
    principal.treeSel.UnselectAll()
    principal.treeview.AllInvisible()
    tx, err := principal.database.Database.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare("SELECT performers.name, albums.name, rolas.path, rolas.title, rolas.genre, rolas.id_rola FROM rolas INNER JOIN performers ON performers.id_performer = rolas.id_performer INNER JOIN albums ON albums.id_album = rolas.id_album WHERE performers.name LIKE ? OR albums.name LIKE ? OR rolas.title LIKE ?")
	if err != nil {
		log.Fatal("could not prepare query: ", err)
	}
    defer stmt.Close()

    wildCard := "%" + strings.TrimSpace(wildcard) + "%"

    rows, err := stmt.Query(wildCard, wildCard, wildCard)
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()
    for rows.Next() {
        var performer string
        var album string
        var path string
		var title string
        var genre string
        var id int64
        err = rows.Scan(&performer, &album, &path, &title, &genre, &id)
        if err != nil {
            log.Fatal(err)
        }
        iter := principal.treeview.Rows[id]
        principal.treeview.ListStore.SetValue(iter, 5, true)
    }
    err = rows.Err()
    if err != nil {
        log.Fatal(err)
    }
    tx.Commit()
}

func (principal *Principal) PopulateFromExistingDB(database *model.Database) {
    rows, err := database.Database.Query("SELECT performers.name, albums.name, rolas.path, rolas.title, rolas.genre, rolas.id_rola FROM rolas INNER JOIN performers ON performers.id_performer = rolas.id_performer INNER JOIN albums ON albums.id_album = rolas.id_album")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
        var performer string
        var album string
        var path string
		var title string
        var genre string
        var id int64
		err = rows.Scan(&performer, &album, &path, &title, &genre, &id)
		if err != nil {
			log.Fatal(err)
		}
		glib.IdleAdd(principal.treeview.addRowStruct, &RowInfo{title, performer, album, genre, path, true, id})
	}
	err = rows.Err()
	if err != nil {
		log.Fatal("could not scan row:", err)
	}
}

/**
 * Función que crea una nueva ventana principal, y la dibuja.
 * Esta es la función que corre el archivo principal para el
 * cliente (cliente.go)
 */
func MainWindow() {
	principal := NewPrincipal()
    principal.mainWindow.Win.ShowAll()
    principal.Repopulate()
}
