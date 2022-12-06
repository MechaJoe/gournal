package main

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"os"

// 	"github.com/PullRequestInc/go-gpt3"
// 	"github.com/joho/godotenv"
// )

// func main() {
// 	godotenv.Load()

// 	apiKey := os.Getenv("API_KEY")
// 	if apiKey == "" {
// 		log.Fatalln("Missing API KEY")
// 	}

// 	ctx := context.Background()
// 	client := gpt3.NewClient(apiKey)

// 	resp, err := client.Completion(ctx, gpt3.CompletionRequest{
// 		Prompt:    []string{"Hey Gournal. Today I want to solve that bug, and get a haircut."},
// 		MaxTokens: gpt3.IntPtr(40),
// 		Stop:      []string{"."},
// 		Echo:      true,
// 	})
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// 	fmt.Println(resp.Choices[0].Text)
// }