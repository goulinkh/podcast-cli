package audioplayer

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
	"github.com/goulinkh/podcast-cli/config"
	itunesapi "github.com/goulinkh/podcast-cli/itunes-api"
)

type AudioPlayer struct {
	Streamer beep.StreamSeekCloser
	Format   beep.Format
}

var (
	MainCtrl  *beep.Ctrl
	Volume    *effects.Volume
	resampler *beep.Resampler

	Streamer beep.StreamSeekCloser
	Format   beep.Format
)

func init() {
	Volume = &effects.Volume{Base: 2}
}

func fetchContent(URL string, filepath string, directory string) error {

	_, err := ioutil.ReadFile(filepath)
	if err == nil {
		return nil
	}
	response, err := http.Get(URL)
	if err != nil {
		return err
	}
	audio, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	os.MkdirAll(directory, 0755)
	// download content if not in .cache
	err = ioutil.WriteFile(filepath, audio, 0755)
	if err != nil {
		return err
	}
	return nil
}

// PlaySound play the given audio url, supported Formats: mp3, wav
func PlaySound(e *itunesapi.Episode) error {
	if Streamer != nil {
		speaker.Lock()
		Streamer.Close()
		speaker.Unlock()
	}
	URL := e.AudioURL
	filename := fmt.Sprintf("%s.mp3", e.Id)
	directory := config.CachePath
	filename = url.PathEscape(path.Clean(strings.ReplaceAll(filename, ":", "")))
	filepath := path.Join(directory, filename)
	file, err := os.Open(filepath)
	if err != nil {
		err = fetchContent(URL, filepath, directory)
		if err != nil {
			return err
		}
	}
	file, err = os.Open(filepath)
	if err != nil {
		return err
	}
	Streamer, Format, err = mp3.Decode(file)
	if err != nil {
		Streamer, Format, err = wav.Decode(file)
	}
	if err != nil {
		return errors.New("Unsupported audio format")
	}
	sr := Format.SampleRate * 2
	speaker.Init(sr, sr.N(time.Millisecond*500))

	streamer := beep.Resample(4, Format.SampleRate, sr, Streamer)
	MainCtrl = &beep.Ctrl{Streamer: streamer}
	resampler = beep.ResampleRatio(4, 1, MainCtrl)
	Volume = &effects.Volume{Streamer: resampler, Base: 2}
	speaker.Play(Volume)
	e.DurationInMilliseconds = int(float32(Streamer.Len())/float32(Format.SampleRate)) * 1000
	return nil
}

func PauseSong(state bool) {
	speaker.Lock()
	MainCtrl.Paused = state
	speaker.Unlock()
}

func IncreaseSpeed() {
	speed := resampler.Ratio() * 1.100000e+000
	if speed >= 1.800000e+000 {
		return
	}
	speaker.Lock()
	resampler.SetRatio(speed)
	speaker.Unlock()
}

func DecreaseSpeed() {
	speed := resampler.Ratio() * 0.900000e+000
	if speed <= 1.000000e+000 {
		return
	}
	speaker.Lock()
	resampler.SetRatio(speed)
	speaker.Unlock()
}

func Seek(pos int) error {
	if MainCtrl != nil {
		speaker.Lock()
		err := Streamer.Seek(Format.SampleRate.N(time.Second) * pos)
		speaker.Unlock()
		return err
	}
	return nil
}

func SetVolume(percent int) {
	if percent > 100 {
		return
	}

	if percent == 0 {
		Volume.Silent = true
	} else {
		Volume.Silent = false
		Volume.Volume = -float64(100-percent) / 100.0 * 5
	}
}

func Position() int {
	return int(Format.SampleRate.D(Streamer.Position()).Round(time.Second).Seconds())
}
