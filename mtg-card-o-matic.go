package main

//☀️ 💧💀🔥🌳

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type ImageUris struct {
	Small       string
	Normal      string
	Large       string
	Png         string
	Art_crop    string
	Border_crop string
}

type Legality struct {
	Standard  string
	Future    string
	Frontier  string
	Modern    string
	Legacy    string
	Pauper    string
	Vintage   string
	Penny     string
	Commander string
	OneVOne   string `json:"1v1"`
	Duel      string
	Brawl     string
}

type Related struct {
	Gather          string
	Tcgplayer_decks string
	Edhrec          string
	Mtgtop8         string
}

type Purchase struct {
	Amazon          string
	Ebay            string
	Tcgplayer       string
	magiccardmarket string
	Cardhoarder     string
	Card_kingdome   string
	Mtgo_traders    string
	Coolstuffinc    string
}

//Structure for Json Object returned from api.scryfall.com/card/random
type MtgCard struct {
	Object            string
	Id                string
	Oracle_id         string
	Multiverse_ids    []int
	Name              string
	Lang              string
	Uri               string
	Scryfall_uri      string
	Layout            string
	Highres_image     bool
	Image_uris        ImageUris
	Mana_cost         string
	Cmc               float64
	Type_line         string
	Oracle_text       string
	Power             string
	Toughness         string
	Colors            []string
	Color_identity    []string
	Legalities        Legality
	Reserved          bool
	Foil              bool
	Nonfoil           bool
	Oversized         bool
	Reprint           bool
	Set               string
	Set_name          string
	Set_uri           string
	Set_search_uri    string
	Scryfall_set_uri  string
	Rulings_uri       string
	Prints_search_uri string
	Collector_number  string
	Digital           bool
	Rarity            string
	Flavor_text       string
	Watermark         string
	Illustration_id   string
	Artist            string
	Frame             string
	Full_art          bool
	Border_color      string
	Timeshifted       bool
	Colorshifted      bool
	Futureshifted     bool
	Edhrec_rank       int
	Tix               string
	Related_uris      Related
	Purchase_uris     Purchase
}

var cardWidth = 40
var card_height = 30

var card string = `
+--------------------------------------+
|{THECREATURESVERYLONGNAME}  {MANACOST}| 
+--------------------------------------+
|                                      |
|                                      |
|                                      |
|                                      |
|                                      |
|                                      |
|                                      |
|                                      |
|                                      |
|                                      |
+--------------------------------------+
|{THECREATURETYPELINE}         {R}{SET}|
+--------------------------------------+
|{ORACLETEXT}                          |
|                                      |
|                                      |
|                                      |
|                                      |
|   --------------------------------   |
|{FLAVORTEXT}                          |
|                                      |
+--------------------------------------+
|                         {POWER/TOUGH}|
+--------------------------------------+
`

func spaceFiller(width int) string {
	var spaceString string
	for i := 0; i < width; i++ {
		spaceString = spaceString + " "
	}
	return spaceString
}

//Wraps text on spaces
func textWrapper(wrapText string, width int) string {
	var wrappedString string
	lineLength := 0

	splitText := strings.Split(wrapText, " ")

	for _, i := range splitText {
		if lineLength+len(i+" ") > width {
			wrappedString += "\n" + i
			lineLength = len(i + " ")
		} else {
			wrappedString += i + " "
			lineLength += len(i + " ")
		}

	}

	return wrappedString
}

func (mtgcard MtgCard) String() string {
	var formatedCard string

	formatedCard += mtgcard.Name + spaceFiller(cardWidth-(len(mtgcard.Name)+len((mtgcard.Mana_cost)))) + ManaSymbol(mtgcard.Mana_cost) + "\n"
	formatedCard += "------------" + "\n"
	formatedCard += mtgcard.Type_line + spaceFiller(cardWidth-(len(mtgcard.Type_line)+len(mtgcard.Rarity+" "+mtgcard.Set))) + mtgcard.Rarity + " " + mtgcard.Set + "\n"
	formatedCard += "------------" + "\n"
	formatedCard += textWrapper(ManaSymbol(mtgcard.Oracle_text), cardWidth) + "\n"
	formatedCard += "------------" + "\n"
	formatedCard += textWrapper(mtgcard.Flavor_text, cardWidth) + "\n"
	formatedCard += "------------" + "\n"
	formatedCard += spaceFiller(cardWidth-len(PowerToughnessFormat(mtgcard.Power, mtgcard.Toughness))) + PowerToughnessFormat(mtgcard.Power, mtgcard.Toughness) + "\n"
	return formatedCard
}

//pull from url to get json object and put it into a MtgCard struct
func getJsonCard(url string) (MtgCard, error) {
	response, err := http.Get(url)
	if err != nil {
		return MtgCard{}, err
	}
	defer response.Body.Close()

	respBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return MtgCard{}, err
	}

	mtgcard := MtgCard{}

	jsonErr := json.Unmarshal(respBody, &mtgcard)
	if jsonErr != nil {
		return MtgCard{}, jsonErr
	}

	return mtgcard, nil
}

//Uses SymbolMap to take the Mana Symbols returned and applys an emoji to it
func ManaSymbol(mana string) string {
	var SymbolMap = map[string]string{
		"{W}":  "☀️ ",
		"{U}":  "💧",
		"{B}":  "💀",
		"{R}":  "🔥",
		"{G}":  "🌳",
		"{T}":  "↩️ ",
		"{1}":  "1️⃣ ",
		"{2}":  "2️⃣ ",
		"{3}":  "3️⃣ ",
		"{4}":  "4️⃣ ",
		"{5}":  "5️⃣ ",
		"{6}":  "6️⃣ ",
		"{7}":  "7️⃣ ",
		"{8}":  "8️⃣ ",
		"{9}":  "9️⃣ ",
		"{10}": "🔟",
	}

	for stringSymbol, cardSymbol := range SymbolMap {
		mana = strings.Replace(mana, stringSymbol, cardSymbol, -1)
	}

	return mana
}

//Returns a blank string if empty, or formats it to N/N or X/X where appropriate
func PowerToughnessFormat(power, toughness string) string {
	powTough := ""
	if len(power) > 0 && len(toughness) > 0 {
		powTough = fmt.Sprintf("%s/%s", power, toughness)
	}

	return powTough
}

func main() {
	url := "https://api.scryfall.com/cards/random"

	mtgCard, err := getJsonCard(url)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(mtgCard)
	}
}
