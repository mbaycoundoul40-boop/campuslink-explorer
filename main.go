package main

import (
	"errors"
	"html/template"
	"log"
	"net/http"
	"regexp"
)

// Un code pays valide = 3 lettres (ex : MAR).
var codeValide = regexp.MustCompile(`^[A-Za-z]{3}$`)

// page charge et affiche une page HTML avec ses données.
func page(w http.ResponseWriter, fichier string, donnees any) {
	tmpl, err := template.ParseFiles("templates/base.html", "templates/"+fichier)
	if err != nil {
		http.Error(w, "Erreur de template", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	if err := tmpl.ExecuteTemplate(w, "base", donnees); err != nil {
		log.Println(err)
	}
}

// erreur affiche la page d'erreur avec le bon code HTTP.
func erreur(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	page(w, "erreur.html", map[string]any{"Code": code, "Message": message})
}

// accueil affiche le nombre de pays disponibles.
func accueil(w http.ResponseWriter, r *http.Request) {
	pays, err := TousLesPays()
	if err != nil {
		log.Println(err)
		erreur(w, http.StatusBadGateway, "Les données sont indisponibles pour le moment.")
		return
	}
	page(w, "accueil.html", map[string]any{"NbPays": len(pays)})
}

// catalogue récupère tous les pays pour les afficher.
func catalogue(w http.ResponseWriter, r *http.Request) {
	pays, err := TousLesPays()
	if err != nil {
		log.Println(err)
		erreur(w, http.StatusBadGateway, "Les données sont indisponibles pour le moment.")
		return
	}
	page(w, "catalogue.html", map[string]any{"Pays": pays})
}

// fiche vérifie le code pays puis affiche les informations correspondantes.
func fiche(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
	if !codeValide.MatchString(code) {
		erreur(w, http.StatusBadRequest, "Code pays invalide : il faut 3 lettres.")
		return
	}
	p, err := PaysParCode(code)
	if errors.Is(err, ErrIntrouvable) {
		erreur(w, http.StatusNotFound, "Pays introuvable.")
		return
	}
	if err != nil {
		log.Println(err)
		erreur(w, http.StatusBadGateway, "Les données sont indisponibles pour le moment.")
		return
	}
	page(w, "fiche.html", map[string]any{"Pays": p})
}

// aPropos affiche la page « À propos »
func aPropos(w http.ResponseWriter, r *http.Request) {
	page(w, "apropos.html", nil)
}

// introuvable affiche une erreur 404 si la page n’existe pas.
func introuvable(w http.ResponseWriter, r *http.Request) {
	erreur(w, http.StatusNotFound, "Page introuvable.")
}

// main configure les routes et lance le serveur sur le port 8080.
func main() {
	mux := http.NewServeMux()
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	mux.HandleFunc("GET /{$}", accueil)
	mux.HandleFunc("GET /pays", catalogue)
	mux.HandleFunc("GET /pays/{code}", fiche)
	mux.HandleFunc("GET /a-propos", aPropos)
	mux.HandleFunc("/", introuvable) // toute autre adresse = 404

	log.Println("Serveur lancé sur http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
