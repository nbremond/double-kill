package commands

import (
    "fmt"
    "os"
    "path/filepath"
    "log"
    "strings"

    "github.com/codegangsta/cli"
    "github.com/AlasdairF/File"

    "github.com/nbremond/double-kill/models"
    "github.com/nbremond/double-kill/helpers"
)

var CmdSearch = cli.Command{
    Name:  "search",
    Usage: "search path",
    Description: `Index files in database`,
    Before: runSearch,
    Action: func(ctx *cli.Context) {},
    Flags:  []cli.Flag{
        cli.IntFlag{"remove, r",-1,"The Maximum number of files to delete. -1 means no limit. This is usefull if you use a distant filesystem and if you lose connection",""},
        cli.BoolFlag{"tiny-hash, t", "Compute a tiny hash (on the first bytes) of every file",""},
        cli.StringSliceFlag{"ignore-dir",&ignore_dir_slice,"Files in the listed directories will not be indexed",""},
    },
}

var computeTinyHash bool
var ignore_dir []string
var ignore_dir_slice = cli.StringSlice{}

func runSearch(c *cli.Context) error {
    var err error
    computeTinyHash = c.Bool("tiny-hash")
    ignore_dir = ignore_dir_slice.Value()
    ignore_dir = c.StringSlice("ignore-dir")
    if err = models.InitDB(); err != nil {
        return err
    }
    basePath := filepath.Clean(c.Args()[0])+"/"

    fmt.Println("Removing deleted subfiles of «"+basePath +"»…")
    dbFiles := models.GetSubfiles(basePath)
    //fmt.Println(dbFiles)
    removedFiles := 0
    for pos := range dbFiles {
        if c.Int("remove") == removedFiles {
            fmt.Print(removedFiles)
            fmt.Println(" files removed. abort.")
            break
        }
        forFile := dbFiles[pos]
        if ! file.Exists(filepath.Join(forFile.Dir, forFile.Filename)){
            fmt.Println("Removing «"+forFile.Filename+"» from database.")
            removedFiles++;
            forFile.Delete()
        }
    }
    fmt.Println("Indexing new files…")
    filepath.Walk(c.Args()[0],indexFile)
    return nil
}

func indexFile(path string, info os.FileInfo, err error) error {
    if err != nil {
        log.Println(err)
        return nil
    }
    if info.IsDir() {
        for _,value := range ignore_dir{
            if strings.HasPrefix(path,value) {
                return filepath.SkipDir
            }
        }
    }else{// path isn' t a dir. We must check if it's a regular file.
        dir, filename := filepath.Split(path)
        upToDate := true
        isNew := false
        dbFile := models.GetOrCreateFile(dir, filename)
        if dbFile.Id == 0 {
            upToDate = false
            isNew = true
        }
        if dbFile.ModTime != info.ModTime(){
            dbFile.ModTime = info.ModTime()
            upToDate = false
        }
        if dbFile.Size != info.Size(){
            dbFile.Size = info.Size()
            upToDate = false
        }
        if ! upToDate{
            dbFile.TinyHash = ""
            dbFile.Hash = ""
            if isNew && computeTinyHash{
                dbFile.TinyHash = helpers.ComputeTinyHash(path)
            }
        }
        dbFile.Save()
        //fmt.Println("done"+dbFile.Filename)
    }
    return nil
}

