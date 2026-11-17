package main

import (
	"html/template"
	"log"
	"net/http"
)

// page affiche une page : le gabarit commun (base.html) + le contenu de la page.
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

func accueil(w http.ResponseWriter, r *http.Request) {
	// Toute adresse inconnue arrive ici : on renvoie 404.
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		page(w, "erreur.html", map[string]any{"Code": 404, "Message": "Page introuvable"})
		return
	}

	pays, err := TousLesPays()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadGateway) // 502 : l'API est en panne
		page(w, "erreur.html", map[string]any{"Code": 502, "Message": "Les données sont indisponibles pour le moment."})
		return
	}
	page(w, "accueil.html", map[string]any{"NbPays": len(pays)})
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	mux.HandleFunc("GET /", accueil)

	log.Println("Serveur lancé sur http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
