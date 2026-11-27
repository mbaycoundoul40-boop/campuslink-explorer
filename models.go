package main

import "time"

// countryPayload correspond au JSON brut renvoyé par REST Countries v3.1.
type countryPayload struct {
	Name struct {
		Common   string `json:"common"`
		Official string `json:"official"`
	} `json:"name"`
	CCA3       string            `json:"cca3"`
	CCA2       string            `json:"cca2"`
	Capital    []string          `json:"capital"`
	Region     string            `json:"region"`
	Subregion  string            `json:"subregion"`
	Population int64             `json:"population"`
	Area       float64           `json:"area"`
	Languages  map[string]string `json:"languages"`
	Currencies map[string]struct {
		Name   string `json:"name"`
		Symbol string `json:"symbol"`
	} `json:"currencies"`
	Borders    []string `json:"borders"`
	Timezones  []string `json:"timezones"`
	Continents []string `json:"continents"`
	Flags      struct {
		PNG string `json:"png"`
		SVG string `json:"svg"`
		Alt string `json:"alt"`
	} `json:"flags"`
}

// Pays est la structure propre utilisée par le catalogue, le graphique SVG et les templates HTML.
type Pays struct {
	Code        string
	Code2       string
	Nom         string
	NomOfficiel string
	Capitale    string
	Region      string
	SousRegion  string
	Population  int64
	Superficie  float64
	Langues     []string
	Monnaies    []string
	Frontieres  []string
	Fuseaux     []string
	DrapeauURL  string
	DrapeauAlt  string
}

// FiltresCatalogue regroupe les critères de recherche, filtre, tri et pagination.
type FiltresCatalogue struct {
	Recherche  string
	Region     string
	PopMin     int64
	Tri        string
	Page       int
	TaillePage int
}

// ResultatCatalogue contient les données envoyées au template catalogue.html.
type ResultatCatalogue struct {
	PaysList      []Pays
	TotalFiltres  int
	TotalPages    int
	Filtres       FiltresCatalogue
	RegionsDispo  []string
	SourceSecours bool
	Avertissement string
}

// APIError permet de remonter un statut HTTP précis (400, 404, 502, 504).
type APIError struct {
	StatusCode int
	Message    string
	Err        error
}

func (e *APIError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

// Config contient la configuration chargée au démarrage.
type Config struct {
	Port        string
	APIBaseURL  string
	Timeout     time.Duration
	SecoursPath string
}
