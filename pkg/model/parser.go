package model

import (
    "strings"
)

type Parser struct {
    stmt string
}

var instanceP *Parser

func GetParser() *Parser {
    once.Do(func() {
        instanceP = &Parser{
            "SELECT " +
            " rolas.id_rola " +
            "FROM " +
            " rolas " +
            "INNER JOIN performers ON performers.id_performer = rolas.id_performer " +
            "INNER JOIN albums ON albums.id_album = rolas.id_album " +
            "WHERE ",
        }
    })
    return instanceP
}

func (parser *Parser) Parse(entry string) (string, []interface{}, bool) {
    queryTerms := make([]interface{}, 0)
    if !strings.HasPrefix(entry, "*~*") {
        return entry, queryTerms, false
    }
    entry = strings.TrimPrefix(entry, "*~*")
    statement := "( "
    tempLayer := strings.Split(entry, "*OR*")
    orLayer := make([]string, 0)
    andLayer := make ([]string, 0)
    for _, term := range tempLayer {
        term = strings.TrimSpace(term)
        if isParseable(term) {
            orLayer = append(orLayer, term)
        }
    }
    if len(orLayer) == 0 {
        return "", queryTerms, false
    }
    for i, term := range orLayer {
        for _, command := range strings.Split(term, "*AND*") {
            command = strings.TrimSpace(command)
            if isParseable(command) {
                andLayer = append(andLayer, command)
            }
        }
        if i < len(orLayer) - 1 {
            andLayer = append(andLayer, "*OR*")
        }
    }
    for i, term := range andLayer {
        switch {
        case term == "*OR*":
            if hasNext(i, andLayer) {
                statement = statement + ") OR ( "
            }

        case strings.HasPrefix(term, "*TITLE*~"):
            if hasNext(i, andLayer) {
                statement = statement + "rolas.title LIKE ? AND "
            } else {
                statement = statement + "rolas.title LIKE ? "
            }
            queryTerms = append(queryTerms, wildcard(strings.TrimSpace(strings.TrimPrefix(term, "*TITLE*~"))))

        case strings.HasPrefix(term, "*TITLE*="):
            if hasNext(i, andLayer) {
                statement = statement + "rolas.title = ? AND "
            } else {
                statement = statement + "rolas.title = ? "
            }
            queryTerms = append(queryTerms, strings.TrimSpace(strings.TrimPrefix(term, "*TITLE*=")))

        case strings.HasPrefix(term, "*ARTIST*~"):
            if hasNext(i, andLayer) {
                statement = statement + "performers.name LIKE ? AND "
            } else {
                statement = statement + "performers.name LIKE ? "
            }
            queryTerms = append(queryTerms, wildcard(strings.TrimSpace(strings.TrimPrefix(term, "*ARTIST*~"))))

        case strings.HasPrefix(term, "*ARTIST*="):
            if hasNext(i, andLayer) {
                statement = statement + "performers.name = ? AND "
            } else {
                statement = statement + "performers.name = ? "
            }
            queryTerms = append(queryTerms, strings.TrimSpace(strings.TrimPrefix(term, "*ARTIST*=")))

        case strings.HasPrefix(term, "*ALBUM*~"):
            if hasNext(i, andLayer) {
                statement = statement + "albums.name LIKE ? AND "
            } else {
                statement = statement + "albums.name LIKE ? "
            }
            queryTerms = append(queryTerms, wildcard(strings.TrimSpace(strings.TrimPrefix(term, "*ALBUM*~"))))

        case strings.HasPrefix(term, "*ALBUM*="):
            if hasNext(i, andLayer) {
                statement = statement + "albums.name = ? AND "
            } else {
                statement = statement + "albums.name = ? "
            }
            queryTerms = append(queryTerms, strings.TrimSpace(strings.TrimPrefix(term, "*ALBUM*=")))

        case strings.HasPrefix(term, "*TRACK*~"):
            if hasNext(i, andLayer) {
                statement = statement + "rolas.track LIKE ? AND "
            } else {
                statement = statement + "rolas.track LIKE ? "
            }
            queryTerms = append(queryTerms, wildcard(strings.TrimSpace(strings.TrimPrefix(term, "*TRACK*~"))))

        case strings.HasPrefix(term, "*TRACK*="):
            if hasNext(i, andLayer) {
                statement = statement + "rolas.track = ? AND "
            } else {
                statement = statement + "rolas.track = ? "
            }
            queryTerms = append(queryTerms, strings.TrimSpace(strings.TrimPrefix(term, "*TRACK*=")))

        case strings.HasPrefix(term, "*GENRE*~"):
            if hasNext(i, andLayer) {
                statement = statement + "rolas.genre LIKE ? AND "
            } else {
                statement = statement + "rolas.genre LIKE ? "
            }
            queryTerms = append(queryTerms, wildcard(strings.TrimSpace(strings.TrimPrefix(term, "*GENRE*~"))))

        case strings.HasPrefix(term, "*GENRE*="):
            if hasNext(i, andLayer) {
                statement = statement + "rolas.genre = ? AND "
            } else {
                statement = statement + "rolas.genre = ? "
            }
            queryTerms = append(queryTerms, strings.TrimSpace(strings.TrimPrefix(term, "*GENRE*=")))

        case strings.HasPrefix(term, "*YEAR*~"):
            if hasNext(i, andLayer) {
                statement = statement + "rolas.YEAR LIKE ? AND "
            } else {
                statement = statement + "rolas.YEAR LIKE ? "
            }
            queryTerms = append(queryTerms, wildcard(strings.TrimSpace(strings.TrimPrefix(term, "*YEAR*~"))))

        case strings.HasPrefix(term, "*YEAR*="):
            if hasNext(i, andLayer) {
                statement = statement + "rolas.year = ? AND "
            } else {
                statement = statement + "rolas.year = ? "
            }
            queryTerms = append(queryTerms, strings.TrimSpace(strings.TrimPrefix(term, "*YEAR*=")))

        case strings.HasPrefix(term, "*YEAR*<="):
            if hasNext(i, andLayer) {
                statement = statement + "rolas.year <= ? AND "
            } else {
                statement = statement + "rolas.year <= ? "
            }
            queryTerms = append(queryTerms, strings.TrimSpace(strings.TrimPrefix(term, "*YEAR*<=")))

        case strings.HasPrefix(term, "*YEAR*>="):
            if hasNext(i, andLayer) {
                statement = statement + "rolas.year >= ? AND "
            } else {
                statement = statement + "rolas.year >= ? "
            }
            queryTerms = append(queryTerms, strings.TrimSpace(strings.TrimPrefix(term, "*YEAR*>=")))
        }
    }
    statement = parser.stmt + statement + ")"
    return statement, queryTerms, true
}

