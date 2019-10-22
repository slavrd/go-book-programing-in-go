// Copyright © 2011-12 Qtrac Ltd.
//
// This program or package and any associated files are licensed under the
// Apache License, Version 2.0 (the "License"); you may not use these files
// except in compliance with the License. You can get a copy of the License
// at: http://www.apache.org/licenses/LICENSE-2.0.
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestBigDigits(t *testing.T) {
	log.SetFlags(0)
	log.Println("TEST bigdigits")

	path, _ := os.Getwd()
	executable := filepath.Join(path, "bigdigits")

	type testCase struct {
		command     *exec.Cmd
		expFileName string
	}

	tests := []testCase{
		{
			command:     exec.Command(executable, "0123456789"),
			expFileName: "0123456789.txt",
		},
		{
			command:     exec.Command(executable, "--bar", "0123456789"),
			expFileName: "0123456789-bar.txt",
		},
	}

	for _, test := range tests {

		expected, err := ioutil.ReadFile(filepath.Join(path, test.expFileName))
		if err != nil {
			t.Fatal(err)
		}

		reader, writer, err := os.Pipe()
		if err != nil {
			t.Fatal(err)
		}

		test.command.Stdout = writer
		err = test.command.Run()
		if err != nil {
			t.Fatal(err)
		}
		writer.Close()
		actual, err := ioutil.ReadAll(reader)
		if err != nil {
			t.Fatal(err)
		}
		reader.Close()
		if bytes.Compare(actual, expected) != 0 {
			t.Fatal("actual != expected")
		}
	}
}
