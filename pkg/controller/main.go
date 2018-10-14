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

    "github.com/dhowden/tag"
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

    treeview.TreeView.TreeView.Connect("row-activated", func() {principal.RowActivated()})

    mainWindow.Buttons["populate"].Connect("clicked", func() {
        principal.treeSel.UnselectAll()
        principal.treeSel.SetMode(gtk.SELECTION_NONE)
        principal.database, _ = model.NewDatabase()
        principal.database.CreateDB()
        principal.database.LoadDB()
        mainWindow.Buttons["populate"].SetSensitive(false)
        glib.IdleAdd(principal.treeview.ListStore.Clear, )
        principal.Populate()
    })

    mainWindow.Buttons["edit"].Connect("clicked", func() {
        principal.RowActivated()
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

func (principal *Principal) rowTextValues() []string {
    values := make([]string, 0)
    _, iter, ok := principal.treeSel.GetSelected()
    if !ok {
        return values
    }
    for i := 0; i < 5; i++ {
        cell, _ := principal.treeview.Filter.GetValue(iter, i)
		str, _ := cell.GetString()
		values = append(values, str)
    }
    return values
}

// Handler of "activate" signal of TreeView's selection
func (principal *Principal) SelectionChanged(s *gtk.TreeSelection) {
	items := principal.rowTextValues()
    if len(items) == 0 {
        return
    }
    for i := 0; i < 4; i++ {
        if i < 3 {
            fmt.Print(items[i] + " - ")
        } else {
            fmt.Println(items[i])
        }
    }
    file, err := os.Open(items[4])
    if err != nil {
        log.Fatal("could not open file:", err)
    }
    metadata, err := tag.ReadFrom(file)
    if err != nil {
        log.Fatal("error while reading the tags in file: ", items[4] + " ", err)
    }
    picture := metadata.Picture()
    if picture == nil {
        principal.defaultImage(items[0], items[1], items[2])
    } else {
        pic := picture.Data
        file, _ := os.Create("./image.jpg")
        imageReader := bytes.NewReader(pic)
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

func (principal *Principal) RowActivated() {
    rowValues := principal.rowTextValues()
    if len(rowValues) == 0 {
        return
    }
    performerID, performerType := principal.database.QueryPerformerType(rowValues[1])
    name := principal.rowTextValues()[1]
    switch performerType {
    case 0:
        personPopUp := view.EditPersonWindow()
        personID := principal.database.ExistsPerson(name, "ñ")
        stageName, realName, birth, death := principal.database.QueryPerson(personID)
        personPopUp.PersonContent.StageNameE.SetText(stageName)
        personPopUp.PersonContent.RealNameE.SetText(realName)
        personPopUp.PersonContent.BirthE.SetText(birth)
        personPopUp.PersonContent.DeathE.SetText(death)
        personPopUp.SaveB.Connect("clicked", func() {
            principal.savePersonContent(personPopUp.PersonContent)
            personPopUp.Win.Close()
        })
    case 1:
        groupPopUp := view.EditGroupWindow()
        groupID := principal.database.ExistsGroup(name)
        groupName, start, end := principal.database.QueryGroup(groupID)
        groupPopUp.GroupContent.GroupNameE.SetText(groupName)
        groupPopUp.GroupContent.StartE.SetText(start)
        groupPopUp.GroupContent.EndE.SetText(end)
        groupPopUp.SaveB.Connect("clicked", func() {
            principal.saveGroupContent(groupPopUp.GroupContent)
            groupPopUp.Win.Close()
        })
    case 2:
        performerPopUp := view.EditPerformerWindow()
        performerPopUp.PersonContent.StageNameE.SetText(name)
        performerPopUp.GroupContent.GroupNameE.SetText(name)
        performerPopUp.SaveB.Connect("clicked", func() {
            newType := performerPopUp.Notebook.GetCurrentPage()
            principal.database.UpdatePerformerType(performerID, newType)
            switch newType {
            case 0:
                principal.savePersonContent(performerPopUp.PersonContent)
            case 1:
                principal.saveGroupContent(performerPopUp.GroupContent)
            }
            performerPopUp.Win.Close()
        })
    }
}

func (principal *Principal) saveGroupContent(groupContent *view.GroupContent) {
    newGroupName := view.GetTextEntry(groupContent.GroupNameE)
    newStart := view.GetTextEntry(groupContent.StartE)
    newEnd := view.GetTextEntry(groupContent.EndE)
    principal.saveGroup(newGroupName, newStart, newEnd)
}

func (principal *Principal) savePersonContent(personContent *view.PersonContent) {
    newStageName := view.GetTextEntry(personContent.StageNameE)
    newRealName := view.GetTextEntry(personContent.RealNameE)
    newBirth := view.GetTextEntry(personContent.BirthE)
    newDeath := view.GetTextEntry(personContent.DeathE)
    principal.savePerson(newStageName, newRealName, newBirth, newDeath)
}

func (principal *Principal) saveGroup(groupName, start, end string) {
    groupID := principal.database.ExistsGroup(groupName)
    if groupID > 0 {
        principal.database.UpdateGroup(groupName, start, end, groupID)
    } else {
        principal.database.AddGroup(groupName, start, end)
    }
}

func (principal *Principal) savePerson(stageName, realName, birth, death string) {
    personID := principal.database.ExistsPerson(stageName, realName)
    if personID > 0 {
        principal.database.UpdatePerson(stageName, realName, birth, death, personID)
    } else {
        principal.database.AddPerson(stageName, realName, birth, death)
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
