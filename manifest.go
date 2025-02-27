// Copyright 2020, 2021 Google LLC
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

package runfiles

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ManifestFile specifies the location of the runfile manifest file.  You can
// pass this as an option to New.  If unset or empty, use the value of the
// environmental variable RUNFILES_MANIFEST_FILE.
type ManifestFile string

func (f ManifestFile) new() (*Runfiles, error) {
	m, err := f.parse()
	if err != nil {
		return nil, err
	}
	return &Runfiles{m, manifestFileVar + "=" + string(f)}, nil
}

type manifest map[string]string

func (f ManifestFile) parse() (manifest, error) {
	r, err := os.Open(string(f))
	if err != nil {
		return nil, fmt.Errorf("runfiles: can’t open manifest file: %w", err)
	}
	defer r.Close()
	s := bufio.NewScanner(r)
	m := make(manifest)
	for s.Scan() {
		fields := strings.SplitN(s.Text(), " ", 2)
		if len(fields) != 2 || fields[0] == "" {
			return nil, fmt.Errorf("runfiles: bad manifest line %q in file %s", s.Text(), f)
		}
		m[fields[0]] = filepath.FromSlash(fields[1])
	}
	if err := s.Err(); err != nil {
		return nil, fmt.Errorf("runfiles: error parsing manifest file %s: %w", f, err)
	}
	return m, nil
}

func (m manifest) path(s string) (string, error) {
	r, ok := m[s]
	if !ok {
		return "", os.ErrNotExist
	}
	if r == "" {
		return "", ErrEmpty
	}
	return r, nil
}

const manifestFileVar = "RUNFILES_MANIFEST_FILE"
