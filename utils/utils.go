package utils

import (
	"encoding/binary"
	"github.com/mitchellh/go-homedir"
	"github.com/zcong1993/utils/colors"
	"github.com/zcong1993/utils/terminal"
	"math/rand"
	"os"
	"path"
	"time"
)

func getEnvOrDefault(key, d string) string {
	v := os.Getenv(key)
	if v == "" {
		return d
	}
	return v
}

// MustGetDb is func return db path or abort
func MustGetDb(db string) string {
	home, err := homedir.Dir()
	if err != nil {
		terminal.Fail("An error occurred when get user home!")
		os.Exit(1)
	}
	return getEnvOrDefault("DB_PATH", path.Join(home, db))
}

func random(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

// RandomColor wrap your string with random color
func RandomColor(str string) string {
	cs := []colors.Func{colors.Blue, colors.Cyan, colors.Gray, colors.Green, colors.Purple, colors.Yellow}
	i := random(0, len(cs))
	fn := cs[i]
	return fn(str)
}

// Int64ToBytes is helper func to convert intf64 into []byte
func Int64ToBytes(i int64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(i))
	return buf
}

// BytesToInt64 is helper func to convert []byte into intf64
func BytesToInt64(buf []byte) int64 {
	return int64(binary.BigEndian.Uint64(buf))
}
