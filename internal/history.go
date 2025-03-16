package internal

import (
	"fmt"

	"go.etcd.io/bbolt"
)

var db *bbolt.DB

func InitDB() error {
	db, _ = bbolt.Open("history.db", 0600, nil)
	db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("chat"))
		return err
	})

	return nil
}

func SaveMessage(user, ai string) {
	db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("chat"))
		return b.Put([]byte(user), []byte(ai))
	})
}

func LoadHistory() string {
	history := ""
	db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("chat"))
		b.ForEach(func(k, v []byte) error {
			history += fmt.Sprintf("User: %s\nAI: %s\n\n", k, v)
			return nil
		})

		return nil
	})

	return history
}
