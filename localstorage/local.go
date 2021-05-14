// Copyright 2021 Harran Ali <harran.m@gmail.com>. All rights reserved.
// Use of this source code is governed by MIT-style
// license that can be found in the LICENSE file.

package localstorage

import (
	"errors"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
	"unicode/utf8"
)

// LocalStorage local storage
type LocalStorage struct {
	rootFolder string
}

// FileInfo provides file information
type FileInfo struct {
	Name                 string
	NameWithoutExtension string
	LastModified         time.Time
	Size                 int64
	Extension            string
	Path                 string
	IsDirectory          bool
	FsFileInfo           fs.FileInfo //golang's fs file info
}

var local *LocalStorage

// New initiate local storage
func New(path string) *LocalStorage {
	local = &LocalStorage{
		rootFolder: path,
	}

	return local
}

//FileInfo returns information about the given file or an error incase there is any
func (l *LocalStorage) FileInfo(filepath string) (fileinfo FileInfo, err error) {
	fullpath := path.Join(l.rootFolder, filepath)
	// make sure the file exists
	if _, err := os.Stat(fullpath); os.IsNotExist(err) {
		return FileInfo{}, errors.New("file does not exist")
	} else if err != nil {
		return FileInfo{}, err
	}

	info, err := os.Stat(fullpath)

	fileinfo = FileInfo{
		Name:                 info.Name(),
		Extension:            removeFirstChar(path.Ext(fullpath)),
		NameWithoutExtension: removeExtension(info.Name(), path.Ext(fullpath)),
		Size:                 info.Size(),
		Path:                 path.Dir(fullpath),
		LastModified:         info.ModTime(),
		IsDirectory:          info.IsDir(),
		FsFileInfo:           info,
	}

	return fileinfo, nil
}

// Put helps you copy files into the root directory
// from external locations, filePath is the full path to the file
// you would like to put, it returns error incase there is any
func (l *LocalStorage) Put(filePath string) error {
	// make sure the source file exists
	s, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return err
	}
	if !s.Mode().IsRegular() {
		return errors.New("File is not in regular mode")
	}

	// make sure there is no file with the same name in destFile
	_, err = os.Stat(path.Join(l.rootFolder, s.Name()))
	if err == nil {
		return errors.New("the file is already exists")
	}

	// open the source file
	srcFile, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// create a file in destination with the same name of source fil
	destFile, err := os.Create(path.Join(l.rootFolder, s.Name()))
	if err != nil {
		return err
	}
	defer destFile.Close()

	// copy the file content
	buf := make([]byte, 100)
	for {
		n, err := srcFile.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}
		if _, err := destFile.Write(buf[:n]); err != nil {
			return err
		}
	}
	return err
}

// PutAs helps you copy files into the root directory
// from external directory,  the first param 'filePath' is the full
// path to the file you would like to put,
// the second param 'fileName' is the name you would
// like to give to the file, it returns error incase there is any
func (l *LocalStorage) PutAs(filePath string, filename string) error {
	// make sure the source file exists
	s, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return err
	}
	if !s.Mode().IsRegular() {
		return errors.New("File is not in regular mode")
	}

	// make sure there is no file with the same name in destFile
	_, err = os.Stat(path.Join(l.rootFolder, filename))
	if err == nil {
		return errors.New("the file is already exists")
	}

	// open the source file
	srcFile, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// create a file in destination with the same name of source fil
	destFile, err := os.Create(path.Join(l.rootFolder, filename))
	if err != nil {
		return err
	}
	defer destFile.Close()

	// copy the file content
	buf := make([]byte, 100)
	for {
		n, err := srcFile.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}
		if _, err := destFile.Write(buf[:n]); err != nil {
			return err
		}
	}
	return err
}

// Copy helps you copy files within the root folder,
// Please note that the reference of the paths of these
// files is the root folder, it accepts the source
// file starting from the root folder,
// and the destination folder starting from the root folder,
// it returns an error incase there is any
func (l *LocalStorage) Copy(filePath string, destPath string) error {
	//unify slashes
	filePath = filepath.ToSlash(filePath)
	destPath = filepath.ToSlash(destPath)

	// construct the full paths
	srcFileFullPath := filepath.Join(l.rootFolder, filePath)
	DestFullPath := filepath.Join(l.rootFolder, destPath)

	// make sure the path of dest folder exists
	os.MkdirAll(DestFullPath, 0755)

	// make sure the source file exists
	s, err := os.Stat(srcFileFullPath)
	if os.IsNotExist(err) {
		return err
	}
	if !s.Mode().IsRegular() {
		return errors.New("File is not in regular mode")
	}

	// make sure there is no file with the same name in dest folder
	_, err = os.Stat(path.Join(DestFullPath, s.Name()))
	if err == nil {
		return errors.New("the file already exists in dest folder")
	}

	// open the source file
	srcFile, err := os.Open(srcFileFullPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// create a file in destination with the same name of source fil
	destFile, err := os.Create(path.Join(DestFullPath, s.Name()))
	if err != nil {
		return err
	}
	defer destFile.Close()

	// copy the file content
	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return err
	}

	return nil
}

