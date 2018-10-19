package model

import (
	"strings"
)

// Parser for the search bar of the main application window.
// The parser uses && and || for 'AND' and 'OR', respectively;
// searches can be performed by any of the id3v2 fields in a
// Rola.   Start a search with '*~*', use the first two letters
//in the field between asterisks, followed by '=', '~', '<', or
// '>' for an exact search, a wildcard search, or for certain
// ranges (for numeric fields), respectively, e.g., '*AR*~punk'
// searches for all artists containing 'punk' in their name.
// There are negated versions of the four operators, '!=', etc.
// The parser joins the atomic formulas to get a formula in
// disjunctive normal form.
type Parser struct {
	stmt string
}

var instanceP *Parser

// GetParser returns the singleton instance of Parser. The
// singleton implementation is not thread safe due to conflicts
// with gtk.
func GetParser() *Parser {
	if instanceP == nil {
		instanceP = &Parser{
			`SELECT
           rolas.id_rola
         FROM
           rolas
         INNER JOIN performers ON performers.id_performer = rolas.id_performer
         INNER JOIN albums ON albums.id_album = rolas.id_album
         WHERE `,
		}
	}
	return instanceP
}

// Parse parses a search string into a sqlite query.
func (parser *Parser) Parse(entry string) (string, []interface{}, bool) {
	queryTerms := make([]interface{}, 0)
	if !strings.HasPrefix(entry, "*~*") {
		return entry, queryTerms, false
	}
	entry = strings.TrimPrefix(entry, "*~*")
	statement := "( "
	tempLayer := strings.Split(entry, "||")
	orLayer := make([]string, 0)
	andLayer := make([]string, 0)
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
		for _, command := range strings.Split(term, "&&") {
			command = strings.TrimSpace(command)
			if isParseable(command) {
				andLayer = append(andLayer, command)
			}
		}
		if i < len(orLayer)-1 {
			andLayer = append(andLayer, "||")
		}
	}
	if len(andLayer) == 0 {
		return "", queryTerms, false
	}
	for i, term := range andLayer {
		switch {
		case term == "||":
			if hasNext(i, andLayer) {
				statement = statement + ") OR ( "
			}

		case strings.HasPrefix(term, "*TI*~"):
			if hasNext(i, andLayer) {
				statement = statement + "rolas.title LIKE ? AND "
			} else {
				statement = statement + "rolas.title LIKE ? "
			}
			queryTerms = append(queryTerms, wildcard(strings.TrimSpace(strings.TrimPrefix(term, "*TI*~"))))

		case strings.HasPrefix(term, "*TI*="):
			if hasNext(i, andLayer) {
				statement = statement + "rolas.title = ? AND "
			} else {
				statement = statement + "rolas.title = ? "
			}
			queryTerms = append(queryTerms, strings.TrimSpace(strings.TrimPrefix(term, "*TI*=")))

		case strings.HasPrefix(term, "*TI*!~"):
			if hasNext(i, andLayer) {
				statement = statement + "NOT rolas.title LIKE ? AND "
			} else {
				statement = statement + "NOT rolas.title LIKE ? "
			}
			queryTerms = append(queryTerms, wildcard(strings.TrimSpace(strings.TrimPrefix(term, "*TI*!~"))))

		case strings.HasPrefix(term, "*TI*!="):
			if hasNext(i, andLayer) {
				statement = statement + "NOT rolas.title = ? AND "
			} else {
				statement = statement + "NOT rolas.title = ? "
			}
			queryTerms = append(queryTerms, strings.TrimSpace(strings.TrimPrefix(term, "*TI*!=")))

		case strings.HasPrefix(term, "*AR*~"):
			if hasNext(i, andLayer) {
				statement = statement + "performers.name LIKE ? AND "
			} else {
				statement = statement + "performers.name LIKE ? "
			}
			queryTerms = append(queryTerms, wildcard(strings.TrimSpace(strings.TrimPrefix(term, "*AR*~"))))

		case strings.HasPrefix(term, "*AR*="):
			if hasNext(i, andLayer) {
				statement = statement + "performers.name = ? AND "
			} else {
				statement = statement + "performers.name = ? "
			}
			queryTerms = append(queryTerms, strings.TrimSpace(strings.TrimPrefix(term, "*AR*=")))

		case strings.HasPrefix(term, "*AR*!~"):
			if hasNext(i, andLayer) {
				statement = statement + "NOT performers.name LIKE ? AND "
			} else {
				statement = statement + "NOT performers.name LIKE ? "
			}
			queryTerms = append(queryTerms, wildcard(strings.TrimSpace(strings.TrimPrefix(term, "*AR*!~"))))

		case strings.HasPrefix(term, "*AR*!="):
			if hasNext(i, andLayer) {
				statement = statement + "NOT performers.name = ? AND "
			} else {
				statement = statement + "NOT performers.name = ? "
			}
			queryTerms = append(queryTerms, strings.TrimSpace(strings.TrimPrefix(term, "*AR*!=")))

		case strings.HasPrefix(term, "*AL*~"):
			if hasNext(i, andLayer) {
				statement = statement + "albums.name LIKE ? AND "
			} else {
				statement = statement + "albums.name LIKE ? "
			}
			queryTerms = append(queryTerms, wildcard(strings.TrimSpace(strings.TrimPrefix(term, "*AL*~"))))

		case strings.HasPrefix(term, "*AL*="):
			if hasNext(i, andLayer) {
				statement = statement + "albums.name = ? AND "
			} else {
				statement = statement + "albums.name = ? "
			}
			queryTerms = append(queryTerms, strings.TrimSpace(strings.TrimPrefix(term, "*AL*=")))

		case strings.HasPrefix(term, "*AL*!~"):
			if hasNext(i, andLayer) {
				statement = statement + "NOT albums.name LIKE ? AND "
			} else {
				statement = statement + "NOT albums.name LIKE ? "
			}
			queryTerms = append(queryTerms, wildcard(strings.TrimSpace(strings.TrimPrefix(term, "*AL*!~"))))

		case strings.HasPrefix(term, "*AL*!="):
			if hasNext(i, andLayer) {
				statement = statement + "NOT albums.name = ? AND "
			} else {
				statement = statement + "NOT albums.name = ? "
			}
			queryTerms = append(queryTerms, strings.TrimSpace(strings.TrimPrefix(term, "*AL*!=")))

		case strings.HasPrefix(term, "*TR*~"):
			if hasNext(i, andLayer) {
				statement = statement + "rolas.track LIKE ? AND "
			} else {
				statement = statement + "rolas.track LIKE ? "
			}
			queryTerms = append(queryTerms, wildcard(strings.TrimSpace(strings.TrimPrefix(term, "*TR*~"))))

		case strings.HasPrefix(term, "*TR*="):
			if hasNext(i, andLayer) {
				statement = statement + "rolas.track = ? AND "
			} else {
				statement = statement + "rolas.track = ? "
			}
			queryTerms = append(queryTerms, strings.TrimSpace(strings.TrimPrefix(term, "*TR*=")))

		case strings.HasPrefix(term, "*TR*<"):
			if hasNext(i, andLayer) {
				statement = statement + "rolas.track < ? AND "
			} else {
				statement = statement + "rolas.track < ? "
			}
			queryTerms = append(queryTerms, strings.TrimSpace(strings.TrimPrefix(term, "*TR*<")))

		case strings.HasPrefix(term, "*TR*>"):
			if hasNext(i, andLayer) {
				statement = statement + "rolas.track > ? AND "
			} else {
				statement = statement + "rolas.track > ? "
			}
			queryTerms = append(queryTerms, strings.TrimSpace(strings.TrimPrefix(term, "*TR*>")))

		case strings.HasPrefix(term, "*TR*!~"):
			if hasNext(i, andLayer) {
				statement = statement + "NOT rolas.track LIKE ? AND "
			} else {
				statement = statement + "NOT rolas.track LIKE ? "
			}
			queryTerms = append(queryTerms, wildcard(strings.TrimSpace(strings.TrimPrefix(term, "*TR*!~"))))

		case strings.HasPrefix(term, "*TR*!="):
			if hasNext(i, andLayer) {
				statement = statement + "NOT rolas.track = ? AND "
			} else {
				statement = statement + "NOT rolas.track = ? "
			}
			queryTerms = append(queryTerms, strings.TrimSpace(strings.TrimPrefix(term, "*TR*!=")))

		case strings.HasPrefix(term, "*TR*!<"):
			if hasNext(i, andLayer) {
				statement = statement + "rolas.track >= ? AND "
			} else {
				statement = statement + "rolas.track >= ? "
			}
			queryTerms = append(queryTerms, strings.TrimSpace(strings.TrimPrefix(term, "*TR*!<")))

		case strings.HasPrefix(term, "*TR*!>"):
			if hasNext(i, andLayer) {
				statement = statement + "rolas.track <= ? AND "
			} else {
				statement = statement + "rolas.track <= ? "
			}
			queryTerms = append(queryTerms, strings.TrimSpace(strings.TrimPrefix(term, "*TR*!>")))

		case strings.HasPrefix(term, "*GE*~"):
			if hasNext(i, andLayer) {
				statement = statement + "rolas.genre LIKE ? AND "
			} else {
				statement = statement + "rolas.genre LIKE ? "
			}
			queryTerms = append(queryTerms, wildcard(strings.TrimSpace(strings.TrimPrefix(term, "*GE*~"))))

		case strings.HasPrefix(term, "*GE*="):
			if hasNext(i, andLayer) {
				statement = statement + "rolas.genre = ? AND "
			} else {
				statement = statement + "rolas.genre = ? "
			}
			queryTerms = append(queryTerms, strings.TrimSpace(strings.TrimPrefix(term, "*GE*=")))

		case strings.HasPrefix(term, "*GE*!~"):
			if hasNext(i, andLayer) {
				statement = statement + "NOT rolas.genre LIKE ? AND "
			} else {
				statement = statement + "NOT rolas.genre LIKE ? "
			}
			queryTerms = append(queryTerms, wildcard(strings.TrimSpace(strings.TrimPrefix(term, "*GE*!~"))))

		case strings.HasPrefix(term, "*GE*!="):
			if hasNext(i, andLayer) {
				statement = statement + "NOT rolas.genre = ? AND "
			} else {
				statement = statement + "NOT rolas.genre = ? "
			}
			queryTerms = append(queryTerms, strings.TrimSpace(strings.TrimPrefix(term, "*GE*!=")))

		case strings.HasPrefix(term, "*YE*~"):
			if hasNext(i, andLayer) {
				statement = statement + "rolas.YEAR LIKE ? AND "
			} else {
				statement = statement + "rolas.YEAR LIKE ? "
			}
			queryTerms = append(queryTerms, wildcard(strings.TrimSpace(strings.TrimPrefix(term, "*YE*~"))))

		case strings.HasPrefix(term, "*YE*="):
			if hasNext(i, andLayer) {
				statement = statement + "rolas.year = ? AND "
			} else {
				statement = statement + "rolas.year = ? "
			}
			queryTerms = append(queryTerms, strings.TrimSpace(strings.TrimPrefix(term, "*YE*=")))

		case strings.HasPrefix(term, "*YE*<"):
			if hasNext(i, andLayer) {
				statement = statement + "rolas.year < ? AND "
			} else {
				statement = statement + "rolas.year < ? "
			}
			queryTerms = append(queryTerms, strings.TrimSpace(strings.TrimPrefix(term, "*YE*<")))

		case strings.HasPrefix(term, "*YE*>"):
			if hasNext(i, andLayer) {
				statement = statement + "rolas.year > ? AND "
			} else {
				statement = statement + "rolas.year > ? "
			}
			queryTerms = append(queryTerms, strings.TrimSpace(strings.TrimPrefix(term, "*YE*>")))

		case strings.HasPrefix(term, "*YE*!~"):
			if hasNext(i, andLayer) {
				statement = statement + "NOT rolas.YEAR LIKE ? AND "
			} else {
				statement = statement + "NOT rolas.YEAR LIKE ? "
			}
			queryTerms = append(queryTerms, wildcard(strings.TrimSpace(strings.TrimPrefix(term, "*YE*!~"))))

		case strings.HasPrefix(term, "*YE*!="):
			if hasNext(i, andLayer) {
				statement = statement + "NOT rolas.year = ? AND "
			} else {
				statement = statement + "NOT rolas.year = ? "
			}
			queryTerms = append(queryTerms, strings.TrimSpace(strings.TrimPrefix(term, "*YE*!=")))

		case strings.HasPrefix(term, "*YE*!<"):
			if hasNext(i, andLayer) {
				statement = statement + "rolas.year >= ? AND "
			} else {
				statement = statement + "rolas.year >= ? "
			}
			queryTerms = append(queryTerms, strings.TrimSpace(strings.TrimPrefix(term, "*YE*!<")))

		case strings.HasPrefix(term, "*YE*!>"):
			if hasNext(i, andLayer) {
				statement = statement + "rolas.year <= ? AND "
			} else {
				statement = statement + "rolas.year <= ? "
			}
			queryTerms = append(queryTerms, strings.TrimSpace(strings.TrimPrefix(term, "*YE*!>")))
		}
	}
	statement = parser.stmt + statement + ")"
	return statement, queryTerms, true
}

func wildcard(entry string) string {
	return "%" + entry + "%"
}

func hasNext(i int, slice []string) bool {
	if i < len(slice)-1 {
		return isParseable(slice[i+1])
	}
	return false
}

func isParseable(entry string) bool {
	entry = strings.TrimSpace(entry)
	frames := []string{"*TI*", "*AR*", "*AL*", "*GE*"}
	operators := []string{"~", "=", "!~", "!="}
	for _, frame := range frames {
		for _, operator := range operators {
			if strings.HasPrefix(entry, frame+operator) {
				return true
			}
		}
	}
	frames = []string{"*TR*", "*YE*"}
	operators = []string{"~", "=", "<", ">", "!~", "!=", "!<", "!>"}
	for _, frame := range frames {
		for _, operator := range operators {
			if strings.HasPrefix(entry, frame+operator) {
				return true
			}
		}
	}
	return false
}
