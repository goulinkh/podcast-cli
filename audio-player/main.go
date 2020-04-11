package audioplayer

import (
	"errors"
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
)

type AudioPlayer struct {
	Streamer beep.StreamSeekCloser
	Format   beep.Format
}

var (
	MainCtrl *beep.Ctrl
	Volume   *effects.Volume
	Streamer beep.StreamSeekCloser
	Format   beep.Format
)

func init() {
	Volume = &effects.Volume{Base: 2}
}

func downloadContent(URL string, filename string, directory string) (string, error) {
	response, err := http.Get(URL)
	if err != nil {
		return "", err
	}
	audio, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	filename = url.PathEscape(path.Clean(strings.ReplaceAll(filename, ":", "")))
	os.MkdirAll(directory, 0755)
	tmpFile, err := ioutil.TempFile(directory, filename)
	if err != nil {
		return "", err
	}
	// download content if not is .cache
	err = ioutil.WriteFile(tmpFile.Name(), audio, 0755)
	if err != nil {
		return "", err
	}
	return tmpFile.Name(), nil
}

// PlaySound play the given audio url, supported Formats: mp3, wav
func PlaySound(filename, directory, URL string) (int, error) {
	if Streamer != nil {
		Streamer.Close()
	}
	audioFile, err := downloadContent(URL, filename, directory)
	if err != nil {
		return 0, err
	}
	file, err := os.Open(audioFile)
	Streamer, Format, err = mp3.Decode(file)
	if err != nil {
		Streamer, Format, err = wav.Decode(file)
	}
	if err != nil {
		return 0, errors.New("Unsupported audio format")
	}
	sr := Format.SampleRate * 2
	speaker.Init(sr, sr.N(time.Second/10))

	Volume.Streamer = beep.Resample(4, Format.SampleRate, sr, Streamer)
	MainCtrl = &beep.Ctrl{Streamer: Volume}
	speaker.Play(MainCtrl)

	return int(float32(Streamer.Len()) / float32(Format.SampleRate)), nil
}

func PauseSong(state bool) {
	speaker.Lock()
	MainCtrl.Paused = state
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
