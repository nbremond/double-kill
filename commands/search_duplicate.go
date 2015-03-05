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
    var size int64
    potentialFiles := make([]*models.File,0,10)
    files := models.GetFilesBySize()
    if len(files) == 0{
        fmt.Println("No file inedexed. Run search before search_duplicate")
        return nil
    }

    for pos := range files {
        forFile := &files[pos]
        potentialFiles = append(potentialFiles, forFile)
        size = forFile.Size
        if pos+1 == len(files) || files[pos+1].Size != size {
            analyseSameSize(potentialFiles)

            potentialFiles = make([]*models.File,0,10)
        }
    }
    return nil
}

func analyseSameSize(files []*models.File) {
    if len(files) < 2 {
        return
    }
    /*
    fmt.Print(len(files))
    fmt.Print(" fichiers de taile ")
    fmt.Print(files[0].Size)
    fmt.Println()
    */
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
        //fmt.Println(forFile.TinyHash)
        if pos+1 == len(files) || files[pos+1].TinyHash != currentHash{
            if currentHash != "" && len(potentialFiles) > 1{
                /////Maybe set a verbosity level ?
                //fmt.Print(len(potentialFiles))
                //fmt.Print(" files with the same TinyHash «"+currentHash+"» (")
                //fmt.Print(potentialFiles[0].Size)
                //fmt.Println(" bytes)")

                ///// are the set realy usefull ?
                //m := models.MatchingFilesSet{
                //    Level:  models.TinyHash,
                //    Files:  make([]models.File,0,10),
                //}
                //for _,e := range potentialFiles {
                //m.Files = append(m.Files, *e)
                //}
                //m.Save()
                //////
                analyseSameTinyHash(potentialFiles)
            }
            potentialFiles = make([]*models.File,0,10)
        }
    }
}



func analyseSameTinyHash(files []*models.File) {
    if len(files) < 2 {
        return
    }
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
        //fmt.Println(forFile.TinyHash)
        if pos+1 == len(files) || files[pos+1].Hash != currentHash{
            if currentHash != "" && len(potentialFiles) > 1{
                fmt.Print(len(potentialFiles))
                fmt.Print(" files with the same Hash «"+currentHash+"» (")
                fmt.Print(potentialFiles[0].Size)
                fmt.Println(" bytes)")
                sortedFiles,sortErr := helpers.SortFilesByteByByte(potentialFiles)
                if sortErr == nil {
                    for _,fileSet := range sortedFiles {
                        if len(fileSet) > 1 {
                            fmt.Print(len(fileSet))
                            fmt.Println(" files are strictly identical")
                            for _,file := range fileSet {
                                fmt.Println(filepath.Join(file.Dir, file.Filename))
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

type fileSorter struct {
    files  []*models.File
    by     func(f1,f2 *models.File) bool
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
    return s.by(s.files[i], s.files[j])
}
