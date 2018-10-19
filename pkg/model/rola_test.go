package model

import (
	"testing"
)

func TestNewRola(t *testing.T) {
	rola := NewRola()
	expecting := "Unknown"
	if rola.artist != expecting {
		t.Errorf("expecting %v, received %v", expecting, rola.artist)
	}
	if rola.title != expecting {
		t.Errorf("expecting %v, received %v", expecting, rola.title)
	}
	if rola.album != expecting {
		t.Errorf("expecting %v, received %v", expecting, rola.album)
	}
	if rola.track != 0 {
		t.Errorf("expecting %v, received %v", expecting, rola.track)
	}
	if rola.year != 2018 {
		t.Errorf("expecting %v, received %v", expecting, rola.year)
	}
	if rola.genre != expecting {
		t.Errorf("expecting %v, received %v", expecting, rola.genre)
	}
	if rola.artist != expecting {
		t.Errorf("expecting %v, received %v", expecting, rola.path)
	}
}

func TestGetters(t *testing.T) {
	rola := NewRola()
	expecting := "Unknown"
	if rola.Artist() != expecting {
		t.Errorf("expecting %v, received %v", expecting, rola.Artist())
	}
	if rola.Title() != expecting {
		t.Errorf("expecting %v, received %v", expecting, rola.Title())
	}
	if rola.Album() != expecting {
		t.Errorf("expecting %v, received %v", expecting, rola.Album())
	}
	if rola.Track() != 0 {
		t.Errorf("expecting %v, received %v", expecting, rola.Track())
	}
	if rola.Year() != 2018 {
		t.Errorf("expecting %v, received %v", expecting, rola.Year())
	}
	if rola.Genre() != expecting {
		t.Errorf("expecting %v, received %v", expecting, rola.Genre())
	}
	if rola.Path() != expecting {
		t.Errorf("expecting %v, received %v", expecting, rola.Path())
	}
}

func TestSetters(t *testing.T) {
	rola := NewRola()
	artist := "Mark Ronson"
	title := "Lose It (In the End)"
	album := "Record Collection"
	track := 2
	year := 2010
	genre := "Alternative"
	path := "/Music/Mark Ronson/Record Collection/"

	rola.SetArtist(artist)
	rola.SetTitle(title)
	rola.SetAlbum(album)
	rola.SetTrack(track)
	rola.SetYear(year)
	rola.SetGenre(genre)
	rola.SetPath(path)

	if rola.Artist() != artist {
		t.Errorf("expecting %v, received %v", artist, rola.Artist())
	}
	if rola.Title() != title {
		t.Errorf("expecting %v, received %v", title, rola.Title())
	}
	if rola.Album() != album {
		t.Errorf("expecting %v, received %v", album, rola.Album())
	}
	if rola.Track() != track {
		t.Errorf("expecting %v, received %v", track, rola.Track())
	}
	if rola.Year() != year {
		t.Errorf("expecting %v, received %v", year, rola.Year())
	}
	if rola.Genre() != genre {
		t.Errorf("expecting %v, received %v", genre, rola.Genre())
	}
	if rola.Path() != path {
		t.Errorf("expecting %v, received %v", path, rola.Path())
	}
}
