package utils

import (
	utils2 "github.com/gost-c/gost-cli/utils"
	"github.com/mitchellh/go-homedir"
	"os"
	"path"
)

var DB_PATH string

func init() {
	home, err := homedir.Dir()
	if err != nil {
		utils2.Fail("An error occurred when get user home!")
		os.Exit(1)
	}
	DB_PATH = path.Join(home, ".note.db")
}
