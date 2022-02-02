# google-squawk

You've heard of Google Chat, Google Voice, Google Talk? Well, this is Google Squawk! :-D
It's a commandline application that connects to the Google Cloud TTS API and generates audio from text.

Google cloud account credentials are required. You can specify the filename in an environment variable:

```
export GOOGLE_APPLICATION_CREDENTIALS=/path/to/serviceaccount/credentials_file.json
```

Run ```gsquawk -h``` to see the help.

```
Usage of ./gsquawk:
  --db float
        Volume gain in dB. [-96 to 16]
  -f string
        Audio format selection. MP3 is 32k [mp3,opus,pcm,ulaw,alaw] (default "mp3")
  -g string
        Gender selection. [m,f,n] 'n' means neutral/don't care. (default "m")
  -i string
        Input file path. Defaults to stdin. (default "-")
  -l string
        Language selection. 'en-US', 'en-GB', 'en-AU', 'en-IN',
        'el-GR', 'ru-RU', etc. (default "en-US")
  -o string
        Output file path. Use '-' for stdout. (default "./tts.mp3")
  -p float
        Pitch. E.g. '0.0' is normal. '20.0' is highest,
        '-20.0' is lowest. (default 1)
  -r int
        Samplerate in Hz. [8000,11025,16000,22050,24000,32000,44100,48000] (default 24000)
  -s float
        Speed. E.g. '1.0' is normal. '2.0' is double
        speed, '0.25' is quarter speed, etc. (default 1)
  -ssml
        Input is SSML format, rather than plain text.
  -v string
        Voice. If specified, this overrides language & gender. (default "unspecified")
```

