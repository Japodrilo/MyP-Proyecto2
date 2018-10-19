package model

import (
    "log"
    "os"
    "os/user"
    "path/filepath"
    "strings"

    "github.com/dhowden/tag"
)


// A Miner searches for mp3 files in the /home/user/Music directory
// along the file tree, gathers their information, and puts it in a
// Rola object, which is then loaded into a channel for external use.
type Miner struct {
    paths     []string
    ore       chan *Rola
    TrackList chan *Rola
}

// NewMiner returns a new Miner with an empty paths slice.
func NewMiner() *Miner {
    return &Miner{
        paths: make([]string, 0),
    }
}

// Traverse walks the file tree looking for mp3 files and saving their
// paths into the paths slice.
func (miner *Miner) Traverse() {
    home, err := user.Current()
    if err != nil {
        log.Fatal("could not retrieve the current user:", err)
    }
	homePath := home.HomeDir
    start := homePath + "/Music"

    err = filepath.Walk(start, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            log.Fatal("failure accessing the path: ", err)
            return err
        }
        if !info.IsDir() && strings.HasSuffix(info.Name(), ".mp3") {
            miner.paths = append(miner.paths, path)
        }
        return nil
    })
    if err != nil {
        log.Fatal("error walking the path: ", err)
    }
}

// Extract traverses the paths slice, opens each of the files whose
// paths are in the slice, reads the ID3v2 tag, saves the information
// into a new Rola, and puts it in the ore channel of the miner.
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

// Populate takes the Rolas in the ore channel of the miner,
// adds them to the database, and if it was a new Rola, it is
// put in the TrackList channel.
// TODO: Maybe this method should be in the controller package.
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
