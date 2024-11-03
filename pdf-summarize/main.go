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
	"rsc.io/pdf"
)

func extractTextFromPDF(filename string) (string, error) {
	// Open the PDF file
	file, err := pdf.Open(filename)
	if err != nil {
		return "", err
	}

	var text bytes.Buffer
	// Iterate over each page
	for i := 1; i <= file.NumPage(); i++ {
		page := file.Page(i)
		content := page.Content()

		// Go through the page content and collect text
		for _, textObj := range content.Text {
			text.WriteString(textObj.S)
			text.WriteString(" ")
		}
	}

	return text.String(), nil
}

func main() {
	client, err := ollama.ClientFromEnvironment()
	if err != nil {
		log.Fatal(err)
	}

	// Extract text from the PDF file
	pdfText, err := extractTextFromPDF("pdf/tax.pdf")
	if err != nil {
		log.Fatal(err)
	}

	// Limit text size if necessary (e.g., API might have a character limit)
	if len(pdfText) > 15000 { // Adjust this limit as needed
		pdfText = pdfText[:15000] + "..."
	}

	// Initial context for the model, explaining what it's reading
	conversationContext := "I'm using a tool to extract text from a PDF. It might not always come out perfect, but I need your help to parse and summarize the contents. If it's in a different language, please summarize in English. Here's the content extracted:\n" + pdfText + "\n\n"

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("PDF content loaded. Start your conversation with the model or type 'exit' to quit:")

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
		var modelResponse strings.Builder
		respFunc := func(resp ollama.GenerateResponse) error {
			modelResponse.WriteString(resp.Response)
			return nil
		}

		// Send request
		err = client.Generate(ctx, req, respFunc)
		if err != nil {
			log.Fatal(err)
		}

		// Display and save the model's response
		fmt.Println("Model:", modelResponse.String())
		conversationContext += modelResponse.String() + "\n\n"
	}
}
