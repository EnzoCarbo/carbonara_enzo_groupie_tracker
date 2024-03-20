package backend

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

type Deck struct {
	Cards []Card
}

type PageInfo struct {
	TotalPages   int
	CurrentPage  int
	PreviousPage int
	NextPage     int
}
