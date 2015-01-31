package helpers


import (
    "os"
    "crypto/sha256"
    "encoding/hex"
    "log"
    "io"
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
        log.Println("Unable to compute TinyHash for \""+path+"\"")
    }else{
        tinyHash := sha256.New()
        tinyHash.Write(data[:numberRead])
        hash = hex.EncodeToString(tinyHash.Sum(nil))
    }
    return hash
}

func ComputeHash(path string) string {
    hash := ""
    hashObject := sha256.New()
    hashError := false
    if file, fileErr := os.Open(path); fileErr != nil {
        hashError =true
    }else{
        defer file.Close()
        _,err := io.Copy(hashObject, file)
        if err != nil {
            hashError = true
        }
    }
    if hashError {
        log.Println("Unable to compute hash for \""+path+"\"")
    }else{
        hash = hex.EncodeToString(hashObject.Sum(nil))
    }
    return hash
}
