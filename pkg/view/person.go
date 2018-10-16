package view

import (
	"github.com/gotk3/gotk3/gtk"
)

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

type AddToGroup struct {
    CurrentGroupLB *gtk.ListBox
    grid           *gtk.Grid
    NewGroupCBT    *gtk.ComboBoxText
}

func NewAddToGroup() *AddToGroup {
    grid := SetupGrid(gtk.ORIENTATION_VERTICAL)

	cornerNW := SetupLabel("    ")
    currentGroupL := SetupLabel("Current Groups:")
    currentGroupLB := SetupListBox()
    newGroupL := SetupLabel("New Group:")
	newGroupCBT := SetupComboBoxText()
    cornerSE := SetupLabel("    ")

    newGroupCBT.SetHExpand(true)
    currentGroupL.SetVExpand(true)

	grid.Add(cornerNW)
	grid.Attach(currentGroupL, 1, 2, 1, 1)
    grid.Attach(currentGroupLB, 2, 2, 1, 1)
    grid.Attach(newGroupL, 1, 1, 1, 1)
    grid.Attach(newGroupCBT, 2, 1, 1, 1)
    grid.Attach(cornerSE, 3, 3, 1, 1)

    return &AddToGroup{
        CurrentGroupLB: currentGroupLB,
        grid:           grid,
        NewGroupCBT:    newGroupCBT,
    }
}

type EditPerson struct {
    CurrentGroupLB *gtk.ListBox
    NewGroupCBT    *gtk.ComboBoxText
    Notebook       *gtk.Notebook
    PersonContent  *PersonContent
    SaveB          *gtk.ToolButton
    Win            *gtk.Window
}

func EditPersonWindow() *EditPerson{
    win := SetupPopupWindow("Edit Person", 350, 216)
    box := SetupBox()
    nb := SetupNotebook()
    tb := SetupToolbar()
    save := SetupToolButtonLabel("Save")

    personContent := NewPersonContent()
    addToGroup := NewAddToGroup()

    save.SetExpand(true)
    save.SetVExpand(true)

    tb.Add(save)
    tb.SetHExpand(true)

    nb.AppendPage(personContent.grid, SetupLabel("Edit Person"))
    nb.AppendPage(addToGroup.grid, SetupLabel("Add to Group"))

    box.Add(nb)
    box.Add(tb)

	win.Add(box)
	win.ShowAll()

	return &EditPerson{
        CurrentGroupLB: addToGroup.CurrentGroupLB,
        NewGroupCBT:    addToGroup.NewGroupCBT,
        Notebook:       nb,
        PersonContent:  personContent,
        SaveB:          save,
		Win:            win,
	}
}
