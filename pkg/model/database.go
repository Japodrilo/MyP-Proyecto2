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

func (database *Database) PrepareStatement(statement string) (*sql.Tx, *sql.Stmt) {
    tx, err := database.Database.Begin()
	if err != nil {
		log.Fatal("could not begin transaction: ", err)
	}
	stmt, err := tx.Prepare(statement)
	if err != nil {
		log.Fatal("could not prepare statement: ", err)
	}
    return tx, stmt
}

func (database *Database) PreparedQuery(statement string, args ...interface{}) (*sql.Tx, *sql.Stmt, *sql.Rows) {
    tx, stmt := database.PrepareStatement(statement)
    rows, err := stmt.Query(args...)
    if err != nil {
        log.Fatal("could not performe query: ", err )
    }
    return tx, stmt, rows
}

func (database *Database) AddAlbum(rola *Rola) int64 {
    idalbum := database.ExistsAlbum(filepath.Dir(rola.Path()))
    if idalbum > 0 {
        return idalbum
    }

    stmtStr := "INSERT INTO albums (path, name, year) SELECT ?, ?, ? " +
               "WHERE NOT EXISTS(SELECT 1 FROM albums WHERE path = ?)"

    tx, stmt := database.PrepareStatement(stmtStr)
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

func (database *Database) AddGroup(groupName, start, end string) {
    stmtStr := "INSERT INTO groups (" +
               " name, " +
               " start_date, " +
               " end_date) "  +
               "SELECT ?, ?, ?"

    tx, stmt := database.PrepareStatement(stmtStr)
    defer stmt.Close()

    _, err := stmt.Exec(groupName, start, end)
	if err != nil {
		log.Fatal(err)
	}
    tx.Commit()
}

func (database *Database) AddPerformer(rola *Rola) int64 {
    idp := database.ExistsPerformer(rola.Artist())
    if idp > 0 {
        return idp
    }

    stmtStr := "INSERT INTO performers(id_type, name) SELECT ?, ? " +
               "WHERE NOT EXISTS(SELECT 1 FROM performers WHERE name = ?)"

    tx, stmt := database.PrepareStatement(stmtStr)
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

func (database *Database) AddPerson(stageName, realName, birth, death string) {
    stmtStr := "INSERT INTO persons (" +
               " stage_name, " +
               " real_name, " +
               " birth_date, " +
               " death_date) "  +
               "SELECT ?, ?, ?, ?"

    tx, stmt := database.PrepareStatement(stmtStr)
    defer stmt.Close()

    _, err := stmt.Exec(stageName, realName, birth, death)
	if err != nil {
		log.Fatal(err)
	}
    tx.Commit()
}

func (database *Database) AddRola(rola *Rola, idperformer, idalbum int64) int64 {
    stmtStr := "INSERT INTO rolas (id_performer, id_album, path, title, " +
               "track, year, genre) SELECT ?, ?, ?, ?, ?, ?, ? WHERE NOT " +
               "EXISTS(SELECT 1 FROM rolas WHERE title = ? " +
               "AND id_performer = ? AND id_album = ? AND genre = ?)"

    tx, stmt := database.PrepareStatement(stmtStr)
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

func (database *Database) ExistsAlbum(albumPath string) int64 {
    stmtStr := "SELECT id_album FROM albums WHERE albums.path = ? LIMIT 1"
    tx, stmt, rows := database.PreparedQuery(stmtStr, albumPath)
    defer stmt.Close()
    defer rows.Close()

    var id int64
    for rows.Next() {
        err := rows.Scan(&id)
        if err != nil {
            log.Fatal(err)
        }
    }
    err := rows.Err()
    if err != nil {
        log.Fatal(err)
    }
    tx.Commit()
    return id
}

func (database *Database) ExistsGroup(groupName string) int64 {
    stmtStr := "SELECT " +
               " id_group " +
               "FROM " +
               " groups " +
               "WHERE " +
               " groups.name = ? " +
               "LIMIT 1"
    tx, stmt, rows := database.PreparedQuery(stmtStr, groupName)
    defer stmt.Close()
    defer rows.Close()

    var id int64
    for rows.Next() {
        err := rows.Scan(&id)
        if err != nil {
            log.Fatal(err)
        }
    }
    err := rows.Err()
    if err != nil {
        log.Fatal(err)
    }
    tx.Commit()
    return id
}

func (database *Database) ExistsPerformer(performerName string) int64 {
    stmtStr := "SELECT id_performer FROM performers WHERE performers.name = ? LIMIT 1"
    tx, stmt, rows := database.PreparedQuery(stmtStr, performerName)
    defer stmt.Close()
    defer rows.Close()

    var id int64
    for rows.Next() {
        err := rows.Scan(&id)
        if err != nil {
            log.Fatal(err)
        }
    }
    err := rows.Err()
    if err != nil {
        log.Fatal(err)
    }
    tx.Commit()
    return id
}

func (database *Database) ExistsPerson(stageName, realName string) int64 {
    stmtStr := "SELECT " +
               " id_person " +
               "FROM " +
               " persons " +
               "WHERE " +
               "persons.stage_name = ? OR persons.real_name = ? " +
               "LIMIT 1"
    tx, stmt, rows := database.PreparedQuery(stmtStr, stageName, realName)
    defer stmt.Close()
    defer rows.Close()

    var id int64
    for rows.Next() {
        err := rows.Scan(&id)
        if err != nil {
            log.Fatal(err)
        }
    }
    err := rows.Err()
    if err != nil {
        log.Fatal(err)
    }
    tx.Commit()
    return id
}

func (database *Database) QueryGroup(groupID int64) (string, string, string){
    stmtStr := "SELECT " +
               " name, " +
               " start_date, " +
               " end_date " +
               "FROM " +
               " groups " +
               "WHERE " +
               " groups.id_group = ?"

    tx, stmt := database.PrepareStatement(stmtStr)
    defer stmt.Close()

    rows, err := stmt.Query(groupID)
    if err != nil {
        log.Fatal("could not execute query: ", err)
    }
    defer rows.Close()

    var name string
    var start string
    var end string
    for rows.Next() {
        err = rows.Scan(&name, &start, &end)
        if err != nil {
            log.Fatal(err)
        }
    }
    err = rows.Err()
    if err != nil {
        log.Fatal(err)
    }
    tx.Commit()
    return name, start, end
}

func (database *Database) QueryPerformerType(name string) (int64, int) {
    stmtStr := "SELECT " +
               " id_performer, " +
               " performers.id_type " +
               "FROM " +
               " performers " +
               "INNER JOIN types ON types.id_type = performers.id_type " +
               "WHERE " +
               " performers.name = ?"

    tx, stmt := database.PrepareStatement(stmtStr)
    defer stmt.Close()

    rows, err := stmt.Query(name)
    if err != nil {
        log.Fatal("could not execute query: ", err)
    }
    defer rows.Close()

    var performerID int64
    var performerType int
    for rows.Next() {
        err = rows.Scan(&performerID, &performerType)
        if err != nil {
            log.Fatal(err)
        }
    }
    err = rows.Err()
    if err != nil {
        log.Fatal(err)
    }
    tx.Commit()
    return performerID, performerType
}

func (database *Database) QueryPerson(personID int64) (string, string, string, string){
    stmtStr := "SELECT " +
               " stage_name, " +
               " real_name, " +
               " birth_date, " +
               " death_date " +
               "FROM " +
               " persons " +
               "WHERE " +
               " persons.id_person = ?"

    tx, stmt := database.PrepareStatement(stmtStr)
    defer stmt.Close()

    rows, err := stmt.Query(personID)
    if err != nil {
        log.Fatal("could not execute query: ", err)
    }
    defer rows.Close()

    var stageName string
    var realName string
    var birth string
    var death string
    for rows.Next() {
        err = rows.Scan(&stageName, &realName, &birth, &death)
        if err != nil {
            log.Fatal(err)
        }
    }
    err = rows.Err()
    if err != nil {
        log.Fatal(err)
    }
    tx.Commit()
    return stageName, realName, birth, death
}

func (database *Database) UpdateGroup(name, start, end string, groupID int64) {
    stmtStr := "UPDATE groups " +
               "SET name = ?, " +
               "    start_date = ?, " +
               "    end_date = ? " +
               "WHERE id_group = ?"

    tx, stmt := database.PrepareStatement(stmtStr)
    defer stmt.Close()

    _, err := stmt.Exec(name, start, end, groupID)
	if err != nil {
		log.Fatal(err)
	}
    tx.Commit()
}

func (database *Database) UpdatePerformerType(performerID int64, performerType int) {
    stmtStr := "UPDATE performers " +
               "SET id_type = ? " +
               "WHERE id_performer = ?"

    tx, stmt := database.PrepareStatement(stmtStr)
    defer stmt.Close()

    _, err := stmt.Exec(performerType, performerID)
	if err != nil {
		log.Fatal("could not execute update: ", err)
	}

    tx.Commit()
}

func (database *Database) UpdatePerson(stageName, realName, birth, death string, personID int64) {
    stmtStr := "UPDATE persons " +
               "SET stage_name = ?, " +
               "    real_name = ?, " +
               "    birth_date = ?, " +
               "    death_date = ? " +
               "WHERE id_person = ?"

    tx, stmt := database.PrepareStatement(stmtStr)
    defer stmt.Close()

    _, err := stmt.Exec(stageName, realName, birth, death, personID)
	if err != nil {
		log.Fatal(err)
	}
    tx.Commit()
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
