package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
)

type CardSet struct {
	SetName   string `json:"set_name"`
	SetCode   string `json:"set_code"`
	SetRarity string `json:"set_rarity"`
	SetPrice  string `json:"set_price"`
}

type CardImage struct {
	ID              int    `json:"id"`
	ImageURL        string `json:"image_url"`
	ImageURLSmall   string `json:"image_url_small"`
	ImageURLCropped string `json:"image_url_cropped"`
}

type CardPrice struct {
	CardmarketPrice   string `json:"cardmarket_price"`
	TcgplayerPrice    string `json:"tcgplayer_price"`
	EbayPrice         string `json:"ebay_price"`
	AmazonPrice       string `json:"amazon_price"`
	CoolstuffincPrice string `json:"coolstuffinc_price"`
}

type Card struct {
	ID            int         `json:"id"`
	Name          string      `json:"name"`
	Type          string      `json:"type"`
	FrameType     string      `json:"frameType"`
	Description   string      `json:"desc"`
	Atk           int         `json:"atk"`
	Def           int         `json:"def"`
	Level         int         `json:"level"`
	Race          string      `json:"race"`
	Attribute     string      `json:"attribute"`
	Archetype     string      `json:"archetype"`
	YgoProDeckURL string      `json:"ygoprodeck_url"`
	CardSets      []CardSet   `json:"card_sets"`
	CardImages    []CardImage `json:"card_images"`
	CardPrices    []CardPrice `json:"card_prices"`
}

type CardResponse struct {
	Data []Card `json:"data"`
}

type SearchRequest struct {
	Query string `json:"query"`
}

func main() {
	temp, err := template.ParseGlob("./templates/*.html")
	if err != nil {
		fmt.Printf("ERREUR => %s", err.Error())
		return
	}

	http.HandleFunc("/main", func(w http.ResponseWriter, r *http.Request) {
		temp.ExecuteTemplate(w, "main", nil)
	})

	http.HandleFunc("/liste", func(w http.ResponseWriter, r *http.Request) {
		var cardsResponse CardResponse

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

		if len(cardsResponse.Data) > 54 {
			cardsResponse.Data = cardsResponse.Data[:54]
		}

		err = temp.ExecuteTemplate(w, "liste", cardsResponse)
		if err != nil {
			fmt.Println("Error rendering HTML:", err)
			return
		}
	})

	http.HandleFunc("/recherche", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("query")

		if query == "" {
			fmt.Fprint(w, "Veuillez fournir un terme de recherche.")
			return
		}

		apiURL := fmt.Sprintf("https://db.ygoprodeck.com/api/v7/cardinfo.php?fname=%s", query)

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

		var cardsResponse CardResponse
		err = json.Unmarshal(body, &cardsResponse)
		if err != nil {
			fmt.Println("Error unmarshalling JSON:", err)
			return
		}

		err = temp.ExecuteTemplate(w, "recherche", cardsResponse)
		if err != nil {
			fmt.Println("Error rendering HTML:", err)
			return
		}
	})

	http.HandleFunc("/info/", func(w http.ResponseWriter, r *http.Request) {
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

		var cardsResponse CardResponse
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
	})

	http.HandleFunc("/categorie", func(w http.ResponseWriter, r *http.Request) {
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

		var cardsResponse CardResponse
		err = json.Unmarshal(body, &cardsResponse)
		if err != nil {
			http.Error(w, "Erreur lors de la désérialisation JSON de la réponse de l'API", http.StatusInternalServerError)
			return
		}

		if len(cardsResponse.Data) > 54 {
			cardsResponse.Data = cardsResponse.Data[:54]
		}

		fmt.Println(categories, levels, attributes)

		err = temp.ExecuteTemplate(w, "categorie", cardsResponse)
		if err != nil {
			fmt.Println("Error rendering HTML:", err)
			return
		}
	})

	rootDoc, _ := os.Getwd()
	fileserver := http.FileServer(http.Dir(rootDoc + "/assets"))
	http.Handle("/static/", http.StripPrefix("/static/", fileserver))
	http.ListenAndServe(":8080", nil)
}
