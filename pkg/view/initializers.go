package view

import(
	"log"

	"github.com/gotk3/gotk3/gtk"
)

/**
 * Este archivo contiene funciones para inicializar los
 * objetos de gtk con manejo de errores integrado.   También
 * se añaden algunas funciones para obtener información de
 * los mismos, también con manejo de erores.
 * No se comenta cada una de las funciones, pues lo que hacen
 * puede deducirse directamente de su nombre.
 */

func GetBufferEntry(entry *gtk.SearchEntry) *gtk.EntryBuffer {
	buffer, err := entry.GetBuffer()
	if err != nil {
		log.Fatal("Unable to get buffer:", err)
	}
	return buffer
}

func GetLabelText(label *gtk.Label) string {
	text, err := label.GetText()
	if err != nil {
		log.Fatal("Unable to get text from label:", err)
	}
	return text
}

func GetTextEntry(entry *gtk.Entry) string {
	text, err := entry.GetText()
	if err != nil {
		log.Fatal("Unable to get text from buffer:", err)
	}
	return text
}

func GetTextEntryClean(entry *gtk.SearchEntry) string {
	text, err := entry.GetText()
	if err != nil {
		log.Fatal("Unable to get text from buffer:", err)
	}
	buffer := GetBufferEntry(entry)
	buffer.DeleteText(0,-1)
	return text + "\n"
}

func SetupBox() *gtk.Box {
	box, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		log.Fatal("Unable to create box:", err)
	}
	box.SetHomogeneous(false)
	return box
}

// func SetupButtonClick(label string, onClick func()) *gtk.Button {
// 	btn := SetupButton(label)
// 	btn.Connect("clicked", onClick)
// 	return btn
// }

func SetupComboBoxText() *gtk.ComboBoxText {
	cb, err := gtk.ComboBoxTextNew()
	if err != nil {
		log.Fatal("Unable to create Combo Box:", err)
	}
	return cb
}

func SetupEntry() *gtk.Entry {
	entry, err := gtk.EntryNew()
	if err != nil {
		log.Fatal("Unable to create Entry:", err)
	}
	return entry
}

func SetupGrid(orient gtk.Orientation) *gtk.Grid {
	grid, err := gtk.GridNew()
	if err != nil {
		log.Fatal("Unable to create grid:", err)
	}
	grid.SetOrientation(orient)
	return grid
}

func SetupLabel(text string) *gtk.Label {
	label, err := gtk.LabelNew(text)
	if err != nil {
		log.Fatal("Unable to create label:", err)
	}
    label.SetXAlign(0)
	return label
}

func SetupNotebook() *gtk.Notebook {
	nb, err := gtk.NotebookNew()
	if err != nil {
		log.Fatal("Unable to create notebook:", err)
	}
	return nb
}

func SetupPopupWindow(title string, width, height int) *gtk.Window {
	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("Unable to create window:", err)
	}
	win.SetTitle(title)
	win.Connect("destroy", func() {
		win.Close()
	})
	win.SetDefaultSize(width, height)
	win.SetPosition(gtk.WIN_POS_CENTER)
	return win
}

func SetupScrolledWindow() *gtk.ScrolledWindow {
	scrwin, err := gtk.ScrolledWindowNew(nil, nil)
	if err != nil {
		log.Fatal("Unable to create scrolled window:", err)
	}
	scrwin.SetPolicy(1,1)
	scrwin.SetHExpand(true)
	return scrwin
}

func SetupSearchBar() *gtk.SearchBar {
	sb, err := gtk.SearchBarNew()
	if err != nil {
		log.Fatal("unable to create searchbar:", err)
	}
	return sb
}

func SetupSearchEntry() *gtk.SearchEntry {
	se, err := gtk.SearchEntryNew()
	if err != nil {
		log.Fatal("unable to create searchentry:", err)
	}
	return se
}

func SetupToolbar() *gtk.Toolbar {
	tb, err := gtk.ToolbarNew()
	if err != nil {
		log.Fatal("unable to create toolbar:", err)
	}
	return tb
}

func SetupToolButtonIcon(iconName string) *gtk.ToolButton {
    image,err := gtk.ImageNewFromIconName(iconName, gtk.ICON_SIZE_BUTTON)
    if err != nil {
        log.Fatal("Unable to create image:", err)
    }
	btn, err := gtk.ToolButtonNew(image, "")
	if err != nil {
		log.Fatal("Unable to create button:", err)
	}
	return btn
}

func SetupToolButtonLabel(text string) *gtk.ToolButton {
	btn, err := gtk.ToolButtonNew(nil, text)
	if err != nil {
		log.Fatal("Unable to create button:", err)
	}
	return btn
}

func SetupWindow(title string) *gtk.Window {
	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("Unable to create window:", err)
	}
	win.SetTitle(title)
	win.Connect("destroy", func() {
		gtk.MainQuit()
	})
	win.SetDefaultSize(1000, 700)
	win.SetPosition(gtk.WIN_POS_CENTER)
	return win
}
