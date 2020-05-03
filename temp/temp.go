package temp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/PuerkitoBio/goquery"
	"github.com/tidwall/gjson"
)

const applePodcastsMainPage = "https://podcasts.apple.com/genre/podcasts/id26"

var (
	authRegEx      = regexp.MustCompile("22privateKeyPath%22%3A%22ssl%2Fwebplayer.p8%22%2C%22token%22%3A%22(?P<Bearer>.+)%22%7D%2C%22routerScroll")
	podcastIdRegEx = regexp.MustCompile(`\d+$`)
)

type Genre struct {
	Text     string   `json:"text"`
	URL      string   `json:"url"`
	SubGenre []*Genre `json:"sub-genres"`
}

type Podcast struct {
	Title string `json:"title"`
	URL   string `json:"url`
	Id    string `json:"id"`
}
type Episode struct {
	Id                     string `json:"id"`
	Artwork                string `json:"artwork"`
	Title                  string `json:"title"`
	AudioURL               string `json:"audiourl"`
	releaseDate            string `json:"releasedate"`
	DurationInMilliseconds string `json:"duratioInMilliseconds"`
	Description            string `json:"description"`
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
			Id:    podcastIdRegEx.FindString(href),
		})
	})
	return podcasts, nil
}
func (p *Podcast) getEpisodes() ([]*Episode, error) {
	response, err := http.Get(p.URL)
	if err != nil {
		return nil, err
	}
	podcastPageBytes, err := ioutil.ReadAll(response.Body)

	match := authRegEx.FindStringSubmatch(string(podcastPageBytes))
	if len(match) != 2 {
		return nil, err

	}
	authorization := "Bearer " + match[1]
	fmt.Println(authorization)
	req, err := http.NewRequest("GET", fmt.Sprintf("https://amp-api.podcasts.apple.com/v1/catalog/us/podcasts/%s/episodes?offset=0&limit=300", p.Id), nil)
	if err != nil {
		return nil, err

	}
	req.Header.Add("Authorization", authorization)
	response, err = http.DefaultClient.Do(req)
	if err != nil {
		return nil, err

	}
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err

	}
	episodesJson := gjson.Get(string(data), "data").Array()
	episodes := make([]*Episode, len(episodesJson))
	for i, episode := range episodesJson {
		episodes[i] = &Episode{
			Id:                     episode.Get(`id`).String(),
			Artwork:                episode.Get(`attributes.artwork.url`).String(),
			Title:                  episode.Get(`attributes.name`).String(),
			AudioURL:               episode.Get(`attributes.assetUrl`).String(),
			releaseDate:            episode.Get(`attributes.releaseDateTime`).String(),
			DurationInMilliseconds: episode.Get(`attributes.durationInMilliseconds`).String(),
			Description:            episode.Get(`attributes.description.standard`).String(),
		}
	}
	return episodes, nil
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

		genre := &Genre{Text: genreSelection.Text(), URL: href, SubGenre: make([]*Genre, 0)}

		s.Find(".top-level-subgenres > li > a").Each(func(i int, s *goquery.Selection) {
			href, exists := genreSelection.Attr("href")
			if !exists {
				return
			}
			genre.SubGenre = append(genre.SubGenre, &Genre{Text: s.Text(), URL: href})
		})
		genres = append(genres, genre)

	})
	return genres, nil
}

func Main() {
	/* genres, _ := getGenres()
	podcasts, _ := genres[0].SubGenre[0].getPodcasts()
	res, _ := json.Marshal(podcasts)
	fmt.Println(string(res))
	fmt.Println((genres[0].SubGenre[0])) */
	genres, _ := GetGenres()
	podcasts, _ := genres[0].SubGenre[0].getPodcasts()
	episodes, _ := podcasts[0].getEpisodes()
	fmt.Println(len(episodes))
}
