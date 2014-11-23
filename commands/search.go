package commands

import (
    "fmt"
    "os"
    "path/filepath"
    "strconv"

    "github.com/codegangsta/cli"
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
    fmt.Println("It's working!")
    filepath.Walk(c.Args()[0],printFile)
    return nil
}

func printFile(path string, info os.FileInfo, err error) error {
    fmt.Println(path + strconv.FormatInt(info.Size(), 10))
    return nil
}
