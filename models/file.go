package models

import (
    "path/filepath"
    "time"
    "sort"
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


type fileSorter struct {
    files  []File
    by     func(f1,f2 *File) bool
}
// Len is part of sort.Interface.
func (s *fileSorter) Len() int {
    return len(s.files)
}

// Swap is part of sort.Interface.
func (s *fileSorter) Swap(i, j int) {
    s.files[i], s.files[j] = s.files[j], s.files[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *fileSorter) Less(i, j int) bool {
    return s.by(&s.files[i], &s.files[j])
}

func GetFilesBySize() []File{
    var files []File
    db.Find(&files)
    fs := &fileSorter{
        files : files,
        by :    func (f1,f2 *File) bool{
            return f1.Size < f2.Size
        },
    }
    sort.Sort(fs)
    return files
}
