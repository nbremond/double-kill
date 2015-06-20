package commands

import (
    "fmt"
    "os"
    //"path/filepath"
    //"log"
    "strings"

    "github.com/codegangsta/cli"
    //"github.com/AlasdairF/File"

    "github.com/nbremond/double-kill/models"
    //"github.com/nbremond/double-kill/helpers"
)

var CmdRemove = cli.Command{
    Name:  "remove",
    Usage: "remove", //remove [keep-dir <path> | remove-dir <path>]",
    Description: `remove duplicates files`,
    Before: runRemove,
    Action: func(ctx *cli.Context) {},
    Flags:  []cli.Flag{
        cli.BoolFlag{"no-byte-comparison", "Only verify hash before removing a file",""},
        cli.StringSliceFlag{"keep-dir",&keep_dir_slice,"Files in the listed directories will not be removed",""},
        cli.StringSliceFlag{"remove-dir",&remove_dir_slice,"Files in the listed directories will be removed",""},
    },
}

var keep_dir_slice = cli.StringSlice{}
var remove_dir_slice = cli.StringSlice{}

func runRemove(c *cli.Context) error {
    var err error
    if err = models.InitDB(); err != nil {
        return err
    }
    groupedFiles := models.GetMatchingHashFiles()
    for _,group := range groupedFiles{
        refFiles := make([]models.File, 0, 10)
        toDeleteFiles := make([]models.File, 0, 10)
        othersFiles := make([]models.File, 0, 10)
        for _,f := range group{
            info,err := os.Stat(f.Dir + f.Filename)
            if err != nil { continue }
            if ! f.IsUpToDate(info) { continue }
            isRefFile := false
            mustBeRemoved := false
            for _,path := range keep_dir_slice{
                if strings.HasPrefix(f.Dir, path){
                    isRefFile = true
                }
            }
            if isRefFile {
                refFiles = append(refFiles, f)
            } else {
                for _,path := range remove_dir_slice{
                    if strings.HasPrefix(f.Dir, path){
                        mustBeRemoved = true
                    }
                }
                if mustBeRemoved {
                    toDeleteFiles = append(toDeleteFiles, f)
                } else {
                    othersFiles = append(othersFiles, f)
                }
            }

            //fmt.Print(f.Size)
            //fmt.Print("\t"+f.TinyHash)
            //fmt.Print(" "+f.Hash+" ")
            //fmt.Print(isRefFile)
            //fmt.Print("  "+f.Dir)
            //fmt.Print(" "+f.Filename)
            //fmt.Println()
        }
        if len(refFiles) == 0 {
            refFiles = othersFiles
            othersFiles = make([]models.File, 0, 0)
        }
        if len(refFiles) == 0 {
            fmt.Println("The files :")
            for _,f := range toDeleteFiles{
                fmt.Println("\t"+f.Dir+f.Filename)
            }
            fmt.Println("have no copie in any kept dir. The files are not removed")
            continue
        }
        for _,f := range othersFiles{
            toDeleteFiles = append(toDeleteFiles, f)
        }
        if len(toDeleteFiles) == 0{
            fmt.Println("The files :")
            for _,f := range refFiles{
                fmt.Println("\t"+f.Dir+f.Filename)
            }
            fmt.Println("are not in a remove directory. All copies are kept.")
        continue
        }
        //TODO a last byte-to-byte check
        fmt.Print("rm")
        for _,f := range toDeleteFiles{
            fmt.Print(" "+f.Dir+f.Filename)
        }
        fmt.Println()

    }
    return nil
}

