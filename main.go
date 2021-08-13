/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"github.com/curiosinauts/platformctl/cmd"
	_ "github.com/lib/pq"
	"os"
)

func init() {
	//connStr := GetEnvWithDefault("DB_CONN", "host=192.168.0.105 user=curiosity password=Goldstar114$ dbname=curiosityworks sslmode=disable search_path=curiosity")
	//
	//newdb, err := sqlx.Connect("postgres", connStr)
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//db = newdb

	//db.MustExec(database.Schema)
}

func main() {
	cmd.Execute()
}

// GetEnvWithDefault attempts to retrieve from env. default calculated based on stage if env value empty.
func GetEnvWithDefault(env, defaultV string) string {
	v := os.Getenv(env)
	if v == "" {
		return defaultV
	}
	return v
}
