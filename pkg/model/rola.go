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

// NewRola creates a Rola with default values; text fields are "Unknown"
// and numeric fields are 0.
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

// Artist returns the performer of the Rola.
func (rola *Rola) Artist() string {
	return rola.artist
}

// Title returns the title of the Rola.
func (rola *Rola) Title() string {
	return rola.title
}

// Album returns the album where the Rola is included.
func (rola *Rola) Album() string {
	return rola.album
}

// Track returns the track number of the Rola as an int.
func (rola *Rola) Track() int {
	return rola.track
}

// Year returns the year of the Rola as an int.
func (rola *Rola) Year() int {
	return rola.year
}

// Genre returns the genre of the Rola as a string.
func (rola *Rola) Genre() string {
	return rola.genre
}

// Path returns the path of the song file where the Rola was mined.
func (rola *Rola) Path() string {
	return rola.path
}

// ID returns the ID assigned to the Rola by the database at insertion.
func (rola *Rola) ID() int64 {
	return rola.id
}

// SetArtist sets the Rola performer.
func (rola *Rola) SetArtist(artist string) {
	rola.artist = strings.TrimSpace(artist)
}

// SetTitle sets the Rola title.
func (rola *Rola) SetTitle(title string) {
	rola.title = strings.TrimSpace(title)
}

// SetAlbum sets the album title of the Rola.
func (rola *Rola) SetAlbum(album string) {
	rola.album = strings.TrimSpace(album)
}

// SetTrack sets the track number of the Rola. Should be an int.
func (rola *Rola) SetTrack(track int) {
	rola.track = track
}

// SetYear sets the year of the Rola. Should be an int.
func (rola *Rola) SetYear(year int) {
	rola.year = year
}

// SetGenre sets the genre of the Rola.
func (rola *Rola) SetGenre(genre string) {
	rola.genre = strings.TrimSpace(genre)
}

// SetPath sets the path of the file where the song represented by the Rola is.
func (rola *Rola) SetPath(path string) {
	rola.path = strings.TrimSpace(path)
}

// SetID sets the ID of the Rola. This value should not be changed unless
// the corresponding value changes in the Database.
func (rola *Rola) SetID(id int64) {
	rola.id = id
}
