package localstorage

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
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

// Put files in the root directory from external locations,
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

// PutAs  puts files in the root directory from external directory with the given name in the second parameter
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

// Copy helps you copy files within the root folder
// Please note that the refrence of the paths of these files is the root folder
// it accepts the source file starting from the root folder
// and the destination folder starting from the root folder
// it returns an error incase there any
func (l *LocalStorage) Copy(srcfile string, destfolder string) error {
	//unify slashes
	srcfile = filepath.ToSlash(srcfile)
	destfolder = filepath.ToSlash(destfolder)

	// construct the full paths
	srcFileFullPath := filepath.Join(l.rootFolder, srcfile)
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
	_, err = os.Stat(path.Join(DestFileFullPath, s.Name()))
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
	return err
}

// Copy helps you copy files within the root folder
// Please note that the refrence of the paths of these files is the root folder
// it accepts the source file starting from the root folder
// and the destination folder starting from the root folder
// and the new file name
// it returns an error incase there any
func (l *LocalStorage) CopyAs(srcfile string, destfolder string, newfilename string) error {
	//unify slashes
	srcfile = filepath.ToSlash(srcfile)
	destfolder = filepath.ToSlash(destfolder)

	// construct the full paths
	srcFileFullPath := filepath.Join(l.rootFolder, srcfile)
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
	_, err = os.Stat(path.Join(DestFileFullPath, newfilename))
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
	destFile, err := os.Create(path.Join(DestFileFullPath, newfilename))
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

// Move helps you Move files within the root folder
// Please note that the refrence of the paths of these files is the root folder
// it accepts the source file starting from the root folder
// and the destination folder starting from the root folder
// it returns an error incase there any
func (l *LocalStorage) Move(srcfile string, destfolder string) error {
	//unify slashes
	srcfile = filepath.ToSlash(srcfile)
	destfolder = filepath.ToSlash(destfolder)

	// construct the full paths
	srcFileFullPath := filepath.Join(l.rootFolder, srcfile)
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
	_, err = os.Stat(path.Join(DestFileFullPath, s.Name()))
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
	err = os.Remove(srcFileFullPath)
	return err
}

// Move helps you Move files within the root folder
// Please note that the refrence of the paths of these files is the root folder
// it accepts the source file starting from the root folder
// and the destination folder starting from the root folder
// and the new file name
// it returns an error incase there any
func (l *LocalStorage) MoveAs(srcfile string, destfolder string, newfilename string) error {
	//unify slashes
	srcfile = filepath.ToSlash(srcfile)
	destfolder = filepath.ToSlash(destfolder)

	// construct the full paths
	srcFileFullPath := filepath.Join(l.rootFolder, srcfile)
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
	_, err = os.Stat(path.Join(DestFileFullPath, newfilename))
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
	destFile, err := os.Create(path.Join(DestFileFullPath, newfilename))
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
	err = os.Remove(srcFileFullPath)
	return err
}

func removeFirstChar(s string) string {
	_, i := utf8.DecodeRuneInString(s)
	return s[i:]
}

func removeExtension(fileName string, ext string) string {
	return strings.TrimSuffix(fileName, ext)
}
