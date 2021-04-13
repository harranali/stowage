package localstorage

import (
	"io/fs"
	"os"
	"path"
	"strings"
	"unicode/utf8"
)

// Local local storage
type Local struct {
	rootFolder string
}

// FileInfo provides file information
type FileInfo struct {
	Name                 string
	NameWithoutExtension string
	Size                 int64
	Extension            string
	Path                 string
	FsFileInfo           fs.FileInfo
}

var local *Local

// New initiate local storage
func New(path string) *Local {
	local = &Local{
		rootFolder: path,
	}

	return local
}

// File returns information about the given file or an error when there is any
func (l *Local) FileInfo(filepath string) (fileinfo FileInfo, err error) {
	fullpath := path.Join(l.rootFolder, filepath)

	info, err := os.Stat(fullpath)
	fileinfo = FileInfo{
		Name:                 info.Name(),
		Extension:            removeFirstChar(path.Ext(fullpath)),
		NameWithoutExtension: removeExtension(info.Name(), path.Ext(fullpath)),
		Size:                 info.Size(),
		Path:                 path.Dir(fullpath),
		FsFileInfo:           info,
	}

	return
}

func removeFirstChar(s string) string {
	_, i := utf8.DecodeRuneInString(s)
	return s[i:]
}

func removeExtension(fileName string, ext string) string {
	return strings.TrimSuffix(fileName, ext)
}
