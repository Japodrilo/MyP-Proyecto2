package main

import (
    "github.com/gotk3/gotk3/gtk"
    "github.com/Japodrilo/MyP-Proyecto2/pkg/controller"
)

func main() {
    gtk.Init(nil)
    controller.MainWindow()
    gtk.Main()
}
