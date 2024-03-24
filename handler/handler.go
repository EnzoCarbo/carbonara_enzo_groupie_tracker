package handler

import (
	"apiperso/backend"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"text/template"
)

func HandlerMain(w http.ResponseWriter, r *http.Request) {
	temp, err := template.ParseGlob("./templates/*.html")
	if err != nil {
		fmt.Printf("ERREUR => %s", err.Error())
		return
	}
	temp.ExecuteTemplate(w, "main", nil)
}

func HandlerListe(w http.ResponseWriter, r *http.Request) {
	backend.DisplayCardListe(w, r)
}

func HandlerInfo(w http.ResponseWriter, r *http.Request) {
	backend.DisplayCardInfo(w, r)
}

func HandlerCategorie(w http.ResponseWriter, r *http.Request) {
	backend.DisplayCategorie(w, r)
}

// Variable globale pour stocker le deck de l'utilisateur
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

func HandlerRecherche(w http.ResponseWriter, r *http.Request) {
	backend.DisplayRecherche(w, r)
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/404", http.StatusSeeOther)
}

func Handler404(w http.ResponseWriter, r *http.Request) {
	temp, err := template.ParseGlob("./templates/*.html")
	if err != nil {
		fmt.Printf("ERREUR => %s", err.Error())
		return
	}
	err = temp.ExecuteTemplate(w, "404", nil)
	if err != nil {
		fmt.Println("Error rendering HTML:", err)
		return
	}
}

func HandlerDeckRemove(w http.ResponseWriter, r *http.Request) {
	// Extraire l'ID de la carte à supprimer à partir de l'URL
	cardID := r.URL.Path[len("/deck/remove/"):]

	// Convertir l'ID de la carte en entier
	id, err := strconv.Atoi(cardID)
	if err != nil {
		http.Error(w, "Invalid card ID", http.StatusBadRequest)
		return
	}

	// Recherche de la carte correspondante dans le deck de l'utilisateur
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

// Nombre maximum de fois où une carte peut être ajoutée au deck
const maxCardCount = 3

func HandlerDeckAdd(w http.ResponseWriter, r *http.Request) {
	// Extraire l'ID de la carte à ajouter à partir de l'URL
	cardID := r.URL.Path[len("/deck/add/"):]

	// Convertir l'ID de la carte en entier
	id, err := strconv.Atoi(cardID)
	if err != nil {
		http.Error(w, "Invalid card ID", http.StatusBadRequest)
		return
	}

	// Récupérer les données de la carte depuis l'API
	apiURL := fmt.Sprintf("https://db.ygoprodeck.com/api/v7/cardinfo.php?id=%d", id)
	response, err := http.Get(apiURL)
	if err != nil {
		http.Error(w, "Error fetching card data from the API", http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	var cardResponse backend.CardResponse
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		http.Error(w, "Error reading card data from the API response", http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, &cardResponse)
	if err != nil {
		http.Error(w, "Error unmarshalling card data from the API response", http.StatusInternalServerError)
		return
	}

	if len(cardResponse.Data) == 0 {
		http.Error(w, "Card data not found in API response", http.StatusNotFound)
		return
	}

	// Sélectionner la première carte de la réponse
	selectedCard := cardResponse.Data[0]

	// Vérifier si la carte est déjà dans le deck
	var cardCountInDeck int
	for _, card := range userDeck.Cards {
		if card.ID == selectedCard.ID {
			cardCountInDeck++
		}
	}

	// Si la carte est déjà dans le deck 3 fois, rediriger sans l'ajouter
	if cardCountInDeck >= maxCardCount {
		http.Redirect(w, r, "/deck", http.StatusSeeOther)
		return
	}

	// Ajouter la carte sélectionnée au deck
	userDeck.Cards = append(userDeck.Cards, selectedCard)

	// Redirection vers la page /deck après l'ajout de la carte
	http.Redirect(w, r, "/deck", http.StatusSeeOther)
}

func HandlerAboutUs(w http.ResponseWriter, r *http.Request) {
	temp, err := template.ParseGlob("./templates/*.html")
	if err != nil {
		fmt.Printf("ERREUR => %s", err.Error())
		return
	}
	err = temp.ExecuteTemplate(w, "aboutus", userDeck)
	if err != nil {
		fmt.Println("Error rendering HTML:", err)
		return
	}
}
