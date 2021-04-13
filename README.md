# Stowage
## (under development)

## What is stowage?
Stowage is a go package for working with filesystems, it supports local, Amazon S3, Google Cloud Storage, Alicloud OSS

## Install
To install stowage run the following command: 
```bash
go get github.com/harranali/stowage
```
## Getting Started 

### Local storage
let's get a local file information
```go
s := stowage.New()

// Initiate the local storage
storePath, _ := filepath.Abs("./my-base-folder")
s.InitLocalStorage(stowage.LocalStorageOpts{
    RootFolder: storePath,
})

// Get file information
info, _ := s.LocalStorage.FileInfo("testfile.txt")
fmt.Println(info.Name)
fmt.Println(info.NameWithoutExtension)
fmt.Println(info.Extension)
fmt.Println(info.Size)
fmt.Println(info.Path)
```