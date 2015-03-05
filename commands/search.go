package commands

import (
    "fmt"
    "os"
    "path/filepath"
    "log"

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
    Flags:  []cli.Flag{},
}


func runSearch(c *cli.Context) error {
    var err error

    if err = models.InitDB(); err != nil {
        return err
    }
    basePath := filepath.Clean(c.Args()[0])+"/"

    fmt.Println("Removing deleted subfiles of «"+basePath +"»…")
    dbFiles := models.GetSubfiles(basePath)
    //fmt.Println(dbFiles)
    for pos := range dbFiles {
        forFile := dbFiles[pos]
        if ! file.Exists(filepath.Join(forFile.Dir, forFile.Filename)){
            fmt.Println("Removing «"+forFile.Filename+"» from database.")
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
        return err
    }
    if info.IsDir() {
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
        if ! upToDate {
            dbFile.TinyHash = ""
            dbFile.Hash = ""
            if isNew {
                dbFile.TinyHash = helpers.ComputeTinyHash(path)
            }
        }
        dbFile.Save()
        //fmt.Println("done"+dbFile.Filename)
    }
    return nil
}

