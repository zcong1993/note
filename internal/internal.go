package internal

import (
	"fmt"
	"github.com/zcong1993/note/utils"
	"github.com/zcong1993/utils/colors"
	"github.com/zcong1993/utils/terminal"
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

// DB is db interface
type DB interface {
	Insert(txt string) (int64, error)
	GetAll() ([]Note, error)
	GetNotes(limit, offset int) ([]Note, error)
	Update(id int64, txt string) (int64, error)
	Delete(id int64) (int64, error)
	DeleteAll() (int64, error)
	Flush() error
}

// ShowNotes is helper func for showing notes in terminal
func ShowNotes(notes []Note) {
	if len(notes) == 0 {
		terminal.LogPad(fmt.Sprintf("%s %s", colors.Green("INFO"), "No note now. "))
		return
	}
	s := ""
	for _, note := range notes {
		s += fmt.Sprintf("%s %s\n", colors.Blue(fmt.Sprintf("%d.", note.Id)), utils.RandomColor(note.Txt))
	}
	terminal.LogPad(s)
}
