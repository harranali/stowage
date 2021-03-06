// Copyright 2021 Harran Ali <harran.m@gmail.com>. All rights reserved.
// Use of this source code is governed by MIT-style
// license that can be found in the LICENSE file.

package localstorage_test

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
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
	info, err := l.FileInfo("filetotestinfo.md")
	if err != nil {
		t.Error("failed asserting file info")
	}

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

	if info.FsFileInfo.Name() != "filetotestinfo.md" {
		t.Error("failed asserting file info: FsFileInfo.Name()")
	}
}

func TestPut(t *testing.T) {
	//create full path to the root folder
	root, _ := filepath.Abs("./testdata/root")
	// initiate the loal storage
	l := New(root)

	// execute Put
	err := l.Put("./testdata/filetobeput.md")
	if err != nil {
		t.Error("failed assert putting file. ", err)
	}
	// assert
	_, err = os.Stat("testdata/root/filetobeput.md")
	if os.IsNotExist(err) {
		t.Error("failed assert putting file. ", err)
	}
	//cleanup
	err = os.Remove("testdata/root/filetobeput.md")
	if err != nil {
		t.Error("failed assert putting file. ", err)
	}
}

func TestPutAs(t *testing.T) {
	//create full path to the root folder
	root, _ := filepath.Abs("./testdata/root")
	// initiate the loal storage
	l := New(root)
	//execute PutAs
	err := l.PutAs("./testdata/filetobeput.md", "filetobeputnewname.md")
	if err != nil {
		t.Error("failed assert putting file. ", err)
	}
	// assert
	_, err = os.Stat("testdata/root/filetobeputnewname.md")
	if os.IsNotExist(err) {
		t.Error("failed assert putting file. ", err)
	}
	//cleanup
	err = os.Remove("testdata/root/filetobeputnewname.md")
	if err != nil {
		t.Error("failed assert putting file. ", err)
	}
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
	err = os.RemoveAll("testdata/root/sub1")
	if err != nil {
		t.Error("failed assert copy. ", err)
	}
}

func TestCopyAs(t *testing.T) {
	//create full path to the root folder
	root, _ := filepath.Abs("./testdata/root")
	// initiate the loal storage
	l := New(root)

	//execute copy to the root
	err := l.CopyAs("filetocopy.md", "/", "filetocopynew.md")
	if err != nil {
		t.Error("failed assert file copyAs. ", err)
	}
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
	err = os.RemoveAll("testdata/root/sub1")
	if err != nil {
		t.Error("failed assert file copyAs. ", err)
	}
}

func TestMove(t *testing.T) {
	//create full path to the root folder
	root, _ := filepath.Abs("./testdata/root")
	// initiate the loal storage
	l := New(root)

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
	if err != nil {
		t.Error("failed assert move file back ", err)
	}
	// delete the sub dirs
	err = os.RemoveAll("/sub1")
	if err != nil {
		t.Error("failed assert remove sub1 ", err)
	}
}

func TestMoveAs(t *testing.T) {
	//create full path to the root folder
	root, _ := filepath.Abs("./testdata/root")
	// initiate the loal storage
	l := New(root)

	// 1- test moving to same dir with the same name
	err := l.MoveAs("filetomove.md", "/", "filetomove.md")
	if err == nil {
		t.Error("error when moving file to same dir with the same name")
	}

	// 2- test moving to same dir with different name
	err = l.MoveAs("filetomove.md", "/", "filetomovenew.md")
	if err != nil {
		t.Error("failed asserting error when moving to same dir with different name")
	}
	// assert the source file
	_, err = os.Stat(path.Join(root, "filetomove.md"))
	if err == nil {
		t.Error("failed assert moving to the same dir: source file still exist")
	}

	//assert the dest file
	_, err = os.Stat(path.Join(root, "filetomovenew.md"))
	if err != nil {
		t.Error("failed assert moving to the same dir: new file not exist")
	}

	// rename file back to original name
	err = l.MoveAs("filetomovenew.md", "/", "filetomove.md")
	if err != nil {
		t.Error("failed assert rename file back to original name")
	}
	// 3- test moving to sub folder
	err = l.MoveAs("filetomove.md", "/sub1/sub2", "filetomove.md")
	if err != nil {
		t.Error("failed assert  moving to sub folder")
	}
	// assert the source file
	_, err = os.Stat(path.Join(root, "filetomove.md"))
	if err == nil {
		t.Error("failed assert moving to the same dir: source file still exist")
	}
	// assert the dest file
	_, err = os.Stat(path.Join(root, "sub1/sub2/filetomove.md"))
	if err != nil {
		t.Error("failed assert moving to the same dir: dest file not exist")
	}

	// move the file back to original dir
	err = l.MoveAs("sub1/sub2/filetomove.md", "/", "filetomove.md")
	if err != nil {
		t.Error("failed move the file back to original dir")
	}
	// remove the sub dirs
	err = os.RemoveAll(path.Join(root, "sub1"))
	if err != nil {
		t.Error("failed remove the sub dirs")
	}
}

