// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0
package jsonarray // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/operator/parser/jsonarray"

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/valyala/fastjson"
	"go.opentelemetry.io/collector/featuregate"
	"go.uber.org/zap"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/entry"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/operator"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/operator/helper"
)

const operatorType = "json_array_parser"
const headerDelimiter = ","

var jsonArrayParserFeatureGate = featuregate.GlobalRegistry().MustRegister(
	"logs.jsonParserArray",
	featuregate.StageAlpha,
	featuregate.WithRegisterDescription("When enabled, allows usage of `json_array_parser`."),
	featuregate.WithRegisterReferenceURL("https://github.com/open-telemetry/opentelemetry-collector-contrib/issues/30321"),
)

func init() {
	if jsonArrayParserFeatureGate.IsEnabled() {
		operator.Register(operatorType, func() operator.Builder { return NewConfig() })
	}
}

// NewConfig creates a new json array parser config with default values
func NewConfig() *Config {
	return NewConfigWithID(operatorType)
}

// NewConfigWithID creates a new json array parser config with default values
func NewConfigWithID(operatorID string) *Config {
	return &Config{
		ParserConfig: helper.NewParserConfig(operatorID, operatorType),
	}
}

// Config is the configuration of a json array parser operator.
type Config struct {
	helper.ParserConfig `mapstructure:",squash"`
	Header              string `mapstructure:"header"`
}

// Build will build a json array parser operator.
func (c Config) Build(logger *zap.SugaredLogger) (operator.Operator, error) {
	parserOperator, err := c.ParserConfig.Build(logger)
	if err != nil {
		return nil, err
	}

	if c.Header != "" {
		return &Parser{
			ParserOperator: parserOperator,
			parse:          generateParseToMapFunc(new(fastjson.ParserPool), strings.Split(c.Header, headerDelimiter)),
		}, nil
	}

	return &Parser{
		ParserOperator: parserOperator,
		parse:          generateParseToArrayFunc(new(fastjson.ParserPool)),
	}, nil
}

// Parser is an operator that parses json array in an entry.
type Parser struct {
	helper.ParserOperator
	parse parseFunc
}

type parseFunc func(any) (any, error)

// Process will parse an entry for json array.
func (r *Parser) Process(ctx context.Context, e *entry.Entry) error {
	return r.ParserOperator.ProcessWith(ctx, e, r.parse)
}

func generateParseToArrayFunc(pool *fastjson.ParserPool) parseFunc {
	return func(value any) (any, error) {
		jArrayLine, err := valueAsString(value)
		if err != nil {
			return nil, err
		}

		p := pool.Get()
		v, err := p.Parse(jArrayLine)
		pool.Put(p)
		if err != nil {
			return nil, errors.New("failed to parse entry")
		}

		jArray := v.GetArray() // a is a []*Value slice
		parsedValues := make([]any, len(jArray))
		for i := range jArray {
			switch jArray[i].Type() {
			case fastjson.TypeNumber:
				parsedValues[i] = jArray[i].GetInt64()
			case fastjson.TypeString:
				parsedValues[i] = string(jArray[i].GetStringBytes())
			case fastjson.TypeTrue:
				parsedValues[i] = true
			case fastjson.TypeFalse:
				parsedValues[i] = false
			case fastjson.TypeNull:
				parsedValues[i] = nil
			case fastjson.TypeObject:
				// Nested objects handled as a string since this parser doesn't support nested headers
				parsedValues[i] = jArray[i].String()
			default:
				return nil, errors.New("failed to parse entry: " + string(jArray[i].MarshalTo(nil)))
			}
		}

		return parsedValues, nil
	}
}

func generateParseToMapFunc(pool *fastjson.ParserPool, header []string) parseFunc {
	return func(value any) (any, error) {
		jArrayLine, err := valueAsString(value)
		if err != nil {
			return nil, err
		}

		p := pool.Get()
		v, err := p.Parse(jArrayLine)
		pool.Put(p)
		if err != nil {
			return nil, errors.New("failed to parse entry")
		}

		jArray := v.GetArray() // a is a []*Value slice
		if len(header) != len(jArray) {
			return nil, fmt.Errorf("wrong number of fields: expected %d, found %d", len(header), len(jArray))
		}
		parsedValues := make(map[string]any, len(jArray))
		for i := range jArray {
			switch jArray[i].Type() {
			case fastjson.TypeNumber:
				parsedValues[header[i]] = jArray[i].GetInt64()
			case fastjson.TypeString:
				parsedValues[header[i]] = string(jArray[i].GetStringBytes())
			case fastjson.TypeTrue:
				parsedValues[header[i]] = true
			case fastjson.TypeFalse:
				parsedValues[header[i]] = false
			case fastjson.TypeNull:
				parsedValues[header[i]] = nil
			case fastjson.TypeObject:
				// Nested objects handled as a string since this parser doesn't support nested headers
				parsedValues[header[i]] = jArray[i].String()
			default:
				return nil, errors.New("failed to parse entry: " + string(jArray[i].MarshalTo(nil)))
			}
		}

		return parsedValues, nil
	}
}

// valueAsString interprets the given value as a string.
func valueAsString(value any) (string, error) {
	var s string
	switch t := value.(type) {
	case string:
		s += t
	case []byte:
		s += string(t)
	default:
		return s, fmt.Errorf("type '%T' cannot be parsed as json array", value)
	}

	return s, nil
}
