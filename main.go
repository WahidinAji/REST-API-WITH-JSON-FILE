package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type Data struct {
	Data []Film `json:"data"`
}

type Film struct {
	Title    string   `json:"title"`
	Genre    []string `json:"genre"`
	Rating   string   `json:"rating"`
	Duration string   `json:"duration"`
	Quality  string   `json:"quality"`
	Trailer  string   `json:"trailer"`
	Watch    string   `json:"watch"`
}

type Dependency struct {
	dummy    *os.File
	Filename string
}

func (d *Dependency) jsonToFilm(w http.ResponseWriter, r *http.Request) error {
	//readJson, err := ioutil.ReadAll(d.dummy)
	readJson, err := ioutil.ReadFile(d.Filename)
	if err != nil {
		return err
	}
	//d.dummy.Read(readJson)
	defer d.dummy.Close()

	var data Data
	_ = json.Unmarshal([]byte(readJson), &data)
	var films []Film
	for _, d := range data.Data {
		var film Film
		film.Title = strings.Trim(d.Title, "\n\t")
		film.Genre = d.Genre
		film.Rating = strings.TrimSpace(d.Rating)
		film.Duration = strings.TrimSpace(d.Duration)
		film.Quality = d.Quality
		film.Trailer = d.Trailer
		film.Watch = d.Watch
		films = append(films, film)
	}
	return json.NewEncoder(w).Encode(Data{
		films,
	})
}

//func AppendToJson(w http.ResponseWriter, r *http.Request) error {
//	readJson, err
//}

func handleRequest() {

	file, err := os.Open("dummy.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	//fileName, _ := ioutil.ReadAll(file)
	//deps := Dependency{dummy: file}

	fileName := "dummy.json"

	deps := Dependency{dummy: file, Filename: string(fileName)}

	http.HandleFunc("/", JsonReader)
	http.HandleFunc("/dummies", jsonToFilmWithNoDependency)

	//with dependencies
	http.HandleFunc("/films", func(w http.ResponseWriter, r *http.Request) {
		deps.jsonToFilm(w, r)
	})
	log.Fatal(http.ListenAndServe(":8001", nil))
}

func main() {

	fmt.Println("Sample read json file")
	handleRequest()

}

func jsonToFilmWithNoDependency(w http.ResponseWriter, r *http.Request) {
	//jsonFile, err := os.Open("dummy.json")
	//if err != nil {
	//	http.Error(w, err.Error(), 500)
	//}
	//fmt.Println("Success opened dummy.json")
	//defer jsonFile.Close()
	//readJson, _ := ioutil.ReadAll(jsonFile)

	readJson, err := ioutil.ReadFile("dummy.json")
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
	var data Data
	_ = json.Unmarshal([]byte(readJson), &data)
	var films []Film
	for _, d := range data.Data {
		var film Film
		film.Title = strings.Trim(d.Title, "\n\t")
		film.Genre = d.Genre
		film.Rating = d.Rating
		film.Duration = d.Duration
		film.Quality = d.Quality
		film.Trailer = d.Trailer
		film.Watch = d.Watch
		films = append(films, film)
	}
	json.NewEncoder(w).Encode(Data{
		films,
	})
}

func JsonReader(w http.ResponseWriter, r *http.Request) {
	jsonFile, err := os.Open("dummy.json")
	if err != nil {
		panic(err)
	}
	fmt.Println("Success opened dummy.json")
	defer jsonFile.Close()

	jsonData, _ := ioutil.ReadFile("dummy.json")
	var data Data
	err = json.Unmarshal([]byte(jsonData), &data)
	/*
		//fmt.Println(data)
		//for i := 0; i < len(data.Data); i++ {
		//	fmt.Println("Title : ", data.Data[i].Title)
		//}
		//fmt.Println(data)
	*/
	json.NewEncoder(w).Encode(data)
}

type Deps struct {
	FileName string
	Dummy    *os.File
}

func (d *Deps) LoadAndTransform(w http.ResponseWriter, r *http.Request) {
	readJson, _ := ioutil.ReadFile(d.FileName)
	d.Dummy.Read(readJson)
	defer d.Dummy.Close()

	var data Data
	_ = json.Unmarshal([]byte(readJson), &data)
	var films []Film
	for _, d := range data.Data {
		var film Film
		film.Title = strings.Trim(d.Title, "\n\t")
		film.Genre = d.Genre
		film.Rating = d.Rating
		film.Duration = d.Duration
		film.Quality = d.Quality
		film.Trailer = d.Trailer
		film.Watch = d.Watch
		films = append(films, film)
		data.Data = films
	}
}
