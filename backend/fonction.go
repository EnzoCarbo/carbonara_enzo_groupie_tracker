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
	startIndex := (currentPage - 1) * cardsToDisplay
	endIndex := startIndex + cardsToDisplay
	if endIndex > len(cards) {
		endIndex = len(cards)
	}
	cardsToRender := cards[startIndex:endIndex]

	totalPages := int(math.Ceil(float64(len(cards)) / float64(cardsToDisplay)))

	pageInfo := PageInfo{
		TotalPages:   totalPages,
		CurrentPage:  currentPage,
		PreviousPage: currentPage - 1,
		NextPage:     currentPage + 1,
	}

	return cardsToRender, pageInfo
}

// Fonction pour gérer la requête de la liste des cartes
func DisplayCardListe(w http.ResponseWriter, r *http.Request) {
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

	if nbCartes != "" {
		cardsToDisplay, _ = strconv.Atoi(nbCartes)
	}

	if page != "" {
		currentPage, _ = strconv.Atoi(page)
	}

	// Récupérer les cartes depuis l'API
	cardsResponse, err := GetAllCards()
	if err != nil {
		http.Error(w, fmt.Sprintf("error fetching cards: %v", err), http.StatusInternalServerError)
		return
	}

	// Paginer les cartes
	cardsToRender, pageInfo := PaginatePage(cardsResponse.Data, currentPage, cardsToDisplay)

	// Afficher la liste des cartes paginée
	err = temp.ExecuteTemplate(w, "liste", struct {
		Cards        []Card
		PageInfo     PageInfo
		CardsPerPage int
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
	// Récupérer les paramètres de la requête
	categories := r.URL.Query()["categorie"]
	levels := r.URL.Query()["level"]
	attributes := r.URL.Query()["attribute"]

	if len(categories) == 0 && len(levels) == 0 && len(attributes) == 0 {
		http.Redirect(w, r, "liste", http.StatusSeeOther)
	}

	temp, err := template.ParseGlob("./templates/*.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("error parsing templates: %v", err), http.StatusInternalServerError)
		return
	}

	// Construire l'URL de l'API en fonction des paramètres de requête
	apiURL := buildAPIURL(categories, levels, attributes)

	// Récupérer les cartes depuis l'API
	cardsResponse, err := GetCards(apiURL)
	if err != nil {
		http.Error(w, fmt.Sprintf("error fetching cards: %v", err), http.StatusInternalServerError)
		return
	}

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
	var cardsResponse CardResponse

	apiURL := fmt.Sprintf("https://db.ygoprodeck.com/api/v7/cardinfo.php?fname=%s", query)
	response, err := http.Get(apiURL)
	if err != nil {
		return cardsResponse, errors.New("error fetching card data from the API")
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return cardsResponse, errors.New("error reading card data from the API response")
	}

	err = json.Unmarshal(body, &cardsResponse)
	if err != nil {
		return cardsResponse, nil
	}

	return cardsResponse, nil
}

func DisplayRecherche(w http.ResponseWriter, r *http.Request) {
	temp, err := template.ParseGlob("./templates/*.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("error parsing templates: %v", err), http.StatusInternalServerError)
		return
	}

	query := extractSearchQuery(r)
	if query == "" {
		http.Error(w, "Veuillez fournir un terme de recherche.", http.StatusBadRequest)
		return
	}

	cardsResponse, _ := GetCardsByQuery(query)

	err = temp.ExecuteTemplate(w, "recherche", cardsResponse)
	if err != nil {
		http.Error(w, fmt.Sprintf("error rendering HTML: %v", err), http.StatusInternalServerError)
		return
	}
}
