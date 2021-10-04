package main

import (
	"log"
	"fmt"
	"net/http"
	"os"
	"io/ioutil"
	"encoding/json"
	"time"
	"math/rand"
	"html/template"
	"path"
)

type Status struct {
	Water int `json:"water"`	
	Wind int `json:"wind"`	
}

type WaterAndWind struct {
	Status Status `json:"status"`
}

func UpdateWaterAndWind() {	
	var waterAndWind WaterAndWind
	rand.Seed(time.Now().UnixNano())
	waterAndWind.Status.Wind = int(rand.Intn(100 - 1) + 1)
	waterAndWind.Status.Water = int(rand.Intn(100 - 1) + 1)
	waterAndWindByte, err := json.MarshalIndent(waterAndWind,"", "  ")
	if err != nil {
		log.Println(err)
	}
	ioutil.WriteFile("data.json", waterAndWindByte, 0644)
}

func main(){
	go func() {
		for {
			time.Sleep(time.Second * 15)
			UpdateWaterAndWind()
		}
	}()

	mux := http.DefaultServeMux
	mux.HandleFunc("/", func (w http.ResponseWriter, r *http.Request){
		jsonFile, err := os.Open("data.json")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer jsonFile.Close()
		byteValue, err := ioutil.ReadAll(jsonFile)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var waterAndWind WaterAndWind
		json.Unmarshal(byteValue, &waterAndWind)
		filePath := path.Join("views", "status.html")
		tmplt, err := template.ParseFiles(filePath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		water := waterAndWind.Status.Water
		wind := waterAndWind.Status.Wind
		var responseData = make(map[string]interface{})
		responseData["water"] = water
		responseData["wind"] = wind

		switch {
			case wind <= 6:
				responseData["windStatus"] = "aman"
			case wind <= 15:
				responseData["windStatus"] = "siaga"
			default:
				responseData["windStatus"] = "bahaya"
		}
	
		switch {
			case water <= 5:
				responseData["waterStatus"] = "aman"
			case water <= 8:
				responseData["waterStatus"] = "siaga"
			default:
				responseData["waterStatus"] = "bahaya"
		}
		tmplt.Execute(w, responseData)
	})

	handler := mux
	server := new(http.Server)
	server.Addr = ":8000"
	server.Handler = handler
	fmt.Println("server running on port 8000")
	server.ListenAndServe()
}






