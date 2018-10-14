package model

import (
    "fmt"
    "log"
    "os"
    "os/user"
    "path/filepath"
    "strings"

    "github.com/dhowden/tag"
)

/**
 * This class represents a miner that searches for mp3 files in the
 * /home/user/Music directory, along the file tree, gathers their
 * information, and puts it in a Rola object, which is then loaded
 * into a channel for external use.
 */
type Miner struct {
    paths     []string
    ore       chan *Rola
    TrackList chan *Rola
}

func NewMiner() *Miner {
    return &Miner{
        paths: make([]string, 0),
    }
}

func (miner *Miner) Traverse() {
    home, err := user.Current()
    if err != nil {
        log.Fatal("could not retrieve the current user:", err)
    }
	homePath := home.HomeDir
    start := homePath + "/Music"

    err = filepath.Walk(start, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            fmt.Printf("failure accessing the path %q: %v\n", path, err)
            return err
        }
        if !info.IsDir() && strings.HasSuffix(info.Name(), ".mp3") {
            miner.paths = append(miner.paths, path)
        }
        return nil
    })
    if err != nil {
        fmt.Printf("error walking the path %q: %v\n", start, err)
    }
}


func (miner *Miner) Extract() {
    miner.ore = make(chan *Rola)
    genreConverter := GetGenre()
    for _, path := range miner.paths {
        file, err := os.Open(path)
        if err != nil {
            log.Fatal("could not open file: ", path, err)
        }
        metadata, err := tag.ReadFrom(file)
        if err != nil {
            log.Fatal("could not read the tag:", err)
        }

        rola := NewRola()
        if metadata.Artist() != "" {
            rola.SetArtist(metadata.Artist())
        }
        if metadata.Title() != "" {
            rola.SetTitle(metadata.Title())
        }
        if metadata.Album() != "" {
            rola.SetAlbum(metadata.Album())
        }
        track, _ := metadata.Track()
        if track != 0 {
            rola.SetTrack(track)
        }
        if metadata.Year() != 0 {
            rola.SetYear(metadata.Year())
        }
        if metadata.Genre() != "" {
            rola.SetGenre(genreConverter.Get(metadata.Genre()))
        }
        rola.SetPath(path)
        miner.ore <- rola
    }
    close(miner.ore)
}

func (miner *Miner) Populate(database *Database) {
    miner.TrackList = make(chan *Rola)
    for rola := range miner.ore {
        idperformer := database.AddPerformer(rola)
        idalbum := database.AddAlbum(rola)
        id := database.AddRola(rola, idperformer, idalbum)
        if id > 0 {
            rola.SetID(id)
            miner.TrackList <- rola
        }
    }
    close(miner.TrackList)
}
