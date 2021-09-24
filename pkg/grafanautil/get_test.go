package grafanautil

import (
	"fmt"
	"testing"

	"github.com/curiosinauts/platformctl/pkg/testutil"
	"github.com/google/uuid"
)

func init() {
	testutil.InitConfig()
}

func TestDownload(t *testing.T) {
	id := uuid.NewString()
	err := DownloadPanel(15, 600, 300, 1, id)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(id)
}
