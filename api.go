package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Adresse de base de l'API REST Countries.
const baseURL = "https://restcountries.com/v3.1"

// Client HTTP avec un délai maximal : si l'API ne répond pas en 5 secondes, on abandonne.
var client = &http.Client{Timeout: 5 * time.Second}

// Pays : les champs qu'on garde de la réponse JSON de l'API.
type Pays struct {
	Nom struct {
		Commun string `json:"common"`
	} `json:"name"`
	Code       string   `json:"cca3"`
	Capitale   []string `json:"capital"`
	Region     string   `json:"region"`
	Population int      `json:"population"`
	Superficie float64  `json:"area"`
	Drapeau    struct {
		PNG string `json:"png"`
	} `json:"flags"`
}

// getJSON appelle une adresse de l'API et décode la réponse dans "cible".
func getJSON(url string, cible any) error {
	resp, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("API injoignable : %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("l'API a répondu %d", resp.StatusCode)
	}
	return json.NewDecoder(resp.Body).Decode(cible)
}

// TousLesPays : route de la collection (environ 250 pays).
func TousLesPays() ([]Pays, error) {
	var liste []Pays
	url := baseURL + "/all?fields=name,cca3,capital,region,population,area,flags"
	err := getJSON(url, &liste)
	return liste, err
}
