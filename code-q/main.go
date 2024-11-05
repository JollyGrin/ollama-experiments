package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	ollama "github.com/ollama/ollama/api"
)

// Function to read a code file and include line numbers
func extractTextFromCodeFile(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var text bytes.Buffer
	scanner := bufio.NewScanner(file)
	lineNumber := 1
	for scanner.Scan() {
		// Write each line with its line number
		text.WriteString(fmt.Sprintf("%d: %s\n", lineNumber, scanner.Text()))
		lineNumber++
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return text.String(), nil
}

func main() {
	// Check if the filename is provided as an argument
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go -- <filename>")
		os.Exit(1)
	}

	// Get the filename from the argument
	filename := os.Args[2]

	client, err := ollama.ClientFromEnvironment()
	if err != nil {
		log.Fatal(err)
	}

	// Extract text from the specified code file
	codeText, err := extractTextFromCodeFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	// Limit text size if necessary (e.g., API might have a character limit)
	if len(codeText) > 15000 { // Adjust this limit as needed
		codeText = codeText[:15000] + "..."
	}

	// Initial context for the model, explaining the purpose
	conversationContext := "You are a code assistant helping to analyze code files. Please provide concise answers, focused mainly on the code itself. When the user references specific line numbers, respond directly regarding the relevant lines.\n\nHere is the code file content:\n" + codeText + "\n\n"

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Code file content loaded. Start your conversation with the model or type 'exit' to quit:")

	for {
		// Get user input
		fmt.Print("You: ")
		if !scanner.Scan() {
			break
		}
		userInput := scanner.Text()
		if strings.ToLower(userInput) == "exit" {
			break
		}

		// Append user input to the context
		conversationContext += "User: " + userInput + "\nModel: "

		// Set up the request
		req := &ollama.GenerateRequest{
			Model:  "llama3.2",
			Prompt: conversationContext,
		}

		ctx := context.Background()

		// Stream the response in real-time
		fmt.Print("Model: ")
		respFunc := func(resp ollama.GenerateResponse) error {
			fmt.Print(resp.Response) // Print each chunk immediately as it arrives
			return nil
		}

		// Send request
		err = client.Generate(ctx, req, respFunc)
		if err != nil {
			log.Fatal(err)
		}

		// Add a newline after the model's response completes
		fmt.Println()
		// Update the conversation context with the model's response
		conversationContext += "\n\n"
	}
}
