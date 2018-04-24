package utils

import (
	utils2 "github.com/gost-c/gost-cli/utils"
	"github.com/mitchellh/go-homedir"
	"os"
	"path"
)

func MustGetDb() string {
	home, err := homedir.Dir()
	if err != nil {
		utils2.Fail("An error occurred when get user home!")
		os.Exit(1)
	}
	return utils2.GetEnvOrDefault("DB_PATH", path.Join(home, ".note.db"))
}
