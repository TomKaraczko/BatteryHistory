package routes

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/Plaenkler/BatteryHistory/pkg/handler"
	"github.com/Plaenkler/BatteryHistory/pkg/rtls"
	"github.com/Plaenkler/BatteryHistory/pkg/rtls/model"
)

type showPage struct {
	Title string
	MAC   string
	Data  string
}

func ProvideShowPage(writer http.ResponseWriter, request *http.Request) {
	defer handler.HandlePanic("routes")

	log.Printf("[provideshowpage] request on /show/")

	err := request.ParseForm()
	if err != nil {
		fmt.Fprintf(writer, "[provideshowpage] could not parse form - error: %s", err)
		return
	}

	mac := request.Form.Get("device")
	if len(mac) == 0 {
		fmt.Fprintf(writer, "[provideshowpage] device is empty")
		return
	}

	if len(mac) != 17 {
		fmt.Fprintf(writer, "[provideshowpage] malformed mac address")
		return
	}

	manager := rtls.GetManager()

	var history model.BatteryResponse
	err = manager.GetBattery(&history, mac)
	if err != nil {
		fmt.Fprintf(writer, "[provideshowpage] could not get battery history - error: %s", err)
		return
	}

	json, err := json.Marshal(history)
	if err != nil {
		fmt.Fprintf(writer, "[provideshowpage] could not marshal data history - error: %s", err)
		return
	}

	showPage := showPage{
		Title: "Showpage",
	}
	showPage.Data = string(json)

	template, err := template.New("show").ParseFS(static, "static/html/show.html")
	if err != nil {
		fmt.Fprintf(writer, "[provideshowpage] could not provide template - error: %s", err)
		return
	}

	writer.Header().Add("Content-Type", "text/html")
	template.Execute(writer, showPage)
}
