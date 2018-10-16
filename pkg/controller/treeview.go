package controller

import (
    "log"

    "github.com/gotk3/gotk3/gtk"
    "github.com/Japodrilo/MyP-Proyecto2/pkg/model"
    "github.com/Japodrilo/MyP-Proyecto2/pkg/view"
)

const (
	COLUMN_TITLE = iota
	COLUMN_ARTIST
    COLUMN_ALBUM
    COLUMN_GENRE
    COLUMN_PATH
    COLUMN_VISIBLE
    COLUMN_ID
)

type TreeView struct {
    *view.TreeView
    Rows map[int64]*gtk.TreeIter
}

func NewTreeView(treeView *view.TreeView) *TreeView {
    Rows := make(map[int64]*gtk.TreeIter)
	return &TreeView{
        treeView,
        Rows,
	}
}

type RowInfo struct {
    title   string
    artist  string
    album   string
    genre   string
    path    string
    visible bool
    id      int64
}

//Method to append a row to the list store for the tree view
func (treeview *TreeView) addRow(title, artist, album, genre, path string, visible bool, id int64) {
	iter := treeview.ListStore.Append()

	err := treeview.ListStore.Set(iter,
		[]int{COLUMN_TITLE, COLUMN_ARTIST, COLUMN_ALBUM, COLUMN_GENRE, COLUMN_PATH, COLUMN_VISIBLE, COLUMN_ID},
		[]interface{}{title, artist, album, genre, path, visible, id})

	if err != nil {
		log.Fatal("Unable to add row:", err)
	}
    treeview.Rows[id] = iter
}

//Method to append a row to the list store for the tree view
func (treeview *TreeView) addRowStruct(rowInfo *RowInfo) {
	treeview.addRow(rowInfo.title, rowInfo.artist, rowInfo.album, rowInfo.genre, rowInfo.path, rowInfo.visible, rowInfo.id)
}

func (treeview *TreeView) addRowFromRola(rola *model.Rola) {
    // Get an iterator for a new row at the end of the list store
	iter := treeview.ListStore.Append()

	// Set the contents of the list store row that the iterator represents
	err := treeview.ListStore.Set(iter,
		[]int{COLUMN_TITLE, COLUMN_ARTIST, COLUMN_ALBUM, COLUMN_GENRE, COLUMN_PATH, COLUMN_VISIBLE, COLUMN_ID},
		[]interface{}{rola.Title(), rola.Artist(), rola.Album(), rola.Genre(), rola.Path(), true, rola.ID()})

	if err != nil {
		log.Fatal("Unable to add row:", err)
	}
    treeview.Rows[rola.ID()] = iter
}

func (treeview *TreeView) AllVisible() {
    iter, ok := treeview.ListStore.GetIterFirst()
    for ok {
        treeview.ListStore.SetValue(iter, 5, true)
        ok = treeview.ListStore.IterNext(iter)
    }
}

func (treeview *TreeView) AllInvisible() {
    sel, err := treeview.TreeView.TreeView.GetSelection()
    if err != nil {
        log.Fatal("could not get tree selection:", err)
    }
    sel.SetMode(gtk.SELECTION_NONE)
    iter, ok := treeview.ListStore.GetIterFirst()
    for ok {
        treeview.ListStore.SetValue(iter, 5, false)
        ok = treeview.ListStore.IterNext(iter)
    }
    sel.SetMode(gtk.SELECTION_SINGLE)
}

func (treeview *TreeView) updatePerformer(rola *model.Rola) {
    iter := treeview.Rows[rola.ID()]
    treeview.ListStore.SetValue(iter, 1, rola.Artist())
}

func (treeview *TreeView) updateRow(rola *model.Rola) {
    iter := treeview.Rows[rola.ID()]
    treeview.ListStore.SetValue(iter, 0, rola.Title())
    treeview.ListStore.SetValue(iter, 1, rola.Artist())
    treeview.ListStore.SetValue(iter, 2, rola.Album())
    treeview.ListStore.SetValue(iter, 3, rola.Genre())
}
