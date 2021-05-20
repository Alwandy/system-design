package url

import (
	"bytes"
	"encoding/json"
	"github.com/Alwandy/system-design/models"
	"github.com/Alwandy/system-design/pkg/dynamodb"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
)

type Bitly struct {
	CreatedAt      string        `json:"created_at"`
	ID             string        `json:"id"`
	Link           string        `json:"link"`
	CustomBitlinks []interface{} `json:"custom_bitlinks"`
	LongURL        string        `json:"long_url"`
	Archived       bool          `json:"archived"`
	Tags           []interface{} `json:"tags"`
	Deeplinks      []interface{} `json:"deeplinks"`
	References     struct {
		Group string `json:"group"`
	} `json:"references"`
}

var bitlyToken = "5a4f5d9332f3eb753d19dcc5bf7fc636942dc4b9"

func NewUrlHandler(w http.ResponseWriter, r *http.Request) {
	u := models.Url{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&u); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	response, err := callBitly(u.Url)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	m := Bitly{}
	err = json.Unmarshal(response, &m)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	u.ShortenUrl = m.Link
	err = db.CreateItem(u)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, u)
	return
}


func GetUrlHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	res := models.Url{}
	res, err := db.GetItem(id)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	http.Redirect(w, r, res.Url, 307)
	return
}

func callBitly(url string) ([]byte, error) {
	var bearer = "Bearer " + bitlyToken
	requestBody, err := json.Marshal(map[string]string{
		"group_guid": "Bl5k9OfPBBf",
		"domain": "bit.ly",
		"long_url": url,
	})

	if err != nil {
		return nil, errors.New("Error with payload")
	}

	req, err := http.NewRequest("POST", "https://api-ssl.bitly.com/v4/shorten", bytes.NewBuffer(requestBody))
	req.Header.Set("Authorization", bearer)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New("Error with payload")
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	return body, nil

}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}