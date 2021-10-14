package grafanautil

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/curiosinauts/platformctl/pkg/io"
	"github.com/grafana-tools/sdk"
	"github.com/spf13/viper"
)

// DownloadPanel retrieves the resource using given URL
func DownloadPanel(panelID, width, height, from int, uuid string, debug bool) error {
	url := fmt.Sprintf("https://grafana.int.curiosityworks.org/render/d-solo/7UdvG-Mnk/base-system-health?"+
		"orgId=1&panelId=%d&width=%d&height=%d&tz=America/New_York"+
		"&from=now-%dh&to=now", panelID, width, height, from)

	req, _ := http.NewRequest("GET", url, nil)
	grafanaAPIKey := viper.Get("grafana_api_key").(string)
	req.Header.Set("Authorization", "Bearer "+grafanaAPIKey)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode > 200 {
		return errors.New("issue with downloading")
	}

	filename := fmt.Sprintf("%s.png", uuid)

	message, err := io.ByteStreamFileUpload("http://192.168.0.118:8080/stream-upload",
		"/var/www/fileserver.curiosityworks.org/htdocs", filename, res.Body)

	if debug {
		log.Println(message)
	}

	return err
}

// ListDashboards lists dashboards
func ListDashboards(query string) ([]sdk.FoundBoard, error) {
	grafanaAPIKey := viper.Get("grafana_api_key").(string)

	api, _ := sdk.NewClient("https://grafana.int.curiosityworks.org", grafanaAPIKey, http.DefaultClient)

	ctx := context.Background()
	foundBoards, err := api.SearchDashboards(ctx, query, false, []string{}...)
	if err != nil {
		return []sdk.FoundBoard{}, err
	}
	return foundBoards, err
}

// ListPanels lists panels
func ListPanels(uid string, partialPanelTitle string) ([]*sdk.Panel, error) {
	grafanaAPIKey := viper.Get("grafana_api_key").(string)

	api, _ := sdk.NewClient("https://grafana.int.curiosityworks.org", grafanaAPIKey, http.DefaultClient)

	ctx := context.Background()
	board, _, _ := api.GetDashboardByUID(ctx, uid)

	panelsMatched := []*sdk.Panel{}
	for _, panel := range board.Panels {
		if strings.Contains(strings.ToLower(panel.Title), strings.ToLower(partialPanelTitle)) {
			panelsMatched = append(panelsMatched, panel)
		}
	}

	return panelsMatched, nil
}
