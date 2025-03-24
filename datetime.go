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
	"strings"

	"github.com/dromara/carbon/v2"
)

type DateTimeExpr struct {
	Now          bool       `( @"now"`
	DateTime     *string    `| @Date "||"`
	Math         []MathExpr `) @@*`
	DateTimeOnly *string    `| @Date`
	carbon       *carbon.Carbon
}

func (e *DateTimeExpr) String() string {
	buf := strings.Builder{}
	if e.Instant() == e.Instant().StartOfDay() {
		buf.WriteString(e.Instant().String())
	} else {
		buf.WriteString(e.Instant().String())
	}
	if len(e.Math) == 0 {
		return buf.String()
	}

	buf.WriteString("||")
	buf.WriteString(e.Math[0].String())
	for i := 1; i < len(e.Math); i++ {
		buf.WriteString(e.Math[i].String())
	}
	return buf.String()
}

func (e *DateTimeExpr) Instant() *carbon.Carbon {
	if e.carbon != nil {
		return e.carbon
	}
	if e.Now {
		e.carbon = carbon.Now(carbon.Local)
	} else if e.DateTime != nil {
		e.carbon = carbon.Parse(*e.DateTime, carbon.Local)
	} else {
		e.carbon = carbon.Parse(*e.DateTimeOnly, carbon.Local)
	}
	return e.carbon
}

type MathExpr struct {
	Truncate string  `  "/" @TimeUnit`
	Summand  Summand `| @@`
}

func (m MathExpr) String() string {
	if m.Truncate != "" {
		return "/" + m.Truncate
	}
	return m.Summand.String()
}

type Summand struct {
	Sign     string `@("+"|"-")?`
	Number   int    `@Integer`
	TimeUnit string `@TimeUnit`
}

func (s Summand) String() string {
	return fmt.Sprintf("%s%d%s", s.Sign, s.Number, s.TimeUnit)
}

func (s Summand) Value() int {
	if s.Sign == "-" {
		return -s.Number
	}
	return s.Number
}
