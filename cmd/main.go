package main

import (
    "github.com/Japodrilo/MyP-Proyecto2/pkg/modelo"
)

func main() {
    miner := model.NewMiner()
    miner.Traverse()
    go miner.Extract()

    database := model.NewDatabase()
    database.StartDB()

    for rola := range miner.Ore() {
        database.AddPerformer(rola)
        database.AddAlbum(rola)
        database.AddRola(rola)
    }

    database.Display()
}
