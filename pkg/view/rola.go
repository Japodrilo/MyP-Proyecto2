package view

import (
	"github.com/gotk3/gotk3/gtk"
)

// An EditRola represents the window used by the 'Edit Rola' menu
// in the main application window.
type EditRola struct {
	RolaContent *RolaContent
	SaveB       *gtk.ToolButton
	Win         *gtk.Window
}

// RolaContent contains all the entries (gtk.Entry) used in the 'Edit Rola'
// menu in the main application window, as well as the grid holding them
// together.   It is meant to be used inside a gtk.Container.
type RolaContent struct {
	grid    *gtk.Grid
	TitleE  *gtk.Entry
	ArtistE *gtk.Entry
	AlbumE  *gtk.Entry
	GenreE  *gtk.Entry
	TrackE  *gtk.Entry
	YearE   *gtk.Entry
}

// NewRolaContent creates and returns a new RolaContent.
func NewRolaContent() *RolaContent {
	grid := SetupGrid(gtk.ORIENTATION_VERTICAL)

	cornerNW := SetupLabel("    ")
	titleL := SetupLabel("Title:")
	titleE := SetupEntry()
	artistL := SetupLabel("Artist:")
	artistE := SetupEntry()
	albumL := SetupLabel("Album:")
	albumE := SetupEntry()
	genreL := SetupLabel("Genre:")
	genreE := SetupEntry()
	trackL := SetupLabel("Track:")
	trackE := SetupEntry()
	yearL := SetupLabel("Year:")
	yearE := SetupEntry()
	cornerSE := SetupLabel("    ")

	titleE.SetHExpand(true)
	artistE.SetHExpand(true)
	albumE.SetHExpand(true)
	genreE.SetHExpand(true)
	trackE.SetHExpand(true)
	yearE.SetHExpand(true)

	grid.Add(cornerNW)
	grid.Attach(titleL, 1, 1, 1, 1)
	grid.Attach(artistL, 1, 2, 1, 1)
	grid.Attach(albumL, 1, 3, 1, 1)
	grid.Attach(genreL, 1, 4, 1, 1)
	grid.Attach(trackL, 1, 5, 1, 1)
	grid.Attach(yearL, 1, 6, 1, 1)
	grid.Attach(titleE, 2, 1, 1, 1)
	grid.Attach(artistE, 2, 2, 1, 1)
	grid.Attach(albumE, 2, 3, 1, 1)
	grid.Attach(genreE, 2, 4, 1, 1)
	grid.Attach(trackE, 2, 5, 1, 1)
	grid.Attach(yearE, 2, 6, 1, 1)

	grid.Attach(cornerSE, 3, 7, 1, 1)

	return &RolaContent{
		grid:    grid,
		TitleE:  titleE,
		ArtistE: artistE,
		AlbumE:  albumE,
		GenreE:  genreE,
		TrackE:  trackE,
		YearE:   yearE,
	}
}

// EditRolaWindow creates an EditRola and draws the corresponding
// window.
func EditRolaWindow() *EditRola {
	win := SetupPopupWindow("Edit Rola", 350, 216)
	box := SetupBox()
	tb := SetupToolbar()
	save := SetupToolButtonLabel("Save")

	rolaContent := NewRolaContent()

	save.SetExpand(true)
	save.SetVExpand(true)

	tb.Add(save)
	tb.SetHExpand(true)

	box.Add(rolaContent.grid)
	box.Add(tb)

	win.Add(box)
	win.ShowAll()

	return &EditRola{
		RolaContent: rolaContent,
		SaveB:       save,
		Win:         win,
	}
}
