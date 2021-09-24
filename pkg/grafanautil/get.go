package grafanautil

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/spf13/viper"
)

// DownloadPanel retrieves the resource using given URL
func DownloadPanel(panelID, width, height, from int, uuid string) error {
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

	filename := fmt.Sprintf("/tmp/%s.png", uuid)
	err = ioutil.WriteFile(filename, data, 0755)

	return err
}
