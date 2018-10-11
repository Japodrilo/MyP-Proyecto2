package view

import (
    "log"

    "github.com/gotk3/gotk3/glib"
    "github.com/gotk3/gotk3/gtk"
)

type TreeView struct {
    TreeView  *gtk.TreeView
    ListStore *gtk.ListStore
    Filter    *gtk.TreeModelFilter
}

func NewTreeView() *TreeView {
    tv, ls, filter := setupTreeView()
    return &TreeView{
        TreeView:        tv,
        ListStore:       ls,
        Filter:          filter,
    }

}

// IDs to access the tree view columns by
const (
	COLUMN_TITLE = iota
	COLUMN_ARTIST
    COLUMN_ALBUM
    COLUMN_GENRE
    COLUMN_PATH
    COLUMN_VISIBLE
    COLUMN_ID
)

// Add a column to the tree view (during the initialization of the tree view)
func createColumn(title string, id int) *gtk.TreeViewColumn {
	cellRenderer, err := gtk.CellRendererTextNew()
	if err != nil {
		log.Fatal("Unable to create text cell renderer:", err)
	}

	column, err := gtk.TreeViewColumnNewWithAttribute(title, cellRenderer, "text", id)
	if err != nil {
		log.Fatal("Unable to create cell column:", err)
	}

	return column
}

// Add a column to the tree view (during the initialization of the tree view)
func createInvisibleColumn(title string, id int) *gtk.TreeViewColumn {
	column := createColumn(title, id)
    column.SetVisible(false)
	return column
}

/**
 * Creates a tree view and the list store that holds its data
 */
func setupTreeView() (*gtk.TreeView, *gtk.ListStore, *gtk.TreeModelFilter) {
	treeView, err := gtk.TreeViewNew()
	if err != nil {
		log.Fatal("Unable to create tree view:", err)
	}

	treeView.AppendColumn(createColumn("Title", COLUMN_TITLE))
	treeView.AppendColumn(createColumn("Artist", COLUMN_ARTIST))
    treeView.AppendColumn(createColumn("Album", COLUMN_ALBUM))
    treeView.AppendColumn(createColumn("Genre", COLUMN_GENRE))
    treeView.AppendColumn(createInvisibleColumn("Path", COLUMN_PATH))
    treeView.AppendColumn(createInvisibleColumn("Visible", COLUMN_VISIBLE))
    treeView.AppendColumn(createInvisibleColumn("ID", COLUMN_ID))


	listStore, err := gtk.ListStoreNew(glib.TYPE_STRING, glib.TYPE_STRING, glib.TYPE_STRING, glib.TYPE_STRING, glib.TYPE_STRING, glib.TYPE_BOOLEAN, glib.TYPE_INT)
	if err != nil {
		log.Fatal("Unable to create list store:", err)
	}
    filter, err := listStore.FilterNew(nil)
    if err != nil {
		log.Fatal("Unable to create filter:", err)
	}
    filter.SetVisibleColumn(5)
    treeView.SetModel(filter)

	return treeView, listStore, filter
}
