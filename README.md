# google-squawk

You've heard of Google Chat, Google Voice, Google Talk? Well, this is Google Squawk! :-D
It's a commandline application that connects to the Google Cloud TTS API and generates audio from text.

Google cloud account credentials are required. You can specify the filename in an environment variable:

```
export GOOGLE_APPLICATION_CREDENTIALS=/path/to/serviceaccount/credentials_file.json
```
By default, input is supplied via standard input, but can also be specified from a file with the `-i` option.
Both plain text and SSML input are supported. Specify `-ssml` to tell the program to expect SSML input from your text source.

Output is to a file, but you can set the filename to `-` to send to stdout. This can be useful if you want to convert the file on the fly to a format that gsquawk doesn't support, with **sox** or **ffmpeg** or the like.

Instead of performing TTS, you can use the `-listvoices` option to have gsquawk send you back a list of all available voices. If you run with `-listvoices` and no other options it limits to `en-US` voices, which is the default language selection. You may use the `-l` option along with `-listvoices` to set a specific language code, or use `-l ALL` to give you a list of every voice available in all languages.

Run ```gsquawk -h``` to see the help:

```
Usage of ./gsquawk:
  -db float
    	Volume gain in dB. [-96 to 16]
    	 (default 0)
  -f string
    	Audio format selection. PCM is uncompressed best quality. Opus is
    	excellent quality. MP3 is 32kb bitrate. [pcm,opus,mp3,ulaw,alaw]
    	 (default "mp3")
  -g string
    	Gender selection. [m,f,n] 'n' means neutral/don't care.
    	 (default "m")
  -i string
    	Input file path. Defaults to stdin.
    	 (default "-")
  -l string
    	Language selection. 'en-US', 'en-GB', 'en-AU', 'en-IN',
    	'el-GR', 'ru-RU', etc.
    	 (default "en-US")
  -listvoices
    	List available voices, rather than generate TTS. Use in
    	combination with '-l ALL' to show voices from all languages.
  -o string
    	Output file path. Use '-' for stdout.
    	 (default "./tts.mp3")
  -p float
    	Pitch. E.g. '0.0' is normal. '20.0' is highest,
    	'-20.0' is lowest.
    	 (default 0)
  -r int
    	Samplerate in Hz. [8000,11025,16000,22050,24000,32000,44100,48000]
    	 (default 24000)
  -s float
    	Speed. E.g. '1.0' is normal. '2.0' is double
    	speed, '0.25' is quarter speed, etc.
    	 (default 1)
  -ssml
    	Input is SSML format, rather than plain text.
  -v string
    	Voice. If specified, this overrides language & gender.
    	 (default "unspecified")
```


