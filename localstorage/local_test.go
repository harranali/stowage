package localstorage_test

import (
	"fmt"
	"os"
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
	//create full path to the root folder
	root, _ := filepath.Abs("./testdata/root")
	// initiate the loal storage
	l := New(root)
	// execute file info
	info, _ := l.FileInfo("filetotestinfo.md")

	// start asserting
	if info.Name != "filetotestinfo.md" {
		t.Error("failed asserting file info: Name")
	}
	if info.Extension != "md" {
		t.Error("failed asserting file info: Extension")
	}
	if info.Path != root {
		t.Error("failed asserting file info: Path")
	}
	if info.NameWithoutExtension != "filetotestinfo" {
		t.Error("failed asserting file info: NameWithoutExtension")
	}
	if info.IsDirectory != false {
		t.Error("failed asserting file info: IsDirectory")
	}
	if info.Size != 14 {
		t.Error("failed asserting file info: Size")
	}
}

func TestPut(t *testing.T) {
	//create full path to the root folder
	root, _ := filepath.Abs("./testdata/root")
	// initiate the loal storage
	l := New(root)

	// execute Put
	l.Put("./testdata/filetobeput.md")
	// assert
	_, err := os.Stat("testdata/root/filetobeput.md")
	if os.IsNotExist(err) {
		t.Error("failed assert putting file. ", err)
	}
	//cleanup
	os.Remove("testdata/root/filetobeput.md")
}

func TestPutAs(t *testing.T) {
	//create full path to the root folder
	root, _ := filepath.Abs("./testdata/root")
	// initiate the loal storage
	l := New(root)
	//execute PutAs
	l.PutAs("./testdata/filetobeput.md", "filetobeputnewname.md")
	// assert
	_, err := os.Stat("testdata/root/filetobeputnewname.md")
	if os.IsNotExist(err) {
		t.Error("failed assert putting file. ", err)
	}
	//cleanup
	os.Remove("testdata/root/filetobeputnewname.md")
}

func TestCopy(t *testing.T) {
	//create full path to the root folder
	root, _ := filepath.Abs("./testdata/root")
	// initiate the loal storage
	l := New(root)

	// execute copy to same dir
	err := l.Copy("filetocopy.md", "/")
	// assert copy to the same dir
	if err == nil {
		t.Error("failed assert copy error: file already exists. ", err)
	}

	// execute copy to different dir
	err = l.Copy("filetocopy.md", "/sub1/sub2")
	// check if for error
	if err != nil {
		t.Error("failed assert copy. ", err)
	}
	// assert file exesist
	_, err = os.Stat("testdata/root/sub1/sub2/filetocopy.md")
	if os.IsNotExist(err) {
		t.Error("failed assert copy file. ", err)
	}
	// cleanup
	os.RemoveAll("testdata/root/sub1")
}

func TestCopyAs(t *testing.T) {
	//create full path to the root folder
	root, _ := filepath.Abs("./testdata/root")
	// initiate the loal storage
	l := New(root)

	//execute copy to the root
	err := l.CopyAs("filetocopy.md", "/", "filetocopynew.md")
	//assert
	_, err = os.Stat("testdata/root/filetocopynew.md")
	if os.IsNotExist(err) {
		t.Error("failed assert file copyAs. ", err)
	}
	//cleanup
	os.RemoveAll("testdata/root/filetocopynew.md")

	// execute copy to sub dir
	err = l.CopyAs("filetocopy.md", "/sub1/sub2", "filetocopynew.md")
	if err != nil {
		t.Error("failed assert copyAs. ", err)
	}
	// assert
	_, err = os.Stat("testdata/root/sub1/sub2/filetocopynew.md")
	if os.IsNotExist(err) {
		t.Error("failed assert file copyAs. ", err)
	}
	// cleanup
	os.RemoveAll("testdata/root/sub1")
}

func TestMove(t *testing.T) {
	//create full path to the root folder
	root, _ := filepath.Abs("./testdata/root")
	// initiate the loal storage
	l := New(root)

	// Move(srcfile string, destfolder string) error

	// test moving to same
	err := l.Move("filetomove.md", "/")
	// assert
	if err == nil {
		t.Error("failed asserting error when moving file to same dir")
	}

	// test move file to sub dir
	err = l.Move("filetomove.md", "/sub1/sub2")
	// assert
	if err != nil {
		t.Error("failed asserting moving file to sub dir. ", err)
	}
	// assert soure file not there
	_, err = os.Stat("testdata/root/filetomove.md")
	if err == nil {
		t.Error("source file still present after moving. ")
	}

	// assert file exist in new dest
	_, err = os.Stat("testdata/root/sub1/sub2/filetomove.md")
	if os.IsNotExist(err) {
		t.Error("source file still present after moving. ")
	}

	// move the file back
	err = l.Move("/sub1/sub2/filetomove.md", "/")
	// delete the sub dirs
	os.RemoveAll("/sub1")
}