func wildcard(entry string) string {
    return "%" + entry + "%"
}

func hasNext(i int, slice []string) bool {
    if i < len(slice) - 1 {
        return isParseable(slice[i+1])
    }
    return false
}

func isParseable(entry string) bool {
    entry = strings.TrimSpace(entry)
    switch {
    case strings.HasPrefix(entry, "*TITLE*~"):
        return true
    case strings.HasPrefix(entry, "*TITLE*="):
        return true
    case strings.HasPrefix(entry, "*ARTIST*~"):
        return true
    case strings.HasPrefix(entry, "*ARTIST*="):
        return true
    case strings.HasPrefix(entry, "*ALBUM*~"):
        return true
    case strings.HasPrefix(entry, "*ALBUM*="):
        return true
    case strings.HasPrefix(entry, "*TRACK*~"):
        return true
    case strings.HasPrefix(entry, "*TRACK*="):
        return true
    case strings.HasPrefix(entry, "*GENRE*~"):
        return true
    case strings.HasPrefix(entry, "*GENRE*="):
        return true
    case strings.HasPrefix(entry, "*YEAR*~"):
        return true
    case strings.HasPrefix(entry, "*YEAR*="):
        return true
    case strings.HasPrefix(entry, "*YEAR*<="):
        return true
    case strings.HasPrefix(entry, "*YEAR*>="):
        return true
    }
    return false
}
