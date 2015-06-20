package models

import (
    "path/filepath"
    "time"
    "os"
    //    "fmt"

    _ "github.com/jinzhu/gorm"

    //    "github.com/nbremond/double-kill/modules/settings"
)

type File struct {
    Id          int64
    Dir         string `sql:"size:255"`
    Filename    string `sql:"size:255"`
    Size        int64
    TinyHash    string `sql:"size:65"`
    Hash        string `sql:"size:65"`
    ModTime     time.Time
}

func init(){
    register(&File{})
}

func standardizePath(path string) string{
    basePath := "."
    //var err error
    path, _ = filepath.Rel(basePath, path)
    //    fmt.Println(err)
    path = filepath.Clean(path)
    if path == "." || path == "" {
        path = ""
    }else{
        path = path+"/"
    }
    return path
}

func (f *File) Save() {
    f.Dir = standardizePath(f.Dir)
    db.Save(f)
}

func GetSubfiles(path string) []File {
    path = standardizePath(path)
    var files []File
    db.Where("dir LIKE ?", path+"%").Find(&files)
    return files
}

func GetOrCreateFile(dir string, filename string) File{
    dir = standardizePath(dir)
    file := File{
        Dir:        dir,
        Filename:   filename,
    }
    db.Where("dir = ? AND filename = ?", dir, filename).Find(&file)
    return file
}

func (file *File) Delete() {
    db.Delete(file)
}


func GetFilesBySize() []File{
    var files []File
    db.Order("size desc").Find(&files)
    return files
}

func GetMatchingHashFiles() [][]File{
    var sortedFiles []File
    var groupedFiles [][]File
    db.Order("size desc").Order("tiny_hash desc").Order("hash").Find(&sortedFiles)
    currentFiles := make([]File,0,10)
    currentSize := int64(0)
    currentTinyHash := ""
    currentHash := ""
    for _,f := range sortedFiles{
        if (f.Size != currentSize) ||
        (f.TinyHash != currentTinyHash) ||
        (f.Hash != currentHash) {
            if len(currentFiles) > 1{
                groupedFiles = append(groupedFiles, currentFiles)
            }
            currentFiles = make([]File,0,10)
            currentSize = f.Size
            currentTinyHash = f.TinyHash
            currentHash = f.Hash
        }
        currentFiles = append(currentFiles, f)
    }
    return groupedFiles
}

func (dbFile *File) IsUpToDate(info os.FileInfo) bool {
    upToDate := true
    if dbFile.Id == 0 {
        upToDate = false
    }
    if dbFile.ModTime != info.ModTime().UTC(){
        dbFile.ModTime = info.ModTime().UTC()
        upToDate = false
    }
    if dbFile.Size != info.Size(){
        dbFile.Size = info.Size()
        upToDate = false
    }
    if ! upToDate{
        dbFile.TinyHash = ""
        dbFile.Hash = ""
    }
    dbFile.Save()
    return upToDate
}
