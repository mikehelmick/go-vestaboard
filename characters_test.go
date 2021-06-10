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
	"strings"
	"testing"
)

func TestCharIndex(t *testing.T) {
	t.Parallel()

	for want, inC := range PrintableChars {
		if c := string(inC); c == " " && want > 0 {
			continue
		}
		got, err := CharToCode(string(inC))
		if err != nil {
			t.Errorf("error on %q, err: %v", inC, err)
		}

		if want != got {
			t.Errorf("wrong char code for %q, want: %v, got %v", inC, want, got)
		}

		if inC >= 'A' && inC <= 'Z' {
			invalid := strings.ToLower(string(inC))
			_, err := CharToCode(string(invalid))
			if err == nil {
				t.Errorf("expected error, got nil")
			}
			if !errors.Is(err, ErrInvalidCharacter) {
				t.Errorf("wrong error: %v", err)
			}
		}
	}
}
