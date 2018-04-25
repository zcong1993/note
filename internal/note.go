package internal

import (
	"fmt"
	"github.com/go-xorm/xorm"
	"github.com/gost-c/gost-cli/colors"
	utils2 "github.com/gost-c/gost-cli/utils"
	_ "github.com/mattn/go-sqlite3"
	"github.com/zcong1993/note/utils"
	"os"
	"time"
)

// Note is db model
type Note struct {
	// Id is primary key
	Id int64
	// Txt is note content
	Txt string
	// Created is created_at
	Created time.Time `xorm:"created"`
	// Updated is updated_at
	Updated time.Time `xorm:"updated"`
}

var orm *xorm.Engine

func init() {
	o, err := xorm.NewEngine("sqlite3", utils.MustGetDb())
	if err != nil {
		utils2.Fail(fmt.Sprintf("xorm error %s. ", err.Error()))
		return
	}
	err = o.CreateTables(&Note{})
	if err != nil {
		utils2.Fail(err.Error())
		return
	}
	orm = o
}

// Insert add a note into db
func Insert(txt string) (int64, error) {
	f, err := orm.Insert(&Note{Txt: txt})
	return f, err
}

// GetAll get all the notes from db
func GetAll() ([]Note, error) {
	var notes []Note
	err := orm.Find(&notes)
	return notes, err
}

// GetNotes can get notes by limit and offset
func GetNotes(limit, offset int) ([]Note, error) {
	var notes []Note
	err := orm.Limit(limit, offset).Find(&notes)
	return notes, err
}

// Update can update a note by id
func Update(id int64, txt string) (int64, error) {
	f, err := orm.ID(id).Update(&Note{Txt: txt})
	return f, err
}

// Delete can delete a note by id
func Delete(id int64) (int64, error) {
	f, err := orm.ID(id).Delete(new(Note))
	return f, err
}

// DeleteAll can delete all the note
func DeleteAll() (int64, error) {
	return orm.Where("id > 0").Delete(new(Note))
}

// Flush will remove db file
func Flush() error {
	defer func() {
		utils2.LogPad(fmt.Sprintf("%s %s", colors.Yellow("WARN"), "Flush db. "))
		os.Exit(0)
	}()
	return os.Remove(utils.MustGetDb())
}

// ShowNotes is helper func for showing notes in terminal
func ShowNotes(notes []Note) {
	if len(notes) == 0 {
		utils2.LogPad(fmt.Sprintf("%s %s", colors.Green("INFO"), "No note now. "))
		return
	}
	s := ""
	for _, note := range notes {
		s += fmt.Sprintf("%s %s\n", colors.Blue(fmt.Sprintf("%d.", note.Id)), note.Txt)
	}
	utils2.LogPad(s)
}
