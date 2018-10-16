package view

import (
	"github.com/gotk3/gotk3/gtk"
)

type EditPerformer struct {
    GroupContent  *GroupContent
    Notebook      *gtk.Notebook
    PersonContent *PersonContent
    SaveB         *gtk.ToolButton
    Win           *gtk.Window
}


func EditPerformerWindow() *EditPerformer{
    win := SetupPopupWindow("Edit Performer", 350, 216)
    box := SetupBox()
    nb := SetupNotebook()
    tb := SetupToolbar()
    save := SetupToolButtonLabel("Save")

    groupContent := NewGroupContent()
    personContent := NewPersonContent()

    save.SetExpand(true)
    save.SetVExpand(true)

    tb.Add(save)
    tb.SetHExpand(true)

    nb.AppendPage(personContent.grid, SetupLabel("Person"))
    nb.AppendPage(groupContent.grid, SetupLabel("Group"))

    box.Add(nb)
    box.Add(tb)

	win.Add(box)
	win.ShowAll()

	return &EditPerformer{
        GroupContent:  groupContent,
        Notebook:      nb,
        PersonContent: personContent,
        SaveB:         save,
		Win:           win,
	}
}
