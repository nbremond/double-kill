package commands

import (
    "fmt"
    "path/filepath"
    "sort"

    "github.com/codegangsta/cli"

    "github.com/nbremond/double-kill/models"
    "github.com/nbremond/double-kill/helpers"
)


var CmdSearchDuplicate = cli.Command{
    Name:  "search_duplicate",
    Usage: "search_duplicate",
    Description: `Search for duplicates`,
    Before: runSearchDuplicate,
    Action: func(ctx *cli.Context) {},
    Flags:  []cli.Flag{},
}


func runSearchDuplicate(c *cli.Context) error {
    var err error
    if err = models.InitDB(); err != nil {
        return err
    }

    files := models.GetFilesBySize()
    if len(files) == 0{
        fmt.Println("No file inedexed. Run " + CmdSearch.Name + " before search_duplicate")
        return nil
    }

    potentialFiles := make([]*models.File,0,10) //will contain files with the same tinyHash
    var fileSize int64 //size of files curently in potentialFiles
    for i := range files {
        forFile := &files[i]
        potentialFiles = append(potentialFiles, forFile)
        fileSize = forFile.Size
        if i == (len(files) - 1) || files[i+1].Size != fileSize {
            if len(potentialFiles) > 1 {
                analyseSameSize(potentialFiles)
            }
            potentialFiles = make([]*models.File,0,10)
        }
    }
    return nil
}


func analyseSameSize(files []*models.File) {
    // first compute all tinyHash
    allTinyHashesKnown := true
    for pos := range files {
        forFile := files[pos]
        if forFile.TinyHash == "" {
            forFile.TinyHash = helpers.ComputeTinyHash(filepath.Join(forFile.Dir, forFile.Filename))
            forFile.Save()
            if forFile.TinyHash == "" {
                allTinyHashesKnown = false
            }
        }
    }
    if ! allTinyHashesKnown {
        fmt.Print(len(files))
        fmt.Print(" files with the same size but which cannot be compared. (")
        fmt.Print(files[0].Size)
        fmt.Println(" bytes)")

    }
    //then we sort by their tinyhash
    fs := &fileSorter{
        files : files,
        by :    func (f1,f2 *models.File) bool{
            return f1.TinyHash < f2.TinyHash
        },
    }
    sort.Sort(fs)
    currentHash := ""
    potentialFiles := make([]*models.File,0,10)
    for pos := range files {
        forFile := files[pos]
        potentialFiles = append(potentialFiles, forFile)
        currentHash = forFile.TinyHash
        if pos+1 == len(files) || files[pos+1].TinyHash != currentHash{
            if currentHash != "" && len(potentialFiles) > 1{
                /////TODO Maybe set a verbosity level ?
                //fmt.Print(len(potentialFiles))
                //fmt.Print(" files with the same TinyHash «"+currentHash+"» (")
                //fmt.Print(potentialFiles[0].Size)
                //fmt.Println(" bytes)")
                analyseSameTinyHash(potentialFiles)
            }
            potentialFiles = make([]*models.File,0,10)
        }
    }
}



func analyseSameTinyHash(files []*models.File) {
    // first compute all Hash
    allHashesKnown := true
    for _,forFile := range files {
        if forFile.Hash == "" {
            forFile.Hash = helpers.ComputeHash(filepath.Join(forFile.Dir, forFile.Filename))
            forFile.Save()
            if forFile.Hash == "" {
                allHashesKnown = false
            }
        }
    }
    if ! allHashesKnown {
        fmt.Print(len(files))
        fmt.Print(" files with the same TinyHash but which cannot be compared. (")
        fmt.Print(files[0].Size)
        fmt.Println(" bytes)")

    }
    //then we sort by their hash
    fs := &fileSorter{
        files : files,
        by :    func (f1,f2 *models.File) bool{
            return f1.Hash < f2.Hash
        },
    }
    sort.Sort(fs)
    currentHash := ""
    potentialFiles := make([]*models.File,0,10)
    for pos := range files {
        forFile := files[pos]
        potentialFiles = append(potentialFiles, forFile)
        currentHash = forFile.Hash
        if pos+1 == len(files) || files[pos+1].Hash != currentHash{
            if currentHash != "" && len(potentialFiles) > 1{
                //TODO set a verbosity level ?
                // fmt.Print(len(potentialFiles))
                // fmt.Print(" files with the same Hash «"+currentHash+"» (")
                // fmt.Print(potentialFiles[0].Size)
                // fmt.Println(" bytes)")
                sortedFiles,sortErr := helpers.SortFilesByteByByte(potentialFiles)
                if sortErr == nil {
                    for _,fileSet := range sortedFiles {
                        if len(fileSet) > 1 {
                            fmt.Print(len(fileSet))
                            fmt.Println(" files are strictly identical")
                            for _,file := range fileSet {
                                fmt.Println(filepath.Join(file.Dir, file.Filename))
                                //TODO keep this result in db
                            }
                        }
                    }
                } else {
                    fmt.Println(sortErr)
                }
            }
            potentialFiles = make([]*models.File,0,10)
        }
    }
}


//helpers for sorting files

type fileSorter struct {
    files  []*models.File
    by     func(f1,f2 *models.File) bool
}
func (s *fileSorter) Len() int {
    return len(s.files)
}
func (s *fileSorter) Swap(i, j int) {
    s.files[i], s.files[j] = s.files[j], s.files[i]
}
func (s *fileSorter) Less(i, j int) bool {
    return s.by(s.files[i], s.files[j])
}
