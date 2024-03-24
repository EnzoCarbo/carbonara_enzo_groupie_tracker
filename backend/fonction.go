package backend

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"text/template"
)

func UnmarshalData(response *http.Response, target interface{}) error {
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("error reading API response: %v", err)
	}

	err = json.Unmarshal(body, target)
	if err != nil {
		return fmt.Errorf("error unmarshalling JSON response: %v", err)
	}

	return nil
}

// Récupère toutes les cartes de l'API
func GetAllCards() (CardResponse, error) {
	var cardsResponse CardResponse

	apiURL := "https://db.ygoprodeck.com/api/v7/cardinfo.php"

	response, err := http.Get(apiURL)
	if err != nil {
		return cardsResponse, fmt.Errorf("error fetching cards from API: %v", err)
	}
	defer response.Body.Close()

	UnmarshalData(response, &cardsResponse)
	return cardsResponse, nil

}

// Récupère les informations d'une carte spécifique de l'API
func GetInfoCards(cardID string) (CardResponse, error) {
	var cardsResponse CardResponse

	apiURL := fmt.Sprintf("https://db.ygoprodeck.com/api/v7/cardinfo.php?id=%s", cardID)

	response, err := http.Get(apiURL)
	if err != nil {
		return cardsResponse, fmt.Errorf("error fetching card information from API: %v", err)
	}
	defer response.Body.Close()

	UnmarshalData(response, &cardsResponse)

	return cardsResponse, nil
}

func GetCards(apiURL string) (CardResponse, error) {
	var cardsResponse CardResponse

	response, err := http.Get(apiURL)
	if err != nil {
		return cardsResponse, fmt.Errorf("error fetching cards from API: %v", err)
	}
	defer response.Body.Close()

	UnmarshalData(response, &cardsResponse)

	return cardsResponse, nil
}

// Fonction pour construire l'URL de l'API en fonction des paramètres de requête
func buildAPIURL(categories, levels, attributes []string) string {
	apiURL := "https://db.ygoprodeck.com/api/v7/cardinfo.php?"

	for _, category := range categories {
		apiURL += "&type=" + category
	}

	for _, level := range levels {
		apiURL += "&level=" + level
	}

	for _, attribute := range attributes {
		apiURL += "&attribute=" + attribute
	}

	return apiURL
}

// Fonction pour paginer les cartes
func PaginatePage(cards []Card, currentPage, cardsToDisplay int) ([]Card, PageInfo) {
	// Calcule l'indice de départ et l'indice de fin des cartes à afficher sur la page actuelle
	startIndex := (currentPage - 1) * cardsToDisplay
	endIndex := startIndex + cardsToDisplay
	if endIndex > len(cards) {
		endIndex = len(cards)
	}
	// Sélectionne les cartes à afficher sur la page actuelle en fonction des indices calculés
	cardsToRender := cards[startIndex:endIndex]

	// Calcule le nombre total de pages en fonction du nombre total de cartes et du nombre de cartes à afficher par page
	totalPages := int(math.Ceil(float64(len(cards)) / float64(cardsToDisplay)))

	// Crée une structure PageInfo contenant des informations sur la pagination
	pageInfo := PageInfo{
		TotalPages:   totalPages,      // Nombre total de pages
		CurrentPage:  currentPage,     // Page actuelle
		PreviousPage: currentPage - 1, // Page précédente
		NextPage:     currentPage + 1, // Page suivante
	}

	// Retourne les cartes à afficher sur la page actuelle et les informations sur la pagination
	return cardsToRender, pageInfo
}

func DisplayCardListe(w http.ResponseWriter, r *http.Request) {
	// Parse les fichiers de modèle HTML dans le répertoire "templates"
	temp, err := template.ParseGlob("./templates/*.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("error parsing templates: %v", err), http.StatusInternalServerError)
		return
	}

	// Récupérer les paramètres de la requête
	nbCartes := r.URL.Query().Get("nb_cartes")
	page := r.URL.Query().Get("page")
	cardsToDisplay := 20
	currentPage := 1

	// Convertir les paramètres de requête en entiers
	if nbCartes != "" {
		cardsToDisplay, _ = strconv.Atoi(nbCartes)
	}

	if page != "" {
		currentPage, _ = strconv.Atoi(page)
	}

	// Récupérer toutes les cartes depuis l'API
	cardsResponse, err := GetAllCards()
	if err != nil {
		http.Error(w, fmt.Sprintf("error fetching cards: %v", err), http.StatusInternalServerError)
		return
	}

	// Paginer les cartes en fonction des paramètres de requête
	cardsToRender, pageInfo := PaginatePage(cardsResponse.Data, currentPage, cardsToDisplay)

	// Afficher la liste paginée des cartes
	err = temp.ExecuteTemplate(w, "liste", struct {
		Cards        []Card   // Les cartes à afficher sur la page actuelle
		PageInfo     PageInfo // Informations sur la pagination
		CardsPerPage int      // Nombre de cartes à afficher par page
	}{
		Cards:        cardsToRender,
		PageInfo:     pageInfo,
		CardsPerPage: cardsToDisplay,
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("error rendering HTML: %v", err), http.StatusInternalServerError)
		return
	}
}

