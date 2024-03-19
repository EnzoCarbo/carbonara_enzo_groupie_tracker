package handler

import (
	"apiperso/backend"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"text/template"
)

var cardsResponse backend.CardResponse

func HandlerMain(w http.ResponseWriter, r *http.Request) {
	temp, err := template.ParseGlob("./templates/*.html")
	if err != nil {
		fmt.Printf("ERREUR => %s", err.Error())
		return
	}
	temp.ExecuteTemplate(w, "main", nil)
}

func HandlerListe(w http.ResponseWriter, r *http.Request) {
	temp, err := template.ParseGlob("./templates/*.html")
	if err != nil {
		fmt.Printf("ERREUR => %s", err.Error())
		return
	}
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

	apiURL := "https://db.ygoprodeck.com/api/v7/cardinfo.php"

	response, err := http.Get(apiURL)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	err = json.Unmarshal(body, &cardsResponse)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}

	// Paginer les cartes
	startIndex := (currentPage - 1) * cardsToDisplay
	endIndex := startIndex + cardsToDisplay
	if endIndex > len(cardsResponse.Data) {
		endIndex = len(cardsResponse.Data)
	}
	cardsToRender := cardsResponse.Data[startIndex:endIndex]

	// Calculer le nombre total de pages
	totalPages := int(math.Ceil(float64(len(cardsResponse.Data)) / float64(cardsToDisplay)))

	// Passer les informations de pagination au modèle
	pageInfo := struct {
		TotalPages   int
		CurrentPage  int
		PreviousPage int
		NextPage     int
	}{
		TotalPages:   totalPages,
		CurrentPage:  currentPage,
		PreviousPage: currentPage - 1,
		NextPage:     currentPage + 1,
	}

	// Afficher la liste des cartes paginée
	err = temp.ExecuteTemplate(w, "liste", struct {
		Cards        []backend.Card
		PageInfo     interface{}
		CardsPerPage int
	}{
		Cards:        cardsToRender,
		PageInfo:     pageInfo,
		CardsPerPage: cardsToDisplay,
	})
	if err != nil {
		fmt.Println("Error rendering HTML:", err)
		return
	}
}

func HandlerInfo(w http.ResponseWriter, r *http.Request) {
	temp, err := template.ParseGlob("./templates/*.html")
	if err != nil {
		fmt.Printf("ERREUR => %s", err.Error())
		return
	}
	cardID := r.URL.Path[len("/info/"):]

	apiURL := fmt.Sprintf("https://db.ygoprodeck.com/api/v7/cardinfo.php?id=%s", cardID)

	response, err := http.Get(apiURL)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	err = json.Unmarshal(body, &cardsResponse)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}

	err = temp.ExecuteTemplate(w, "info", cardsResponse)
	if err != nil {
		fmt.Println("Error rendering HTML:", err)
		return
	}
}

func HandlerCategorie(w http.ResponseWriter, r *http.Request) {
	temp, err := template.ParseGlob("./templates/*.html")
	if err != nil {
		fmt.Printf("ERREUR => %s", err.Error())
		return
	}
	categories := r.URL.Query()["categorie"]
	levels := r.URL.Query()["level"]
	attributes := r.URL.Query()["attribute"]

	if len(categories) == 0 && len(levels) == 0 && len(attributes) == 0 {
		http.Error(w, "Veuillez fournir au moins une option de tri.", http.StatusBadRequest)
		return
	}

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

	response, err := http.Get(apiURL)
	if err != nil {
		http.Error(w, "Erreur lors de la requête à l'API", http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		http.Error(w, "Erreur lors de la lecture de la réponse de l'API", http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, &cardsResponse)
	if err != nil {
		http.Error(w, "Erreur lors de la désérialisation JSON de la réponse de l'API", http.StatusInternalServerError)
		return
	}

	fmt.Println(categories, levels, attributes)

	err = temp.ExecuteTemplate(w, "categorie", cardsResponse)
	if err != nil {
		fmt.Println("Error rendering HTML:", err)
		return
	}
}

var userDeck backend.Deck

func HandlerDeck(w http.ResponseWriter, r *http.Request) {
	temp, err := template.ParseGlob("./templates/*.html")
	if err != nil {
		fmt.Printf("ERREUR => %s", err.Error())
		return
	}
	err = temp.ExecuteTemplate(w, "deck", userDeck)
	if err != nil {
		fmt.Println("Error rendering HTML:", err)
		return
	}
}

func HandlerDeckRemove(w http.ResponseWriter, r *http.Request) {
	cardID := r.URL.Path[len("/deck/remove/"):]

	// Convertir cardID en entier
	id, err := strconv.Atoi(cardID)
	if err != nil {
		http.Error(w, "Invalid card ID", http.StatusBadRequest)
		return
	}

	// Recherche de la carte correspondante dans le deck
	var cardIndex int = -1
	for i, card := range userDeck.Cards {
		if card.ID == id {
			cardIndex = i
			break
		}
	}

	// Si la carte est trouvée dans le deck, la supprimer
	if cardIndex != -1 {
		userDeck.Cards = append(userDeck.Cards[:cardIndex], userDeck.Cards[cardIndex+1:]...)
	}

	// Redirection vers la page /deck après la suppression de la carte du deck
	http.Redirect(w, r, "/deck", http.StatusSeeOther)
}

func HandlerDeckAdd(w http.ResponseWriter, r *http.Request) {
	cardID := r.URL.Path[len("/deck/add/"):]

	// Recherche de la carte correspondante dans les données récupérées
	var selectedCard backend.Card
	for _, card := range cardsResponse.Data {
		if strconv.Itoa(card.ID) == cardID {
			selectedCard = card
			break
		}
	}

	// Ajout de la carte au deck uniquement si elle n'est pas déjà présente
	cardAlreadyInDeck := false
	for _, card := range userDeck.Cards {
		if card.ID == selectedCard.ID {
			cardAlreadyInDeck = true
			break
		}
	}
	if !cardAlreadyInDeck {
		userDeck.Cards = append(userDeck.Cards, selectedCard)
	}

	// Redirection vers la page /deck après l'ajout au deck
	http.Redirect(w, r, "/deck", http.StatusSeeOther)
}
