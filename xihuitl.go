// Copyright 2025 The Xihuitl authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package xihuitl

import (
	"fmt"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/dromara/carbon/v2"
)

func Parse(expr string) (*carbon.Carbon, error) {
	ast, err := createParser().ParseString("", expr)
	if err != nil {
		return nil, err
	}
	return applyMath(ast)
}

func createParser() *participle.Parser[DateTimeExpr] {
	lex := lexer.MustSimple([]lexer.SimpleRule{
		{"Date", `\d{4}-\d{2}-\d{2}`},
		{"Time", `\d\d:\d\d(:\d\d(\.\d+)?)?`},
		{"DateTime", `\d{4}-\d{2}-\d{2}T\d\d:\d\d(:\d\d(\.\d+)?)?`},

		{"TimeUnit", `[smHhdwMqQyY]`},
		{"RangeUnit", `(d|h|m|s|ms|micros|nanos)`},

		{"Ident", `[a-zA-Z_][a-zA-Z_0-9]*`},
		{"Integer", `[0-9]+`},
		{"Special", `(\|\|)|[/+-]`},
		{"Whitespace", `\s+`},
	})

	p, _ := participle.Build[DateTimeExpr](
		participle.Lexer(lex),
		participle.Elide("Whitespace"),
	)
	return p
}

func applyMath(e *DateTimeExpr) (c *carbon.Carbon, err error) {
	c = e.Instant()
	for _, m := range e.Math {
		if m.Truncate != "" {
			c, err = truncate(c, m.Truncate)
		} else {
			c, err = calc(c, m.Summand.Value(), m.Summand.TimeUnit)
		}
		if err != nil {
			return c, err
		}
	}
	return
}

func calc(c *carbon.Carbon, val int, unit string) (*carbon.Carbon, error) {
	switch unit {
	case "s":
		return c.AddSeconds(val), nil
	case "m":
		return c.AddMinutes(val), nil
	case "H", "h":
		return c.AddHours(val), nil
	case "d":
		return c.AddDays(val), nil
	case "w":
		return c.AddDays(val * 7), nil
	case "M":
		return c.AddMonths(val), nil
	case "q", "Q":
		return c.AddQuarters(val), nil
	case "y", "Y":
		return c.AddYears(val), nil
	default:
		return c, fmt.Errorf("unknown unit %q", unit)
	}
}

func truncate(c *carbon.Carbon, unit string) (*carbon.Carbon, error) {
	switch unit {
	case "s":
		return c.StartOfSecond(), nil
	case "m":
		return c.StartOfMinute(), nil
	case "H", "h":
		return c.StartOfHour(), nil
	case "d":
		return c.StartOfDay(), nil
	case "w":
		return c.StartOfWeek(), nil
	case "M":
		return c.StartOfMonth(), nil
	case "q", "Q":
		return c.StartOfQuarter(), nil
	case "y", "Y":
		return c.StartOfYear(), nil
	default:
		return c, fmt.Errorf("unknown unit %q", unit)
	}
}
