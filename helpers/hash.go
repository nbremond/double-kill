package helpers


import (
    "os"
    "crypto/sha256"
    "encoding/hex"
    "log"
)


const tinyHashSize = 10 * 1024

func ComputeTinyHash(path string) string {
    hash := ""
    hashError := false
    data := make([]byte, tinyHashSize)
    var numberRead int
    if file, fileErr := os.Open(path); fileErr != nil {
        hashError =true
    }else{
        defer file.Close()
        var readErr error
        numberRead, readErr = file.Read(data)
        if  readErr != nil {
            hashError = true
        }
    }
    if hashError {
        log.Println("Unable to compute hash for \""+path+"\"")
    }else{
        tinyHash := sha256.New()
        tinyHash.Write(data[:numberRead])
        hash = hex.EncodeToString(tinyHash.Sum(nil))
    }
    return hash
}
