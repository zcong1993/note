package main

import (
	"errors"
	"fmt"
	"github.com/zcong1993/note/internal"
	"github.com/zcong1993/utils/colors"
	"github.com/zcong1993/utils/terminal"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"strings"
)

var (
	// GitCommit is commit hash for version
	GitCommit = ""
	// Version is app version
	Version = "v0.1.0"
)

var (
	app     = kingpin.New("note", "Command line tool for note.")
	addCmd  = app.Command("add", "add a note.")
	addTxts = addCmd.Arg("content", "not content").Required().Strings()

	listCmd = app.Command("list", "show all notes.").Default()

	deleteCmd = app.Command("delete", "delete a note by id.")
	deleteID  = deleteCmd.Arg("id", "note id").Required().Int64()

	getCmd = app.Command("get", "get notes by limit and offset.")
	limit  = getCmd.Flag("limit", "query limit").Short('l').Int()
	offset = getCmd.Flag("offset", "query offset").Short('o').Int()

	updateCmd = app.Command("update", "update a note.")
	updateID  = updateCmd.Flag("id", "note id for updating.").Short('i').Required().Int64()
	updateTxt = updateCmd.Arg("content", "not content").Required().String()

	deleteAllCmd = app.Command("delete-all", "delete all notes.")

	flushCmd = app.Command("flush", "flush note db.")

	version = app.Command("version", "Show note cli version.")
)

func list() {
	notes, err := internal.GetAll()
	if err != nil {
		terminal.LogErrPad(err)
		return
	}
	internal.ShowNotes(notes)
}

func add() {
	txt := strings.Join(*addTxts, " ")
	_, err := internal.Insert(txt)
	if err != nil {
		terminal.LogErrPad(err)
		return
	}
	terminal.LogSuccessPad("Add success.")
}

func d() {
	f, err := internal.Delete(*deleteID)
	if err != nil {
		terminal.LogErrPad(err)
		return
	}
	if f != 1 {
		terminal.LogErrPad(errors.New("Delete failed, maybe note not exists. "))
		return
	}
	terminal.LogSuccessPad("Delete success.")
}

func update() {
	f, err := internal.Update(*updateID, *updateTxt)
	if err != nil {
		terminal.LogErrPad(err)
		return
	}
	if f != 1 {
		terminal.LogErrPad(errors.New("Update failed, maybe note not exists. "))
		return
	}
	terminal.LogSuccessPad("Update success.")
}

func deleteAll() {
	f, err := internal.DeleteAll()
	if err != nil {
		terminal.LogErrPad(err)
		return
	}
	if f < 1 {
		terminal.LogErrPad(errors.New("Delete all failed, maybe no note now. "))
		return
	}
	terminal.LogSuccessPad("Delete all success.")
}

func get() {
	notes, err := internal.GetNotes(*limit, *offset)
	if err != nil {
		terminal.LogErrPad(err)
		return
	}
	internal.ShowNotes(notes)
}

func main() {
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case version.FullCommand():
		showVersion()
	case addCmd.FullCommand():
		add()
	case listCmd.FullCommand():
		list()
	case deleteCmd.FullCommand():
		d()
	case updateCmd.FullCommand():
		update()
	case flushCmd.FullCommand():
		internal.Flush()
	case deleteAllCmd.FullCommand():
		deleteAll()
	case getCmd.FullCommand():
		get()
	default:
		list()
	}
}

func showVersion() {
	version := fmt.Sprintf("%s version %s", colors.Cyan(app.Name), colors.Purple(Version))
	if len(GitCommit) != 0 {
		version += colors.Gray(fmt.Sprintf(" (%s)", GitCommit))
	}
	terminal.LogPad(version)
}
