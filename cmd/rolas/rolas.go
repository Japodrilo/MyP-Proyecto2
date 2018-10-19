package main

// #cgo LDFLAGS: -lX11
// #include <X11/Xlib.h>
import "C"

import (
	"github.com/Japodrilo/MyP-Proyecto2/pkg/controller"
	"github.com/gotk3/gotk3/gtk"
)

func main() {
	C.XInitThreads()
	gtk.Init(nil)
	controller.MainWindow()
	gtk.Main()
}
