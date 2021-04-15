# Stowage
## (under development)

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

    // copy a file from outside the root dir to the root dir with the same name
	Put(filepath string) error
	PutAs(filepath string, filename string) error // copy a file from outside the root dir to the root dir with different
	Copy(srcfile string, destfolder string) error
	CopyAs(srcfile string, destfolder string, newfilename string) error
	Move(srcfile string, destfile string) error
	MoveAs(srcfile string, destfile string, newfilename string) error
	Rename(filename string, newfilename string) error
	Delete(filepath string) error
	DeleteMultiple(filepaths []string) error
	Create(filepath string, content []byte) error
	Append(filepath string, content []byte) error
	Exists(filepath string) (bool, error)
	Missing(filepath string) (bool, error)
	Read(filepath string) ([]byte, error)
```