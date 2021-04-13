package localstorage

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"path"
	"strings"
	"unicode/utf8"
)

// Local local storage
type LocalStorage struct {
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

var local *LocalStorage

// New initiate local storage
func New(path string) *LocalStorage {
	local = &LocalStorage{
		rootFolder: path,
	}

	return local
}

// FileInfo returns information about the given file or an error incase there is
func (l *LocalStorage) FileInfo(filepath string) (fileinfo FileInfo, err error) {
	// make sure the file exists
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return FileInfo{}, err
	}

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

	return fileinfo, nil
}

// Put files in the root directory,
// filepath is the full path to the file you would like to put
// it returns error incase there was
func (l *LocalStorage) Put(filepath string) error {
	// make sure the source file exists
	s, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		return err
	}
	if !s.Mode().IsRegular() {
		return errors.New("File is not in regular mode")
	}

	// make sure there is no file with the same name in destFile
	_, err = os.Stat(path.Join(l.rootFolder, s.Name()))
	if !os.IsNotExist(err) {
		return errors.New("the file is already exists")
	}

	// open the source file
	srcFile, err := os.Open(filepath)
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

// PutAs  puts files in the root directory with the given name
// the first arg 'filepath' is the file path
// the second arg 'filename' is the name of the new file
// filepath is the full path to the file you would like to put
// it returns error incase there was
func (l *LocalStorage) PutAs(filepath string, filename string) error {
	// make sure the source file exists
	s, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		return err
	}
	if !s.Mode().IsRegular() {
		return errors.New("File is not in regular mode")
	}

	// make sure there is no file with the same name in destFile
	_, err = os.Stat(path.Join(l.rootFolder, filename))
	if !os.IsNotExist(err) {
		return errors.New("the file is already exists")
	}

	// open the source file
	srcFile, err := os.Open(filepath)
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

func removeFirstChar(s string) string {
	_, i := utf8.DecodeRuneInString(s)
	return s[i:]
}

func removeExtension(fileName string, ext string) string {
	return strings.TrimSuffix(fileName, ext)
}