// CopyAs helps you copy files within the root folder,
// Please note that the reference of the paths of these
// files is the root folder, it accepts the source file starting
// from the root folder and the destination folder
// starting from the root folder, and the new file name,
// it returns an error incase there is any
func (l *LocalStorage) CopyAs(filePath string, destfolder string, newFilePath string) error {
	//unify slashes
	filePath = filepath.ToSlash(filePath)
	destfolder = filepath.ToSlash(destfolder)

	// construct the full paths
	srcFileFullPath := filepath.Join(l.rootFolder, filePath)
	DestFileFullPath := filepath.Join(l.rootFolder, destfolder)

	// make sure the path of dest folder exists
	os.MkdirAll(DestFileFullPath, 0755)

	// make sure the source file exists
	s, err := os.Stat(srcFileFullPath)
	if os.IsNotExist(err) {
		return err
	}
	if !s.Mode().IsRegular() {
		return errors.New("File is not in regular mode")
	}

	// make sure there is no file with the same name in dest folder
	_, err = os.Stat(path.Join(DestFileFullPath, newFilePath))
	if err == nil {
		return errors.New("the file is already exists")
	}

	// open the source file
	srcFile, err := os.Open(srcFileFullPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// create a file in destination with the same name of source fil
	destFile, err := os.Create(path.Join(DestFileFullPath, newFilePath))
	if err != nil {
		return err
	}
	defer destFile.Close()

	// copy the file content
	buf := make([]byte, 100)
	for {
		n, err := srcFile.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}
		if _, err := destFile.Write(buf[:n]); err != nil {
			return err
		}
	}
	return err
}

// Move helps you Move files within the root folder,
// Please note that the  of the paths of these
// files is the root folder, it accepts the source file
// starting from the root folder, and the destination
// folder starting from the root folder,
// it returns an error incase there any
func (l *LocalStorage) Move(filePath string, destFolder string) error {
	//unify slashes
	filePath = filepath.ToSlash(filePath)
	destFolder = filepath.ToSlash(destFolder)

	// construct the full paths
	srcFileFullPath := filepath.Join(l.rootFolder, filePath)
	DestFileFullPath := filepath.Join(l.rootFolder, destFolder)

	// make sure the path of dest folder exists
	os.MkdirAll(DestFileFullPath, 0755)

	// make sure the source file exists
	s, err := os.Stat(srcFileFullPath)
	if os.IsNotExist(err) {
		return err
	}
	if !s.Mode().IsRegular() {
		return errors.New("File is not in regular mode")
	}

	// make sure there is no file with the same name in dest folder
	_, err = os.Stat(path.Join(DestFileFullPath, s.Name()))
	if err == nil {
		return errors.New("the file is already exists")
	}

	// open the source file
	srcFile, err := os.Open(srcFileFullPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// create a file in destination with the same name of source fil
	destFile, err := os.Create(path.Join(DestFileFullPath, s.Name()))
	if err != nil {
		return err
	}
	defer destFile.Close()

	// copy the file content
	buf := make([]byte, 100)
	for {
		n, err := srcFile.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}
		if _, err := destFile.Write(buf[:n]); err != nil {
			return err
		}
	}

	// remove the source file
	srcFile.Close()
	os.Remove(srcFileFullPath)
	return err
}

// MoveAs helps you Move files within the root folder,
// Please note that the reference of the paths of these files
// is the root folder, it accepts the source
// file starting from the root folder and the destination
// folder starting from the root folder,
// and the new file name, it returns an error incase there any
func (l *LocalStorage) MoveAs(filePath string, destFolder string, newFilePath string) error {
	//unify slashes
	filePath = filepath.ToSlash(filePath)
	destFolder = filepath.ToSlash(destFolder)

	// construct the full paths
	srcFileFullPath := filepath.Join(l.rootFolder, filePath)
	DestFileFullPath := filepath.Join(l.rootFolder, destFolder)

	// make sure the path of dest folder exists
	os.MkdirAll(DestFileFullPath, 0755)

	// make sure the source file exists
	s, err := os.Stat(srcFileFullPath)
	if os.IsNotExist(err) {
		return err
	}
	if !s.Mode().IsRegular() {
		return errors.New("File is not in regular mode")
	}

	// make sure there is no file with the same name in dest folder
	_, err = os.Stat(path.Join(DestFileFullPath, newFilePath))
	if !os.IsNotExist(err) {
		return errors.New("the file is already exists")
	}

	// open the source file
	srcFile, err := os.Open(srcFileFullPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// create a file in destination with the same name of source fil
	destFile, err := os.Create(path.Join(DestFileFullPath, newFilePath))
	if err != nil {
		return err
	}
	defer destFile.Close()

	// copy the file content
	buf := make([]byte, 100)
	for {
		n, err := srcFile.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}
		if _, err := destFile.Write(buf[:n]); err != nil {
			return err
		}
	}

	// remove the source file
	srcFile.Close()
	os.Remove(srcFileFullPath)
	return err
}

