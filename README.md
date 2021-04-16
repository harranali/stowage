# Stowage
A Go package for working with files and directories

![Build Status](https://github.com/harranali/stowage/actions/workflows/build-master.yml/badge.svg)
![Test Status](https://github.com/harranali/stowage/actions/workflows/test-master.yml/badge.svg)

## Install
To install stowage run the following command: 
```bash
go get github.com/harranali/stowage
```
## Getting Started 
```go
// get the absolute path to the root directory
rootFolder, _ := filepath.Abs("./my-base-folder")

// first create the package variable 
s := stowage.New()

// Initiate the storage with root directory
s.InitLocalStorage(stowage.LocalStorageOpts{
    RootFolder: rootFolder,
})

// get file information
info, _ := s.LocalStorage.FileInfo("testfile.txt")
fmt.Println(info.Name) // the file name with extension
fmt.Println(info.NameWithoutExtension) // the file name without extension
fmt.Println(info.Extension) // the file extension

// copy file 
err := s.LocalStorage.Copy("myfile.txt", "files/backup/txt")

```
####  Root folder
The `rootFolder` acts as the reference for all operations, it has to be a full absolute path
here is how you can get the absolute path for your root folder:
```go
// get the absolute path to the root directory
rootFolder, _ := filepath.Abs("./my-base-folder")
```

#### Package initiation
you need first to create the package variable with by calling the method New `s := stowage.New()`, next you need to initiate the storage engine by calling the method `s.InitLocalStorage(opts)` and pass to it the options, the code below shows how you can create the variable and initiate the storage engine
```go
// first create the package variable 
s := stowage.New()

rootFolder, _ := filepath.Abs("./my-base-folder")
// Initiate the storage enginey
s.InitLocalStorage(stowage.LocalStorageOpts{
    RootFolder: rootFolder,
})
```


#### File information 
Here is how you can get information about a file such as name, extension, size, and more.
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
#### Example operations
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

// rename a file 
s.LocalStorage.Rename(filePath string, newFilePath string) error
```


#### All operations 
Here is a list of all supported operations
```go
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
Directories(DirectoryPath string) (directoryPaths []string, err error)
AllDirectories(DirectoryPath string) (directoryPaths []string, err error)
MakeDirectory(DirectoryPath string, perm int) error
RenameDirectory(DirectoryPath string, NewDirectoryPath string) (err error)
DeleteDirectory(DirectoryPath string) (err error)
```