package utils

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestMustGetDb(t *testing.T) {
	mock := ".note.db"
	p := MustGetDb(mock)
	assert.NotEmpty(t, p, "should get db path")
	fakePath := "test.db"
	os.Setenv("DB_PATH", fakePath)
	p = MustGetDb(mock)
	assert.Equal(t, p, fakePath)
	os.Unsetenv("DB_PATH")
}
