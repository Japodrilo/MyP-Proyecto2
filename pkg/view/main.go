package view

import (
    "github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

/**
 * Estructura que representa a la ventana principal y
 * sus campos de interés para el controlador.
 */
type MainWindow struct {
    Buttons        map[string]*gtk.ToolButton
    Grid           *gtk.Grid
    ScrolledWindow *gtk.ScrolledWindow
    SearchBar      *gtk.SearchBar
    SearchEntry    *gtk.SearchEntry
    SongInfo       []*gtk.Label
    TreeView       *TreeView
    Win		       *gtk.Window
}

/**
 * Constructor, función que dibuja la ventana principal
 * e inicializa los campos que necesita el controlador.
 */
func SetupMainWindow() *MainWindow {
    buttons := make(map[string]*gtk.ToolButton)
	win := SetupWindow("Rolas")
	box := SetupBox()
    boxtop := SetupBox()
    boxinfo:= SetupBox()
    albumLabel := SetupLabel("\tAlbum")
    artistLabel := SetupLabel("\tArtist\n\n\n")
    titleLabel := SetupLabel("\n\n\n\tTitle\n\n\n")
    tb := SetupToolbar()
    sb := SetupSearchBar()
    se := SetupSearchEntry()
    play := SetupToolButtonIcon("media-playback-start")
    edit := SetupToolButtonIcon("stock_edit")
    search := SetupToolButtonIcon("list-add")
    populate := SetupToolButtonIcon("reload")
    treeview := NewTreeView()
    scrwin := SetupScrolledWindow()
	grid := SetupGrid(gtk.ORIENTATION_HORIZONTAL)

    pix, _ := gdk.PixbufNewFromFileAtScale("./noimage.png", 250, 250, true)
    defaultImage, _ := gtk.ImageNewFromPixbuf(pix)

    boxtop.SetOrientation(gtk.ORIENTATION_HORIZONTAL)
    boxtop.Add(tb)
    boxtop.Add(sb)

    boxinfo.Add(titleLabel)
    boxinfo.Add(artistLabel)
    boxinfo.Add(albumLabel)

    sb.SetSearchMode(true)
    sb.Add(se)
    sb.ConnectEntry(se)
    sb.SetShowCloseButton(false)

    tb.Add(populate)
    tb.Add(play)
    tb.Add(edit)
    tb.Add(search)
    tb.SetStyle(gtk.TOOLBAR_ICONS)

    buttons["populate"] = populate
    buttons["play"] = play
    buttons["edit"] = edit
    buttons["search"] = search

    box.Add(boxtop)
    box.Add(scrwin)
    box.Add(grid)
    //grid.SetVExpand(true)

    grid.Attach(defaultImage, 0, 0, 1, 1)
    grid.Attach(boxinfo, 2, 0, 1, 1)

    scrwin.SetVExpand(true)
    scrwin.Add(treeview.TreeView)
	//box.Add(grid)

	win.Add(box)
	win.ShowAll()

    songInfo := []*gtk.Label{titleLabel, artistLabel, albumLabel}

	return &MainWindow{
		Buttons:        buttons,
        Grid:           grid,
        ScrolledWindow: scrwin,
        SearchBar:      sb,
        SearchEntry:    se,
        SongInfo:       songInfo,
        TreeView:       treeview,
        Win:            win,
	}
}
