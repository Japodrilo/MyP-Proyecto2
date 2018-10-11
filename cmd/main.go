package main

import (
    //"time"

    "github.com/gotk3/gotk3/gtk"
    "github.com/Japodrilo/MyP-Proyecto2/pkg/controller"
    //"github.com/Japodrilo/MyP-Proyecto2/pkg/model"
)

func main() {
    //miner := model.NewMiner()
    //miner.Traverse()
    //go miner.Extract()

    //database := model.NewDatabase()
    //database.CreateDB()

    //miner.Populate(database)

    //database.Display()

    gtk.Init(nil)

	controller.MainWindow()

	gtk.Main()
}
