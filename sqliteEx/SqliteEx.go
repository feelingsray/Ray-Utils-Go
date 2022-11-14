package sqliteEx

import (
	"database/sql"

	"github.com/feelingsray/Ray-Utils-Go/tools"
)

type SqliteEx struct {
	DBPath string
	DBObj  *sql.DB
	DBStmt *sql.Stmt
}

func NewSqliteEx(dbPath string) *SqliteEx {
	db := SqliteEx{
		DBPath: dbPath,
		DBObj:  nil,
		DBStmt: nil,
	}
	return &db
}

func (s *SqliteEx) Open() (bool, error) {
	isExits, err := tools.PathExists(s.DBPath)
	if err != nil {
		return false, err
	}
	if !isExits {
		s.DBObj, err = sql.Open("sqlite3", s.DBPath)
		if err != nil {
			return false, err
		} else {
			return true, nil
		}
	}
	return true, nil

}
