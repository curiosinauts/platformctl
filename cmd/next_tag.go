package cmd

import (
	"github.com/spf13/cobra"
)

// tagCmd represents the tag command
var tagCmd = &cobra.Command{
	Use:   "tag",
	Short: "Generates next docker tag",
	Long:  `Generates next docker tag`,
	Run: func(cmd *cobra.Command, args []string) {
		// r = requests.get(f'https://docker-registry.int.curiosityworks.org/v2/{image}/tags/list' )

		// if r.status_code != 200:
		// print('1.0.0')
		// sys.exit(0)

		// repo = json.loads(r.text)

		// if not repo['tags']:
		// print('1.0.0')
		// sys.exit(0)

		// tags = repo['tags']
		// # tags = ['30.0.12', '30.0.13', '1.0.10', '1.0.11', 'latest']

		// tt = []
		// for tag in tags:
		// if tag == "latest":
		// 	continue
		// tag_list = tag.split('.')
		// tag_int_list = [int(i) for i in tag_list]
		// tag_tuple = tuple(tag_int_list)
		// tt.append(tag_tuple)

		// tt = sorted(tt, key = lambda x: (x[0], x[1], x[2]))

		// current_tag_tuple = tt[-1]

		// next_tag = f'{current_tag_tuple[0]}.{current_tag_tuple[1]}.{current_tag_tuple[2]+1}'

		// print(next_tag)%
	},
}

func init() {
	addCmd.AddCommand(tagCmd)
}
