package itunesapi

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/tidwall/gjson"
)

type Podcast struct {
	Title       string `json:"title"`
	URL         string `json:"url`
	Id          string `json:"id"`
	Description string `json:"description"`
	Author      string `json:"author"`
}

func FindPodcasts(query string) ([]*Podcast, error) {
	authorization, err := getAuthorization()
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("GET", fmt.Sprintf("https://itunes.apple.com/search?country=us&entity=podcast&term=%s", query), nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Authorization", authorization)
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	podcastsJSON := gjson.Get(string(data), "results").Array()
	podcasts := make([]*Podcast, 0)
	for _, podcast := range podcastsJSON {
		if podcast.Get("kind").String() == "podcast" && podcast.Get("wrapperType").String() == "track" {
			podcasts = append(podcasts, &Podcast{
				Author:      podcast.Get("artistName").String(),
				Description: "",
				Id:          podcast.Get("trackId").String(),
				Title:       podcast.Get("collectionName").String(),
				URL:         regexp.MustCompile(`\?.*$`).ReplaceAllString(podcast.Get("collectionViewUrl").String(), ""),
			})
		}
	}
	return podcasts, nil
}
func (p *Podcast) GetEpisodes() ([]*Episode, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://amp-api.podcasts.apple.com/v1/catalog/us/podcasts/%s/episodes?offset=0&limit=300", p.Id), nil)
	if err != nil {
		return nil, err

	}
	authorization, err := getAuthorization()
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", authorization)
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err

	}
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err

	}
	episodesJSON := gjson.Get(string(data), "data").Array()
	episodes := make([]*Episode, 0)
	for _, episode := range episodesJSON {
		if episode.Get(`attributes.mediaKind`).String() == "audio" {

			episodes = append(episodes, &Episode{
				Id:                     episode.Get(`id`).String(),
				Artwork:                episode.Get(`attributes.artwork.url`).String(),
				Title:                  episode.Get(`attributes.name`).String(),
				AudioURL:               episode.Get(`attributes.assetUrl`).String(),
				ReleaseDate:            episode.Get(`attributes.releaseDateTime`).String(),
				DurationInMilliseconds: int(episode.Get(`attributes.durationInMilliseconds`).Int()),
				Description:            episode.Get(`attributes.description.standard`).String(),
			})
		}
	}
	return episodes, nil
}

func (g *Genre) GetPodcasts() ([]*Podcast, error) {
	request, err := http.NewRequest("GET", "https://amp-api.podcasts.apple.com/v1/catalog/us/charts?types=podcasts&limit=200&genre="+g.Id, nil)
	if err != nil {
		return nil, err
	}
	authorization, err := getAuthorization()
	if err != nil {
		return nil, err
	}
	request.Header.Add("Authorization", authorization)
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	podcastsJSON := gjson.Get(string(data), "results.podcasts.0.data").Array()
	podcasts := make([]*Podcast, len(podcastsJSON))
	for i, podcast := range podcastsJSON {
		podcasts[i] = &Podcast{
			Id:          podcast.Get("id").String(),
			Description: podcast.Get("attributes.description.standard").String(),
			Title:       podcast.Get("attributes.name").String(),
			URL:         "https://amp-api.podcasts.apple.com" + podcast.Get("href").String(),
			Author:      podcast.Get("attributes.artistName").String(),
		}
	}
	return podcasts, nil
}

func getAuthorization() (string, error) {
	if authorization != "" {
		return authorization, nil
	}
	authRegEx := regexp.MustCompile("privateKeyPath.+token%22%3A%22(?P<Bearer>.*?)%22%7D%2C%22")
	resp, err := http.Get("https://podcasts.apple.com/us/podcast/the-joe-rogan-experience/id360084272")
	if err != nil {
		return "", err
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	match := authRegEx.FindStringSubmatch(string(data))
	if len(match) != 2 {
		return "", errors.New("Authorization access token is not found")
	}
	authorization = "Bearer " + match[1]
	return authorization, nil
}
