# Stowage
## (under development)
![Build Status](https://github.com/harranali/stowage/actions/workflows/build-master.yml/badge.svg)
![Test Status](https://github.com/harranali/stowage/actions/workflows/test-master.yml/badge.svg)

## What is stowage?
A Go package for working with local filesystems

## Install
To install stowage run the following command: 
```bash
go get github.com/harranali/stowage
```
## Getting Started 

#### Initiate the package
```go
// get the full path absolute path to the root dir
rootFolder, _ := filepath.Abs("./my-base-folder")

// first create the package variable 
s := stowage.New()

// Initiate the storage
s.InitLocalStorage(stowage.LocalStorageOpts{
    RootFolder: rootFolder,
})
```

#### File information 
```go 
rootFolder, _ := filepath.Abs("./my-base-folder")

s := stowage.New()
s.InitLocalStorage(stowage.LocalStorageOpts{
    RootFolder: rootFolder,
})

info, _ := s.LocalStorage.FileInfo("testfile.txt")
fmt.Println(info.Name) // the file name with extension
fmt.Println(info.NameWithoutExtension) // the file name without extension
fmt.Println(info.Extension) // the file extension
fmt.Println(info.Size) // the file size 
fmt.Println(info.Path) // the file full path
```
#### File operations
All file operations are performed with respect to the root directory
```go 
rootFolder, _ := filepath.Abs("./my-base-folder")

s := stowage.New()
s.InitLocalStorage(stowage.LocalStorageOpts{
    RootFolder: rootFolder,
})

// copy file to the root directory from external directory
s.LocalStorage.Put(filePath string) error

// copy file to the root directory from external directory with the given name
s.LocalStorage.PutAs(filePath string, filename string) error

// copy file around within the root directory
s.LocalStorage.Copy(filePath string, destfolder string) error

// copy file around within the root directory with given name as third param
s.LocalStorage.CopyAs(filePath string, destfolder string, newFilePath string) error

// move file around within the root directory
s.LocalStorage.Move(filePath string, destfolder string) error

// move file around within the root directory with given name as third param
s.LocalStorage.MoveAs(filePath string, destFolder string, newFilePath string) error

// rename a file 
s.LocalStorage.Rename(filePath string, newFilePath string) error

// delete a file
s.LocalStorage.Delete(filePath string) error

// delete multipe files
s.LocalStorage.DeleteMultiple(filePaths []string) error

// create a file 
s.LocalStorage.Create(filePath string, content []byte) error

// append to a file
s.LocalStorage.Append(filePath string, content []byte) error

// check if file exists
s.LocalStorage.Exists(filePath string) (bool, error)

// check if file is missing
s.LocalStorage.Missing(filePath string) (bool, error)

// read content of a file
s.LocalStorage.Read(filePath string) ([]byte, error)
```