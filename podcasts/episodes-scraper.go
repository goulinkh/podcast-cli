package podcasts

import (
	"fmt"
	"math"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const episodesPerPage = 20

func (podcast *Podcast) GetEpisodes(index int, count int) []*Episode {
	startPage := index/episodesPerPage + 1
	pagesToFetch := int(math.Ceil(float64(count+index%episodesPerPage) / episodesPerPage))

	var wg sync.WaitGroup
	episodes := make([]*Episode, 0)
	for i := 0; i < pagesToFetch; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			episodesPerPage, err := podcast.scrapEpisodesByPodcast(fmt.Sprintf("%s/-episodes?page=%d", podcast.URL, startPage+index))
			if err != nil {
				fmt.Println(err)
			} else {
				episodes = append(episodes, episodesPerPage...)
			}
		}(i)
	}
	wg.Wait()
	// TODO: sort by release date
	if len(episodes) < count {
		return episodes
	}

	return episodes[:count]
}

func (podcast *Podcast) scrapEpisodesByPodcast(URL string) ([]*Episode, error) {
	response, err := http.Get(URL)
	if err != nil || response.StatusCode != 200 {
		return nil, fmt.Errorf("Invalid podcast URL")
	}
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to get episodes list")
	}
	episodes := make([]*Episode, 0)
	doc.Find(".episode").Each(func(i int, s *goquery.Selection) {
		title := s.Find(".header > .title > a").Text()
		url, success := s.Find(".header > .title > a").Attr("href")
		if !success {
			return
		}
		description := strings.Trim(s.Find(".description").Text(), " \n")
		releaseDateString := strings.Trim(s.Find(".header > .released").Text(), " \n")
		releaseDate := parseReleaseDate(releaseDateString)
		episodes = append(episodes, &Episode{
			URL: "https://gpodder.net" + url, Title: title, Description: description, ReleaseDate: releaseDate, Podcast: podcast,
		})
	})

	return episodes, nil
}

func (e *Episode) AudioURL() (string, error) {
	if e.audioURL != "" {
		fmt.Println("hey there")
		return e.audioURL, nil
	}
	response, err := http.Get(e.URL)
	if err != nil || response.StatusCode != 200 {
		fmt.Println(err)
		return "", fmt.Errorf("Invalid episode URL")
	}

	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return "", err
	}

	link, success := doc.Find(".description > a:nth-child(2)").Attr("href")
	if !success {
		return "", fmt.Errorf("Failed to find audi link")
	}
	e.audioURL = link
	return link, nil
}

func parseReleaseDate(date string) (releaseDate time.Time) {
	formats := []string{
		"Jan. 2, 2006", "January 2, 2006",
	}
	date = strings.Replace(date, "Sept.", "September", -1)
	for _, format := range formats {
		releaseDate, err := time.Parse(format, date) //Sept. 30, 2019
		if err == nil {
			return releaseDate
		}
	}
	return releaseDate
}
