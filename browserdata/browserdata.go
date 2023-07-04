package browserdata

import (
	"database/sql"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"

	"github.com/moond4rk/hackbrowserdata/utils/fileutil"
)

type Extractor interface {
	Extract() (interface{}, error)
}

type RowsHandler func([]byte, interface{}) (interface{}, error)

type ExtractorHandler func([]byte, string, string, RowsHandler) (interface{}, error)

func DefaultDBHandler(masterKey []byte, dbpath, dbQuery string, rowsHandler RowsHandler) (interface{}, error) {
	tempFile := filepath.Join(os.TempDir(), filepath.Base(dbpath))
	if err := fileutil.CopyFile(dbpath, tempFile); err != nil {
		return nil, err
	}
	defer os.Remove(tempFile)
	db, err := sql.Open("sqlite3", tempFile)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	rows, err := db.Query(dbQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return rowsHandler(masterKey, rows)
}

func DefaultJSONHandler(masterKey []byte, dbpath, dbQuery string, rowsHandler RowsHandler) (interface{}, error) {
	tempFile := filepath.Join(os.TempDir(), filepath.Base(dbpath))
	if err := fileutil.CopyFile(dbpath, tempFile); err != nil {
		return nil, err
	}
	defer os.Remove(tempFile)
	s, err := os.ReadFile(tempFile)
	if err != nil {
		return nil, err
	}
	return rowsHandler(masterKey, s)
}
