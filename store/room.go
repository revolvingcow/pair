package store

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/boltdb/bolt"
)

type Room struct {
	Name          string    `json:"name"`
	LocalAddress  string    `json:"local"`
	RemoteAddress string    `json:"remote"`
	Connections   int       `json:"connections"`
	Created       time.Time `json:"created"`
	Touched       time.Time `json:"touched"`
}

func (room *Room) Get() error {
	db, err := bolt.Open("router.db", 0644, nil)
	if err != nil {
		return err
	}
	defer db.Close()

	world := []byte("rooms")
	key := []byte(room.Name)

	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(world)
		if bucket == nil {
			return fmt.Errorf("Bucket %q not found!", world)
		}

		b := bucket.Get(key)
		return json.Unmarshal(b, room)
	})

	return err
}

func (room *Room) Update() error {
	log.Printf("Room %s active with %d connections", room.Name, room.Connections)

	db, err := bolt.Open("router.db", 0644, nil)
	if err != nil {
		return err
	}
	defer db.Close()

	world := []byte("rooms")
	key := []byte(room.Name)
	room.Touched = time.Now()
	b, err := json.Marshal(room)
	if err != nil {
		return err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(world)
		if err != nil {
			return err
		}

		return bucket.Put(key, b)
	})

	return err
}

func (room *Room) Eject() error {
	log.Printf("Room %s ejected!", room.Name)

	db, err := bolt.Open("router.db", 0644, nil)
	if err != nil {
		return err
	}
	defer db.Close()

	world := []byte("rooms")
	key := []byte(room.Name)
	err = db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(world)
		if err != nil {
			return err
		}

		return bucket.Delete(key)
	})

	return err
}
