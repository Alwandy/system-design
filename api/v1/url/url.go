package url

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type urlJson struct {
	Url string `json:"url"`
}

func NewUrlHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var url urlJson
	err := json.NewDecoder(r.Body).Decode(url)
	if err != nil {
		fmt.Fprintf(w, "%s", err)
	}
	fmt.Fprintf(w, "%s", url)
}


func GetUrl(w http.ResponseWriter, r *http.Request) {

}
