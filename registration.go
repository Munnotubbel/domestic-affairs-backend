package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// User contains user information
type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Token contains all information to calculate the hash
type Token struct {
	Warrantor string `json:"warrantor"`
	Seed      int    `json:"seed"`
	Hash      string `json:"hash"`
	Valid     bool   `json:"valid"`
	Applicant User   `json:"applicant"`
}

// ReadTokens bliblablub
func ReadTokens(file string) {

	// Open our jsonFile
	jsonFile, err := os.Open(file)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened " + file)

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var Token []Token
	json.Unmarshal(byteValue, &Token)

}

// GenerateTokens generates
func GenerateTokens(warrantor string, seedStart int, count int) map[string]Token {

	ret := make(map[string]Token)

	for i := seedStart; i < count; i++ {
		var id string = fmt.Sprintf("%s:%08d\n", warrantor, i)

		h := sha256.New()
		h.Write([]byte(id))

		cur := Token{
			Hash:      fmt.Sprintf("%x", string(h.Sum(nil))),
			Seed:      i,
			Valid:     true,
			Warrantor: warrantor,
		}

		ret[cur.Hash] = cur
	}

	return ret
}
