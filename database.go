package main

import (
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis"
)

// Database is the main container
type Database struct {
	Client *redis.Client
}

// NewDatabase creates a redis client
func NewDatabase(dbIndex int) Database {

	d := Database{
		Client: redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "",      // no password set
			DB:       dbIndex, // use default DB
		}),
	}

	_, err := d.Client.Ping().Result()
	if err == nil {
		fmt.Println("Database is up and running")
	}

	return d
}

// Insert adds a new Token to the db
func (d *Database) Insert(t Token) error {

	v, err := json.Marshal(t)
	if err != nil {
		return err
	}

	err = d.Client.Set(t.Hash, v, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

// Query returns a Token, if existing
func (d *Database) Query(hash string) (Token, error) {

	val, err := d.Client.Get(hash).Result()
	if err != nil {
		return Token{}, err
	}

	var ret Token
	err = json.Unmarshal([]byte(val), &ret)
	if err != nil {
		return Token{}, err
	}

	return ret, nil
}

// List returns a map of all Tokens in DB
func (d *Database) List() (map[string]Token, error) {

	var ret = make(map[string]Token)

	keys := d.Client.Keys("*").Val()

	var item Token
	var err error

	for _, key := range keys {
		item, err = d.Query(key)
		if err != nil {
			return ret, err
		}

		ret[key] = item
	}
	return ret, nil
}

// Emails returns a map of all k=Emails v=name in DB
func (d *Database) Emails() (map[string]string, error) {

	var ret = make(map[string]string)

	keys := d.Client.Keys("*").Val()

	var item Token
	var err error

	for _, key := range keys {
		item, err = d.Query(key)
		if err != nil {
			return ret, err
		}

		ret[item.Applicant.Email] = item.Applicant.Name
	}
	return ret, nil
}
