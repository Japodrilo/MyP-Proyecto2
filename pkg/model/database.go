// Package model of the application, mainly the database, miner and parser.
package model

import (
	"database/sql"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/gchaincl/dotsql"
	// SQLite driver
	_ "github.com/mattn/go-sqlite3"
)

// A Database is the intermediary between the sql database and
// the rest of the model and the view.
type Database struct {
	Database *sql.DB
	cache    string
}

// NewDatabase creates a new Database object, checking whether a .db
// file already exists and opening a connection to it, or creating a
// new one.   The database is saved in the file "~/.cache/rolas/rolas.db"
func NewDatabase() (*Database, bool) {
	home, err := user.Current()
	if err != nil {
		log.Fatal("could not retrieve the current user:", err)
	}
	cache := home.HomeDir + "/.cache/rolas"
	os.Mkdir(cache, 0700)

	fileExists := true
	if _, err := os.Stat(cache + "/rolas.db"); os.IsNotExist(err) {
		fileExists = false
	}
	db, err := sql.Open("sqlite3", cache+"/rolas.db")
	if err != nil {
		log.Fatal("could not open the database: ", err)
	}
	DB := &Database{
		Database: db,
		cache:    cache,
	}
	return DB, fileExists
}

// AddAlbum takes a Rola as a parameter, adds its album to the database
// and returns the ID number of the album in the database.   If the album
// was already in the database, this method does nothing and returns the
// ID of the album in the database.
func (database *Database) AddAlbum(rola *Rola) int64 {
	idalbum := database.ExistsAlbum(filepath.Dir(rola.Path()), rola.Album())
	if idalbum > 0 {
		return idalbum
	}

	stmtStr := `INSERT
                INTO albums (
                  path,
                  name,
                  year)
                SELECT ?, ?, ?
                WHERE NOT EXISTS
                (SELECT 1 FROM albums WHERE path = ? AND name = ?)`

	tx, stmt := database.PrepareStatement(stmtStr)
	defer stmt.Close()

	id, err := stmt.Exec(filepath.Dir(rola.Path()), rola.Album(), rola.Year(), filepath.Dir(rola.Path()), rola.Album())
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

// AddGroup takes a Rola as a parameter, adds its performer, which has been
// already verified to be a group, to the database, and returns the ID
// of the group in the database. If the group is already in the database, this
// method does nothing and returns the group ID in the database.
func (database *Database) AddGroup(groupName, start, end string) int64 {
	stmtStr := `INSERT INTO groups (
                 name,
                 start_date,
                 end_date)
                SELECT ?, ?, ?`

	tx, stmt := database.PrepareStatement(stmtStr)
	defer stmt.Close()

	id, err := stmt.Exec(groupName, start, end)
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

// AddPerformer takes a Rola as a parameter, adds its performer to the
// database, and returns the ID number of the performer in the database. The
// performer is added with type 2 (not known if it is a person of a group).
// If the performer is already in the database, this method does nothing
// and returns the performer ID in the database.
func (database *Database) AddPerformer(rola *Rola) int64 {
	idp := database.ExistsPerformer(rola.Artist())
	if idp > 0 {
		return idp
	}

	stmtStr := `INSERT
                INTO performers (
                  id_type,
                  name)
                SELECT ?, ?
                WHERE NOT EXISTS
                (SELECT 1 FROM performers WHERE name = ?)`

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

// AddPerson takes a Rola as a parameter, adds its performer, which has been
// already verified to be a person, to the database, and returns the ID
// of the person in the database. If the person is already in the database,
// this method does nothing and returns the person ID in the database.
func (database *Database) AddPerson(stageName, realName, birth, death string) {
	stmtStr := `INSERT INTO persons (
                  stage_name,
                  real_name,
                  birth_date,
                  death_date)
                SELECT ?, ?, ?, ?`

	tx, stmt := database.PrepareStatement(stmtStr)
	defer stmt.Close()

	_, err := stmt.Exec(stageName, realName, birth, death)
	if err != nil {
		log.Fatal(err)
	}
	tx.Commit()
}

// AddPersonToGroup receives the ID of a person and group in the database,
// respectively, and adds them to the in_group table, i.e., adds the person
// to the group.
func (database *Database) AddPersonToGroup(personID, groupID int64) {
	stmtStr := "INSERT INTO in_group (" +
		" id_person, " +
		" id_group) " +
		"SELECT ?, ?"

	tx, stmt := database.PrepareStatement(stmtStr)
	defer stmt.Close()

	_, err := stmt.Exec(personID, groupID)
	if err != nil {
		log.Fatal(err)
	}
	tx.Commit()
}

// AddRola takes a Rola and the IDs of the performer and album of the Rola
// as parameters, and attempts to add the Rola to the database.   If it was
// already in the database, it does nothing and returns -1.   Otherwise it
// returns the ID asigned to the Rola by the database.
func (database *Database) AddRola(rola *Rola, idperformer, idalbum int64) int64 {
	stmtStr := `INSERT
                INTO rolas (
                  id_performer,
                  id_album,
                  path,
                  title,
                  track,
                  year,
                  genre)
                SELECT ?, ?, ?, ?, ?, ?, ?
                WHERE NOT EXISTS
                (SELECT 1 FROM rolas WHERE (title = ?
                  AND id_performer = ?
                  AND id_album = ?
                  AND genre = ?)
				  OR path = ?)`

	tx, stmt := database.PrepareStatement(stmtStr)
	defer stmt.Close()

	result, err := stmt.Exec(idperformer, idalbum, rola.Path(), rola.Title(), rola.Track(), rola.Year(), rola.Genre(), rola.Title(), idperformer, idalbum, rola.Genre(), rola.Path())
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

// AllGroups queries the database and returns a map whose keys are the names
// of all the groups in the database, and the corresponding values are the
// IDs of the groups.
func (database *Database) AllGroups() map[string]int64 {
	groups := make(map[string]int64)
	rows, err := database.Database.Query("SELECT id_group, name FROM groups")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int64
		var name string
		err = rows.Scan(&id, &name)
		if err != nil {
			log.Fatal(err)
		}
		groups[name] = id
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return groups
}

// AllPersons queries the database and returns a map whose keys are the names
// of all the persons in the database, and the corresponding values are the
// IDs of the persons.
func (database *Database) AllPersons() map[string]int64 {
	persons := make(map[string]int64)
	rows, err := database.Database.Query("SELECT id_person, stage_name FROM persons")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int64
		var name string
		err = rows.Scan(&id, &name)
		if err != nil {
			log.Fatal(err)
		}
		persons[name] = id
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return persons
}

// CreateDB creates the tables specified in the rolas.sql file.
func (database *Database) CreateDB() {
	home, err := user.Current()
	if err != nil {
		log.Fatal("could not retrieve the current user:", err)
	}
	cache := home.HomeDir + "/.cache/rolas"
	fileExists := true
	if _, err := os.Stat(cache + "/rolas.sql"); os.IsNotExist(err) {
		fileExists = false
	}
	if !fileExists {
		RestoreAsset(cache, "rolas.sql")
	}

	dot, err := dotsql.LoadFromFile(cache + "/data/rolas.sql")
	if err != nil {
		log.Fatal("could not load rolas.sql: ", err)
	}

	CREATE := "create-"
	TABLE := "-table"

	setup := make([]string, 0)

	setup = append(setup, CREATE+"types-table")
	setup = append(setup, CREATE+"type0")
	setup = append(setup, CREATE+"type1")
	setup = append(setup, CREATE+"type2")
	setup = append(setup, CREATE+"performers"+TABLE)
	setup = append(setup, CREATE+"persons"+TABLE)
	setup = append(setup, CREATE+"groups"+TABLE)
	setup = append(setup, CREATE+"albums"+TABLE)
	setup = append(setup, CREATE+"rolas"+TABLE)
	setup = append(setup, CREATE+"in_group"+TABLE)

	for _, query := range setup {
		_, err = dot.Exec(database.Database, query)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// ExistsAlbum takes an album's path and name, and returns the album ID in
// the database, or 0 if the album is not in the database.
func (database *Database) ExistsAlbum(albumPath, name string) int64 {
	stmtStr := "SELECT id_album FROM albums WHERE albums.path = ? AND albums.name = ? LIMIT 1"
	tx, stmt, rows := database.PreparedQuery(stmtStr, albumPath, name)
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

// ExistsGroup takes a group's name as an argument and returns the group
// ID in the database, or 0 if the group is not in the database.
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

// ExistsPerformer takes a performer's name as an argument and returns the
// performer ID in the database, or 0 if the performer is not in the database.
func (database *Database) ExistsPerformer(performerName string) int64 {
	stmtStr := `SELECT
                  id_performer
                FROM performers
                WHERE performers.name = ?
                LIMIT 1`
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

// ExistsPerson takes a person's name and returns the person ID in the
// database, or 0 if the album is not in the database.
func (database *Database) ExistsPerson(stageName string) int64 {
	stmtStr := "SELECT " +
		" id_person " +
		"FROM " +
		" persons " +
		"WHERE " +
		"persons.stage_name = ? " +
		"LIMIT 1"
	tx, stmt, rows := database.PreparedQuery(stmtStr, stageName)
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

// LoadDB pings the database to verify if the connection is active.
func (database *Database) LoadDB() {
	err := database.Database.Ping()
	if err != nil {
		log.Fatal("connection is dead:", err)
	}
}

// PreparedQuery executes a prepared query and returns the resulting rows,
// it handles the errors and returns the context and prepared statement
// for the user to close them.
func (database *Database) PreparedQuery(statement string, args ...interface{}) (*sql.Tx, *sql.Stmt, *sql.Rows) {
	tx, stmt := database.PrepareStatement(statement)
	rows, err := stmt.Query(args...)
	if err != nil {
		log.Fatal("could not perform query: ", err)
	}
	return tx, stmt, rows
}

// PrepareStatement initializes an sqlite prepared statement from a string
// and returns the corresponding sql context and prepared statement.
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

// QueryCustom receives a parsed string in disjunctive normal form and
// executes the corresponding query.
func (database *Database) QueryCustom(stmtStr string, terms ...interface{}) []int64 {
	result := make([]int64, 0)

	tx, stmt, rows := database.PreparedQuery(stmtStr, terms...)
	defer stmt.Close()
	defer rows.Close()

	for rows.Next() {
		var id int64
		err := rows.Scan(&id)
		if err != nil {
			log.Fatal(err)
		}
		result = append(result, id)
	}
	err := rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	tx.Commit()
	return result
}

// QueryGroup receives a group ID and returns its name, start_date and
// end_date. It is assumed that the group is in the database.
func (database *Database) QueryGroup(groupID int64) (string, string, string) {
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

// QueryGroupMembers receives a group ID as a parameter and returns a map
// having the group members' names as keys and the value is true if the
// person is a member of the group.   It is assumed that the group is in
// the database.
func (database *Database) QueryGroupMembers(groupID int64) map[string]bool {
	members := make(map[string]bool)
	stmtStr := "SELECT " +
		" persons.stage_name " +
		"FROM " +
		" persons " +
		"INNER JOIN in_group ON in_group.id_person = persons.id_person " +
		"WHERE " +
		" in_group.id_group = ?"

	tx, stmt := database.PrepareStatement(stmtStr)
	defer stmt.Close()

	rows, err := stmt.Query(groupID)
	if err != nil {
		log.Fatal("could not execute query: ", err)
	}
	defer rows.Close()

	var member string
	for rows.Next() {
		err = rows.Scan(&member)
		if err != nil {
			log.Fatal(err)
		}
		members[member] = true
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	tx.Commit()
	return members
}

// QueryPerformerType receives a performer's ID as an argument and returns
// its type and name.   It is assumed that the performer is in the database.
func (database *Database) QueryPerformerType(id int64) (int, string) {
	stmtStr := "SELECT " +
		" performers.id_type, " +
		" performers.name " +
		"FROM " +
		" performers " +
		"INNER JOIN types ON types.id_type = performers.id_type " +
		"WHERE " +
		" performers.id_performer = ?"

	tx, stmt := database.PrepareStatement(stmtStr)
	defer stmt.Close()

	rows, err := stmt.Query(id)
	if err != nil {
		log.Fatal("could not execute query: ", err)
	}
	defer rows.Close()

	var performerType int
	var name string
	for rows.Next() {
		err = rows.Scan(&performerType, &name)
		if err != nil {
			log.Fatal(err)
		}
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	tx.Commit()
	return performerType, name
}

// QueryPerson receives a person's ID as an argument and returns its
// stage_name, real_name, birth_date and death_date, all as strings.   It
// is assumed that the person is in the database.
func (database *Database) QueryPerson(personID int64) (string, string, string, string) {
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

// QueryPersonGroups takes a person's ID as an argument and returns a map
// whose keys are the groups where the person is a member, and the values
// are all true.   It is assumed that the person is in the database.
func (database *Database) QueryPersonGroups(personID int64) map[string]bool {
	groups := make(map[string]bool)
	stmtStr := "SELECT " +
		" groups.name " +
		"FROM " +
		" groups " +
		"INNER JOIN in_group ON in_group.id_group = groups.id_group " +
		"WHERE " +
		" in_group.id_person = ?"

	tx, stmt := database.PrepareStatement(stmtStr)
	defer stmt.Close()

	rows, err := stmt.Query(personID)
	if err != nil {
		log.Fatal("could not execute query: ", err)
	}
	defer rows.Close()

	var group string
	for rows.Next() {
		err = rows.Scan(&group)
		if err != nil {
			log.Fatal(err)
		}
		groups[group] = true
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	tx.Commit()
	return groups
}

// QueryRola receives a Rola's ID as an argument and returns the correspoding
// rola, but with an empty path.   It is assumed that the rola is in the
// database.
func (database *Database) QueryRola(rolaID int64) *Rola {
	stmtStr := "SELECT " +
		" performers.name, " +
		" albums.name, " +
		" rolas.title, " +
		" rolas.track, " +
		" rolas.year, " +
		" rolas.genre " +
		"FROM rolas " +
		"INNER JOIN performers ON performers.id_performer = rolas.id_performer " +
		"INNER JOIN albums ON albums.id_album = rolas.id_album " +
		"WHERE " +
		" rolas.id_rola = ?"

	tx, stmt := database.PrepareStatement(stmtStr)
	defer stmt.Close()

	rows, err := stmt.Query(rolaID)
	if err != nil {
		log.Fatal("could not execute query: ", err)
	}
	defer rows.Close()

	var performer string
	var album string
	var title string
	var track int
	var year int
	var genre string
	for rows.Next() {
		err = rows.Scan(&performer, &album, &title, &track, &year, &genre)
		if err != nil {
			log.Fatal(err)
		}
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	tx.Commit()
	return &Rola{artist: performer,
		title: title,
		album: album,
		track: track,
		year:  year,
		genre: genre,
		path:  "",
		id:    rolaID,
	}
}

// QueryRolaForeign takes a Rola's ID as an argument and returns the
// IDs associated to its performer and album.
func (database *Database) QueryRolaForeign(rolaID int64) (int64, int64) {
	stmtStr := "SELECT " +
		" id_performer, " +
		" id_album " +
		"FROM " +
		" rolas " +
		"WHERE " +
		" rolas.id_rola = ?"

	tx, stmt := database.PrepareStatement(stmtStr)
	defer stmt.Close()

	rows, err := stmt.Query(rolaID)
	if err != nil {
		log.Fatal("could not execute query: ", err)
	}
	defer rows.Close()

	var performerID int64
	var albumID int64
	for rows.Next() {
		err = rows.Scan(&performerID, &albumID)
		if err != nil {
			log.Fatal(err)
		}
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	tx.Commit()
	return performerID, albumID
}

// QuerySimple receives a string as an argument, and returns a slice with
// the IDs of all the Rolas containing the string in its performer name,
// album name, title, or genre.
func (database *Database) QuerySimple(wildcard string) []int64 {
	result := make([]int64, 0)
	stmtStr := "SELECT " +
		" rolas.id_rola " +
		"FROM " +
		" rolas " +
		"INNER JOIN performers ON performers.id_performer = rolas.id_performer " +
		"INNER JOIN albums ON albums.id_album = rolas.id_album " +
		"WHERE " +
		" performers.name LIKE ? " +
		" OR albums.name LIKE ? " +
		" OR rolas.title LIKE ? " +
		" OR rolas.genre LIKE ?"

	wildCard := "%" + strings.TrimSpace(wildcard) + "%"
	tx, stmt, rows := database.PreparedQuery(stmtStr, wildCard, wildCard, wildCard, wildCard)
	defer stmt.Close()
	defer rows.Close()

	for rows.Next() {
		var id int64
		err := rows.Scan(&id)
		if err != nil {
			log.Fatal(err)
		}
		result = append(result, id)
	}
	err := rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	tx.Commit()
	return result
}

// UpdateGroup receives new values for the fields of a group, together with the
// group's ID, and updates the information. It is assumed that the group is
// in the database.
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

// UpdatePerformerType receives a performer's ID and a performer's type
// (0, 1, 2), to set the new peformer's type.   It is assumed that the
// performer is in the database.
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

// UpdatePerson receives new values for the fields of a person,
// together with the person's ID, and updates the information.   It is
// assumed that the person is in the database.
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

// UpdateRola takes a Rola as an argument and updates all its fields in
// the database.   It is assumed that the Rola taken as argument has the
// same ID as the rola we want to update.
func (database *Database) UpdateRola(rola *Rola) {
	stmtStr := "UPDATE rolas " +
		"SET title = ?, " +
		"    track = ?, " +
		"    year = ?, " +
		"    genre = ? " +
		"WHERE id_rola = ?"

	tx, stmt1 := database.PrepareStatement(stmtStr)
	defer stmt1.Close()

	_, err := stmt1.Exec(rola.title, rola.track, rola.year, rola.genre, rola.id)
	if err != nil {
		log.Fatal(err)
	}
	tx.Commit()

	performerID := database.AddPerformer(rola)
	albumID := database.AddAlbum(rola)

	stmtStr = "UPDATE rolas " +
		"SET id_performer = ?, " +
		"    id_album = ? " +
		"WHERE id_rola = ?"

	tx, stmt2 := database.PrepareStatement(stmtStr)
	defer stmt1.Close()

	_, err = stmt2.Exec(performerID, albumID, rola.id)
	if err != nil {
		log.Fatal(err)
	}
	tx.Commit()
}
