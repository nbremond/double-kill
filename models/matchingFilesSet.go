package models


type MatchingFilesSet struct {
    Id          int64
    Size        int64
    TinyHash    string `sql:"size:65"`
    Hash        string `sql:"size:65"`
    Files       []File
}

func init() {
    register(&MatchingFilesSet{})
}

func (m *MatchingFilesSet) Save() {
    db.Save(m)
}
