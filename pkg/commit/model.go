package commit

import "time"

type FileStat struct {
	Insert int
	Delete int
	Ext    string
}

type Diff struct {
	Changes []Change
}

type Change struct {
	Dir       string
	Filename  string
	Text      string
	Insertion int
	Deletion  int
}

type Meta struct {
	Time time.Time
	Hash string
}
