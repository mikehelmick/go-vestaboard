// Copyright 2021 Mike Helmick
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package vestaboard

import (
	"errors"
	"fmt"
)

// PrintableChars is a string of all of the Vestaboard accepted chrs
const PrintableChars = " ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890!@#$() - +&=;: '\"%,.  /? °"

// Color represents constants for supported colors.
type Color int

const (
	PoppyRed Color = iota + 63
	Orange
	Yellow
	Green
	ParisBlue
	Violet
	White
)

var (
	ErrInvalidCharacter = errors.New("invalid character")
	ErrInvalidColor     = errors.New("invalid color")

	charNumbers map[string]int
)

func init() {
	charNumbers = make(map[string]int)
	for i, c := range PrintableChars {
		if _, ok := charNumbers[string(c)]; ok {
			// skip spaces.
			continue
		}
		charNumbers[string(c)] = i
	}
}

func CharToCode(c string) (int, error) {
	i, ok := charNumbers[c]
	if !ok {
		return -1, ErrInvalidCharacter
	}

	return i, nil
}

func ValidText(t string, newlineAccepted bool) error {
	for i, c := range t {
		if newlineAccepted && c == '\n' {
			continue
		}

		if _, err := CharToCode(string(c)); err != nil {
			return fmt.Errorf("invalid character %q at position %d, %w", string(c), i, err)
		}
	}
	return nil
}
