package model

import (
    "database/sql"
    "fmt"
    "log"
    "os"
    "path/filepath"
    "strings"

    "github.com/gchaincl/dotsql"
    _ "github.com/mattn/go-sqlite3"
)

type Database struct {
    Database *sql.DB
}

func NewDatabase() (*Database, bool) {
    fileExists := true
    if _, err := os.Stat("./database.db"); os.IsNotExist(err) {
        fileExists = false
    }
    db, err := sql.Open("sqlite3", "./database.db")
    if err != nil {
		log.Fatal("could not open the database: ", err)
	}
    DB := &Database{
        Database: db,
    }
    return DB, fileExists
}

func(database *Database) CreateDB() {
    os.Remove("./database.db")

    dot, err := dotsql.LoadFromFile("rolas.sql")
    if err != nil {
		log.Fatal("could not load rolas.sql: ", err)
	}

    CREATE := "create-"
    TABLE := "-table"

    setup := make([]string, 0)

    setup = append(setup, CREATE + "types-table")
    setup = append(setup, CREATE + "type0")
    setup = append(setup, CREATE + "type1")
    setup = append(setup, CREATE + "type2")
    setup = append(setup, CREATE + "performers" + TABLE)
    setup = append(setup, CREATE + "persons" + TABLE)
    setup = append(setup, CREATE + "groups" + TABLE)
    setup = append(setup, CREATE + "albums" + TABLE)
    setup = append(setup, CREATE + "rolas" + TABLE)
    setup = append(setup, CREATE + "in_group" + TABLE)

    for _, query := range setup {
        _, err = dot.Exec(database.Database, query)
        if err != nil {
            log.Fatal(err)
        }
    }
}

func (database *Database) LoadDB() {
    err := database.Database.Ping()
    if err != nil {
        log.Fatal("connection is dead:", err)
    }
}

func (database *Database) PerformerExists(performerName string) int64 {
    var id int64
    tx, err := database.Database.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare("SELECT id_performer FROM performers WHERE performers.name = ?")
	if err != nil {
		log.Fatal("could not prepare query: ", err)
	}
    defer stmt.Close()

    rows, err := stmt.Query(performerName)
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()
    for rows.Next() {
        err = rows.Scan(&id)
        if err != nil {
            log.Fatal(err)
        }
    }
    err = rows.Err()
    if err != nil {
        log.Fatal(err)
    }
    tx.Commit()
    return id
}

func (database *Database) AddPerformer(rola *Rola) int64 {
    idp := database.PerformerExists(rola.Artist())
    if idp > 0 {
        return idp
    }

    tx, err := database.Database.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare("INSERT INTO performers(id_type, name) " +
                            "SELECT ?, ? " +
                            "WHERE NOT EXISTS(SELECT 1 FROM performers WHERE name = ?)")
	if err != nil {
		log.Fatal("could not prepare insert: ", err)
	}
    defer stmt.Close()

    id, err := stmt.Exec(2, strings.TrimSpace(rola.Artist()), rola.Artist())
	if err != nil {
		log.Fatal(err)
	}
    tx.Commit()
    lastId, err := id.LastInsertId()
    if err != nil {
        log.Fatal("could not retrieve the last insert ID:", err)
    }
    return lastId
}

func (database *Database) AlbumExists(albumPath string) int64 {
    var id int64
    tx, err := database.Database.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare("SELECT id_album FROM albums WHERE albums.path = ?")
	if err != nil {
		log.Fatal("could not prepare insert: ", err)
	}
    defer stmt.Close()

    rows, err := stmt.Query(albumPath)
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()
    for rows.Next() {
        err = rows.Scan(&id)
        if err != nil {
            log.Fatal(err)
        }
    }
    err = rows.Err()
    if err != nil {
        log.Fatal(err)
    }
    tx.Commit()
    return id
}

func (database *Database) AddAlbum(rola *Rola) int64 {
    idalbum := database.AlbumExists(filepath.Dir(rola.Path()))
    if idalbum > 0 {
        return idalbum
    }

    tx, err := database.Database.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare("INSERT INTO albums (path, name, year) " +
                            "SELECT ?, ?, ? " +
                            "WHERE NOT EXISTS(SELECT 1 FROM albums WHERE path = ?)")
	if err != nil {
		log.Fatal("could not prepare query: ", err)
	}
    defer stmt.Close()

    id, err := stmt.Exec(filepath.Dir(rola.Path()), rola.Album(), rola.Year(), filepath.Dir(rola.Path()))
	if err != nil {
		log.Fatal(err)
	}
    tx.Commit()
    lastId, err := id.LastInsertId()
    if err != nil {
        log.Fatal("could not retrieve the last insert ID:", err)
    }
    return lastId
}

func (database *Database) AddRola(rola *Rola, idperformer, idalbum int64) int64 {
    tx, err := database.Database.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("INSERT INTO rolas (id_performer, id_album, path, title, track, year, genre) " +
                            "SELECT ?, ?, ?, ?, ?, ?, ? " +
                            "WHERE NOT EXISTS(SELECT 1 FROM rolas WHERE title = ? AND id_performer = ? AND id_album = ? AND genre = ?)")
	if err != nil {
		log.Fatal("could not prepare insert: ", err)
	}
    defer stmt.Close()

    result, err := stmt.Exec(idperformer, idalbum, rola.Path(), rola.Title(), rola.Track(), rola.Year(), rola.Genre(), rola.Title(), idperformer, idalbum, rola.Genre())
	if err != nil {
		log.Fatal("could not execute insert:", err)
	}
    rowsAdded, err := result.RowsAffected()
    if err != nil {
        log.Fatal("could not retrieve number of affected rows:", err)
    }
    tx.Commit()
    if rowsAdded > 0 {
        id, err := result.LastInsertId()
        if err != nil {
            log.Fatal("could not retrieve last inserted id:", err)
        }
        return id
    }
    return -1
}

func (database *Database) Display() {
    rows, err := database.Database.Query("SELECT id_performer, name FROM performers")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var name string
		err = rows.Scan(&id, &name)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(id, name)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

    rows, err = database.Database.Query("SELECT id_album, path, name, year FROM albums")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
        var path string
		var name string
        var year int
		err = rows.Scan(&id, &path, &name, &year)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(id, path, name, year)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

    rows, err = database.Database.Query("SELECT id_rola, performers.name, albums.name, rolas.path, rolas.title, rolas.track, rolas.year, rolas.genre FROM rolas INNER JOIN performers ON performers.id_performer = rolas.id_performer INNER JOIN albums ON albums.id_album = rolas.id_album")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var idr int
        var performer string
        var album string
        var path string
		var title string
        var track string
        var year int
        var genre string
		err = rows.Scan(&idr, &performer, &album, &path, &title, &track, &year, &genre)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(idr, title, performer, album, track, year, genre, path)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}