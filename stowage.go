// Copyright 2021 Harran Ali <harran.m@gmail.com>. All rights reserved.
// Use of this source code is governed by MIT-style
// license that can be found in the LICENSE file.

package stowage

import (
	"github.com/harranali/stowage/localstorage"
)

type LocalStorageOpts struct {
	RootFolder string
}

type S3Opts struct {
	token string
}

type GoogleCloudStorageOpts struct {
	token string
}

type OSSOpts struct {
	Token string
}

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
}

type Stowage struct {
	LocalStorage       Disk
	S3                 Disk
	GoogleCloudStorage Disk
	OSS                Disk
}

var stowage *Stowage

func New() *Stowage {
	stowage = &Stowage{}

	return stowage
}

// InitLocalStorage initializes local storage
func (s *Stowage) InitLocalStorage(opts LocalStorageOpts) {
	s.LocalStorage = localstorage.New(opts.RootFolder)
}

// InitS3 initializes Amazon S3 storage
func (s *Stowage) InitS3(opts S3Opts) {
}

// InitGoogleCloudStorage initializes google cloud storage
func (s *Stowage) InitGoogleCloudStorage(opts GoogleCloudStorageOpts) {
}

// InitOSS initializes Alicloud OSS
func (s *Stowage) InitOSS(opts OSSOpts) {
}
