# Stowage
A Go filesystem package for working with files and directories, it features a simple API with support for the common files and directories operations such as copy, move, list files, file exists check, and more.

![Build Status](https://github.com/harranali/stowage/actions/workflows/build-master.yml/badge.svg) ![Test Status](https://github.com/harranali/stowage/actions/workflows/test-master.yml/badge.svg) [![GoDoc](https://godoc.org/github.com/harranali/stowage?status.svg)](https://godoc.org/github.com/harranali/stowage)
[![reportcard](https://goreportcard.com/badge/harranali/stowage)](https://goreportcard.com/report/harranali/stowage) [![codecov](https://codecov.io/gh/harranali/stowage/branch/master/graph/badge.svg?token=EHLAZVHY55)](https://codecov.io/gh/harranali/stowage)


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
## Root folder
The `rootFolder` acts as the reference for all operations, it has to be a full absolute path
here is how you can get the absolute path for your root folder:
```go
// get the absolute path to the root directory
rootFolder, _ := filepath.Abs("./my-base-folder")
```

## Package initiation
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


## Getting File information 
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
## Example operations
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


## List of supported operations 
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

## docs
Here are the details of each operation

##### FileInfo(filePath string) (fileinfo localstorage.FileInfo, err error)
`FileInfo` returns information about the given file or an error incase there is any
```go
info, _ := s.LocalStorage.FileInfo("testfile.txt")
info.Name // the file name with extension
info.NameWithoutExtension // the file name without extension
info.Extension // the file extension
info.Size // the file size 
info.Path // the file full path
```

#### Put(filePath string) error
`Put` helps you copy files into the root directory from external locations, filePath is the full path to the file you would like to put, it returns error incase there is any
```go
err := s.LocalStorage.Put("testfile.txt")
```

####  PutAs(filePath string, filename string) error
`PutAs` helps you copy files into the root directory from external directory,  the first param 'filePath' is the full path to the file you would like to put, the second param 'fileName' is the name you would like to give to the file, it returns error incase there is any
```go
err := s.LocalStorage.PutAs("testfile.txt", "newtestfile.txt")
```

#### Copy(filePath string, destPath string) error
`Copy` helps you copy files within the root folder, Please note that the reference of the paths of these files is the root folder, it accepts the source file starting from the root folder, and the destination folder starting from the root folder, it returns an error incase there is any
```go
err := s.LocalStorage.Copy("testfile.txt", "/folder/subfolder")
```

#### CopyAs(filePath string, destfolder string, newFilePath string) error
`CopyAs` helps you copy files within the root folder, Please note that the reference of the paths of these files is the root folder, it accepts the source file starting from the root folder and the destination folder starting from the root folder, and the new file name, it returns an error incase there is any
```go
err := s.LocalStorage.CopyAs("testfile.txt", "/folder/subfolder", "newtestfile.txt")
```

####  Move(filePath string, destFolder string) error
`Move` helps you Move files within the root folder, Please note that the reference of the paths of these files is the root folder, it accepts the source file starting from the root folder, and the destination folder starting from the root folder, it returns an error incase there any
```go
err := s.LocalStorage.Move("testfile.txt", "/folder/subfolder")
```

####  MoveAs(filePath string, destFolder string, newFilePath string) error
`MoveAs` helps you Move files within the root folder, Please note that the reference of the paths of these files is the root folder, it accepts the source file starting from the root folder and the destination folder starting from the root folder, and the new file name, it returns an error incase there any
```go
err := s.LocalStorage.Move("testfile.txt", "/folder/subfolder", "newtestfile.txt")
```

#### Rename(filePath string, newFilePath string) error
`Rename` renames the given file as first parameter to the name given as a second parameter, it returns error incase there is any
```go
err := s.LocalStorage.Rename("testfile.txt",  "newtestfile.txt")
```

#### Delete(filePath string) error
`Delete` deletes the given file it returns error incase there is any
```go
err := s.LocalStorage.Delete("testfile.txt")
```

#### DeleteMultiple(filePaths []string) (err error)
`DeleteMultiple` deltes multiple files given as slice of strings of file paths, it returns error incase there is any
```go
files := []string{"testfile1.txt", "testfile1.txt"}
err := s.LocalStorage.DeleteMultiple(files)
```

#### Create(filePath string, content []byte) error
`Create` helps you create new a file and add content to it, it returns error incase there is any
```go
err := s.LocalStorage.Create("newfile.txt", []byte("this is a sample text"))
```

#### Append(filePath string, content []byte) error
`Append` helps you append content to a file, it returns error incase there is any
```go
err := s.LocalStorage.Append("newfile.txt",[]byte("this is a sample text"))
```

#### Exists(filePath string) (bool, error)
`Exists` checks if a file exists withn the root folder, it returns a bool and an error incase any
```go
exists, err := s.LocalStorage.Exists("newfile.txt")
```

#### Missing(filePath string) (bool, error)
`Missing` checks if a file is missing in the root folder, it returns a bool and an error incase any
```go
missing, err := s.LocalStorage.Missing("newfile.txt")
```
#### Read(filePath string) ([]byte, error)
`Read` helps you grap the content of a file, it returns the data in a slice of bytes and an error incase there is any
```go
content, err := s.LocalStorage.Read("newfile.txt")
```

#### Files(DirectoryPath string) (files []FileInfo, err error)
`Files` returns a list of files in a given directory, the file type is LocalStorage.FileInfo  NOT the standard library fs.FileInfo, and it returns an error incase any occurred, if you want a list of files including the files in sub directories, consider using the method `AllFiles(DirectoryPath string)`
```go
files, err := s.LocalStorage.Files("mydir")
```

####  AllFiles(DirectoryPath string) (files []FileInfo, err error)
`AllFiles` returns a list of files in the given directory including files in sub directories, the file type in the list is LocalStorage.FileInfo NOT the standard library fs.FileInfo
```go
files, err := s.LocalStorage.Files("mydir")
```

####  Directories(DirectoryPath string) (directoryPaths []string, err error)
`Directories` returns a slice of string containing the paths of the directories, if you want the list of directories including subdirectories, consider using the method "AllDirectories(DirectoryPath string)", it returns an error incase is any
```go
directories, err := s.LocalStorage.Directories("mydir")
```

#### AllDirectories(DirectoryPath string) (SubDirectoryPath []string, err error)
`AllDirectories` returns a list of directories including sub directories, it returns an error incase is any
```go
directories, err := s.LocalStorage.AllDirectories("mydir")
```

#### MakeDirectory(DirectoryPath string, perm int) (err error)
`MakeDirectory` creates a new directory and the necessary parent directories with the given permissions, permissions could be (example: 0777) or any linux based permissions, it returns an error incase is any
```go
err := s.LocalStorage.MakeDirectory("mydir")
```

#### RenameDirectory(DirectoryPath string, NewDirectoryPath string) (err error)
`RenameDirectory` changes the name of directory to new name, it returns an error incase there is any
```go
err := s.LocalStorage.RenameDirectory("mydir",  "new-dir-name")
```


#### DeleteDirectory(DirectoryPath string) (err error)
`DeleteDirectory` deletes the given directory
```go
err := s.LocalStorage.DeleteDirectory("mydir")
```

