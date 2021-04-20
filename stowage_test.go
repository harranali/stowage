package stowage_test

import (
	"fmt"
	"path/filepath"
	"testing"

	. "github.com/harranali/stowage"
)

func TestNew(t *testing.T) {
	s := New()
	if fmt.Sprintf("%T", s) != "*stowage.Stowage" {
		t.Error("faild assert var type")
	}
}

func TestInitLocalStorage(t *testing.T) {
	s := New()
	root, _ := filepath.Abs("./localstorage/testdata/root")
	s.InitLocalStorage(LocalStorageOpts{
		RootFolder: root,
	})

	info, err := s.LocalStorage.FileInfo("filetotestinfo.md")
	if err != nil {
		t.Error("failed assert getting state: ", err)
	}
	if info.Name != "filetotestinfo.md" {
		t.Error("failed assert reading file name")
	}
}
