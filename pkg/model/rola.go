package model

import (
    "strings"
)

/**
 * Una rola representa una canción lista con información suficente para ser
 * ingresada a la base de datos.   La información mínima que contiene es else {
 * nombre del intérprete, título de la rola, álbum donde aparece la rola,
 * número de la rola en el álbum, year de publicación del álbum, género de
 * la rola y path del archivo en disco.
 */
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