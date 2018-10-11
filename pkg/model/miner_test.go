package model

import (
    "testing"
)

func TestTraverse(t *testing.T) {
    miner := NewMiner()
    miner.Traverse()
    if (miner.paths[0] != "/home/cesar/Music/Camelia Jordana/Camelia Jordana - Non Non Non.mp3") {
        t.Errorf("expecting %v, received %v", "/home/cesar/Music/Camelia Jordana/Camelia Jordana - Non Non Non.mp3", miner.paths[0])
    }
    if len(miner.paths) != 4 {
        t.Errorf("expecting %v, received %v", 4, len(miner.paths))
    }
}

func TestExtract(t *testing.T) {
    miner := NewMiner()
    paths := make([]string, 0)
    miner.Traverse()
    go miner.Extract()
    go func() {
        select{
        case rola, ok := <- miner.ore:
            if !ok {
                return
            }
            paths = append(paths, rola.Path())
        }
    }()
    for i, path := range paths {
        if path != miner.paths[i] {
            t.Errorf("expecting %v, received %v", miner.paths[i], path)
        }
    }
}

func TestExtract2(t *testing.T) {
    miner := NewMiner()
    paths := make([]string, 0)
    miner.Traverse()
    go miner.Extract()
    for rola := range miner.Ore() {
        paths = append(paths, rola.Path())
    }
    for i, path := range paths {
        if path != miner.paths[i] {
            t.Errorf("expecting %v, received %v", miner.paths[i], path)
        }
    }
}
