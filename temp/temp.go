package temp

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/PuerkitoBio/goquery"
	"github.com/tidwall/gjson"
)

const applePodcastsMainPage = "https://podcasts.apple.com/genre/podcasts/id26"

var (
	authorization string
)

type Genre struct {
	Text     string   `json:"text"`
	URL      string   `json:"url"`
	Id       string   `json:"url"`
	SubGenre []*Genre `json:"sub-genres"`
}

type Podcast struct {
	Title       string `json:"title"`
	URL         string `json:"url`
	Id          string `json:"id"`
	Description string `json:"description"`
	Author      string `json:"author"`
}
type Episode struct {
	Id                     string `json:"id"`
	Artwork                string `json:"artwork"`
	Title                  string `json:"title"`
	AudioURL               string `json:"audiourl"`
	ReleaseDate            string `json:"releasedate"`
	DurationInMilliseconds string `json:"duratioInMilliseconds"`
	Description            string `json:"description"`
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
	episodes := make([]*Episode, len(episodesJSON))
	for i, episode := range episodesJSON {
		episodes[i] = &Episode{
			Id:                     episode.Get(`id`).String(),
			Artwork:                episode.Get(`attributes.artwork.url`).String(),
			Title:                  episode.Get(`attributes.name`).String(),
			AudioURL:               episode.Get(`attributes.assetUrl`).String(),
			ReleaseDate:            episode.Get(`attributes.releaseDateTime`).String(),
			DurationInMilliseconds: episode.Get(`attributes.durationInMilliseconds`).String(),
			Description:            episode.Get(`attributes.description.standard`).String(),
		}
	}
	return episodes, nil
}

func getGenreId(url string) string {
	idRegExp := regexp.MustCompile(`\d+$`)
	return idRegExp.FindString(url)
}

func GetGenres() ([]*Genre, error) {
	resp, err := http.Get(applePodcastsMainPage)
	if err != nil {
		return nil, err
	}
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return nil, err
	}
	genres := make([]*Genre, 0)
	doc.Find(".list.column > li").Each(func(i int, s *goquery.Selection) {
		genreSelection := s.Find(".top-level-genre")
		href, exists := genreSelection.Attr("href")
		if !exists {
			return
		}

		genre := &Genre{Text: genreSelection.Text(), URL: href, SubGenre: make([]*Genre, 0), Id: getGenreId(href)}

		s.Find(".top-level-subgenres > li > a").Each(func(i int, s *goquery.Selection) {
			href, exists := genreSelection.Attr("href")
			if !exists {
				return
			}
			genre.SubGenre = append(genre.SubGenre, &Genre{Text: s.Text(), URL: href, Id: getGenreId(href)})
		})
		genres = append(genres, genre)

	})
	return genres, nil
}

func getAuthorization() (string, error) {
	if authorization != "" {
		return authorization, nil
	}
	authRegEx := regexp.MustCompile("22privateKeyPath%22%3A%22ssl%2Fwebplayer.p8%22%2C%22token%22%3A%22(?P<Bearer>.+)%22%7D%2C%22routerScroll")
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
