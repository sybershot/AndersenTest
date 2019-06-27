package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	rand2 "math/rand"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/", getNames)
	log.Fatal(http.ListenAndServe(":80", nil))
}

var myClient = &http.Client{Timeout: 10 * time.Second}

type ChildNames struct {
	Names []string `json:"names"`
}

func getNames(w http.ResponseWriter, r *http.Request) {
	req, err := myClient.Get("https://data.cityofnewyork.us/api/views/25th-nujf/rows.json?accessType=DOWNLOAD")

	if err != nil {
		log.Panic(err)
	}

	defer req.Body.Close()

	byteValue, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Panic(err)
	}
	var result map[string]interface{}
	err = json.Unmarshal([]byte(byteValue), &result)

	if err != nil {
		log.Panic(err)
	}

	delete(result, "meta") //We need only the data.

	//All names go here! I have no idea how to work with the names directly, so ¯\_(ツ)_/¯, it's working at least.
	var Names []string
	for k, v := range result {
		switch vv := v.(type) {
		case []interface{}:
			for _, u := range vv {
				currentName := fmt.Sprintf("%v", u.([]interface{})[11])
				Names = append(Names, currentName)
			}
		default:
			fmt.Println(k, "is not a JSON resp!")
		}
	}

	var acc []string //Creating accumulator for JSON.
	i := 0
	for i < 5 {
		acc = append(acc, Names[rand2.Int31n(int32(len(Names)))])
		i++
	}
	resp := &ChildNames{Names: acc}
	payload, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(payload)
	if err != nil {
		log.Panic(err)
	}
}
