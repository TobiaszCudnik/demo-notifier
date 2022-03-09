package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	max := 3
	min := 0
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		jsonString, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			return
		}
		// randomly fail
		randErr := rand.Intn(max-min) + min
		if randErr > max-2 {
			log.Println("Failing...")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(nil)
			return
		}
		// parse msgs
		var msgs []string
		err = json.Unmarshal(jsonString, &msgs)
		if err != nil {
			log.Println(err)
			return
		}
		log.Printf("Received %d msgs", len(msgs))
	})
	http.ListenAndServe(":5050", nil)
}
