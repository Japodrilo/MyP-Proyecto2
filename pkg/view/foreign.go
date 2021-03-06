package view

import (
	"github.com/gotk3/gotk3/gtk"
)

// EditForeignPerformer holds the gtk.Window object corresponding to the
// 'edit performer' window, and its necessary elements to connect it to
// the model.
type EditForeignPerformer struct {
	GroupCBT      *gtk.ComboBoxText
	GroupContent  *GroupContent
	Notebook      *gtk.Notebook
	PersonCBT     *gtk.ComboBoxText
	PersonContent *PersonContent
	SaveB         *gtk.ToolButton
	Win           *gtk.Window
}

// EditForeignPerformerWindow draws the window coresponding to the
// 'Edit Performer' menu accessed through the 'performers' button in the
// main window.   It is intended to edit any existing person or group,
// not necessarily appearing as a performer in the database.
func EditForeignPerformerWindow() *EditForeignPerformer {
	win := SetupPopupWindow("Edit Performer", 350, 216)
	box := SetupBox()
	nb := SetupNotebook()
	tb := SetupToolbar()
	save := SetupToolButtonLabel("Save")
	personCBT := SetupComboBoxText()
	groupCBT := SetupComboBoxText()

	groupContent := NewGroupContent()
	personContent := NewPersonContent()

	save.SetExpand(true)
	save.SetVExpand(true)

	tb.Add(save)
	tb.SetHExpand(true)

	personCBT.SetHExpand(true)
	personContent.grid.InsertRow(1)
	personContent.grid.Attach(personCBT, 1, 1, 2, 1)
	groupCBT.SetHExpand(true)
	groupContent.grid.InsertRow(1)
	groupContent.grid.Attach(groupCBT, 1, 1, 2, 1)

	nb.AppendPage(personContent.grid, SetupLabel("Person"))
	nb.AppendPage(groupContent.grid, SetupLabel("Group"))

	box.Add(nb)
	box.Add(tb)

	win.Add(box)
	win.ShowAll()

	return &EditForeignPerformer{
		GroupCBT:      groupCBT,
		GroupContent:  groupContent,
		Notebook:      nb,
		PersonCBT:     personCBT,
		PersonContent: personContent,
		SaveB:         save,
		Win:           win,
	}
}
