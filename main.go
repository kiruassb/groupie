package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"groupie/models"
)

func main() {
	http.HandleFunc("/", HomeHandler)
	http.HandleFunc("/artist", ArtistHandler)

	// fichiers statiques (CSS)
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Println("Serveur lancé sur http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}


func GetArtists() ([]models.Artist, error) {
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var artists []models.Artist
	err = json.NewDecoder(resp.Body).Decode(&artists)
	if err != nil {
		return nil, err
	}


	spotifyMap := map[string]string{
		"Queen":      "1dfeR4HaWDbWqFHLkxsg1d",
		"Pink Floyd": "0k17h0D3J5VfsdmQ1iZtE9",
		"SOJA":       "6Fx1cjY6uJqB3FqomkLzXU",
		"Scorpions":  "27T030eWyCQRmDyuvr1kxY",
	}

	for i := range artists {
		if id, ok := spotifyMap[artists[i].Name]; ok {
			artists[i].SpotifyID = id
		}
	}


	return artists, nil
}




func HomeHandler(w http.ResponseWriter, r *http.Request) {
	artists, err := GetArtists()
	if err != nil {
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Erreur template", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, artists)
	if err != nil {
		http.Error(w, "Erreur affichage", http.StatusInternalServerError)
	}
}

func ArtistHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "ID manquant", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID invalide", http.StatusBadRequest)
		return
	}

	artists, err := GetArtists()
	if err != nil {
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}

	var selectedArtist models.Artist
	found := false

	for _, artist := range artists {
		if artist.ID == id {
			selectedArtist = artist
			found = true
			break
		}
	}

	if !found {
		http.Error(w, "Artiste introuvable", http.StatusNotFound)
		return
	}

	tmpl, err := template.ParseFiles("templates/artist.html")
	if err != nil {
		http.Error(w, "Erreur template", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, selectedArtist)
	if err != nil {
		http.Error(w, "Erreur affichage", http.StatusInternalServerError)
	}
}