package instructor

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/binarycraft007/instructor-go/pkg/instructor/googleai"
)

type UsageSum struct {
	InputTokens  int
	OutputTokens int
	TotalTokens  int
}

func chatHandler(i Instructor, ctx context.Context, request interface{}, response any) (interface{}, error) {

	var err error
	var schema interface{}

	t := reflect.TypeOf(response)

	if i.Provider() == ProviderGoogleAI {
		schema, err = googleai.GenerateSchemaFromType(t)
		if err != nil {
			return nil, err
		}
	} else {
		schema, err = NewSchema(t)
		if err != nil {
			return nil, err
		}
	}

	// keep a running total of usage
	usage := &UsageSum{}

	for attempt := 0; attempt <= i.MaxRetries(); attempt++ {

		text, resp, err := i.chat(ctx, request, schema)
		if err != nil {
			// no retry on non-marshalling/validation errors
			return i.emptyResponseWithResponseUsage(resp), err
		}

		text = extractJSON(&text)

		err = json.Unmarshal([]byte(text), &response)
		if err != nil {
			// TODO:
			// add more sophisticated retry logic (send back json and parse error for model to fix).
			//
			// Currently, its just recalling with no new information
			// or attempt to fix the error with the last generated JSON

			i.countUsageFromResponse(resp, usage)
			continue
		}

		if i.Validate() {
			validate = validator.New()
			// Validate the response structure against the defined model using the validator
			err = validate.Struct(response)

			if err != nil {
				// TODO:
				// add more sophisticated retry logic (send back validator error and parse error for model to fix).

				i.countUsageFromResponse(resp, usage)
				continue
			}
		}

		return i.addUsageSumToResponse(resp, usage)
	}

	return i.emptyResponseWithUsageSum(usage), errors.New("hit max retry attempts")
}
