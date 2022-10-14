package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"io"
	"strings"
	"log"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	port := ":3000"

	r := chi.NewRouter()

	r.Use(
		middleware.Logger,
		middleware.Recoverer,
	)
	// r.Get("/searchMail", )
	r.Options("/searchMail", func(w http.ResponseWriter, r *http.Request) {
		// cors
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Origin", "*")
	})
	r.Post("/searchMail", handlerPost)

	fmt.Println("Serving on " + port)
	http.ListenAndServe(port, r)
}

func handlerPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// var query string
	// err:= json.NewDecoder((r.Body).Decode(&query))
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// bodyJson, err := json.Marshal(query.Query)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}
	var data map[string]interface{}

	err = json.Unmarshal(body, &data)
	fmt.Println(data["input"])
	
	query := `{
		"query":{
			"bool":{
				"must":[
					{
						"query_string":{
							"query": "`+ data["input"].(string) +`"
						}
					}
				]
			}
		}

	}`
	
	fmt.Println(strings.NewReader(query))

	req, err := http.NewRequest("POST", "http://localhost:4080/es/Emails/_search", strings.NewReader(query))
	if err != nil {
		fmt.Println(err)
	}
	req.SetBasicAuth("admin", "123")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	log.Println(resp.StatusCode)
	results, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	w.Write(results)
}
