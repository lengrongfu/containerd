/*
   Copyright The containerd Authors.

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package platforms

import (
	"bufio"
	"io"
	"os"
	"runtime"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"golang.org/x/sys/unix"
)

func TestCPUVariant(t *testing.T) {
	if !isArmArch(runtime.GOARCH) {
		t.Skip("only relevant on linux/arm")
	}

	variants := []string{"v8", "v7", "v6", "v5", "v4", "v3"}

	p := getCPUVariant()
	for _, variant := range variants {
		if p == variant {
			t.Logf("got valid variant as expected: %#v = %#v\n", p, variant)
			return
		}
	}

	t.Fatalf("could not get valid variant as expected: %v\n", variants)
}

//go:noinline
func TestCPUVariant_MonkeyNoCpuArchitecture(t *testing.T) {
	if !isArmArch(runtime.GOARCH) {
		t.Skip("only relevant on linux/arm")
	}
	fackFile := &os.File{}
	patches := gomonkey.ApplyFunc(os.Open, func(name string) (*os.File, error) {
		return fackFile, nil
	})
	defer patches.Reset()
	patches.ApplyMethodFunc(fackFile, "Close", func() error {
		return nil
	})
	scan := bufio.NewScanner(fackFile)
	patches.ApplyFunc(bufio.NewScanner, func(r io.Reader) *bufio.Scanner {
		return scan
	})
	patches.ApplyMethodFunc(scan, "Scan", func() bool {
		return false
	})
	patches.ApplyFunc(unix.Uname, func(uname *unix.Utsname) error {
		// value is aarch64
		uname.Machine = [256]byte{97, 97, 114, 99, 104, 54, 52, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		return nil
	})
	variants := []string{"v8", "v7", "v6", "v5", "v4", "v3"}

	p := getCPUVariant()
	for _, variant := range variants {
		if p == variant {
			t.Logf("got valid variant as expected: %#v = %#v\n", p, variant)
			return
		}
	}

	t.Fatalf("could not get valid variant as expected: %v\n", variants)
}

//go:noinline
func TestCPUVariant_MonkeyCpuArchitecture(t *testing.T) {
	if !isArmArch(runtime.GOARCH) {
		t.Skip("only relevant on linux/arm")
	}
	fackFile := &os.File{}
	patches := gomonkey.ApplyFunc(os.Open, func(name string) (*os.File, error) {
		return fackFile, nil
	})
	defer patches.Reset()
	patches.ApplyMethodFunc(fackFile, "Close", func() error {
		return nil
	})
	scan := bufio.NewScanner(fackFile)
	patches.ApplyFunc(bufio.NewScanner, func(r io.Reader) *bufio.Scanner {
		return scan
	})
	patches.ApplyMethodFunc(scan, "Scan", func() bool {
		return true
	})
	patches.ApplyMethodFunc(scan, "Text", func() string {
		return "Cpu architecture: aarch64"
	})

	variants := []string{"v8", "v7", "v6", "v5", "v4", "v3"}

	p := getCPUVariant()
	for _, variant := range variants {
		if p == variant {
			t.Logf("got valid variant as expected: %#v = %#v\n", p, variant)
			return
		}
	}

	t.Fatalf("could not get valid variant as expected: %v\n", variants)
}
