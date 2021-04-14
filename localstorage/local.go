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

// Rename renames a given file as first parameter to a new name passed as a second parameter
// it returns error incase there is any
func (l *LocalStorage) Rename(filename string, newfilename string) error {
	srcFileFullPath := filepath.Join(l.rootFolder, filename)
	destFileFullPath := filepath.Join(l.rootFolder, newfilename)

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

// Delete removes the given file
// it returns error incase there is any
func (l *LocalStorage) Delete(filename string) error {
	srcFileFullPath := filepath.Join(l.rootFolder, filename)

	// make sure the source file exists
	s, err := os.Stat(srcFileFullPath)
	if os.IsNotExist(err) {
		return err
	}
	if !s.Mode().IsRegular() {
		return errors.New("File is not in regular mode")
	}

	err = os.Remove(srcFileFullPath)

	return err
}

// Delete removes multiple files given as slice of strings of file paths
// it returns error incase there is any
func (l *LocalStorage) DeleteMultiple(filepaths []string) (err error) {
	for _, file := range filepaths {
		srcFileFullPath := filepath.Join(l.rootFolder, file)

		// make sure the source file exists
		s, err := os.Stat(srcFileFullPath)
		if os.IsNotExist(err) {
			continue
		}
		if !s.Mode().IsRegular() {
			continue
		}

		err = os.Remove(srcFileFullPath)
	}

	return err
}

// Create helps you create new file and add content to it
// it returns error incase there is any
func (l *LocalStorage) Create(filepath string, content []byte) error {
	// make sure the path of dest folder exists
	os.MkdirAll(path.Dir(filepath), 0755)

	fileFullPath := path.Join(l.rootFolder, filepath)

	// check if the file exists
	_, err := os.Stat(fileFullPath)

	if !os.IsNotExist(err) {
		return errors.New("file already exists")
	}
	// create the file
	file, err := os.Create(fileFullPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// add the content
	_, err = file.Write(content)

	return err
}

// Append helps you append content to a file
// it returns error incase there is any
func (l *LocalStorage) Append(filepath string, content []byte) error {
	fileFullPath := path.Join(l.rootFolder, filepath)

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

// Exists checks if a file exists withn the root folder
// it returns a bool and an error incase any
func (l *LocalStorage) Exists(filepath string) (bool, error) {
	fileFullPath := path.Join(l.rootFolder, filepath)

	_, err := os.Stat(fileFullPath)
	if os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}

// Missing checks if a file is missing in the root folder
// it returns a bool and an error incase any
func (l *LocalStorage) Missing(filepath string) (bool, error) {
	fileFullPath := path.Join(l.rootFolder, filepath)

	_, err := os.Stat(fileFullPath)
	if err != nil {
		// there is an error
		if os.IsNotExist(err) {
			// not exist error
			return true, nil
		} else {
			// another errors
			return false, err
		}
	}

	return false, nil
}

func removeFirstChar(s string) string {
	_, i := utf8.DecodeRuneInString(s)
	return s[i:]
}

func removeExtension(fileName string, ext string) string {
	return strings.TrimSuffix(fileName, ext)
}
