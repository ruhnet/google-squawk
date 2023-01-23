//Copyright 2022 Ruel Tmeizeh - All Rights Reserved
package main

import (
	"bufio"
	"context"
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	texttospeechpb "google.golang.org/genproto/googleapis/cloud/texttospeech/v1"
)

type CommandlineOptions struct {
	ListVoices *bool    `json:"listvoices,omitempty"`
	Ssml       *bool    `json:"ssml,omitempty"`
	Output     *string  `json:"output,omitempty"`
	Input      *string  `json:"input,omitempty"`
	Language   *string  `json:"language,omitempty"`
	Gender     *string  `json:"gender,omitempty"`
	Voice      *string  `json:"voice,omitempty"`
	Format     *string  `json:"format,omitempty"`
	Speed      *float64 `json:"speed,omitempty"`
	Pitch      *float64 `json:"pitch,omitempty"`
	SampleRate *int     `json:"samplerate,omitempty"`
	VolumeGain *float64 `json:"volume,omitempty"`
}

func main() {
	//check commandline args:
	opts := &CommandlineOptions{
		ListVoices: flag.Bool("listvoices", false, "List available voices, rather than generate TTS. Use in\ncombination with '-l ALL' to show voices from all languages."),
		Ssml:       flag.Bool("ssml", false, "Input is SSML format, rather than plain text."),
		Input:      flag.String("i", "-", "Input file path. Defaults to stdin.\n"),
		Output:     flag.String("o", "./tts.mp3", "Output file path. Use '-' for stdout.\n"),
		Language:   flag.String("l", "en-US", "Language selection. 'en-US', 'en-GB', 'en-AU', 'en-IN',\n'el-GR', 'ru-RU', etc.\n"),
		Gender:     flag.String("g", "m", "Gender selection. [m,f,n] 'n' means neutral/don't care.\n"),
		Format:     flag.String("f", "mp3", "Audio format selection. PCM is uncompressed best quality. Opus is\nexcellent quality. MP3 is 32kb bitrate. [pcm,opus,mp3,ulaw,alaw]\n"),
		Voice:      flag.String("v", "unspecified", "Voice. If specified, this overrides language & gender.\n"),
		Speed:      flag.Float64("s", 1.0, "Speed. E.g. '1.0' is normal. '2.0' is double\nspeed, '0.25' is quarter speed, etc.\n"),
		Pitch:      flag.Float64("p", 0.0, "Pitch. E.g. '0.0' is normal. '20.0' is highest,\n'-20.0' is lowest.\n (default 0)"),
		SampleRate: flag.Int("r", 24000, "Samplerate in Hz. [8000,11025,16000,22050,24000,32000,44100,48000]\n"),
		VolumeGain: flag.Float64("-db", 0.0, "Volume gain in dB. [-96 to 16]\n (default 0)"),
	}
	flag.Parse()

	var audioFormat texttospeechpb.AudioEncoding
	var fileExtension string
	switch *opts.Format {
	case "mp3":
		audioFormat = texttospeechpb.AudioEncoding_MP3
		fileExtension = "mp3"
	case "opus":
		audioFormat = texttospeechpb.AudioEncoding_OGG_OPUS
		fileExtension = "ogg"
	case "ogg":
		audioFormat = texttospeechpb.AudioEncoding_OGG_OPUS
		fileExtension = "ogg"
	case "pcm":
		audioFormat = texttospeechpb.AudioEncoding_LINEAR16
		fileExtension = "pcm"
	case "ulaw":
		audioFormat = texttospeechpb.AudioEncoding_MULAW
		fileExtension = "ulaw"
	case "alaw":
		audioFormat = texttospeechpb.AudioEncoding_ALAW
		fileExtension = "alaw"
	default:
		audioFormat = texttospeechpb.AudioEncoding_MP3
		fileExtension = "mp3"
	}

	filename := "tts." + fileExtension
	if *opts.Output != "./tts.mp3" {
		filename = *opts.Output
	}

	///////////////////////////////////////
	//Instantiates a Google Cloud client
	ctx := context.Background()
	client, err := texttospeech.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	if *opts.ListVoices {
		fmt.Println("Available Voices:")
		bufStdout := bufio.NewWriter(os.Stdout)
		listVoices(bufStdout, ctx, client, *opts.Language)
		bufStdout.Flush()
		os.Exit(0)
	}

	var inputFile *os.File
	if *opts.Input == "-" {
		//read input from stdin
		inputFile = os.Stdin
	} else {
		//read input from file
		var err error
		inputFile, err = os.Open(*opts.Input)
		if err != nil {
			log.Fatal(err)
		}
		defer inputFile.Close()
	}

	var input string

	scanner := bufio.NewScanner(inputFile)
	for scanner.Scan() {
		//fmt.Println(scanner.Text())
		input = input + scanner.Text()
	}

	//Start building TTS request things
	synthInput := &texttospeechpb.SynthesisInput{}
	synthInput.InputSource = &texttospeechpb.SynthesisInput_Text{Text: input}
	if *opts.Ssml {
		synthInput.InputSource = &texttospeechpb.SynthesisInput_Ssml{Ssml: input}
	}

	//Voice Gender
	var gender texttospeechpb.SsmlVoiceGender
	switch *opts.Gender {
	case "m":
		gender = texttospeechpb.SsmlVoiceGender_MALE
	case "f":
		gender = texttospeechpb.SsmlVoiceGender_FEMALE
	default:
		gender = texttospeechpb.SsmlVoiceGender_NEUTRAL
	}

	voice := &texttospeechpb.VoiceSelectionParams{
		LanguageCode: *opts.Language,
		SsmlGender:   gender,
		//Name:         *opts.Voice, //Name overrides LanguageCode and SsmlGender
		//Name: "en-US-Wavenet-B",
	}
	if *opts.Voice != "unspecified" {
		voice.Name = *opts.Voice
	}

	//the request parameters
	req := texttospeechpb.SynthesizeSpeechRequest{
		Input: synthInput,
		Voice: voice,
		AudioConfig: &texttospeechpb.AudioConfig{
			AudioEncoding:   audioFormat,
			SpeakingRate:    *opts.Speed,
			SampleRateHertz: int32(*opts.SampleRate),
			Pitch:           *opts.Pitch,
			VolumeGainDb:    *opts.VolumeGain,
		},
	}

	resp, err := client.SynthesizeSpeech(ctx, &req)
	if err != nil {
		log.Fatal(err)
	}

	if *opts.Output == "-" { //write to stdout
		//binary.Write(os.Stdout, binary.LittleEndian, resp.AudioContent)
		bufStdout := bufio.NewWriter(os.Stdout) //add a buffer
		defer bufStdout.Flush()
		binary.Write(bufStdout, binary.LittleEndian, resp.AudioContent)
	} else { //write to file
		err = ioutil.WriteFile(filename, resp.AudioContent, 0644)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Audio content written to file: %v\n", filename)
	}

}

func listVoices(w io.Writer, ctx context.Context, client *texttospeech.Client, lang string) error {
	resp, err := client.ListVoices(ctx, &texttospeechpb.ListVoicesRequest{})
	if err != nil {
		return err
	}

	for _, voice := range resp.Voices {
		for _, languageCode := range voice.LanguageCodes {
			if lang == languageCode || lang == "ALL" {
				fmt.Fprintln(w, "___________________________________")
				fmt.Fprintf(w, "Name: %v\n", voice.Name)
				fmt.Fprintf(w, "  Language: %v\n", languageCode)
				fmt.Fprintf(w, "  Gender: %v\n", voice.SsmlGender.String())
				fmt.Fprintf(w, "  Native Sample Rate (in Hz): %v\n", voice.NaturalSampleRateHertz)
			}
		}
	}
	fmt.Fprintln(w, "------------------------------------")

	return nil
}
