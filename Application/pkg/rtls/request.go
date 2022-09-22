package rtls

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/bobziuchkovski/digest"
	"gorm.io/gorm"

	"github.com/Plaenkler/BatteryHistory/pkg/config"
	"github.com/Plaenkler/BatteryHistory/pkg/handler"
	"github.com/Plaenkler/BatteryHistory/pkg/rtls/model"
)

var instance *Manager

type Manager struct {
	DB     *gorm.DB
	config *config.Config
}

func GetManager() *Manager {
	defer handler.HandlePanic("request")

	if instance == nil {
		instance = &Manager{
			config: config.GetConfig(),
		}
	}

	return instance
}

func (m *Manager) GetTags(response *model.TagResponse) error {
	defer handler.HandlePanic("request")

	_, byteValue, err := m.call("/epe/pos/taglist?fields=all")
	if err != nil {
		return fmt.Errorf("[rtls] get tags failed: %s", err.Error())
	}

	err = xml.Unmarshal(byteValue, response)
	if err != nil {
		return fmt.Errorf("[rtls] xml unmarshal failed: %s", err.Error())
	}

	return nil
}

func (m *Manager) GetBattery(response *model.BatteryResponse, mac string) error {
	defer handler.HandlePanic("request")

	_, byteValue, err := m.call("/epe/cfg/batteryhistory?mac=" + mac)
	if err != nil {
		return fmt.Errorf("[rtls] get battery failed: %s", err.Error())
	}

	err = xml.Unmarshal(byteValue, response)
	if err != nil {
		return fmt.Errorf("[rtls] xml unmarshal failed: %s", err.Error())
	}

	return nil
}

func (m *Manager) call(url string) (*http.Response, []byte, error) {
	defer handler.HandlePanic("request")

	t := digest.NewTransport(m.config.ServerUser, m.config.ServerPassword)

	prefix := "http://"

	req, err := http.NewRequest("GET", prefix+m.config.ServerAddress+":"+m.config.ServerPort+url, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("[rtls] creating http request failed: %s", err.Error())
	}

	resp, err := t.RoundTrip(req)
	if err != nil {
		return nil, nil, fmt.Errorf("[rtls] request_failed: %s", err.Error())
	}

	if resp.StatusCode != 200 {
		return resp, nil, fmt.Errorf("[rtls] http_get_request_failed: %s", resp.Status)
	}

	byteValue, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp, nil, fmt.Errorf("[rtls] reading_data_failed: %s", err.Error())
	}

	return resp, byteValue, nil
}
