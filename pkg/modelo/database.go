package model

import (
    "database/sql"
    "fmt"
    "log"
    "os"
    "path/filepath"

    "github.com/gchaincl/dotsql"
    _ "github.com/mattn/go-sqlite3"
)

type Database struct {
    rolas    map[string]*Rola
    database *sql.DB
}

func NewDatabase() *Database {
    return &Database{
        rolas: make(map[string]*Rola),
    }
}

func(database *Database) StartDB() {
    os.Remove("./database.db")

    db, err := sql.Open("sqlite3", "./database.db")
    if err != nil {
		log.Fatal("could not open the database: ", err)
	}

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
        _, err = dot.Exec(db, query)
        if err != nil {
            log.Fatal(err)
        }
    }
    database.database = db
}

func (database *Database) AddPerformer(rola *Rola) {
    tx, err := database.database.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("INSERT into performers (id_type, name)" +
                            "VALUES (?, ?)")
	if err != nil {
		log.Fatal("could not prepare query: ", err)
	}
    defer stmt.Close()

    _, err = stmt.Exec(2, rola.Artist())
	if err != nil {
		log.Fatal(err)
	}
    tx.Commit()
}

func (database *Database) AddAlbum(rola *Rola) {
    tx, err := database.database.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("INSERT into albums (path, name, year)" +
                            "VALUES (?, ?, ?)")
	if err != nil {
		log.Fatal("could not prepare query: ", err)
	}
    defer stmt.Close()

    _, err = stmt.Exec(filepath.Dir(rola.Path()), rola.Album(), rola.Year())
	if err != nil {
		log.Fatal(err)
	}
    tx.Commit()
}

func (database *Database) AddRola(rola *Rola) {
    tx, err := database.database.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("INSERT into rolas (id_performer, id_album, path, title, track, year, genre)" +
                            "VALUES (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal("could not prepare query: ", err)
	}
    defer stmt.Close()

    _, err = stmt.Exec(0, 0, rola.Path(), rola.Title(), rola.Track(), rola.Year(), rola.Genre())
	if err != nil {
		log.Fatal(err)
	}
    tx.Commit()
}

func (database *Database) Display() {
    rows, err := database.database.Query("SELECT id_performer, name FROM performers")
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

    rows, err = database.database.Query("SELECT id_album, path, name, year FROM albums")
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

    rows, err = database.database.Query("SELECT id_rola, id_performer, id_album, path, title, track, year, genre FROM rolas")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var idr int
        var idp int
        var ida int
        var path string
		var title string
        var track string
        var year int
        var genre string
		err = rows.Scan(&idr, &idp, &ida, &path, &title, &track, &year, &genre)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(idr, idp, ida, path, title, track, year, genre)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}
