package view

import (
	"github.com/gotk3/gotk3/gtk"
)

type EditPerson struct {
    AddGroupB     *gtk.ToolButton
    PersonContent *PersonContent
    SaveB         *gtk.ToolButton
    Win           *gtk.Window
}

type PersonContent struct {
    grid       *gtk.Grid
    StageNameE *gtk.Entry
    RealNameE  *gtk.Entry
    BirthE     *gtk.Entry
    DeathE     *gtk.Entry
}

func NewPersonContent() *PersonContent {
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
    cornerSE := SetupLabel("    ")

    stageNameE.SetHExpand(true)
    realNameE.SetHExpand(true)
    birthE.SetHExpand(true)
    deathE.SetHExpand(true)

	grid.Add(cornerNW)
	grid.Attach(stageNameL, 1, 1, 1, 1)
	grid.Attach(realNameL, 1, 2, 1, 1)
    grid.Attach(birthL, 1, 3, 1, 1)
	grid.Attach(deathL, 1, 4, 1, 1)
    grid.Attach(stageNameE, 2, 1, 1, 1)
	grid.Attach(realNameE, 2, 2, 1, 1)
    grid.Attach(birthE, 2, 3, 1, 1)
	grid.Attach(deathE, 2, 4, 1, 1)
    grid.Attach(cornerSE, 3, 5, 1, 1)

    return &PersonContent{
        grid:       grid,
        StageNameE: stageNameE,
        RealNameE:  realNameE,
        BirthE:     birthE,
        DeathE:     deathE,
    }
}

func EditPersonWindow() *EditPerson{
    win := SetupPopupWindow("Edit Person", 350, 216)
    box := SetupBox()
    tb := SetupToolbar()
    save := SetupToolButtonLabel("Save")
    addGroup := SetupToolButtonLabel("Add to group")

    personContent := NewPersonContent()

    save.SetExpand(true)
    save.SetVExpand(true)
    addGroup.SetExpand(true)
    addGroup.SetVExpand(true)

    tb.Add(save)
    tb.Add(addGroup)
    tb.SetHExpand(true)

    box.Add(personContent.grid)
    box.Add(tb)

	win.Add(box)
	win.ShowAll()

	return &EditPerson{
        AddGroupB:     addGroup,
        PersonContent: personContent,
        SaveB:         save,
		Win:           win,
	}
}
