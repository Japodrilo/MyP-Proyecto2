package model

import (
    "strings"
)

// A Rola represents a song, it contains the information present in
// various frames from the id3v2 tag, namely, artist, title, album
// track number, year, genre, and additionally, the path of the song
// file, and the id assigned by the database to the song.
type Rola struct {
    artist string
    title  string
    album  string
    track  int
    year   int
    genre  string
    path   string
    id     int64
}

func NewRola() *Rola {
    initial := "Unknown"
    return &Rola{
        artist: initial,
        title:  initial,
        album:  initial,
        track:  0,
        year:   2018,
        genre:  initial,
        path:   initial,
        id:     0,
    }
}

func (rola *Rola) Artist() string {
    return rola.artist
}

func (rola *Rola) Title() string {
    return rola.title
}

func (rola *Rola) Album() string {
    return rola.album
}

func (rola *Rola) Track() int {
    return rola.track
}

func (rola *Rola) Year() int {
    return rola.year
}

func (rola *Rola) Genre() string {
    return rola.genre
}

func (rola *Rola) Path() string {
    return rola.path
}

func (rola *Rola) ID() int64 {
    return rola.id
}

func (rola *Rola) SetArtist(artist string) {
    rola.artist = strings.TrimSpace(artist)
}

func (rola *Rola) SetTitle(title string) {
    rola.title = strings.TrimSpace(title)
}

func (rola *Rola) SetAlbum(album string) {
    rola.album = strings.TrimSpace(album)
}

func (rola *Rola) SetTrack(track int) {
    rola.track = track
}

func (rola *Rola) SetYear(year int) {
    rola.year = year
}

func (rola *Rola) SetGenre(genre string) {
    rola.genre = strings.TrimSpace(genre)
}

func (rola *Rola) SetPath(path string) {
    rola.path = strings.TrimSpace(path)
}

func (rola *Rola) SetID(id int64) {
    rola.id = id
}
