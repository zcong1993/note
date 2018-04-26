package bolt

import (
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/zcong1993/note/internal"
	"github.com/zcong1993/note/utils"
	"github.com/zcong1993/utils/colors"
	"github.com/zcong1993/utils/terminal"
	"os"
	"time"
)

const dbName = ".note.bolt.db"

var bucket = []byte("note")

// DB implement db interface via boltdb
type DB struct {
	db *bolt.DB
}

// NewBoltDB return a sqlite driven
func NewBoltDB() *DB {
	db, err := bolt.Open(utils.MustGetDb(dbName), 0600, &bolt.Options{Timeout: 2 * time.Second})
	if err != nil {
		terminal.Fail(fmt.Sprintf("boltdb init error %s. ", err.Error()))
		os.Exit(1)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(bucket)
		return err
	})
	if err != nil {
		terminal.Fail(fmt.Sprintf("boltdb create database error %s. ", err.Error()))
		os.Exit(1)
	}
	return &DB{db: db}
}

func (db *DB) getNextID() int64 {
	nextID := int64(1)
	db.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		c := b.Cursor()
		k, _ := c.Last()
		if k != nil {
			nextID = utils.BytesToInt64(k) + 1
		}
		return nil
	})
	return nextID
}

// Insert add a note into db
func (db *DB) Insert(txt string) (int64, error) {
	nextID := db.getNextID()
	note := internal.Note{Txt: txt}
	n, _ := json.Marshal(note)
	err := db.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		return b.Put(utils.Int64ToBytes(nextID), n)
	})
	if err != nil {
		return 0, err
	}
	return 1, nil
}

// GetAll get all the notes from db
func (db *DB) GetAll() ([]internal.Note, error) {
	var notes []internal.Note
	db.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var n internal.Note
			err := json.Unmarshal(v, &n)
			if err != nil {
				continue
			}
			n.Id = utils.BytesToInt64(k)
			notes = append(notes, n)
		}
		return nil
	})
	return notes, nil
}

// GetNotes can get notes by limit and offset
func (db *DB) GetNotes(limit, offset int) ([]internal.Note, error) {
	if limit == 0 && offset == 0 {
		return db.GetAll()
	}
	var notes []internal.Note
	db.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		c := b.Cursor()
		index := 0
		count := 0
		for k, v := c.First(); k != nil; k, v = c.Next() {
			if count >= limit {
				break
			}
			var n internal.Note
			err := json.Unmarshal(v, &n)
			if err != nil {
				continue
			}
			index++
			n.Id = utils.BytesToInt64(k)
			if index > offset {
				notes = append(notes, n)
				count++
			}
		}
		return nil
	})
	return notes, nil
}

// Update can update a note by id
func (db *DB) Update(id int64, txt string) (int64, error) {
	err := db.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		v, err := json.Marshal(&internal.Note{Id: id, Txt: txt})
		if err != nil {
			return err
		}
		return b.Put(utils.Int64ToBytes(id), v)
	})
	if err != nil {
		return 0, err
	}
	return 1, nil
}

// Delete can delete a note by id
func (db *DB) Delete(id int64) (int64, error) {
	k := utils.Int64ToBytes(id)
	err := db.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		v := b.Get(k)
		if v == nil {
			return fmt.Errorf("key %d not exists. ", id)
		}
		return b.Delete(k)
	})
	if err != nil {
		return 0, err
	}
	return 1, nil
}

// DeleteAll can delete all the note
func (db *DB) DeleteAll() (int64, error) {
	count := int64(0)
	db.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		c := b.Cursor()
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			err := b.Delete(k)
			if err != nil {
				continue
			}
			count++
		}
		return nil
	})
	return count, nil
}

// Flush will remove db file
func (db *DB) Flush() error {
	defer func() {
		terminal.LogPad(fmt.Sprintf("%s %s", colors.Yellow("WARN"), "Flush db. "))
		os.Exit(0)
	}()
	return os.Remove(utils.MustGetDb(dbName))
}
