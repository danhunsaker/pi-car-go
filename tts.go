package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
	"github.com/pqyptixa/tts2media"
	"github.com/xlab/closer"
)

var (
	saveDir string
)

type TTSWriter struct {
	Prefix string
	say func(string) error
}

func NewTTSWriter(prefix, engine string) (*TTSWriter, error) {
	var w TTSWriter

	switch engine {
	case "espeak":
		w = TTSWriter{prefix, sayEspeak}
	case "pico", "picotts":
		w = TTSWriter{prefix, sayPico}
	case "festival":
		w = TTSWriter{prefix, sayFestival}
	default:
		return nil, errors.New("Unsupported TTS engine")
	}

	return &w, nil
}

func (w *TTSWriter) Write(text []byte) (int, error) {
	err := w.say(w.Prefix + " " + string(text))
	return len(text), err
}

func init() {
	var err error
	saveDir, err = ioutil.TempDir(os.TempDir(), "pi-car-go-")
	if err != nil {
		fmt.Printf("Error allocating temporary directory: %v\n", err)
	}
	err = os.Mkdir(saveDir+"/tmp", 0700)
	if err != nil {
		fmt.Printf("Error allocating work subdirectory: %v\n", err)
	}
	tts2media.SetDataDir(saveDir + "/")
	closer.Bind(func() {
		os.RemoveAll(saveDir)
	})
}

func sayEspeak(say string) error {
	media, err := (&tts2media.EspeakSpeech{
		Text:     say,    // text to turn to speech
		Lang:     "en",   // language
		Speed:    "135",  // speed
		Gender:   "f",    // gender
		Altvoice: "0",    // alternative voice
		Quality:  "high", // quality of output mp3/ogg audio
		Pitch:    "50",   // pitch
	}).NewEspeakSpeech()
	if err != nil {
		return err
	}

	return playAudio(media)
}

func sayPico(say string) error {
	media, err := (&tts2media.PicoTTSSpeech{
		Text:    say,     // text to turn to speech
		Lang:    "en-US", // language
		Quality: "high",  // quality of output mp3/ogg audio
	}).NewPicoTTSSpeech()
	if err != nil {
		return err
	}

	return playAudio(media)
}

func sayFestival(say string) error {
	media, err := (&tts2media.FestivalSpeech{
		Text:    say,       // text to turn to speech
		Lang:    "english", // language
		Quality: "high",    // quality of output mp3/ogg audio
	}).NewFestivalSpeech()
	if err != nil {
		return err
	}

	return playAudio(media)
}

func playAudio(media *tts2media.Media) error {
	done := make(chan struct{})
	filename := saveDir + "/" + media.Filename

	defer media.RemoveWAV()
	f, err := os.Open(filename + ".wav")
	if err != nil {
		close(done)
		fmt.Printf("Error while opening WAV file: %v\n", err)
		return err
	}
	s, format, err := wav.Decode(f)
	if err != nil {
		close(done)
		fmt.Printf("Error while decoding WAV file: %v\n", err)
		return err
	}

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	speaker.Play(beep.Seq(s, beep.Callback(func() {
		close(done)
	})))

	<-done
	return nil
}
