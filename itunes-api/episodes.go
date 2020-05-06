package itunesapi

type Episode struct {
	Id                     string `json:"id"`
	Artwork                string `json:"artwork"`
	Title                  string `json:"title"`
	AudioURL               string `json:"audiourl"`
	ReleaseDate            string `json:"releasedate"`
	DurationInMilliseconds int    `json:"duratioInMilliseconds"`
	Description            string `json:"description"`
}
