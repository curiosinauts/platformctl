package cmd

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/curiosinauts/platformctl/pkg/regutil"
	"github.com/spf13/cobra"
)

var nextTagCmdDebug bool

// nextTagCmd represents the tag command
var nextTagCmd = &cobra.Command{
	Use:     "tag",
	Short:   "Generates next docker tag",
	Long:    `Generates next docker tag`,
	PreRunE: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		repository := args[0]

		eh := ErrorHandler{"listing tags"}
		tags, err := regutil.ListTags(repository, nextTagCmdDebug)
		eh.HandleError("listing tags", err)

		//tags = []string{"30.0.12", "30.0.13", "1.0.10", "1.0.11", "latest"}

		if nextTagCmdDebug {
			fmt.Println(tags)
		}

		versions := []SemanticVersion{}
		for _, versionStr := range tags {
			sv, err := NewSemanticVersion(versionStr)
			if err != nil {
				continue
			}
			versions = append(versions, sv)
		}

		sort.Slice(versions, func(i, j int) bool {
			lessMajor := versions[i].Major < versions[j].Major
			if lessMajor {
				return true
			}
			if versions[i].Major == versions[j].Major {
				lessMinor := versions[i].Minor < versions[j].Minor
				if lessMinor {
					return true
				}
				if versions[i].Minor == versions[j].Minor {
					lessPatch := versions[i].Patch < versions[j].Patch
					if lessPatch {
						return true
					}
				}
			}
			return false
		})

		if nextTagCmdDebug {
			fmt.Println(versions)
		}

		var nextTag SemanticVersion
		if len(versions) > 0 {
			lastTag := versions[len(versions)-1]
			nextTag = lastTag
			nextTag.Patch = nextTag.Patch + 1
		}

		fmt.Println(nextTag.String())
	},
}

func init() {
	nextCmd.AddCommand(nextTagCmd)
	nextTagCmd.Flags().BoolVarP(&nextTagCmdDebug, "debug", "d", false, "Debug this command")
}

// SemanticVersion is a struct for semantic version
type SemanticVersion struct {
	Major int
	Minor int
	Patch int
}

// NewSemanticVersion creates a new instance of semantic version
func NewSemanticVersion(version string) (SemanticVersion, error) {
	if strings.HasSuffix(version, "v") {
		version = version[1:]
	}
	terms := strings.Split(version, ".")
	sv := SemanticVersion{}
	if len(terms) != 3 {
		return sv, errors.New("invalid semantic versioning schema")
	}
	major, err := strconv.Atoi(terms[0])
	if err != nil {
		return sv, errors.New("invalid semantic versioning schema")
	}
	minor, err := strconv.Atoi(terms[1])
	if err != nil {
		return sv, errors.New("invalid semantic versioning schema")
	}
	patch, err := strconv.Atoi(terms[2])
	if err != nil {
		return sv, errors.New("invalid semantic versioning schema")
	}
	sv.Major = major
	sv.Minor = minor
	sv.Patch = patch

	return sv, nil
}

func (sv SemanticVersion) String() string {
	return fmt.Sprintf("%d.%d.%d", sv.Major, sv.Minor, sv.Patch)
}
