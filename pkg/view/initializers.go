package view

import (
	"log"

	"github.com/gotk3/gotk3/gtk"
)

// GetBufferEntry takes a gtk.SearchEntry as a parameter and
// returns its gtk.EntryBuffer. It includes error handling.
func GetBufferEntry(entry *gtk.SearchEntry) *gtk.EntryBuffer {
	buffer, err := entry.GetBuffer()
	if err != nil {
		log.Fatal("Unable to get buffer:", err)
	}
	return buffer
}

// GetLabelText takes a gtk.Label as a parameter and
// returns the text it contains. It includes error handling.
func GetLabelText(label *gtk.Label) string {
	text, err := label.GetText()
	if err != nil {
		log.Fatal("Unable to get text from label:", err)
	}
	return text
}

// GetTextEntry receives a gtk.Entry as a parameter and
// returns the text it contains. It includes error handling.
func GetTextEntry(entry *gtk.Entry) string {
	text, err := entry.GetText()
	if err != nil {
		log.Fatal("Unable to get text from buffer:", err)
	}
	return text
}

// GetTextEntryClean receives a gtk.SearchEntry as a parameter,
// returns the text it contains and cleans its buffer. It includes
// error handling.
func GetTextEntryClean(entry *gtk.SearchEntry) string {
	text, err := entry.GetText()
	if err != nil {
		log.Fatal("Unable to get text from buffer:", err)
	}
	buffer := GetBufferEntry(entry)
	buffer.DeleteText(0, -1)
	return text + "\n"
}

// GetTextSearchEntry receives a gtk.SearchEntry as a parameter and
// returns the text it contains. It includes error handling.
func GetTextSearchEntry(entry *gtk.SearchEntry) string {
	text, err := entry.GetText()
	if err != nil {
		log.Fatal("Unable to get text from buffer:", err)
	}
	return text
}

// SetupBox creates a new gtk.Box object, sets its 'Homogeneous' property
// to false, and returns it. It includes error handling.
func SetupBox() *gtk.Box {
	box, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		log.Fatal("Unable to create box:", err)
	}
	box.SetHomogeneous(false)
	return box
}

// SetupComboBoxText creates a new gtk.ComboBoxTex object
// and returns it. It includes error handling.
func SetupComboBoxText() *gtk.ComboBoxText {
	cb, err := gtk.ComboBoxTextNew()
	if err != nil {
		log.Fatal("Unable to create Combo Box:", err)
	}
	return cb
}

// SetupEntry creates a new gtk.Entry object and returns it.
// It includes error handling.
func SetupEntry() *gtk.Entry {
	entry, err := gtk.EntryNew()
	if err != nil {
		log.Fatal("Unable to create Entry:", err)
	}
	return entry
}

// SetupGrid creates a new gtk.Grid object with orientation specified
// by a parameter of type gtk.Orientation, and returns it.
// It includes error handling.
func SetupGrid(orient gtk.Orientation) *gtk.Grid {
	grid, err := gtk.GridNew()
	if err != nil {
		log.Fatal("Unable to create grid:", err)
	}
	grid.SetOrientation(orient)
	return grid
}

// SetupLabel creates a new gtk.Label object with text content given
// by the argument of the function, sets its XAlign property to 0, and
// returns it. It includes error handling.
func SetupLabel(text string) *gtk.Label {
	label, err := gtk.LabelNew(text)
	if err != nil {
		log.Fatal("Unable to create label:", err)
	}
	label.SetXAlign(0)
	return label
}

// SetupListBox creates a new gtk.ListBox object and returns it.
// It includes error handling.
func SetupListBox() *gtk.ListBox {
	lb, err := gtk.ListBoxNew()
	if err != nil {
		log.Fatal("Unable to create ListBox:", err)
	}
	return lb
}

// SetupListBoxRowLabel creates a new gtk.ListBoxRow object, creates a
// new label containing the text given as an argument of the function,
// adds the label to the listbox row, and returns it.
// It includes error handling.
func SetupListBoxRowLabel(text string) *gtk.ListBoxRow {
	lbr, err := gtk.ListBoxRowNew()
	if err != nil {
		log.Fatal("Unable to create List Box Row:", err)
	}
	label := SetupLabel(text)
	lbr.Add(label)
	return lbr
}

// SetupNotebook creates a new gtk.Notebook object and returns it.
// It includes error handling.
func SetupNotebook() *gtk.Notebook {
	nb, err := gtk.NotebookNew()
	if err != nil {
		log.Fatal("Unable to create notebook:", err)
	}
	return nb
}

// SetupPopupWindow creates a new gtk.Window object with title,
// width and height given by its parameters; context its "destroy" signal
// to the window Close() function, and returns it. It includes error
// handling.
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

// SetupScrolled Window creates a new gtk.ScrolledWindow object,
// sets its Policy to (1,1), HExpand to true, and returns it.
// It includes error handling.
func SetupScrolledWindow() *gtk.ScrolledWindow {
	scrwin, err := gtk.ScrolledWindowNew(nil, nil)
	if err != nil {
		log.Fatal("Unable to create scrolled window:", err)
	}
	scrwin.SetPolicy(1, 1)
	scrwin.SetHExpand(true)
	return scrwin
}

// SetupSearchBar creates a new gtk.SearchBar object and returns it.
// It includes error handling.
func SetupSearchBar() *gtk.SearchBar {
	sb, err := gtk.SearchBarNew()
	if err != nil {
		log.Fatal("unable to create searchbar:", err)
	}
	return sb
}

// SetupSearchEntry creates a new gtk.SearchEntry object and returns it.
// It includes error handling.
func SetupSearchEntry() *gtk.SearchEntry {
	se, err := gtk.SearchEntryNew()
	if err != nil {
		log.Fatal("unable to create searchentry:", err)
	}
	return se
}

// SetupToolbar creates a new gtk.Toolbar object and returns it.
// It includes error handling.
func SetupToolbar() *gtk.Toolbar {
	tb, err := gtk.ToolbarNew()
	if err != nil {
		log.Fatal("unable to create toolbar:", err)
	}
	return tb
}

// SetupToolButtonIcon creates a new gtk.ToolButton object with image
// determined by de icon name taken as argument, and returns it.
// It includes error handling.
func SetupToolButtonIcon(iconName string) *gtk.ToolButton {
	image, err := gtk.ImageNewFromIconName(iconName, gtk.ICON_SIZE_BUTTON)
	if err != nil {
		log.Fatal("Unable to create image:", err)
	}
	btn, err := gtk.ToolButtonNew(image, "")
	if err != nil {
		log.Fatal("Unable to create button:", err)
	}
	return btn
}

// SetupToolButtonLabel creates a new gtk.ToolButton object with
// text label given by the entry of the function, and returns it.
// It includes error handling.
func SetupToolButtonLabel(text string) *gtk.ToolButton {
	btn, err := gtk.ToolButtonNew(nil, text)
	if err != nil {
		log.Fatal("Unable to create button:", err)
	}
	return btn
}

// SetupWindow creates a new gtk.Window object with fixed size
// (1000 by 700), connects its "destroy" signal to gtk.MainQuit()
// funcion, and returns it.   It includes error handling.
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