func TestRename(t *testing.T) {
	//create full path to the root folder
	root, _ := filepath.Abs("./testdata/root")
	// initiate the loal storage
	l := New(root)

	err := l.Rename("filetorename.md", "filetorename-newname.md")
	if err != nil {
		t.Error("failed assert renaming file")
	}
	_, err = os.Stat(path.Join(root, "filetorename-newname.md"))
	if err != nil {
		t.Error("failed assert renaming file")
	}

	// rename back the file
	err = l.Rename("filetorename-newname.md", "filetorename.md")
	if err != nil {
		t.Error("failed rename back the  file")
	}
}

func TestDelete(t *testing.T) {
	//create full path to the root folder
	root, _ := filepath.Abs("./testdata/root")
	// initiate the loal storage
	l := New(root)

	l.Delete("filetodelete.md")
	_, err := os.Stat("filetodelete.md")
	if err == nil {
		t.Error("failed asserting file delete")
	}

	// create the file
	l.Create("filetodelete.md", []byte("this is a test file"))
}

func TestDeleteMultiple(t *testing.T) {
	//create full path to the root folder
	root, _ := filepath.Abs("./testdata/root")
	// initiate the loal storage
	l := New(root)
	l.DeleteMultiple([]string{"filetodelete1.md", "filetodelete2"})
	_, err := os.Stat("filetodelete1.md")
	if err == nil {
		t.Error("failed asserting delete multiple")
	}

	_, err = os.Stat("filetodelete2.md")
	if err == nil {
		t.Error("failed asserting delete multiple")
	}

	// create the files
	file, _ := os.Create(path.Join(root, "filetodelete1.md"))
	defer file.Close()
	file.Write([]byte("this is a test file"))

	file, _ = os.Create(path.Join(root, "filetodelete2.md"))
	defer file.Close()
	file.Write([]byte("this is a test file"))
}

func TestCreate(t *testing.T) {
	//create full path to the root folder
	root, _ := filepath.Abs("./testdata/root")
	// initiate the loal storage
	l := New(root)

	err := l.Create("filetocreate.md", []byte("this is a test file"))
	if err != nil {
		t.Error(err)
	}

	_, err = os.Stat(path.Join(root, "filetocreate.md"))
	if err != nil {
		t.Error("failed assert file create: ", err)
	}

	os.Remove(path.Join(root, "filetocreate.md"))
}

func TestAppend(t *testing.T) {
	//create full path to the root folder
	root, _ := filepath.Abs("./testdata/root")
	// initiate the loal storage
	l := New(root)
	l.Append("filetoappend.md", []byte("appended"))

	fileContent, err := ioutil.ReadFile(path.Join(root, "filetoappend.md"))
	if err != nil {
		t.Error("failed asserting append: ", err)
	}
	yes := strings.Contains(string(fileContent), "appended")
	if !yes {
		t.Error("failed assert append")
	}

	// fix the files
	os.Remove(path.Join(root, "filetoappend.md"))
	file, _ := os.Create(path.Join(root, "filetoappend.md"))
	file.Write([]byte("this is a test file\n"))
	file.Close()
}

func TestExists(t *testing.T) {
	//create full path to the root folder
	root, _ := filepath.Abs("./testdata/root")
	// initiate the loal storage
	l := New(root)

	yes, _ := l.Exists("filetocheckexist.md")
	if !yes {
		t.Error("failed assert exist")
	}

}

