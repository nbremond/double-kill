
package commands

import (
    "fmt"
    "os"
    "path/filepath"

    "github.com/codegangsta/cli"

    "github.com/nbremond/double-kill/models"
)

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
        fmt.Println(dir+filename)
        file := models.File{
            Dir:        dir,
            Filename:   filename,
            Size:       info.Size(),
            UpdatedAt:  info.ModTime(),
            Hash:       "",
            TinyHash:   "",
        }
        file.Save()
    }
    return nil
}

