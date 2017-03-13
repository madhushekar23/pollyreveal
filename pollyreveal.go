package main

import (
	"fmt"
  "errors"
  "sync"
  "log"
  "os"
	"bytes"
	"golang.org/x/net/html"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/polly"

)

const POLY = "aside"  // tag used by reveal js for
const DEFAULT_VOICE = "Raveena"

func main() {

	voiceId, inputFileName, outputFileName, err := processFlags()
	if err != nil  {
		log.Fatal(err)
	}
	
	fmt.Println("InputFile: ", inputFileName)
	fmt.Println("OutputFile: ", outputFileName)
	fmt.Println("VoiceId: ", voiceId)

  // Open the Input file
  inputFileHandle, err := os.Open(inputFileName)
  if err != nil {
    log.Fatal(err)
  }
  defer inputFileHandle.Close()

  // Open the output file, create it if needed
  outputFileHandle, err := os.OpenFile(outputFileName, os.O_RDWR|os.O_CREATE, 0755)
  if err != nil {
    log.Fatal(err)
  }
  defer outputFileHandle.Close()

  // Parse the input html doc into memory, this is a heavy memory
  // utilization function
  htmlDoc, err := html.Parse(inputFileHandle)
  if err != nil {
    log.Fatal(err)  // unable to parse the html document
  }

	// initialize the AWS session and service
	awsSession := session.Must(session.NewSession())
	pollySvc := polly.New(awsSession)

  // initialize file name counter
  // audio files will have names of the format
  // outputFileName.polly00x.mp3
  counter := 0
  var wg sync.WaitGroup


  // Define the processing loop
  var processDoc func(*html.Node)
  processDoc = func(node *html.Node)  {
    if node.Type == html.ElementNode && node.Data == POLY  {
      polyNode := node.FirstChild
      if polyNode != nil && len(polyNode.Data) > 0  {

        counter++
        mp3FileName := fmt.Sprintf("%s.%03d.mp3", outputFileName, counter)

        wg.Add(1)
        go func(localNode *html.Node, audioFileName string)  {
          defer wg.Done()
          audioNode, err := generateAudioFile(localNode.Data, audioFileName, voiceId, pollySvc)
          if err != nil {
            log.Print(err)
          }
          node.Parent.InsertBefore(audioNode, node)

        }(polyNode, mp3FileName)
      }
    }
    for childNode := node.FirstChild; childNode != nil; childNode = childNode.NextSibling  {
      processDoc(childNode)
    }
  }

  processDoc(htmlDoc)  // Enter the processing loop
  wg.Wait()  // wait for all concurrent requests to end

	html.Render(outputFileHandle, htmlDoc)  // render the output file

  // and we are done.
}


// Use AWS Polly SDK to generate an MP3 file.  Create a html node
// return with the filename and return back to caller for embedding in
// document tree
func generateAudioFile(speechText, audioFileName, voice string, pollySvc *polly.Polly)(*html.Node, error) {
  if len(audioFileName) == 0 {
    return nil, errors.New("audio file name came in empty\n")
  }

	// Define inputs to Polly for Voice
	pollyInput := &polly.SynthesizeSpeechInput{
										OutputFormat: aws.String("mp3"),
										VoiceId:			aws.String(voice),
										Text:					aws.String(speechText),
	}

	// Call Polly Service to generate Audio file
	pollyResp, err := pollySvc.SynthesizeSpeech(pollyInput)
	if err != nil {
		fmt.Println("AWS Polly Error:", err)
		return nil, err
	}

	// Open the audio file for writing
	audioFileHandle, err := os.OpenFile(audioFileName, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil  {
		fmt.Println("Error Writing file: ", audioFileName, err)
		return nil, err
	}
	defer audioFileHandle.Close()

	// Write the output buffer to file
	audioBuffer := new(bytes.Buffer)
	audioBuffer.ReadFrom(pollyResp.AudioStream)
	_, err = audioFileHandle.Write(audioBuffer.Bytes())
	if err != nil  {
		fmt.Println(err)
		return nil, err
	}

	fmt.Println(pollyResp)   // dump the response to console

	// generate the <audio> node and return back to caller
  node := makeAudioNode(audioFileName)
  return node, nil
}

// This generates a simple <audio data-autoplay src="">
// tag for embedding into a html document
func makeAudioNode(filename string) *html.Node  {
  node := &html.Node{
            Type: 3,
            Data: "audio",
            Attr: []html.Attribute{
												html.Attribute{Key: "data-autoplay"},
												html.Attribute{Key: "src", Val: filename,},
										},
  }

  return node
}

// Commandine processing
func processFlags() (vId, iFile, oFile string, err error)  {
	var inputFileName, outputFileName, voiceId string

	argBase := 1
	empty := ""

	const (
		defaultVoiceId = "Raveena"
		appUsage = "\nUsage: pollyreveal [-v VoiceId] inputFile outputFile\nFor VoiceId refer to AWS Polly documentation.\n"
	)

	if len(os.Args) != 3  &&  len(os.Args) != 5  {
		return empty, empty, empty, errors.New(appUsage)
	}

	voiceId = DEFAULT_VOICE  // just because i am Indian
	if len(os.Args) == 5   {
		if os.Args[1] != "-v"  {
			return empty, empty, empty, errors.New(appUsage)
		}
		voiceId = os.Args[2]
		argBase = 3  // account for -v and voiceId
	}

	inputFileName = os.Args[argBase]
	outputFileName = os.Args[argBase + 1]

	return voiceId, inputFileName, outputFileName, nil
}
