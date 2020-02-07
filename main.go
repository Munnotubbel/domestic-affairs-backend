package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"sync"

	"github.com/julienschmidt/httprouter"
)

const (
	reddisDbIndex int = 0
)

type store struct {
	data map[string]string
	m    sync.RWMutex
}

var (
	addr          = flag.String("addr", ":8081", "http service address")
	generate      = flag.Bool("generate", false, "Regenerate Database")
	listAvailable = flag.Bool("list-available", false, "List available Tokens")
	listUsed      = flag.Bool("list-used", false, "List used Tokens")
	db            Database
)

func main() {

	db = NewDatabase(reddisDbIndex)

	flag.Parse()
	// ############################
	// Parse Command Line Arguments
	// ############################
	if *generate {
		err := generateTokens()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(-1)
		}
		os.Exit(0)
	}

	if *listAvailable {
		err := listAvailableTokens()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(-1)
		}
		os.Exit(0)
	}

	if *listUsed {
		err := listUsedTokens()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(-1)
		}
		os.Exit(0)
	}

	// ############################
	// Real Shit
	// ############################
	r := httprouter.New()

	// curl -X GET  "127.0.0.1:8081/register/8ec7c043a478ec5d7604523f2ff6dac8e2f15d01fb55a4a7fed72b31368bb8f0/pascal+huerst/paso@domo.ch"
	r.GET("/register/:key/:name/:email", httpRegister)
	// curl -X GET  "127.0.0.1:8081/list-available" | jq --color-output
	r.GET("/list-available", httpListAvailable)
	// curl -X GET  "127.0.0.1:8081/list-used" | jq --color-output
	r.GET("/list-used", httpListUsed)

	fmt.Println("Starting server on: localhost" + *addr)

	err := http.ListenAndServe(*addr, r)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func generateTokens() error {
	fmt.Println("Regenerating Database!")

	tokens := GenerateTokens("Hendrik", 0, 20)
	for _, v := range tokens {
		err := db.Insert(v)
		if err != nil {
			return err
		}
	}
	return nil
}

func listAvailableTokens() error {
	tokens, _ := db.List()

	var keys []int
	var reorderedTokens = make(map[int]Token)

	for _, v := range tokens {
		if v.Valid {
			keys = append(keys, v.Seed)
			reorderedTokens[v.Seed] = v
		}
	}
	sort.Ints(keys)

	fmt.Printf("Used Tokens [%v]:\n", len(keys))

	for _, key := range keys {
		cur := reorderedTokens[key]
		fmt.Printf("%08d  %s  %s\n", cur.Seed, cur.Hash, cur.Warrantor)
	}

	return nil
}

func listUsedTokens() error {
	tokens, _ := db.List()

	var keys []int
	var reorderedTokens = make(map[int]Token)

	for _, v := range tokens {
		if !v.Valid {
			keys = append(keys, v.Seed)
			reorderedTokens[v.Seed] = v
		}
	}
	sort.Ints(keys)

	fmt.Printf("Used Tokens [%v]:\n", len(keys))

	for _, key := range keys {
		cur := reorderedTokens[key]
		fmt.Printf("%08d  %s  %s -> %s <%s>\n", cur.Seed, cur.Hash, cur.Warrantor, cur.Applicant.Name, cur.Applicant.Email)
	}

	return nil
}

func httpListAvailable(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Applicant")
	w.Header().Set("Content-Type", "application/json")

	if (*r).Method == "OPTIONS" {
		return
	}

	tokens, _ := db.List()

	var keys []int
	var reorderedTokens = make(map[int]Token)

	for _, v := range tokens {
		if v.Valid {
			keys = append(keys, v.Seed)
			reorderedTokens[v.Seed] = v
		}
	}
	sort.Ints(keys)

	var resp []Token

	for _, key := range keys {
		cur := reorderedTokens[key]
		resp = append(resp, cur)
	}

	jd, _ := json.Marshal(resp)

	w.Write(jd)
}

func httpListUsed(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Applicant")
	w.Header().Set("Content-Type", "application/json")

	if (*r).Method == "OPTIONS" {
		return
	}

	tokens, _ := db.List()

	var keys []int
	var reorderedTokens = make(map[int]Token)

	for _, v := range tokens {
		if !v.Valid {
			keys = append(keys, v.Seed)
			reorderedTokens[v.Seed] = v
		}
	}
	sort.Ints(keys)

	var resp []Token

	for _, key := range keys {
		cur := reorderedTokens[key]
		resp = append(resp, cur)
	}

	jd, _ := json.Marshal(resp)

	w.Write(jd)
}

func httpRegister(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	type response struct {
		success  bool   `json:"success"`
		response string `json:"response"`
	}

	var jd []byte

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Applicant")
	w.Header().Set("Content-Type", "application/json")

	if (*r).Method == "OPTIONS" {
		return
	}

	k := p.ByName("key")
	u := User{
		Name:  p.ByName("name"),
		Email: p.ByName("email"),
	}

	log.Printf("Register Attempt: name[%s], email[%s], token[%s]\n", u.Name, u.Email, k)
	//var resp response

	err := validateCredentials(u)
	if err != nil {
		log.Print(err.Error())

		resp := response{
			success:  false,
			response: err.Error(),
		}

		jd, _ = json.Marshal(resp)
		w.Write(jd)
		return
	}

	err = tryRegister(u, k)
	if err != nil {
		log.Print(err.Error())

		resp := response{
			success:  false,
			response: err.Error(),
		}
		jd, _ = json.Marshal(resp)
		w.Write(jd)
		return
	}

	// Success
	//resp = response{
	//	success:  true,
	//	responce: "Sucessfully registered",
	//}
	//var resp []response
	resppp := response{
		success:  true,
		response: "Tadaaaaa",
	}

	//resp = append(resp, tmp)

	jd, err = json.Marshal(resppp)
	if err != nil {
		log.Println(err.Error())
	}

	fmt.Println(resppp)

	w.Write(jd)
	fmt.Println("Reached End!")
}

func tryRegister(u User, k string) error {

	v, err := db.Query(k)
	if err != nil {
		return fmt.Errorf("Invalid Token: %s", k)
	}

	if !v.Valid {
		return fmt.Errorf("Invalid Token: %s", k)
	}

	// Update: Is this attomic or could it be hacked?
	v.Valid = false
	v.Applicant = u

	err = db.Insert(v)
	if err != nil {
		return fmt.Errorf("Invalid Token: %s", k)
	}

	return nil
}

func validateCredentials(u User) error {

	users, err := db.Emails()
	if err != nil {
		return fmt.Errorf("Database connection down. Retry later")
	}

	if _, exists := users[u.Email]; exists {
		return fmt.Errorf("Email address: %s is already in us", u.Email)
	}

	return nil
}
