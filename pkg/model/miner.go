package model

import (
    "fmt"
    "log"
    "os"
    "os/user"
    "path/filepath"
    "strconv"
    "strings"
    "unicode"

    "github.com/bogem/id3v2"
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
    for _, path := range miner.paths {
        tag, err := id3v2.Open(path, id3v2.Options{Parse: true})
	    if err != nil {
 		    log.Fatal("error while opening mp3 file: ", path + " ", err)
 	    }
	    defer tag.Close()

        rola := NewRola()
        if tag.Artist() != "" {
            rola.SetArtist(tag.Artist())
        }
        if tag.Title() != "" {
            rola.SetTitle(tag.Title())
        }
        if tag.Album() != "" {
            rola.SetAlbum(tag.Album())
        }
        if tag.GetTextFrame("TRCK").Text != "" {
            trackString := tag.GetTextFrame("TRCK").Text
            f := func(c rune) bool {
                return !unicode.IsNumber(c)
            }
            fields := strings.FieldsFunc(trackString, f)
            if len(fields) > 0 {
                track, err := strconv.Atoi(fields[0])
                if err != nil {
                    fmt.Println(err)
                    log.Fatal("error while trying to cast track number into int: ", err)
                }
                rola.SetTrack(track)
            }
        }
        if tag.Year() != "" {
            yearString := tag.Year()
            f := func(c rune) bool {
                return !unicode.IsNumber(c)
            }
            fields := strings.FieldsFunc(yearString, f)
            if len(fields) > 0 {
                year, err := strconv.Atoi(fields[0])
                if err != nil {
                    fmt.Println(err)
                    log.Fatal("error while trying to cast year into int: ", err)
                }
                rola.SetYear(year)
            }
        }
        if tag.Genre() != "" {
            rola.SetGenre(tag.Genre())
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
