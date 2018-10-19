// View of the application, all the windows are built in this package.
package view

import (
	"github.com/gotk3/gotk3/gtk"
)

// About represents the about dialog window in the main application.
// It contains its window for the controller to close it.
type About struct {
	Win *gtk.Window
}

// NewAbout creates, returns, and draws a new About object.
func NewAbout() *About {
	win := SetupPopupWindow("About", 280, 40)
	grid := SetupGrid(gtk.ORIENTATION_VERTICAL)

	cornerNW := SetupLabel("    ")
	VersionL := SetupLabel("Rolas Database Version 0.0.7")
	CopyrightL := SetupLabel("MIT License, 2018, César Hernández Cruz")
	cornerSE := SetupLabel("    ")

	VersionL.SetHAlign(3)
	CopyrightL.SetHAlign(3)

	grid.Add(cornerNW)
	grid.Attach(VersionL, 1, 1, 1, 1)
	grid.Attach(CopyrightL, 1, 2, 1, 1)
	grid.Attach(cornerSE, 2, 3, 1, 1)

	win.Add(grid)
	win.Connect("destroy", func() { win.Close() })
	win.ShowAll()
	return &About{
		Win: win,
	}
}
