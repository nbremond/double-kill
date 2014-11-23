package settings

import ()

func init() {
    DB.Engine = "sqlite3"
    DB.Source = "/tmp/test.db"
}

var (
    Version string
    DB struct {
        Engine string
        Source string
    }
)
