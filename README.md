# pollyreveal
## Summary
A pre-processor for RevealJS that would process speaker notes in `<aside>` tags to audio files using AWS Polly capabilities.  Enables a reveal JS presentation to be hosted with audio content to the viewer.

Usage: <code>pollyreveal [-v voiceId] inputFile outputFile</code>
<br>`voiceId` should be a value supported by [AWS Polly](http://docs.aws.amazon.com/polly/latest/dg/API_Voice.html)<br>

## Installation
Download the pollyreveal binary from the builds directory.  It has been built for OSX.  For other OSes, you need to download the source and compile it using the GO compiler for your OS.  Hope to fix this shortly with mutiple OS builds.

## AWS Credentials
You will need an AWS Account with access to AWS Polly<br>
The simplest way to check is to use the AWS CLI tool to check if you have access and are configured correctly.<br>
If you are using Linux of OSX, then ~/.aws directory should have your credentials setup.<br>
Ensure that you have <code>export AWS_REGION=us-east-1</code> of any region that has AWS Polly available.<br>
Not setting this region will result in Missing Region error and many things will go wrong.

## Details
RevealJS provides `<aside>` tags for a developer to plug in speaker notes into the presentation.  `pollyreveal` can pre-process these tags and add `<audio>` tags to the presentation, with mp3 files generated locally that has natural voices of the text in the speaker notes.  The voice can be modified to suit the style or language based on the [AWS Polly](http://docs.aws.amazon.com/polly/latest/dg/API_Voice.html) documentation.

```html
<section>
  <h1>20th Century Physicists</h1>
  <p></p>
  <aside class="notes">
    Hi!, We will be introducing you to some leading 20th century physicists today during this presentation.
    To move forward when you are ready, please press the space bar or use the arrow links at the bottom to click right.
  </aside>
</section>
```
would be converted to have an additional `<audio>` tag like below
```html
<section>
  <h1>20th Century Physicists</h1>
  <p></p>
  <audio data-autoplay="" src="demo1.html.001.mp3"></audio><aside class="notes">
    Hi!, We will be introducing you to some leading 20th century physicists today during this presentation.
    To move forward when you are ready, please press the space bar or use the arrow links at the bottom to click right.
  </aside>
</section>
```
The `src` attribute of the `audio` tag is generated based on the output filename that is passed as parameter.  If the filename is `demo1.html`, then the mp3 files are generated as `demo1.html.001.mp3`.

## Known Bugs and Issues
* Unable to customize the sample-rate, file format of the audio files
* Unable to use many other features like Lexicon of AWS Polly
* Does not process SSML content in speaker notes
* Speaker notes have to be less than 1500 characters or less than 5 mins of talk time.