// Fonction pour gérer la requête d'informations sur une carte spécifique
func DisplayCardInfo(w http.ResponseWriter, r *http.Request) {
	temp, err := template.ParseGlob("./templates/*.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("error parsing templates: %v", err), http.StatusInternalServerError)
		return
	}

	// Extraire l'ID de la carte du chemin URL
	cardID := GetCardId(r.URL.Path)

	// Récupérer les informations sur la carte depuis l'API
	cardsResponse, err := GetInfoCards(cardID)
	if err != nil {
		http.Error(w, fmt.Sprintf("error fetching card information: %v", err), http.StatusInternalServerError)
		return
	}

	// Rendre le modèle HTML avec les informations de la carte
	err = temp.ExecuteTemplate(w, "info", cardsResponse)
	if err != nil {
		http.Error(w, fmt.Sprintf("error rendering HTML: %v", err), http.StatusInternalServerError)
		return
	}
}

func DisplayCategorie(w http.ResponseWriter, r *http.Request) {
	// Récupère les valeurs des paramètres de requête
	categories := r.URL.Query()["categorie"]
	levels := r.URL.Query()["level"]
	attributes := r.URL.Query()["attribute"]

	// Vérifie si aucun paramètre de requête n'a été fourni, si c'est le cas, redirige vers une autre page
	if len(categories) == 0 && len(levels) == 0 && len(attributes) == 0 {
		http.Redirect(w, r, "liste", http.StatusSeeOther)
		return
	}

	// Parse les fichiers de modèle HTML dans le répertoire "templates"
	temp, err := template.ParseGlob("./templates/*.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("error parsing templates: %v", err), http.StatusInternalServerError)
		return
	}

	// Construit l'URL de l'API en fonction des paramètres de requête
	apiURL := buildAPIURL(categories, levels, attributes)

	// Récupère les cartes depuis l'API
	cardsResponse, err := GetCards(apiURL)
	if err != nil {
		http.Error(w, fmt.Sprintf("error fetching cards: %v", err), http.StatusInternalServerError)
		return
	}

	// Vérifie si aucun en-tête Content-Type n'a été défini dans la réponse
	if w.Header().Get("Content-Type") == "" {
		// Aucune donnée de réponse n'a été écrite, donc nous pouvons afficher les cartes
		err = temp.ExecuteTemplate(w, "categorie", cardsResponse)
		if err != nil {
			http.Error(w, fmt.Sprintf("error rendering HTML: %v", err), http.StatusInternalServerError)
			return
		}
	}
}

func GetCardId(path string) string {
	return path[len("/info/"):]
}

func extractSearchQuery(r *http.Request) string {
	return r.URL.Query().Get("query")
}

func GetCardsByQuery(query string) (CardResponse, error) {
	// Crée une variable pour stocker la réponse des cartes
	var cardsResponse CardResponse

	// Construit l'URL de l'API en utilisant le terme de recherche fourni
	apiURL := fmt.Sprintf("https://db.ygoprodeck.com/api/v7/cardinfo.php?fname=%s", query)

	// Effectue une requête HTTP GET à l'URL de l'API
	response, err := http.Get(apiURL)
	if err != nil {
		// En cas d'erreur lors de la récupération des données de carte depuis l'API, retourne une erreur avec un message approprié
		return cardsResponse, errors.New("error fetching card data from the API")
	}
	defer response.Body.Close()

	// Lit le corps de la réponse HTTP
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		// En cas d'erreur lors de la lecture des données de carte à partir de la réponse HTTP, retourne une erreur avec un message approprié
		return cardsResponse, errors.New("error reading card data from the API response")
	}

	// Décode les données JSON du corps de la réponse dans la structure cardsResponse
	err = json.Unmarshal(body, &cardsResponse)
	if err != nil {
		// En cas d'erreur lors du décodage JSON, retourne une erreur avec un message approprié
		return cardsResponse, err
	}

	// Retourne la réponse des cartes et aucune erreur
	return cardsResponse, nil
}

func DisplayRecherche(w http.ResponseWriter, r *http.Request) {
	// Parse les fichiers de modèle HTML dans le répertoire "templates"
	temp, err := template.ParseGlob("./templates/*.html")
	if err != nil {
		// En cas d'erreur lors de l'analyse des modèles, renvoie une erreur interne au serveur avec un message approprié
		http.Error(w, fmt.Sprintf("error parsing templates: %v", err), http.StatusInternalServerError)
		return
	}

	// Extrait le terme de recherche de la requête HTTP
	query := extractSearchQuery(r)
	if query == "" {
		// Si aucun terme de recherche n'est fourni, renvoie une erreur de mauvaise requête avec un message approprié
		http.Error(w, "Veuillez fournir un terme de recherche.", http.StatusBadRequest)
		return
	}

	// Récupère les cartes en fonction du terme de recherche
	cardsResponse, _ := GetCardsByQuery(query)

	// Exécute le modèle "recherche" avec les données des cartes et écrit le résultat dans la réponse HTTP
	err = temp.ExecuteTemplate(w, "recherche", cardsResponse)
	if err != nil {
		// En cas d'erreur lors du rendu HTML, renvoie une erreur interne au serveur avec un message approprié
		http.Error(w, fmt.Sprintf("error rendering HTML: %v", err), http.StatusInternalServerError)
		return
	}
}
