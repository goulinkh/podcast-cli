package itunesapi

import (
	"net/http"
	"regexp"

	"github.com/PuerkitoBio/goquery"
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
