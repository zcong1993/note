package main

import (
	"errors"
	"fmt"
	"github.com/gost-c/gost-cli/colors"
	"github.com/gost-c/gost-cli/utils"
	"github.com/tj/kingpin"
	"github.com/zcong1993/note/internal"
	"os"
)

var (
	// GitCommit is commit hash for version
	GitCommit = ""
	// Version is app version
	Version = "v0.1.0"
)

var (
	app    = kingpin.New("note", "Command line tool for note.")
	addCmd = app.Command("add", "add a note.")
	addTxt = addCmd.Arg("content", "not content").Required().String()

	listCmd = app.Command("list", "show all notes.")

	deleteCmd = app.Command("delete", "delete a note by id.")
	deleteId  = deleteCmd.Arg("id", "note id").Required().Int64()

	updateCmd = app.Command("update", "update a note.")
	updateId  = updateCmd.Flag("id", "note id for updating.").Short('i').Required().Int64()
	updateTxt = updateCmd.Arg("content", "not content").Required().String()

	deleteAllCmd = app.Command("delete-all", "delete all notes.")

	flushCmd = app.Command("flush", "flush note db.")

	version = app.Command("version", "Show note cli version.")
)

func list() {
	notes, err := internal.GetAll()
	if err != nil {
		utils.LogErrPad(err)
		return
	}
	internal.ShowNotes(notes)
}

func main() {
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case version.FullCommand():
		showVersion()
	case addCmd.FullCommand():
		_, err := internal.Insert(*addTxt)
		if err != nil {
			utils.LogErrPad(err)
			return
		}
		utils.LogSuccessPad("Add success.")
	case listCmd.FullCommand():
		list()
	case deleteCmd.FullCommand():
		f, err := internal.Delete(*deleteId)
		if err != nil {
			utils.LogErrPad(err)
			return
		}
		if f != 1 {
			utils.LogErrPad(errors.New("Delete failed, maybe note not exists. "))
			return
		}
		utils.LogSuccessPad("Delete success.")
	case updateCmd.FullCommand():
		f, err := internal.Update(*updateId, *updateTxt)
		if err != nil {
			utils.LogErrPad(err)
			return
		}
		if f != 1 {
			utils.LogErrPad(errors.New("Update failed, maybe note not exists. "))
			return
		}
		utils.LogSuccessPad("Update success.")
	case flushCmd.FullCommand():
		internal.Flush()
	case deleteAllCmd.FullCommand():
		f, err := internal.DeleteAll()
		if err != nil {
			utils.LogErrPad(err)
			return
		}
		if f < 1 {
			utils.LogErrPad(errors.New("Delete all failed, maybe no note now. "))
			return
		}
		utils.LogSuccessPad("Delete all success.")
	default:
		list()
	}
}

func showVersion() {
	version := fmt.Sprintf("%s version %s", colors.Cyan(app.Name), colors.Purple(Version))
	if len(GitCommit) != 0 {
		version += colors.Gray(fmt.Sprintf(" (%s)", GitCommit))
	}
	utils.LogPad(version)
}
