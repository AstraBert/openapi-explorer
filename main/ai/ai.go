package ai

import (
	"context"
	"encoding/json"

	"github.com/invopop/jsonschema"
	"github.com/openai/openai-go/v2"
)

func generateSchema[T any]() any {
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}
	var v T
	schema := reflector.Reflect(v)
	return schema
}

func StructuredChat[T any](message, systemMessage, schemaName, schemaDescription string) (any, error) {

	client := openai.NewClient()
	ctx := context.Background()

	structuredOutputSchema := generateSchema[T]()

	schemaParam := openai.ResponseFormatJSONSchemaJSONSchemaParam{
		Name:        schemaName,
		Description: openai.String(schemaDescription),
		Schema:      structuredOutputSchema,
		Strict:      openai.Bool(true),
	}

	chat, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(systemMessage),
			openai.UserMessage(message),
		},
		ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
			OfJSONSchema: &openai.ResponseFormatJSONSchemaParam{
				JSONSchema: schemaParam,
			},
		},
		Model: openai.ChatModelGPT4_1,
	})

	if err != nil {
		return nil, err
	}

	// extract into a well-typed struct
	var structuredOutput T
	_ = json.Unmarshal([]byte(chat.Choices[0].Message.Content), &structuredOutput)
	return structuredOutput, nil
}
