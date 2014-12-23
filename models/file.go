package models

import (
    "path/filepath"
    "time"
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
