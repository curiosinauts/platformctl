package grafanautil

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	filepath := fmt.Sprintf("/tmp/%s.png", uuid)
	err = ioutil.WriteFile(filepath, data, 0755)
	if err != nil {
		return err
	}
	defer os.Remove(filepath)

	message, err := io.ByteStreamFileUpload("http://192.168.0.118:8080/stream-upload", "/var/www/fileserver.curiosityworks.org/htdocs", filepath)
	if debug {
		log.Println(message)
	}

	return err
}

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
