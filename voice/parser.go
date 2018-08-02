package voice

import (
	"encoding/binary"
	"bytes"
	"log"
	// "os"
	"unsafe"

	"github.com/xlab/closer"
	"github.com/xlab/pocketsphinx-go/sphinx"
	"github.com/xlab/portaudio-go/portaudio"
	"github.com/zaf/resample"
)

const (
	sampleRate        = 44100
	samplesPerChannel = sampleRate
	channels          = 1
	sampleFormat      = portaudio.PaInt16
)

var (
	hmm     = "/usr/local/share/pocketsphinx/model/en-us/en-us"              // Sets directory containing acoustic model files
	dict    = "/usr/local/share/pocketsphinx/model/en-us/cmudict-en-us.dict" // Sets main pronunciation dictionary (lexicon) input file
	lm      = "/usr/local/share/pocketsphinx/model/en-us/en-us.lm.bin"       // Sets word trigram language model input file
	logfile = "pi-car-go-voice-parser.log"                                   // Log file to write log to
	stdout  = true                                                           // Disables log file and writes everything to stdout
	outraw  = ""                                                             // Specify output dir for RAW recorded sound files (s16le). Directory must exist
)

func StartListening() {
	go parse()
}

func parse() {
	log.SetFlags(0)

	defer closer.Close()
	closer.Bind(func() {
		log.Println("Bye!")
	})
	if err := portaudio.Initialize(); paError(err) {
		log.Fatalln("PortAudio init error:", paErrorText(err))
	}
	closer.Bind(func() {
		if err := portaudio.Terminate(); paError(err) {
			log.Println("PortAudio term error:", paErrorText(err))
		}
	})

	// Init CMUSphinx
	cfg := sphinx.NewConfig(
		sphinx.HMMDirOption(hmm),
		sphinx.DictFileOption(dict),
		sphinx.LMFileOption(lm),
		sphinx.SampleRateOption(16000),
	)
	if len(outraw) > 0 {
		sphinx.RawLogDirOption(outraw)(cfg)
	}
	if stdout == false {
		sphinx.LogFileOption(logfile)(cfg)
	}

	log.Println("Loading CMU PhocketSphinx.")
	log.Println("This may take a while depending on the size of your model.")
	dec, err := sphinx.NewDecoder(cfg)
	if err != nil {
		closer.Fatalln(err)
	}
	closer.Bind(func() {
		dec.Destroy()
	})
	l := &paListener{
		dec: dec,
	}

	var stream *portaudio.Stream
	if err := portaudio.OpenDefaultStream(&stream, channels, 0, sampleFormat,
		sampleRate, samplesPerChannel, l.paCallback, nil); paError(err) {
		log.Fatalln("PortAudio error:", paErrorText(err))
	}
	closer.Bind(func() {
		if err := portaudio.CloseStream(stream); paError(err) {
			log.Println("[WARN] PortAudio error:", paErrorText(err))
		}
	})

	if err := portaudio.StartStream(stream); paError(err) {
		log.Fatalln("PortAudio error:", paErrorText(err))
	}
	closer.Bind(func() {
		if err := portaudio.StopStream(stream); paError(err) {
			log.Fatalln("[WARN] PortAudio error:", paErrorText(err))
		}
	})

	if !dec.StartUtt() {
		closer.Fatalln("[ERR] Sphinx failed to start utterance")
	}
	// log.Println(banner)
	log.Println("Ready..")
	closer.Hold()
}

type paListener struct {
	inSpeech   bool
	uttStarted bool
	dec        *sphinx.Decoder
}

// paCallback: for simplicity reasons we process raw audio with sphinx in the stream callback,
// never do that for any serious applications, use a buffered channel instead.
func (l *paListener) paCallback(input unsafe.Pointer, _ unsafe.Pointer, sampleCount uint,
	_ *portaudio.StreamCallbackTimeInfo, _ portaudio.StreamCallbackFlags, _ unsafe.Pointer) int32 {

	const (
		statusContinue = int32(portaudio.PaContinue)
		statusAbort    = int32(portaudio.PaAbort)
	)

	sampleRatio := 16000.0 / sampleRate
	samplesIn := int(sampleCount) * channels
	sampleInBytes := samplesIn * 2
	samplesOut := int(float64(samplesIn) * sampleRatio)
	sampleOutBytes := samplesOut * 2

	// log.Printf("processing %d samples (%d bytes) at %d to %d samples (%d bytes) at %d...", samplesIn, sampleInBytes, sampleRate, samplesOut, sampleOutBytes, 16000)

	buffer := bytes.NewBuffer(make([]byte, sampleOutBytes))
	resampler, _ := resample.New(buffer, float64(sampleRate), float64(16000), channels, resample.I16, resample.HighQ)
	defer resampler.Close()
	raw := (*(*[1 << 24]byte)(input))[:sampleInBytes]

	resampler.Write(raw)
	// wrote, err := resampler.Write(raw)
	// log.Printf("%d bytes written of %d; error: %v", wrote, len(raw), err)
	// log.Printf("%d samples in (%d bytes); % X -> % X", samplesIn, len(buffer.Bytes()), raw, buffer.Bytes())
	in := make([]int16, len(buffer.Bytes()) / 2)
	binary.Read(buffer, binary.LittleEndian, &in)
	// err = binary.Read(buffer, binary.LittleEndian, &in)
	// log.Printf("read error: %v", err)
	// log.Printf("%d samples out; %X", len(in), in)

	// ProcessRaw with disabled search because callback needs to be realtime
	_, ok := l.dec.ProcessRaw(in, true, false)
	// log.Printf("processed: %d frames, ok: %v", frames, ok)
	if !ok {
		return statusAbort
	}
	if l.dec.IsInSpeech() {
		l.inSpeech = true
		if !l.uttStarted {
			l.uttStarted = true
			log.Println("Listening..")
		}
	} else if l.uttStarted {
		// speech -> silence transition, time to start new utterance
		l.dec.EndUtt()
		l.uttStarted = false
		l.report() // report results
		if !l.dec.StartUtt() {
			closer.Fatalln("[ERR] Sphinx failed to start utterance")
		}
	}
	return statusContinue
}

func (l *paListener) report() {
	hyp, _ := l.dec.Hypothesis()
	if len(hyp) > 0 {
		log.Printf("    > hypothesis: %s", hyp)
	} else {
		log.Println("ah, nothing")
	}
}

func paError(err portaudio.Error) bool {
	return portaudio.ErrorCode(err) != portaudio.PaNoError
}

func paErrorText(err portaudio.Error) string {
	return portaudio.GetErrorText(err)
}
