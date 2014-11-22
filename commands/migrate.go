package commands

import (
    "fmt"

    "github.com/codegangsta/cli"
)

var CmdSearch = cli.Command{
    Name:  "search",
    Usage: "Search for duplicates",
    Description: `Search for duplicates`,
    Before: runMigrate,
    Action: func(ctx *cli.Context) {},
    Flags:  []cli.Flag{},
}

func runMigrate(ctx *cli.Context) error {
    fmt.Println("It's working!")

    return nil
}
