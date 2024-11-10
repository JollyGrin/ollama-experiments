package main

import (
	"context"
	"fmt"
	"log"
	"os/exec"

	ollama "github.com/ollama/ollama/api"
)

func call() {
	diffOutput, err := getStagedDiff()
	if err != nil {
		log.Fatalf("Error getting staged diff: %v", err)
	}

	if diffOutput == "" {
		fmt.Println("No staged changes to display.")
	} else {
		fmt.Println("Staged changes:")
		fmt.Println(diffOutput)
	}
}

// getStagedDiff runs `git diff --cached` to get the diff for all staged changes
func getStagedDiff() (string, error) {
	cmd := exec.Command("git", "diff", "--cached")
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func main() {
	// Extract text from the specified code file
	codeText, err := getStagedDiff()
	if err != nil {
		log.Fatal(err)
	}

	client, err := ollama.ClientFromEnvironment()
	if err != nil {
		log.Fatal(err)
	}

	prompt := "I just made some updates to my code and it generated these git diffs. Could you analyze these git diffs and think of a concise yet informative summary that I can copy and paste into a commit message? Here is the git diffs:\n\n" + codeText
	// By default, GenerateRequest is streaming.
	req := &ollama.GenerateRequest{
		Model:  "llama3.2",
		Prompt: prompt,
	}

	ctx := context.Background()
	respFunc := func(resp ollama.GenerateResponse) error {
		// Only print the response here; GenerateResponse has a number of other
		// interesting fields you want to examine.

		// In streaming mode, responses are partial so we call fmt.Print (and not
		// Println) in order to avoid spurious newlines being introduced. The
		// model will insert its own newlines if it wants.
		fmt.Print(resp.Response)
		return nil
	}

	err = client.Generate(ctx, req, respFunc)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println()
}
