package rtls

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"

	"github.com/Plaenkler/BatteryHistory/pkg/config"
	"github.com/Plaenkler/BatteryHistory/pkg/handler"
	"github.com/Plaenkler/BatteryHistory/pkg/rtls/digest"
	"github.com/Plaenkler/BatteryHistory/pkg/rtls/model"
)

var instance *Manager

type Manager struct {
	config *config.Config
	client *http.Client
}

func GetManager() *Manager {
	defer handler.HandlePanic("request")

	if instance == nil {
		client, err := digest.NewTransport(config.GetConfig().ServerUser, config.GetConfig().ServerPassword, 3).Client()
		if err != nil {
			panic(fmt.Errorf("[rtls] creating client failed: %s", err.Error()))
		}

		instance = &Manager{
			config: config.GetConfig(),
			client: client,
		}
	}

	return instance
}

func (m *Manager) GetTags(response *model.TagResponse) error {
	defer handler.HandlePanic("request")

	resp, byteValue, err := m.call("/epe/pos/taglist?fields=all")
	if err != nil {
		return fmt.Errorf("[rtls] get tags failed: %s", err.Error())
	}
	defer resp.Body.Close()

	err = xml.Unmarshal(byteValue, response)
	if err != nil {
		return fmt.Errorf("[rtls] get tags xml unmarshal failed: %s", err.Error())
	}

	return nil
}

func (m *Manager) GetBattery(response *model.BatteryResponse, mac string) error {
	defer handler.HandlePanic("request")

	resp, byteValue, err := m.call("/epe/cfg/batteryhistory?mac=" + mac)
	if err != nil {
		return fmt.Errorf("[rtls] get battery failed: %s", err.Error())
	}
	defer resp.Body.Close()

	err = xml.Unmarshal(byteValue, response)
	if err != nil {
		return fmt.Errorf("[rtls] get battery xml unmarshal failed: %s", err.Error())
	}

	return nil
}

func (m *Manager) call(url string) (*http.Response, []byte, error) {
	defer handler.HandlePanic("request")

	req, err := http.NewRequest(http.MethodGet, "http://"+m.config.ServerAddress+":"+m.config.ServerPort+url, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("[rtls] creating http request failed: %s", err.Error())
	}

	resp, err := m.client.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("[rtls] request failed: %s", err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		return resp, nil, fmt.Errorf("[rtls] http get request failed: %s", resp.Status)
	}

	byteValue, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp, nil, fmt.Errorf("[rtls] reading data failed: %s", err.Error())
	}

	return resp, byteValue, nil
}
