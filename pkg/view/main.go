package view

import (
	"log"
	"os"
	"os/user"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

// MainWindow represents the view of the main window.  As an object it
// holds the gtk objects used by the controller.
type MainWindow struct {
	Buttons        map[string]*gtk.ToolButton
	Grid           *gtk.Grid
	ScrolledWindow *gtk.ScrolledWindow
	SearchEntry    *gtk.SearchEntry
	SongInfo       []*gtk.Label
	TreeView       *TreeView
	Win            *gtk.Window
}

// SetupMainWindow draws the main window and initializes the
// gtk objects needed by the controller.
func SetupMainWindow() *MainWindow {
	buttons := make(map[string]*gtk.ToolButton)
	win := SetupWindow("Rolas")
	box := SetupBox()
	gridtop := SetupGrid(gtk.ORIENTATION_HORIZONTAL)
	boxinfo := SetupBox()
	albumLabel := SetupLabel("\tAlbum")
	artistLabel := SetupLabel("\tArtist\n\n\n")
	titleLabel := SetupLabel("\n\n\n\tTitle\n\n\n")
	tb := SetupToolbar()
	tb2 := SetupToolbar()
	se := SetupSearchEntry()
	edit := SetupToolButtonIcon("gtk-edit")
	performers := SetupToolButtonIcon("gtk-open")
	new := SetupToolButtonIcon("gtk-new")
	populate := SetupToolButtonIcon("gtk-refresh")
	about := SetupToolButtonIcon("gtk-info")
	treeview := NewTreeView()
	scrwin := SetupScrolledWindow()
	grid := SetupGrid(gtk.ORIENTATION_HORIZONTAL)
	space1 := SetupLabel("                       ")
	space2 := SetupLabel("                       ")
	space3 := SetupLabel("                       ")

	home, err := user.Current()
	if err != nil {
		log.Fatal("could not retrieve the current user:", err)
	}
	cache := home.HomeDir + "/.cache/rolas"
	fileExists := true
	if _, err := os.Stat(cache + "/noimage.png"); os.IsNotExist(err) {
		fileExists = false
	}
	if !fileExists {
		RestoreAsset(cache, "noimage.png")
	}
	pix, _ := gdk.PixbufNewFromFileAtScale(cache + "/noimage.png", 250, 250, true)
	defaultImage, _ := gtk.ImageNewFromPixbuf(pix)

	se.SetHExpand(true)

	gridtop.Add(tb)
	gridtop.Add(space1)
	gridtop.Add(space2)
	gridtop.Attach(se, 2, 0, 3, 1)
	gridtop.Add(space3)
	gridtop.Add(tb2)

	boxinfo.Add(titleLabel)
	boxinfo.Add(artistLabel)
	boxinfo.Add(albumLabel)

	tb.Add(populate)
	tb.Add(edit)
	tb.Add(performers)
	tb.Add(new)
	tb.SetStyle(gtk.TOOLBAR_ICONS)

	tb2.Add(about)

	buttons["populate"] = populate
	buttons["edit"] = edit
	buttons["performers"] = performers
	buttons["new"] = new
	buttons["about"] = about

	box.Add(gridtop)
	box.Add(scrwin)
	box.Add(grid)

	grid.Attach(defaultImage, 0, 0, 1, 1)
	grid.Attach(boxinfo, 2, 0, 1, 1)

	scrwin.SetVExpand(true)
	scrwin.Add(treeview.TreeView)

	win.SetIconName("gtk-media-record")
	win.Add(box)
	win.ShowAll()

	songInfo := []*gtk.Label{titleLabel, artistLabel, albumLabel}

	return &MainWindow{
		Buttons:        buttons,
		Grid:           grid,
		ScrolledWindow: scrwin,
		SearchEntry:    se,
		SongInfo:       songInfo,
		TreeView:       treeview,
		Win:            win,
	}
}
