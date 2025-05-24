package utils

import (
	"context"
	"fmt"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type GeminiClient struct {
	client *genai.Client
}

func NewGeminiClient(apiKey string) (*GeminiClient, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}

	return &GeminiClient{
		client: client,
	}, nil
}

func (g *GeminiClient) GenerateResponse(systemPrompt, userMessage string) (string, error) {
	ctx := context.Background()
	model := g.client.GenerativeModel("gemini-pro")

	// Set system instruction
	model.SystemInstruction = &genai.Content{
		Parts: []genai.Part{
			genai.Text(systemPrompt),
		},
	}

	// Generate response
	resp, err := model.GenerateContent(ctx, genai.Text(userMessage))
	if err != nil {
		return "", err
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no response generated")
	}

	// Extract text from response
	if textPart, ok := resp.Candidates[0].Content.Parts[0].(genai.Text); ok {
		return string(textPart), nil
	}

	return "", fmt.Errorf("unexpected response format")
}

func (g *GeminiClient) GenerateResponseWithImage(systemPrompt, userMessage, imageURL string) (string, error) {
	ctx := context.Background()
	model := g.client.GenerativeModel("gemini-pro-vision")

	// Set system instruction
	model.SystemInstruction = &genai.Content{
		Parts: []genai.Part{
			genai.Text(systemPrompt),
		},
	}

	// Create image part (simplified - in production you'd fetch the image)
	prompt := fmt.Sprintf("%s\n\nUser has shared an image. Please analyze it in the context of their message: %s", systemPrompt, userMessage)

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", err
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no response generated")
	}

	if textPart, ok := resp.Candidates[0].Content.Parts[0].(genai.Text); ok {
		return string(textPart), nil
	}

	return "", fmt.Errorf("unexpected response format")
}

func (g *GeminiClient) Close() error {
	return g.client.Close()
}
