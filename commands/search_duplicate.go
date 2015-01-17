package commands

import (
    "fmt"
  //  "os"
  //  "path/filepath"
  //  "crypto/sha256"
  //  "encoding/hex"
  //  "log"

    "github.com/codegangsta/cli"
 //   "github.com/AlasdairF/File"

    "github.com/nbremond/double-kill/models"
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
    var potentialFiles []*models.File
    files := models.GetFilesBySize()
    for pos := range files {
        forFile := files[pos]
        if forFile.Size != size {
        size = forFile.Size

        }
        //Append(potentialFiles, forFile)
        fmt.Println(forFile.Size)
    }
return nil
}
