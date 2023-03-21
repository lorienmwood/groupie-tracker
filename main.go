package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strconv"
)


type ArtistStruct []struct { 
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    string   `json:"locations"`
	ConcertDates string   `json:"concertDates"`
	Relations    string   `json:"relations"`
}

type SingleArtist struct { 
	UrlArtist string
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    string   `json:"locations"`
	ConcertDates string   `json:"concertDates"`
	Relations    string   `json:"relations"`
}

var Artists ArtistStruct
// var UrlArtist string = "https://groupietrackers.herokuapp.com/api/artists/" + strconv.Itoa(ArtistStruct[i].ID)

var tpl *template.Template
var LocalHost = "8080"

// initialize an HTTP client (the party that will send the request to the server)
var client *http.Client

func main () {
	exec.Command("open", "http://localhost:"+LocalHost).Start()

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	tpl = template.Must(template.ParseGlob("static/*.html"))

	http.HandleFunc("/", homePage)
	// http.HandleFunc("/artist", artistInfo)
	err := http.ListenAndServe(":"+LocalHost, nil)
	log.Fatal(err)
	// client = &http.Client{Timeout: 10 * time.Second} //timeout if no response
}

func homePage(w http.ResponseWriter, r *http.Request) {
	// the switch functions handles all possible http.Status errors
	switch{
	case r.Method != http.MethodGet:
		w.WriteHeader(http.StatusBadRequest)
	case r.URL.Path != "/":
		 http.Error(w, "404 not found", http.StatusNotFound)
	default:
		// var a ArtistStruct
		var all []SingleArtist

	response, err := http.Get("https://groupietrackers.herokuapp.com/api/artists") // http.Get request to the url we've defined
	if err != nil {
	  w.WriteHeader(http.StatusInternalServerError)
	}
	
	defer response.Body.Close()
	responseData, _ := ioutil.ReadAll(response.Body) // read the request body, and then we'll pass this body into the unmarshal func
	// fmt.Println("The JSON body: ", (string(responseData)))
	
	err = json.Unmarshal(responseData, &Artists) 
	if err != nil {
      w.WriteHeader(http.StatusInternalServerError)
	} 
	
	for _, v := range Artists{
     data := SingleArtist{
		UrlArtist: "https://groupietrackers.herokuapp.com/api/artists/"+strconv.Itoa(v.ID),
		ID: v.ID,
		Image: v.Image,
		Name: v.Name,
		Members: v.Members,
		CreationDate: v.CreationDate,
		FirstAlbum: v.FirstAlbum,
		ConcertDates: v.ConcertDates,
		Locations: v.Locations,
		Relations: v.Relations,
	 }
	 all = append(all, data)

	}
		tpl.ExecuteTemplate(w, "index.html", all)
	}
	
}

