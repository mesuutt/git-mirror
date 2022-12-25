package main

import "testing"

func TestRepo_AddChange(t *testing.T) {
	// TODO: use afero
	repo := NewRepo("/tmp/git_activity/")

	err := repo.AddChange(".go", 100)

	if err != nil {
		t.Fatal(err)
	}
}
