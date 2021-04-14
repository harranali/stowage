package localstorage_test

import (
	"fmt"
	"path/filepath"
	"testing"

	. "github.com/harranali/stowage/localstorage"
)

func TestNew(t *testing.T) {
	root, _ := filepath.Abs("./testdata/root")
	l := New(root)

	if fmt.Sprintf("%T", l) != "*localstorage.LocalStorage" {
		t.Error("failed initiating localstorage")
	}
}

func TestFileInfo(t *testing.T) {
	root, _ := filepath.Abs("./testdata/root")
	l := New(root)
	info, _ := l.FileInfo("file2.md")

	if info.Name != "file2.md" {
		t.Error("failed asserting file info: Name")
	}

	if info.Extension != "md" {
		t.Error("failed asserting file info: Extension")
	}

	if info.Path != root {
		t.Error("failed asserting file info: Path")
	}

	if info.NameWithoutExtension != "file2" {
		t.Error("failed asserting file info: NameWithoutExtension")
	}

	if info.IsDirectory != false {
		t.Error("failed asserting file info: IsDirectory")
	}

	if info.Size != 21 {
		t.Error("failed asserting file info: Size")
	}
}
