package helpers

import (
    "fmt"
    "os"
    "path/filepath"
    "log"
    "io"

    "github.com/nbremond/double-kill/models"
)

const blockSize = 1024 //buffer size for reading files

type openFile struct {
    originFile  *models.File
    data        []byte
    osFile      *os.File
}

type fileSet struct {
    ToCheckFiles        []openFile 
    CheckedFiles        []openFile //CheckedFiles[0] is the refence file for the block
    changingBlock       []byte
    changingBlockPos    int //all blocks of checked_files before this pos are the sames. 
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
            ofile.read()
            rootSet.ToCheckFiles = append( rootSet.ToCheckFiles, ofile)
        }
    }
    err = rootSet.compute()
    res := make([][]*models.File,0,10)
    rootSet.formatAndClose(&res)
    return res,err
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
        fmt.Println()//there is no newline after the progress status
        log.Println("Unable to open \"" + path + "\"")
        return file, err
    }
    return file, nil
}

func (f *openFile) read() error {
    f.data = make([]byte, blockSize)
    n,err := f.osFile.Read(f.data)
    f.data = f.data[:n]
    if err == io.EOF {
    err = nil
    }
    if err != nil {
        path := filepath.Join(f.originFile.Dir, f.originFile.Filename)
        fmt.Println()//there is no newline after the progress status
        log.Println("Unable to read \"" + path + "\"")
    }
    return err
}

//b1 and b2 are supposed to have the same size
func compareBlock(b1, b2 []byte) bool {
    for i := range b1 {
        if b1[i] != b2[i] {
            return false
        }
    }
    return true
}

func (fs *fileSet) compute() error {
    pos := fs.changingBlockPos
    for len(fs.CheckedFiles[0].data) > 0 {//iterate on all blocks of the file
        newToCheckFiles := make([]openFile,0,10) //this is used to remove files from the slice fs.ToCheckFiles
        for _,file := range fs.ToCheckFiles {
            if !compareBlock(file.data, fs.CheckedFiles[0].data) {
                if len(fs.childs) > 0 && pos == fs.childs[len(fs.childs)-1].changingBlockPos {
                    newfs := fs.childs[len(fs.childs)-1]
                    newfs.ToCheckFiles = append(newfs.CheckedFiles, file)
                } else {
                    newfs := fileSet{
                        ToCheckFiles:     make([]openFile,0,10),
                        CheckedFiles:     append(make([]openFile,0,10), file),
                        changingBlock:    file.data,
                        changingBlockPos: pos,
                        childs:           make([]fileSet,0,10),
                    }
                    fs.childs = append(fs.childs, newfs)
                }
            } else {
                newToCheckFiles = append(newToCheckFiles,file)
            }
        }
        fs.ToCheckFiles = newToCheckFiles

        for _,file := range fs.ToCheckFiles {
            err := file.read();
            if err != nil {
                return err
            }
        }
        err := fs.CheckedFiles[0].read()
        if err != nil {
            return err
        }
    }
    for _,file := range fs.ToCheckFiles {
        fs.CheckedFiles = append(fs.CheckedFiles, file)
    }
    fs.ToCheckFiles = make([]openFile,0,10)
    for _,toCompute := range fs.childs {
        if err := toCompute.compute(); err != nil{
            return err
        }
    }
    return nil
}

func (fs *fileSet) formatAndClose(res *[][]*models.File) {
    if len(fs.CheckedFiles) > 0 {
        sameFiles := make([]*models.File,0,10)
        for _,file := range fs.CheckedFiles {
            sameFiles = append(sameFiles, file.originFile)
            file.osFile.Close()
        }
        *res = append(*res,sameFiles)
    }
    for _,file := range fs.ToCheckFiles {
        file.osFile.Close()
    }
    for _,child := range fs.childs {
        child.formatAndClose(res)
    }
}

