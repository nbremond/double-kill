package helpers

import (
    "fmt"
    "os"
    "path/filepath"
    "log"

    "github.com/nbremond/double-kill/models"
)

const blockSize = 1024

type openFile struct {
    originFile  *models.File
    data        []byte
    osFile      *os.File
}

type fileSet struct {
    ToCheckFiles        []openFile
    CheckedFiles        []openFile
    changingBlock       []byte
    changingBlockPos    int
    childs              []fileSet
}

func SortFilesByteByByte(files []*models.File) ([][]*models.File, error) {
    firstFile,err := open(files[0])
    if err != nil {
        return nil, err
    }
    if firstFile.read() != nil{
        return nil, err
    }

    rootSet := fileSet{
        ToCheckFiles:   make([]openFile,0,10),
        CheckedFiles:   append(make([]openFile,0,10), firstFile),
        changingBlock:  firstFile.data,
        changingBlockPos: 0,
        childs:         make([]fileSet,0,10),
    }
    for pos,file := range files {
        if pos != 0 {
            ofile,err := open(file)
            if err != nil {
                return nil, err
            }
            rootSet.ToCheckFiles = append( rootSet.ToCheckFiles, ofile)
        }
    }
//calculer
//remettre en forme
//fermer les fichiers
    fmt.Println(len(firstFile.data))
    return make([][]*models.File,0,10),nil
}

func open(originFile *models.File) (openFile, error) {
    file := openFile {
        originFile: originFile,
        data:       make([]byte, blockSize),
    }
    var err error
    path := filepath.Join(originFile.Dir, originFile.Filename)
    file.osFile,err = os.Open(path)
    if err != nil {
        log.Println("Unable to open \"" + path + "\"")
        return file, err
    }
    return file, nil
}

func (f *openFile) read() error {
    f.data = make([]byte, blockSize)
    n,err := f.osFile.Read(f.data)
    f.data = f.data[:n]
    if err != nil {
        path := filepath.Join(f.originFile.Dir, f.originFile.Filename)
        log.Println("Unable to read \"" + path + "\"")
    }
    return err
}


