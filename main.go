package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	rand2 "math/rand"
	"net/http"
	"time"
)

func main() {
	port := flag.String("p", "3000", "Sets the port for server")
	flag.Parse()
	http.HandleFunc("/", getNames)
	log.Println("Server is starting on port", *port)
	log.Fatal(http.ListenAndServe(":"+*port, nil))

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
	for _, v := range result {
		switch vv := v.(type) {
		case []interface{}:
			for _, u := range vv {
				currentName := fmt.Sprintf("%v", u.([]interface{})[11])
				Names = append(Names, currentName)
			}
		default:
			log.Panic("This JSON is not compatible!")
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
