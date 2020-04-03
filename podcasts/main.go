package podcasts

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
)

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
func SearchEpisode(episodes []*Episode, query string) []*Episode {
	episodesMatched := make([]*Episode, 0)
	queryRegex := fmt.Sprintf(`(?i)%s`, query)
	for _, episode := range episodes {
		tb, _ := regexp.Match(queryRegex, []byte(episode.Title))
		db, _ := regexp.Match(queryRegex, []byte(episode.Description))
		if tb || db {
			episodesMatched = append(episodesMatched, episode)
		}
	}
	return episodesMatched
}
