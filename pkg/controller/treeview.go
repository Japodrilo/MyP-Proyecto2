package controller

import (
	"log"

	"github.com/Japodrilo/MyP-Proyecto2/pkg/model"
	"github.com/Japodrilo/MyP-Proyecto2/pkg/view"
	"github.com/gotk3/gotk3/gtk"
)

// Constants corresponding to the column numbers in the tree view.
const (
	COLUMN_TITLE = iota
	COLUMN_ARTIST
	COLUMN_ALBUM
	COLUMN_GENRE
	COLUMN_PATH
	COLUMN_VISIBLE
	COLUMN_ID
)

// TreeView represents the tree view in the main window of
// the application.   It contains the gtk window form the view
// module, and a dictionary with Rola id's as keys, and the rows
// of the tree view as entries (*gtk.TreeIter).
type TreeView struct {
	*view.TreeView
	Rows map[int64]*gtk.TreeIter
}

// NewTreeView takes as an argument a view.TreeView, and creates
// a new map to hold the Rola id's and rows of the tree view. It
// returns a TreeView.
func NewTreeView(treeView *view.TreeView) *TreeView {
	Rows := make(map[int64]*gtk.TreeIter)
	return &TreeView{
		treeView,
		Rows,
	}
}

// RowInfo binds the information contained in a row in a single object.
// This is mainly used to pass a row's information as an argument to
// glib.IdleAdd().
type RowInfo struct {
	title   string
	artist  string
	album   string
	genre   string
	path    string
	visible bool
	id      int64
}

// Unexported method to append a row to the list store for the tree view.
// TODO: use RowInfo as parameter.
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

// Unexported method to append a row to the list store for the tree view.
func (treeview *TreeView) addRowStruct(rowInfo *RowInfo) {
	treeview.addRow(rowInfo.title, rowInfo.artist, rowInfo.album, rowInfo.genre, rowInfo.path, rowInfo.visible, rowInfo.id)
}

// Unexported method to append a row to the list store for the
// tree view directly from a Rola.
func (treeview *TreeView) addRowFromRola(rola *model.Rola) {
	iter := treeview.ListStore.Append()

	err := treeview.ListStore.Set(iter,
		[]int{COLUMN_TITLE, COLUMN_ARTIST, COLUMN_ALBUM, COLUMN_GENRE, COLUMN_PATH, COLUMN_VISIBLE, COLUMN_ID},
		[]interface{}{rola.Title(), rola.Artist(), rola.Album(), rola.Genre(), rola.Path(), true, rola.ID()})

	if err != nil {
		log.Fatal("Unable to add row:", err)
	}
	treeview.Rows[rola.ID()] = iter
}

// Unexported method to update the performer of a Rola in the
// tree view.
func (treeview *TreeView) updatePerformer(rola *model.Rola) {
	iter := treeview.Rows[rola.ID()]
	treeview.ListStore.SetValue(iter, 1, rola.Artist())
}

// Unexported method to update all the entries of a row in
// the tree view.
func (treeview *TreeView) updateRow(rola *model.Rola) {
	iter := treeview.Rows[rola.ID()]
	treeview.ListStore.SetValue(iter, 0, rola.Title())
	treeview.ListStore.SetValue(iter, 1, rola.Artist())
	treeview.ListStore.SetValue(iter, 2, rola.Album())
	treeview.ListStore.SetValue(iter, 3, rola.Genre())
}

// AllVisible makes all the rows of the tree view visible.
func (treeview *TreeView) AllVisible() {
	iter, ok := treeview.ListStore.GetIterFirst()
	for ok {
		treeview.ListStore.SetValue(iter, 5, true)
		ok = treeview.ListStore.IterNext(iter)
	}
}

// AllInvisible hides all the rows of the tree view.
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
