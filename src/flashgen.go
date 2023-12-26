package flashgen

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/Clever/csvlint"
	openai "github.com/sashabaranov/go-openai"
)

type ModelConfig struct {
	Model string
	SystemMessage string
	MaxTokens int
}

const FlashcardCSVDelimiter = ','

const ReaderStride = 2048

type Input struct {
	Name string
	Data string
	Tags []string
}

func IsPDF(filePath string) bool {
	return path.Ext(filePath) == ".pdf"
}


func ConvertPDFToText(filePath string) (string, error) {
	if !IsPDF(filePath) {
		return "", errors.New("file is not a pdf")
	}
	return "", errors.New("not implemented")
}

func applyDefaultConfig(config ModelConfig) (ModelConfig, error) {
	if config.Model == "" {
		config.Model = openai.GPT4
	}
	if config.SystemMessage == "" {
		defaultSystemMessageFilePath := "src/prompts/anki-flashcards-csv.txt"
		b, err := os.ReadFile(defaultSystemMessageFilePath)
		if err != nil {
			return ModelConfig{}, err
		}
		config.SystemMessage = string(b)
	}
	if config.MaxTokens == 0 {
		config.MaxTokens = 256
	}
	return config, nil
}


func RequestFlashcardsFromLLM(config ModelConfig, input *Input) (string , error) {

	config, err := applyDefaultConfig(config)
	if err != nil {
		return "", err
	}

	openAIKey := os.Getenv("OPEN_AI_KEY")
	client := openai.NewClient(openAIKey)

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: config.Model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role: openai.ChatMessageRoleSystem,
					Content: config.SystemMessage,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: input.Data,
				},
			},
			MaxTokens: config.MaxTokens,
		},
	)
	if err != nil {
		return "", fmt.Errorf("ChatCompletion error: %v", err)
	}

	return resp.Choices[0].Message.Content, nil
} 

func Clean(flashcards string) string {
	return strings.TrimSuffix(flashcards, "\n")
}

func GenerateFlashcards(inputFilePath string, outputFilePath string, tags string) error {
	
	if IsPDF(inputFilePath) {
		fp, err := ConvertPDFToText(inputFilePath)
		inputFilePath = fp
		if err != nil {
			return err
		}
	}

	inputFile, err := os.Open(inputFilePath)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	if outputFilePath == "" {
		ext := path.Ext(inputFilePath)
		outputFilePath = inputFilePath[0:len(inputFilePath)-len(ext)] + ".csv" 
		if outputFilePath == inputFilePath {
			outputFilePath = "flashgen-" + outputFilePath
		}
	}
	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	config := ModelConfig{}
	
	r := bufio.NewReader(inputFile)
	w := bufio.NewWriter(outputFile)
	stride := ReaderStride
	buf := make([]byte, 0, stride)
	for {
        n, err := io.ReadFull(r, buf[:cap(buf)])
        buf = buf[:n]
        if err != nil {
            if err == io.EOF {
                break
            }
            if err != io.ErrUnexpectedEOF {
				return err
            }
        }

		input := Input{
			Name: inputFilePath,
			Data: string(buf),
			Tags: []string{},
		}

		flashcards, err := RequestFlashcardsFromLLM(config, &input)
		if err != nil {
			return err
		}

		flashcards = Clean(flashcards)
		
		flashcards += "\n"
		_, err = w.WriteString(flashcards)
		if err != nil {
			return err
		}

		w.Flush()
	}
	return nil	
}

func ValidateFlashcards(flashcardsFilePath string) error {
	flashcardsFile , err := os.Open(flashcardsFilePath) 
	if err != nil {
		return err 
	}
	csvErrors, _, err := csvlint.Validate(flashcardsFile, FlashcardCSVDelimiter, false)		
	if err != nil {
		return err
	}
	
	if len(csvErrors) > 0 {
		fmt.Printf("Flashcards were generated and saved in %s, but there were some errors:\n", flashcardsFilePath)
		for _, csvError := range csvErrors {
			fmt.Println(csvError)
		}
	}
	return nil
}

func ImportFlashcardsIntoAnki(flashcardsFilePath string, deck string) error {
	return errors.New("not implemented")
}