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

type Note struct {
	Id      int64
	Txt     string
	Created time.Time `xorm:"created"`
	Updated time.Time `xorm:"updated"`
}

var orm *xorm.Engine

func init() {
	o, err := xorm.NewEngine("sqlite3", utils.DB_PATH)
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

func Insert(txt string) error {
	_, err := orm.Insert(&Note{Txt: txt})
	return err
}

func GetAll() ([]Note, error) {
	var notes []Note
	err := orm.Find(&notes)
	return notes, err
}

func GetNotes(limit, offset int) ([]Note, error) {
	var notes []Note
	err := orm.Limit(limit, offset).Find(&notes)
	return notes, err
}

func Update(id int64, txt string) error {
	_, err := orm.ID(id).Update(&Note{Txt: txt})
	return err
}

func Delete(id int64) error {
	_, err := orm.ID(id).Delete(new(Note))
	return err
}

func DeleteAll() (int64, error) {
	return orm.Where("id > 0").Delete(new(Note))
}

func Flush() error {
	defer func() {
		utils2.LogPad(fmt.Sprintf("%s %s", colors.Yellow("WARN"), "Flush db. "))
		os.Exit(0)
	}()
	return os.Remove(utils.DB_PATH)
}

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
