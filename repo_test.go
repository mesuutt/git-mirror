package main

import "testing"

func TestRepo_AddChange(t *testing.T) {
	// TODO: use afero
	repo := NewRepo("/tmp/git_activity/")

	err := repo.AddStat(FileStat{
		Insert: 1,
		Delete: 2,
		Ext:    "go",
	})

	if err != nil {
		t.Fatal(err)
	}
}
