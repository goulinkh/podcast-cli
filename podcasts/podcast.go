package podcasts

import (
	"encoding/json"
	"fmt"
	"time"
)

type Podcast struct {
	URL         string `json:"mygpo_link"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	Description string `json:"description"`
	LogoURL     string `json:"logo_url"`
	Website     string `json:"website"`
}

func (p *Podcast) String() string {
	res, err := json.Marshal(p)
	if err != nil {
		return "failed"
	}
	return string(res)
}

type Episode struct {
	URL         string    `json:"url"`
	Title       string    `json:"title"`
	Podcast     *Podcast  `json:"podcast"`
	Description string    `json:"description"`
	audioURL    string    `json:"audio-url"`
	ReleaseDate time.Time `json:"release-date"`
}

func (e *Episode) String() string {
	return fmt.Sprintf("title : %s, url: %s\n", e.Title, e.URL)
}
