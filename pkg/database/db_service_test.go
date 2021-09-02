package database

import (
	"fmt"
	"testing"
)

func TestGetMappingConfigFromSlicePointer(t *testing.T) {
	users := &[]User{
		User{Username: "foo"},
		User{Username: "bar"},
	}

	meta := GetMappingConfigFromSlicePointer(users)
	if meta != nil {
		fmt.Println(meta.TableName)
	}

}