func TestMissing(t *testing.T) {
	//create full path to the root folder
	root, _ := filepath.Abs("./testdata/root")
	// initiate the loal storage
	l := New(root)

	yes, _ := l.Missing("missingfile.md")
	if !yes {
		t.Error("failed assert missing")
	}

}

func TestRead(t *testing.T) {
	//create full path to the root folder
	root, _ := filepath.Abs("./testdata/root")
	// initiate the loal storage
	l := New(root)

	content, err := l.Read("filetoread.md")
	if err != nil {
		t.Error("failed assert reading file")
	}

	if string(content) != "contentToRead" {
		t.Error("failed assert reading file")
	}

}

func TestFiles(t *testing.T) {
	//create full path to the root folder
	root, _ := filepath.Abs("./testdata/root")
	// initiate the loal storage
	l := New(root)
	files, err := l.Files("files")
	if err != nil {
		t.Error("failed asserting list files")
	}
	count := len(files)
	if count != 2 {
		t.Error("failed asserting list files")
	}
}

func TestAllFiles(t *testing.T) {
	//create full path to the root folder
	root, _ := filepath.Abs("./testdata/root")
	// initiate the loal storage
	l := New(root)
	files, err := l.AllFiles("files")
	if err != nil {
		t.Error("failed asserting list all files")
	}
	count := len(files)
	if count != 3 {
		t.Error("failed asserting list all files")
	}
}

func TestDirectories(t *testing.T) {
	//create full path to the root folder
	root, _ := filepath.Abs("./testdata/root")
	// initiate the loal storage
	l := New(root)
	dirs, err := l.Directories("dirs")
	if err != nil {
		t.Error("failed asserting list dirs")
	}
	count := len(dirs)
	if count != 2 {
		t.Error("failed asserting list dirs")
	}
}

func TestAllDirectories(t *testing.T) {
	//create full path to the root folder
	root, _ := filepath.Abs("./testdata/root")
	// initiate the loal storage
	l := New(root)
	dirs, err := l.AllDirectories("dirs")

	if err != nil {
		t.Error("failed asserting list dirs")
	}
	count := len(dirs)
	if count != 3 {
		t.Error("failed asserting list dirs")
	}
}

func TestMakeDirectory(t *testing.T) {
	//create full path to the root folder
	root, _ := filepath.Abs("./testdata/root")
	// initiate the loal storage
	l := New(root)
	err := l.MakeDirectory("dirToMake", 0777)
	if err != nil {
		t.Error("failed asserting make directory")
	}

	dirPath := path.Join(root, "dirToMake")
	s, _ := os.Stat(dirPath)

	if !s.IsDir() {
		t.Error("failed asserting make directory")
	}

	err = os.RemoveAll(dirPath)
	if err != nil {
		log.Fatal(err)
	}
}

func TestRenameDirectory(t *testing.T) {
	//create full path to the root folder
	root, _ := filepath.Abs("./testdata/root")
	// initiate the loal storage
	l := New(root)
	err := l.RenameDirectory("dirtorename", "dirtorenamenew")
	if err != nil {
		t.Error("failed asserting rename directory")
	}

	dirPath := path.Join(root, "dirtorenamenew")
	s, _ := os.Stat(dirPath)

	if s != nil && !s.IsDir() {
		t.Error("failed asserting rename directory")
	}

	os.RemoveAll(dirPath)
	l.MakeDirectory("dirtorename", 0777)
	l.Create("dirtorename/.gitkeep", []byte(""))
}

func TestDeleteDirectory(t *testing.T) {
	//create full path to the root folder
	root, _ := filepath.Abs("./testdata/root")
	// initiate the loal storage
	l := New(root)
	err := l.DeleteDirectory("dirtodelete")
	if err != nil {
		t.Error("failed asserting delete directory")
	}

	dirPath := path.Join(root, "dirtodelete")
	s, _ := os.Stat(dirPath)
	if s != nil {
		t.Error("failed asserting delete directory")
	}

	l.MakeDirectory(path.Join(root, "dirtodelete"), 0777)
	l.Create(path.Join(root, "dirtodelete/.gitkeep"), []byte(""))
}
