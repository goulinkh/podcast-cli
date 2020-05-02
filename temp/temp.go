package temp

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

const applePodcastsMainPage = "https://podcasts.apple.com/genre/podcasts/id26"

type Genre struct {
	Text     string      `json:"text"`
	URL      string      `json:"url"`
	SubGenre []*SubGenre `json:"sub-genres"`
}

type SubGenre struct {
	Text string `json:"text"`
	URL  string `json:"url"`
}

type Podcast struct {
	Title string `json:"title"`
	URL   string `json:"url`
}

func (g *Genre) String() string {
	str, err := json.Marshal(g)
	if err != nil {
		return ""
	}
	return string(str)
}
func (g *Genre) getPodcasts() ([]*Podcast, error) {
	resp, err := http.Get(g.URL)
	if err != nil {
		return nil, err
	}
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return nil, err
	}
	podcasts := make([]*Podcast, 0)
	doc.Find("#selectedcontent li a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if !exists {
			return
		}
		podcasts = append(podcasts, &Podcast{
			Title: s.Text(),
			URL:   href,
		})
	})
	return podcasts, nil
}
func getGenres() ([]*Genre, error) {
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

		genre := &Genre{Text: genreSelection.Text(), URL: href, SubGenre: make([]*SubGenre, 0)}

		s.Find(".top-level-subgenres > li > a").Each(func(i int, s *goquery.Selection) {
			href, exists := genreSelection.Attr("href")
			if !exists {
				return
			}
			genre.SubGenre = append(genre.SubGenre, &SubGenre{Text: s.Text(), URL: href})
		})
		genres = append(genres, genre)

	})
	fmt.Println(len(genres))
	return genres, nil
}

func Main() {
	genres, _ := getGenres()
	podcasts, _ := genres[0].getPodcasts()
	res, _ := json.Marshal(podcasts)
	fmt.Println(string(res))
}
