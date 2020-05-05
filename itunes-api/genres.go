package itunesapi

import (
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
