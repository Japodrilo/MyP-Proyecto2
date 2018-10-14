package view

import (
	"github.com/gotk3/gotk3/gtk"
)

type EditGroup struct {
    Win          *gtk.Window
    GroupContent *GroupContent
    SaveB      *gtk.ToolButton
    AddMemberB *gtk.ToolButton
}

type GroupContent struct {
    grid       *gtk.Grid
    GroupNameE *gtk.Entry
    StartE     *gtk.Entry
    EndE       *gtk.Entry
}

func NewGroupContent() *GroupContent {
    grid := SetupGrid(gtk.ORIENTATION_VERTICAL)

	cornerNW := SetupLabel("    ")
	groupNameL := SetupLabel("Name:")
	groupNameE := SetupEntry()
    startL := SetupLabel("Start date:")
	startE := SetupEntry()
	endL := SetupLabel("End date:")
	endE := SetupEntry()
    cornerSE := SetupLabel("    ")

    groupNameE.SetHExpand(true)
    startE.SetHExpand(true)
    endE.SetHExpand(true)

	grid.Add(cornerNW)
	grid.Attach(groupNameL, 1, 1, 1, 1)
    grid.Attach(startL, 1, 2, 1, 1)
	grid.Attach(endL, 1, 3, 1, 1)
    grid.Attach(groupNameE, 2, 1, 1, 1)
    grid.Attach(startE, 2, 2, 1, 1)
	grid.Attach(endE, 2, 3, 1, 1)
    grid.Attach(cornerSE, 3, 4, 1, 1)

    return &GroupContent{
        grid:       grid,
    	GroupNameE: groupNameE,
        StartE:     startE,
        EndE:       endE,
    }
}

func EditGroupWindow() *EditGroup{
    win := SetupPopupWindow("Edit Group", 350, 216)
    box := SetupBox()
    tb := SetupToolbar()
    save := SetupToolButtonLabel("Save")
    addMember := SetupToolButtonLabel("Add member")

    groupContent := NewGroupContent()

    save.SetExpand(true)
    save.SetVExpand(true)
    addMember.SetExpand(true)
    addMember.SetVExpand(true)

    tb.Add(save)
    tb.Add(addMember)
    tb.SetHExpand(true)

    box.Add(groupContent.grid)
    box.Add(tb)

	win.Add(box)
	win.ShowAll()

	return &EditGroup{
        AddMemberB:   addMember,
        GroupContent: groupContent,
        SaveB:        save,
		Win:          win,
	}
}
