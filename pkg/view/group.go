package view

import (
	"github.com/gotk3/gotk3/gtk"
)

// GroupContent contains the entries needed by the controller to retrieve
// information from an 'Edit Group' menu.   It contains the corresopnding
// gtk.Entry objects, and the grid that holds them together.   It is meant
// to be used inside a gtk.Container.
type GroupContent struct {
	grid       *gtk.Grid
	GroupNameE *gtk.Entry
	StartE     *gtk.Entry
	EndE       *gtk.Entry
}

// NewGroupContent creates and returns a new GroupContent object.
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

// AddMember contains a gtk.ListBox, a gtk.ComboBoxText and the grid
// holding them.   These are used to create the dialog where a member
// is added to a group in the application.   It is intended to be used
// inside a gtk.Container.
type AddMember struct {
	CurrentMemberLB *gtk.ListBox
	grid            *gtk.Grid
	NewMemberCBT    *gtk.ComboBoxText
}

func NewAddMember() *AddMember {
	grid := SetupGrid(gtk.ORIENTATION_VERTICAL)

	cornerNW := SetupLabel("    ")
	currentMemberL := SetupLabel("Current Members:")
	currentMemberLB := SetupListBox()
	newMemberL := SetupLabel("New Member:")
	newMemberCBT := SetupComboBoxText()
	cornerSE := SetupLabel("    ")

	newMemberCBT.SetHExpand(true)
	currentMemberLB.SetVExpand(true)

	grid.Add(cornerNW)
	grid.Attach(currentMemberL, 1, 2, 1, 1)
	grid.Attach(currentMemberLB, 2, 2, 1, 1)
	grid.Attach(newMemberL, 1, 1, 1, 1)
	grid.Attach(newMemberCBT, 2, 1, 1, 1)
	grid.Attach(cornerSE, 3, 3, 1, 1)

	return &AddMember{
		CurrentMemberLB: currentMemberLB,
		grid:            grid,
		NewMemberCBT:    newMemberCBT,
	}
}

// EditGroup represents the window from the 'Edit Group' dialog in the
// main application.   It contains the information of both a GroupContent
// and an AddMember objects.
type EditGroup struct {
	CurrentMemberLB *gtk.ListBox
	GroupContent    *GroupContent
	NewMemberCBT    *gtk.ComboBoxText
	Notebook        *gtk.Notebook
	SaveB           *gtk.ToolButton
	Win             *gtk.Window
}

// EditGroupWindow creates and draws the 'Edit Group' window for the main
// application.   It returns a new EditGroup object.
func EditGroupWindow() *EditGroup {
	win := SetupPopupWindow("Edit Group", 350, 216)
	box := SetupBox()
	nb := SetupNotebook()
	tb := SetupToolbar()
	save := SetupToolButtonLabel("Save")

	addMember := NewAddMember()
	groupContent := NewGroupContent()

	save.SetExpand(true)
	save.SetVExpand(true)

	tb.Add(save)
	tb.SetHExpand(true)

	nb.AppendPage(groupContent.grid, SetupLabel("Edit Group"))
	nb.AppendPage(addMember.grid, SetupLabel("Add Member"))

	box.Add(nb)
	box.Add(tb)

	win.Add(box)
	win.ShowAll()

	return &EditGroup{
		CurrentMemberLB: addMember.CurrentMemberLB,
		GroupContent:    groupContent,
		NewMemberCBT:    addMember.NewMemberCBT,
		Notebook:        nb,
		SaveB:           save,
		Win:             win,
	}
}
