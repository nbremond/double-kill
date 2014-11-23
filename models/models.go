package models

import(
    orm "github.com/jinzhu/gorm"
    _ "github.com/mattn/go-sqlite3"

    "github.com/nbremond/double-kill/modules/settings"
)

var (
    db orm.DB
    tables []interface{} = make([]interface{}, 0)
)

func register(model interface{}) {
    tables = append(tables, model)
}

func InitDB() error {
    db, err := orm.Open(settings.DB.Engine, settings.DB.Source)
    if err != nil {
        return err
    }
    db.DB()

    db.LogMode(true)
    db.AutoMigrate(tables...)
    return nil
}
