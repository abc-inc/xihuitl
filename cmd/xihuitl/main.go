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

package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/abc-inc/xihuitl"
	"github.com/dromara/carbon/v2"
)

func main() {
	if len(os.Args) != 2 {
		usage()
	}

	c, err := xihuitl.Parse(os.Args[1])
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(c.Format(strings.TrimSuffix(carbon.ISO8601Format, "P")))
}

func usage() {
	bin := filepath.Base(os.Args[0])
	_, _ = fmt.Fprintln(os.Stderr, heredoc.Docf(`
		Calculate date and time relative to the given anchor date.

		Usage: %s EXPR

		Examples:
		  current time with seconds precision (2006-01-02T15:04:05)
		  %s now

		  previous day (2006-01-02)
		  %s "now-1d/d"

		  apply date math to a given date
		  %s "2025-01-01||-1w-1d+18h+30m"
	`, bin, bin, bin, bin))

	os.Exit(1)
}