// Rename renames the given file as first parameter to the name
// given as a second parameter,
// it returns error incase there is any
func (l *LocalStorage) Rename(filePath string, newFilePath string) error {
	srcFileFullPath := filepath.Join(l.rootFolder, filePath)
	destFileFullPath := filepath.Join(l.rootFolder, newFilePath)

	// make sure the source file exists
	s, err := os.Stat(srcFileFullPath)
	if os.IsNotExist(err) {
		return err
	}
	if !s.Mode().IsRegular() {
		return errors.New("File is not in regular mode")
	}

	err = os.Rename(srcFileFullPath, destFileFullPath)

	return err
}

// Delete deletes the given file it returns error incase there is any
func (l *LocalStorage) Delete(filePath string) error {
	srcFileFullPath := filepath.Join(l.rootFolder, filePath)

	// make sure the source file exists
	s, err := os.Stat(srcFileFullPath)
	if os.IsNotExist(err) {
		return err
	}
	if !s.Mode().IsRegular() {
		return errors.New("File is not in regular mode")
	}

	os.Remove(srcFileFullPath)

	return err
}

// DeleteMultiple deltes multiple files given as slice of strings
// of file paths, it returns error incase there is any
func (l *LocalStorage) DeleteMultiple(filePaths []string) (err error) {
	for _, file := range filePaths {
		srcFileFullPath := filepath.Join(l.rootFolder, file)

		// make sure the source file exists
		s, err := os.Stat(srcFileFullPath)
		if os.IsNotExist(err) {
			continue
		}
		if !s.Mode().IsRegular() {
			continue
		}

		os.Remove(srcFileFullPath)
	}

	return err
}

// Create helps you create new a file and add content to it,
// it returns error incase there is any
func (l *LocalStorage) Create(filePath string, content []byte) error {
	// make sure the path of dest folder exists
	fileFullPath := path.Join(l.rootFolder, filePath)
	fileFullPath = filepath.ToSlash(fileFullPath)
	os.MkdirAll(path.Dir(fileFullPath), 0755)

	// check if the file exists
	_, err := os.Stat(fileFullPath)
	if err == nil {
		return errors.New("file already exists")
	}

	// create the file
	file, err := os.Create(fileFullPath)
	if err != nil {
		return err
	}
	// add the content
	_, err = file.Write(content)
	file.Close()

	return err
}

// Append helps you append content to a file,
// it returns error incase there is any
func (l *LocalStorage) Append(filePath string, content []byte) error {
	fileFullPath := path.Join(l.rootFolder, filePath)

	// check if the file exists
	_, err := os.Stat(fileFullPath)
	if os.IsNotExist(err) {
		return err
	}

	// create the file
	file, err := os.OpenFile(fileFullPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// add the content
	_, err = file.Write(content)

	return err
}

// Exists checks if a file exists withn the root folder,
// it returns a bool and an error incase any
func (l *LocalStorage) Exists(filePath string) (bool, error) {
	fileFullPath := path.Join(l.rootFolder, filePath)

	_, err := os.Stat(fileFullPath)
	if os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}

// Missing checks if a file is missing in the root folder,
// it returns a bool and an error incase any
func (l *LocalStorage) Missing(filePath string) (bool, error) {
	fileFullPath := path.Join(l.rootFolder, filePath)

	_, err := os.Stat(fileFullPath)
	if err != nil {
		// there is an error
		if os.IsNotExist(err) {
			// not exist error
			return true, nil
		}
		// another errors
		return false, err
	}

	return false, nil
}

// Read helps you grap the content of a file,
// it returns the data in a slice of bytes and an error
// incase there is any
func (l *LocalStorage) Read(filePath string) ([]byte, error) {
	fileFullPath := path.Join(l.rootFolder, filePath)

	_, err := os.Stat(fileFullPath)
	if os.IsNotExist(err) {
		// not exist error
		return nil, err
	}

	content, err := ioutil.ReadFile(fileFullPath)
	if err != nil {
		return nil, err
	}

	return content, nil
}

// Files returns a list of files in a given directory,
// the file type is LocalStorage.FileInfo
//  NOT the standard library fs.FileInfo,
// and it returns an error incase any occurred,
// if you want a list of files including
// the files in sub directories, consider using the method
// `AllFiles(DirectoryPath string)`
func (l *LocalStorage) Files(DirectoryPath string) (files []FileInfo, err error) {
	DirectoryFullPath := path.Join(l.rootFolder, DirectoryPath)

	_, err = os.Stat(DirectoryFullPath)
	if os.IsNotExist(err) {
		// not exist error
		return []FileInfo{}, err
	}

	res, err := ioutil.ReadDir(DirectoryFullPath)
	for _, val := range res {
		if !val.IsDir() {
			// assign the result var
			p := path.Join(DirectoryPath, val.Name())
			f, _ := l.FileInfo(p)
			files = append(files, f)
		}
	}

	return files, err
}

// AllFiles returns a list of files in the given directory
// including files in sub directories,
// the file type in the list is LocalStorage.FileInfo
// NOT the standard library fs.FileInfo
func (l *LocalStorage) AllFiles(DirectoryPath string) (files []FileInfo, err error) {
	DirectoryFullPath := path.Join(l.rootFolder, DirectoryPath)

	_, err = os.Stat(DirectoryFullPath)
	if os.IsNotExist(err) {
		// not exist error
		return []FileInfo{}, err
	}

	err = filepath.Walk(DirectoryFullPath, func(filePath string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			// trim root path
			trimedPath := strings.ReplaceAll(filePath, l.rootFolder, "")
			f, _ := l.FileInfo(trimedPath)
			files = append(files, f)
		}
		return nil
	})
	if err != nil {
		return []FileInfo{}, err
	}

	return files, err
}

