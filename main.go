package main

import (
	"os"

	"github.com/curiosinauts/platformctl/cmd"
	_ "github.com/lib/pq"
)

func init() {
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
