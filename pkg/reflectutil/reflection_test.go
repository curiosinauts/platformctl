package reflectutil

import (
	"testing"
)

func Test(t *testing.T) {
	var person = struct {
		Name    string `db:"name"`
		Age     int    `db:"age"`
		Address string `json:"address"`
	}{
		"John",
		20,
		"1 foo ave",
	}

	dbTags := ListDBTagsFor(&person)
	if len(dbTags) != 2 {
		t.Error("two db tags are expected")
	}
}
