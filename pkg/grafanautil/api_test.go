package grafanautil

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/curiosinauts/platformctl/pkg/testutil"
	"github.com/google/uuid"
	"github.com/grafana-tools/sdk"
	"github.com/spf13/viper"
)

func init() {
	testutil.InitConfig()
}

func TestDownload(t *testing.T) {
	id := uuid.NewString()
	err := DownloadPanel(15, 600, 300, 1, id, true)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(id)
}

func TestSearch(t *testing.T) {
	grafanaAPIKey := viper.Get("grafana_api_key").(string)

	api, _ := sdk.NewClient("https://grafana.int.curiosityworks.org", grafanaAPIKey, http.DefaultClient)

	ctx := context.Background()
	foundBoards, err := api.SearchDashboards(ctx, "", false, []string{""}...)

	if err != nil {
		t.Fatalf("SearchDashboards test failed: %s", err)
	}
	for _, foundboard := range foundBoards {

		board, _, _ := api.GetDashboardByUID(ctx, foundboard.UID)

		data, _ := json.Marshal(&board)
		fmt.Println(string(data))
	}

}
