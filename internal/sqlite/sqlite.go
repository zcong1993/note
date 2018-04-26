package sqlite

import (
	"fmt"
	"github.com/go-xorm/xorm"
	"github.com/zcong1993/utils/colors"
	"github.com/zcong1993/utils/terminal"
	// Import sqlite3 driver
	_ "github.com/mattn/go-sqlite3"
	"github.com/zcong1993/note/internal"
	"github.com/zcong1993/note/utils"
	"os"
)

// DB implement db interface via sqlite3
type DB struct {
	orm *xorm.Engine
}

// NewSqliteDB return a sqlite driven
func NewSqliteDB() *DB {
	o, err := xorm.NewEngine("sqlite3", utils.MustGetDb())
	if err != nil {
		terminal.Fail(fmt.Sprintf("xorm error %s. ", err.Error()))
		os.Exit(1)
	}
	err = o.CreateTables(&internal.Note{})
	if err != nil {
		terminal.Fail(err.Error())
		os.Exit(1)
	}
	return &DB{orm: o}
}

// Insert add a note into db
func (db *DB) Insert(txt string) (int64, error) {
	f, err := db.orm.Insert(&internal.Note{Txt: txt})
	return f, err
}

// GetAll get all the notes from db
func (db *DB) GetAll() ([]internal.Note, error) {
	var notes []internal.Note
	err := db.orm.Find(&notes)
	return notes, err
}

// GetNotes can get notes by limit and offset
func (db *DB) GetNotes(limit, offset int) ([]internal.Note, error) {
	var notes []internal.Note
	err := db.orm.Limit(limit, offset).Find(&notes)
	return notes, err
}

// Update can update a note by id
func (db *DB) Update(id int64, txt string) (int64, error) {
	f, err := db.orm.ID(id).Update(&internal.Note{Txt: txt})
	return f, err
}

// Delete can delete a note by id
func (db *DB) Delete(id int64) (int64, error) {
	f, err := db.orm.ID(id).Delete(new(internal.Note))
	return f, err
}

// DeleteAll can delete all the note
func (db *DB) DeleteAll() (int64, error) {
	return db.orm.Where("id > 0").Delete(new(internal.Note))
}

// Flush will remove db file
func (db *DB) Flush() error {
	defer func() {
		terminal.LogPad(fmt.Sprintf("%s %s", colors.Yellow("WARN"), "Flush db. "))
		os.Exit(0)
	}()
	return os.Remove(utils.MustGetDb())
}
