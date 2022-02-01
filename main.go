//Copyright 2022 Ruel Tmeizeh - All Rights Reserved
package main

import (
	"bufio"
	"context"
	"encoding/binary"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	texttospeechpb "google.golang.org/genproto/googleapis/cloud/texttospeech/v1"
)

type CommandlineOptions struct {
	Ssml       *bool    `json:"ssml,omitempty"`
	Output     *string  `json:"output,omitempty"`
	Language   *string  `json:"lang,omitempty"`
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
		Ssml:       flag.Bool("ssml", false, "Input is SSML format, rather than plain text."),
		Output:     flag.String("o", "./tts.mp3", "Output file path. Use '-' for stdout."),
		Language:   flag.String("l", "en-US", "Language selection. 'en-US', 'en-GB', 'en-AU', 'en-IN', 'el-GR', 'ru-RU', etc."),
		Gender:     flag.String("g", "m", "Gender selection. [m,f,n]"),
		Format:     flag.String("f", "mp3", "Format selection. [mp3,opus,pcm,ulaw,alaw]"),
		Voice:      flag.String("v", "unspecified", "Voice. If specified, this overrides language & gender."),
		Speed:      flag.Float64("s", 1.0, "Speed. E.g. '1.0' is normal. '2.0' is double speed, '0.25' is quarter speed, etc."),
		Pitch:      flag.Float64("p", 1.0, "Pitch. E.g. '0.0' is normal. '20.0' is highest, '-20.0' is lowest."),
		SampleRate: flag.Int("r", 32000, "Samplerate. [8000,11025,16000,22050,24000,32000,44100,48000]"),
		VolumeGain: flag.Float64("db", 0.0, "Volume gain in dB."),
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

	//Instantiates a Google Cloud client
	ctx := context.Background()
	client, err := texttospeech.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	//take input from stdin
	stdinReader := bufio.NewReader(os.Stdin)
	input, _ := stdinReader.ReadString('\n')

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
