// Copyright 2021 Harran Ali <harran.m@gmail.com>. All rights reserved.
// Use of this source code is governed by MIT-style
// license that can be found in the LICENSE file.

package stowage

import (
	"github.com/harranali/stowage/localstorage"
)

// LocalStorageOpts options for initiating local storage
type LocalStorageOpts struct {
	RootFolder string
}

// Disk interface defines all supported operations by local storage
type Disk interface {
	FileInfo(filePath string) (fileinfo localstorage.FileInfo, err error)
	Put(filePath string) error
	PutAs(filePath string, filename string) error
	Copy(filePath string, destfolder string) error
	CopyAs(filePath string, destfolder string, newFilePath string) error
	Move(filePath string, destfolder string) error
	MoveAs(filePath string, destFolder string, newFilePath string) error
	Rename(filePath string, newFilePath string) error
	Delete(filePath string) error
	DeleteMultiple(filePaths []string) error
	Create(filePath string, content []byte) error
	Append(filePath string, content []byte) error
	Exists(filePath string) (bool, error)
	Missing(filePath string) (bool, error)
	Read(filePath string) ([]byte, error)
	Files(DirectoryPath string) ([]localstorage.FileInfo, error)
	AllFiles(DirectoryPath string) ([]localstorage.FileInfo, error)
	Directories(SubDirectoryPath string) (directoryPaths []string, err error)
	AllDirectories(SubDirectoryPath string) (directoryPaths []string, err error)
	MakeDirectory(DirectoryPath string, perm int) error
	RenameDirectory(DirectoryPath string, NewDirectoryPath string) (err error)
	DeleteDirectory(DirectoryPath string) (err error)
}

// Stowage represents all supported storages
type Stowage struct {
	LocalStorage Disk
}

var stowage *Stowage

// New initialize stowage
func New() *Stowage {
	stowage = &Stowage{}

	return stowage
}

// InitLocalStorage initializes local storage
func (s *Stowage) InitLocalStorage(opts LocalStorageOpts) {
	s.LocalStorage = localstorage.New(opts.RootFolder)
}
