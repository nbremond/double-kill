package models

import (
    "time"

    _ "github.com/jinzhu/gorm"
)

type File struct {
    Id          int64
    Dir         string `sql:"size:255"`
    Filename    string `sql:"size:255"`
    Size        int64
    TinyHash    string `sql:"size:65"`
    Hash        string `sql:"size:65"`
    UpdatedAt   time.Time
}
