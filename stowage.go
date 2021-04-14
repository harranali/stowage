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
	FileInfo(filepath string) (fileinfo localstorage.FileInfo, err error)
	Put(filepath string) error
	PutAs(filepath string, filename string) error
	Copy(srcfile string, destfolder string) error
	CopyAs(srcfile string, destfolder string, newfilename string) error
	Move(srcfile string, destfile string) error
	MoveAs(srcfile string, destfile string, newfilename string) error
	Rename(filename string, newfilename string) error
	Delete(filepath string) error
	DeleteMultiple(filepaths []string) error
	Create(filepath string, content []byte) error
	Append(filepath string, content []byte) error
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
