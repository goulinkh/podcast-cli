package audioplayer

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
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
	mainCtrl *beep.Ctrl
	volume   *effects.Volume
	streamer beep.StreamSeekCloser
	format   beep.Format
)

func init() {
	volume = &effects.Volume{Base: 2}
}
func downloadContent(URL string, filename string, directory string) (string, error) {
	filename = path.Clean(strings.ReplaceAll(filename, ":", ""))
	filePath := directory + filename
	_, err := os.Open(filePath)
	// download content if not is .cache
	if err == nil {
		return filePath, nil
	}
	response, err := http.Get(URL)
	if err != nil {
		return "", err
	}
	audio, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	fmt.Println("downloaded file")
	os.MkdirAll(".cache", os.ModeAppend)
	err = ioutil.WriteFile(filePath, audio, os.ModeAppend)
	if err != nil {
		return "", err
	}
	return filePath, nil
}

// PlaySound play the given audio url, supported formats: mp3, wav
func PlaySound(filename, directory, URL string) (int, error) {
	if streamer != nil {
		streamer.Close()
	}
	audioFile, err := downloadContent(URL, filename, directory)
	if err != nil {
		return 0, err
	}
	fmt.Println(audioFile)
	file, err := os.Open(audioFile)
	streamer, format, err = mp3.Decode(file)
	if err != nil {
		streamer, format, err = wav.Decode(file)
	}
	if err != nil {
		return 0, errors.New("Unsupported audio format")
	}
	sr := format.SampleRate * 2
	speaker.Init(sr, sr.N(time.Second/10))

	volume.Streamer = beep.Resample(4, format.SampleRate, sr, streamer)
	mainCtrl = &beep.Ctrl{Streamer: volume}
	speaker.Play(mainCtrl)

	return int(float32(streamer.Len()) / float32(format.SampleRate)), nil
}

func PauseSong(state bool) {
	speaker.Lock()
	mainCtrl.Paused = state
	speaker.Unlock()
}

func Seek(pos int) error {
	speaker.Lock()
	err := streamer.Seek(pos * int(format.SampleRate))
	speaker.Unlock()
	return err
}

func SetVolume(percent int) {
	if percent == 0 {
		volume.Silent = true
	} else {
		volume.Silent = false
		volume.Volume = -float64(100-percent) / 100.0 * 5
	}
}
