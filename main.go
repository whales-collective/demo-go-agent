package main

import (
	"context"
	"fmt"
	"os"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

func main() {

	modelRunnerBaseUrl := os.Getenv("MODEL_RUNNER_BASE_URL")

	if modelRunnerBaseUrl == "" {
		panic("MODEL_RUNNER_BASE_URL environment variable is not set")
	}
	modelRunnerChatModel := os.Getenv("MODEL_RUNNER_CHAT_MODEL")
	fmt.Println("Using Model Runner Chat Model:", modelRunnerChatModel)

	if modelRunnerChatModel == "" {
		panic("MODEL_RUNNER_CHAT_MODEL environment variable is not set")
	}

	ctx := context.Background()

	clientEngine := openai.NewClient(
		option.WithBaseURL(modelRunnerBaseUrl),
		option.WithAPIKey(""),
	)

    agentInstructions := `
    You are Dungeon Master Bob, a Dwarf in a fantasy world. 
    You are friendly and helpful, but you can also be mischievous. 
    `

	// Chat Completion parameters
	chatCompletionParams := openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(agentInstructions),
			openai.UserMessage("Tell mee a story about a Werewolf in the forest."),
		},
		Model:       modelRunnerChatModel,
		Temperature: openai.Opt(0.8),
	}

	stream := clientEngine.Chat.Completions.NewStreaming(ctx, chatCompletionParams)

	for stream.Next() {
		chunk := stream.Current()
		// Stream each chunk as it arrives
		if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != "" {
			fmt.Print(chunk.Choices[0].Delta.Content)
		}
	}

	if err := stream.Err(); err != nil {
		fmt.Printf("😡 Stream error: %v\n", err)
	}

	fmt.Println()

}
