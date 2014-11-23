package commands

import (
    "fmt"
    "os"
    "path/filepath"
    "crypto/sha256"
    "encoding/hex"
    "log"

    "github.com/codegangsta/cli"

    "github.com/nbremond/double-kill/models"
)

const tinyHashSize = 10 * 1024

var CmdSearch = cli.Command{
    Name:  "search",
    Usage: "Search for duplicates",
    Description: `Search for duplicates`,
    Before: runSearch,
    Action: func(ctx *cli.Context) {},
    Flags:  []cli.Flag{},
}


func runSearch(c *cli.Context) error {
    var err error

    if err = models.InitDB(); err != nil {
        return err
    }
    fmt.Println("It's working! or not ...")
    filepath.Walk(c.Args()[0],printFile)
    return nil
}

func printFile(path string, info os.FileInfo, err error) error {
    if ! info.IsDir() {
        dir, filename := filepath.Split(path)
        hashError := false
        tinyHashString := ""
        if file, fileErr := os.Open(path); fileErr != nil {
            hashError =true
        }else{
            defer file.Close()
            data := make([]byte, tinyHashSize)
            n, readErr := file.Read(data)
            if  readErr != nil {
                hashError = true
            }
            if hashError {
                log.Println("Unable to compute hash for \""+path+"\"")
            }else{
                tinyHash := sha256.New()
                tinyHash.Write(data[:n])
                tinyHashString = hex.EncodeToString(tinyHash.Sum(nil))
            }
        }
        dbFile := models.File{
            Dir:        dir,
            Filename:   filename,
            Size:       int64(info.Size()),
            ModTime:    info.ModTime(),
            TinyHash:   tinyHashString,
        }
        dbFile.Save()
    }
    return nil
}

