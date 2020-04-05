package podcasts

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type Podcast struct {
	URL           string `json:"mygpo_link"`
	Title         string `json:"title"`
	Author        string `json:"author"`
	Description   string `json:"description"`
	LogoURL       string `json:"logo_url"`
	Website       string `json:"website"`
	episodes      []*Episode
	episodeIndex  int
	episodesCount int
}

func (p *Podcast) String() string {
	res, err := json.Marshal(p)
	if err != nil {
		return "failed"
	}
	return string(res)
}

func GetTop50Podcats() (podcasts []*Podcast) {
	response, err := http.Get("https://gpodder.net/toplist/50.json")
	if err != nil {
		log.Fatal(err)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	json.Unmarshal(responseData, &podcasts)
	return podcasts
}

func FindPodcasts(query string) (podcasts []*Podcast) {
	URL := fmt.Sprintf("https://gpodder.net/search.json?q=%s", url.QueryEscape(query))
	response, err := http.Get(URL)
	if err != nil {
		log.Fatal(err)
	}
	responseData, err := ioutil.ReadAll(response.Body)
	if response.StatusCode != 200 {
		log.Fatal(responseData)
	}
	if err != nil {
		log.Fatal(err)
	}

	json.Unmarshal(responseData, &podcasts)
	return podcasts
}
