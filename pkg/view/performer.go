package view

import (
	"fmt"

	"github.com/gotk3/gotk3/gtk"
)

type EditPerformer struct {
    Win        *gtk.Window
    StageNameE *gtk.Entry
    RealNameE  *gtk.Entry
    BirthE     *gtk.Entry
    DeathE     *gtk.Entry
    SaveB      *gtk.ToolButton
}


// Class that represents the add/edit performer menu.
func NewEditPerformer() *EditPerformer {
	win := SetupPopupWindow("Edit Performer", 300, 185)
    box := SetupBox()
    tb := SetupToolbar()
    save := SetupToolButtonLabel("Save")
	grid := SetupGrid(gtk.ORIENTATION_VERTICAL)

	cornerNW := SetupLabel("    ")
	stageNameL := SetupLabel("Stage name:")
	stageNameE := SetupEntry()
	realNameL := SetupLabel("Real name:")
	realNameE := SetupEntry()
    birthL := SetupLabel("Date of birth:")
	birthE := SetupEntry()
	deathL := SetupLabel("Date of death:")
	deathE := SetupEntry()
    cornerSW := SetupLabel("    ")

	grid.Add(cornerNW)
	grid.Attach(stageNameL, 1, 1, 1, 1)
	grid.Attach(realNameL, 1, 2, 1, 1)
    grid.Attach(birthL, 1, 3, 1, 1)
	grid.Attach(deathL, 1, 4, 1, 1)
    grid.Attach(stageNameE, 2, 1, 1, 1)
	grid.Attach(realNameE, 2, 2, 1, 1)
    grid.Attach(birthE, 2, 3, 1, 1)
	grid.Attach(deathE, 2, 4, 1, 1)
    grid.Attach(cornerSW, 1, 5, 1, 1)

    save.SetExpand(true)
    save.SetVExpand(true)

    tb.Add(save)
    tb.SetHExpand(true)

    box.Add(grid)
    box.Add(tb)

	win.Add(box)
	win.ShowAll()

	return &EditPerformer{
		Win: win,
		StageNameE: stageNameE,
		RealNameE:  realNameE,
        BirthE:     birthE,
        DeathE:    deathE,
		SaveB:      save,
	}
}

/**
 * Dibuja una ventana nueva para el diálogo de error en la
 * conexión.
 */
func PopUpErrorConexion(servidor, puerto string) {
	win := SetupPopupWindow("Error", 500, 48)
	box := SetupBox()
	grid := SetupGrid(gtk.ORIENTATION_HORIZONTAL)
	mensaje := SetupLabel(fmt.Sprintf("No fue posible establecer la conexión con \"%v:%v\"", servidor, puerto))
	espacio1 := SetupLabel("    ")
	espacio1.SetHExpand(true)
	espacio2 := SetupLabel("    ")
	espacio2.SetHExpand(true)
	// aceptar := SetupButtonClick("Aceptar", func() {
	// 	win.Close()
	// })
	box.Add(mensaje)
	grid.Add(espacio1)
	// grid.Add(aceptar)
	grid.Add(espacio2)
	box.Add(grid)
	win.Add(box)
	win.ShowAll()
}
