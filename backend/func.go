package backend

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
)

func getQueryParameters(r *http.Request) (string, string) {
	nbCartes := r.URL.Query().Get("nb_cartes")
	page := r.URL.Query().Get("page")
	return nbCartes, page
}

func fetchCardData() CardResponse {
	apiURL := "https://db.ygoprodeck.com/api/v7/cardinfo.php"
	response, err := http.Get(apiURL)
	if err != nil {
		fmt.Println("Error:", err)
		return CardResponse{}
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return CardResponse{}
	}

	var cardsResponse CardResponse
	err = json.Unmarshal(body, &cardsResponse)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return CardResponse{}
	}

	return cardsResponse
}

func paginateCards(cardsResponse CardResponse, nbCartes, page string) []Card {
	// Convertir le nombre de cartes par page et la page actuelle en entiers
	cardsToDisplay, _ := strconv.Atoi(nbCartes)
	currentPage, _ := strconv.Atoi(page)

	// Déterminer l'indice de début et de fin pour la pagination
	startIndex := (currentPage - 1) * cardsToDisplay
	endIndex := startIndex + cardsToDisplay

	// Vérifier si l'indice de fin dépasse le nombre total de cartes
	if endIndex > len(cardsResponse.Data) {
		endIndex = len(cardsResponse.Data)
	}

	// Récupérer les cartes à afficher pour la page actuelle
	cardsToRender := cardsResponse.Data[startIndex:endIndex]

	return cardsToRender
}

func calculatePaginationInfo(cardsResponse CardResponse, nbCartes, page string) interface{} {
	// Convertir le nombre de cartes par page en entier
	cardsToDisplay, _ := strconv.Atoi(nbCartes)

	// Calculer le nombre total de pages en fonction du nombre total de cartes et du nombre de cartes par page
	totalPages := int(math.Ceil(float64(len(cardsResponse.Data)) / float64(cardsToDisplay)))

	// Convertir la page actuelle en entier
	currentPage, _ := strconv.Atoi(page)

	// Calculer la page précédente et suivante
	previousPage := currentPage - 1
	nextPage := currentPage + 1

	// Si la page précédente est inférieure à 1, la définir sur 1
	if previousPage < 1 {
		previousPage = 1
	}

	// Si la page suivante est supérieure au nombre total de pages, la définir sur le nombre total de pages
	if nextPage > totalPages {
		nextPage = totalPages
	}

	// Créer une structure contenant les informations de pagination
	pageInfo := struct {
		TotalPages   int
		CurrentPage  int
		PreviousPage int
		NextPage     int
	}{
		TotalPages:   totalPages,
		CurrentPage:  currentPage,
		PreviousPage: previousPage,
		NextPage:     nextPage,
	}

	return pageInfo
}
