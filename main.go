package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	spinner "github.com/briandowns/spinner"

	flashgen "github.com/igoragoli/flashgen/src"
)

func printHelpAndExit(code int) {
  fmt.Println(`flashgen [-anki] [-deck "<deck-name>"] [-tags "<tag1> <tag2> ...<tagn>]" <input-file> [<output-file>]`)
	flag.PrintDefaults()
	os.Exit(code)
}

func main() {
  anki := flag.Bool("anki", false, "(not implemented) whether to automatically import generated flashcards in Anki with the AnkiConnect add-on.") 
  deck := flag.String("deck", "", "(not implemented) the name of the deck to import the flashcards into.")
  tags := flag.String("tags", "", "(not implemented) a comma-separated list of tags to apply to the imported flashcards.")
	help := flag.Bool("help", false, "print help and exit")
	flag.Parse()

	if *help {
		printHelpAndExit(0)
	}

  if len(flag.Args()) < 1 || len(flag.Args()) > 2 {
		fmt.Print("flashgen accepts an input filepath and an optional output filepath\n\n")
		printHelpAndExit(1)
	}

  if *anki || *deck != "" || *tags != "" {
    fmt.Println("Anki integration and tags are not implemented yet.")
    os.Exit(1)
  }

  if _, ok := os.LookupEnv("OPEN_AI_KEY"); !ok {
    fmt.Println("OPEN_AI_KEY environment variable is not set.")
    os.Exit(1)
  }

  inputFilePath := flag.Args()[0]
  var outputFilePath string
  if len(flag.Args()) == 2 {
    outputFilePath = flag.Args()[1]
  } 

  s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
  s.Prefix = "Generating flashcards: " // Prefix text before the spinner
  s.Start()
  flashgen.GenerateFlashcards(inputFilePath, outputFilePath, *tags)
  s.Stop()

  flashgen.ValidateFlashcards(outputFilePath)

  if *anki {
    fmt.Println("Importing flashcards into Anki...")
    err := flashgen.ImportFlashcardsIntoAnki(outputFilePath, *deck)
    if err != nil {
      panic(err)
    }
  }

}