// Directories returns a slice of string containing
// the paths of the directories,
// if you want the list of directories including subdirectories,
// consider using the method "AllDirectories(DirectoryPath string)",
// it returns an error incase is any
func (l *LocalStorage) Directories(DirectoryPath string) (SubDirectoryPaths []string, err error) {
	DirectoryFullPath := path.Join(l.rootFolder, DirectoryPath)

	_, err = os.Stat(DirectoryFullPath)
	if os.IsNotExist(err) {
		// not exist error
		return []string{}, err
	}

	res, err := ioutil.ReadDir(DirectoryFullPath)
	for _, val := range res {
		if val.IsDir() {
			// assign the result var
			p := path.Join(l.rootFolder, DirectoryPath, val.Name())
			p = filepath.ToSlash(p)
			SubDirectoryPaths = append(SubDirectoryPaths, p)
		}
	}

	return SubDirectoryPaths, err
}

// AllDirectories returns a list of directories including
// sub directories, it returns an error incase is any
func (l *LocalStorage) AllDirectories(SubDirectoryPath string) (directoryPaths []string, err error) {
	DirectoryFullPath := path.Join(l.rootFolder, SubDirectoryPath)

	_, err = os.Stat(DirectoryFullPath)
	if os.IsNotExist(err) {
		// not exist error
		return []string{}, err
	}

	err = filepath.Walk(DirectoryFullPath, func(filePath string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && filePath != DirectoryFullPath {
			filePath = filepath.ToSlash(filePath)
			directoryPaths = append(directoryPaths, filePath)
		}
		return nil
	})
	if err != nil {
		return []string{}, err
	}

	return directoryPaths, err
}

// MakeDirectory creates a new directory and the necessary
// parent directories with the given permissions,
// permissions could be (example: 0777) or any linux based permissions,
// it returns an error incase is any
func (l *LocalStorage) MakeDirectory(DirectoryPath string, perm int) (err error) {
	DirectoryFullPath := path.Join(l.rootFolder, DirectoryPath)
	err = os.MkdirAll(DirectoryFullPath, fs.FileMode(perm))
	return err
}

// RenameDirectory changes the name of directory to new name,
// it returns an error incase there is any
func (l *LocalStorage) RenameDirectory(DirectoryPath string, NewDirectoryPath string) (err error) {
	DirectoryFullPath := path.Join(l.rootFolder, DirectoryPath)
	NewDirectoryFullPath := path.Join(l.rootFolder, NewDirectoryPath)
	err = os.Rename(DirectoryFullPath, NewDirectoryFullPath)

	return err
}

// DeleteDirectory deletes the given directory
func (l *LocalStorage) DeleteDirectory(DirectoryPath string) (err error) {
	DirectoryFullPath := path.Join(l.rootFolder, DirectoryPath)
	err = os.RemoveAll(DirectoryFullPath)

	return err
}

func removeFirstChar(s string) string {
	_, i := utf8.DecodeRuneInString(s)
	return s[i:]
}

func removeExtension(fileName string, ext string) string {
	return strings.TrimSuffix(fileName, ext)
}
