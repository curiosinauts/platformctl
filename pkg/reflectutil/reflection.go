package reflectutil

import (
	"reflect"
	"strings"
)

// ListDBTagsFor lists go meta tags with prefix db:
func ListDBTagsFor(i interface{}) []string {
	var dbTags []string

	tags := listTagsFor(i)
	for _, tag := range tags {
		if strings.Contains(tag, "db:") {
			terms := strings.Fields(tag)
			for _, term := range terms {
				if strings.HasPrefix(term, "db:") {
					s := strings.Replace(term, "db:", "", -1)
					s = strings.Replace(s, "\"", "", -1)
					if s == "id" {
						continue
					}
					dbTags = append(dbTags, s)
				}
			}
		}
	}

	return dbTags
}

func listTagsFor(i interface{}) []string {
	var tags []string

	t := reflect.TypeOf(i).Elem()
	v := reflect.ValueOf(i).Elem()

	for i := 0; i < v.NumField(); i++ {
		tag := string(t.Field(i).Tag)
		if len(tag) > 0 {
			tags = append(tags, tag)
		}
	}

	return tags
}
