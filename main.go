package main

import (
	"fmt"
	"os"
	"io/ioutil"
	"encoding/json"
	"time"
	"math/rand"
	"net/http"
	"html/template"
	"log"
)

type Element struct {
    Status Status `json:"status"`
}

type Status struct{
	Water int `json:"water"`
	Wind int `json:"wind"`
}

var channelElement chan Element

func autoUpdate(channelElement chan Element){
	for {
		select {
		case <-time.After(15 * time.Second):
			var element Element
			element.Status.Water = rand.Intn(10);
			element.Status.Wind = rand.Intn(20);

			file, _ := json.MarshalIndent(element, "", " ")
			_ = ioutil.WriteFile("element.json", file, 0644)

			fmt.Println(element.Status.Water)
			fmt.Println(element.Status.Wind)
		}
	}
}

func main(){
	go autoUpdate(channelElement)

	funcMap := template.FuncMap{
        "percentage": func(i int, j int) float32 {
            return float32(i) / float32(j) * 100
        },
    }

	http.HandleFunc("/index", func(w http.ResponseWriter, r *http.Request) {
		jsonFile, err := os.Open("element.json")
		if err != nil {
			fmt.Println(err)
		}
		defer jsonFile.Close()

		fmt.Println("Successfully Opened element.json")

		byteValue, _ := ioutil.ReadAll(jsonFile)
		var element Element

		json.Unmarshal([]byte(byteValue), &element)

		var tmpl = template.Must(template.New("test").Funcs(funcMap).ParseFiles(
			"views/index.html",
		))
   
		err = tmpl.ExecuteTemplate(w, "index", element)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	err := http.ListenAndServe(":9090", nil)
    if err != nil {
        log.Fatal("Error running service: ", err)
    }
}