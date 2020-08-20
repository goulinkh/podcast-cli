package rss

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"strconv"

	itunesapi "github.com/goulinkh/podcast-cli/itunes-api"
)

type Rss struct {
	Channel struct {
		Item []struct {
			Title       string `xml:"title"`
			PubDate     string `xml:"pubDate"`
			Author      string `xml:"author"`
			Description string `xml:"description"`
			Image       struct {
				Href string `xml:"href,attr"`
			} `xml:"image"`
			Enclosure struct {
				URL    string `xml:"url,attr"`
				Length int    `xml:"length,attr"`
				Type   string `xml:"type,attr"`
			} `xml:"enclosure"`
			Duration int `xml:"duration"`
		} `xml:"item"`
	} `xml:"channel"`
}

func ParseEpisodes(rssUrl string) ([]*itunesapi.Episode, error) {
	resp, err := http.Get(rssUrl)
	if err != nil {
		return nil, err
	}

	rss, err := ioutil.ReadAll(resp.Body)
	var podcast Rss
	err = xml.Unmarshal(rss, &podcast)
	episodes := make([]*itunesapi.Episode, len(podcast.Channel.Item))
	for i, e := range podcast.Channel.Item {
		episodes[i] =
			&itunesapi.Episode{
				Artwork:                e.Image.Href,
				AudioURL:               e.Enclosure.URL,
				Description:            e.Description,
				DurationInMilliseconds: e.Duration * 1000,
				Id:                     strconv.Itoa(i),
				ReleaseDate:            e.PubDate,
				Title:                  e.Title,
			}
	}
	return episodes, nil
}
