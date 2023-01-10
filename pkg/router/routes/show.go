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
	MAC  string
	Data string
}

func ProvideShowPage(writer http.ResponseWriter, request *http.Request) {
	defer handler.HandlePanic("routes")

	log.Printf("[provide homepage] request on %s - by address: %s", request.RequestURI, request.RemoteAddr)

	if request.Method != http.MethodPost && request.Method != http.MethodGet {
		writer.WriteHeader(http.StatusForbidden)
		_, err := writer.Write([]byte("405 - Method not allowed"))
		if err != nil {
			log.Printf("[provide showpage] could not write http reply - error: %s", err)
		}
		return
	}

	err := request.ParseForm()
	if err != nil {
		fmt.Fprintf(writer, "[provide showpage] could not parse form - error: %s", err)
		return
	}

	mac := request.Form.Get("device")
	if len(mac) == 0 {
		fmt.Fprintf(writer, "[provide showpage] device is empty")
		return
	}

	if len(mac) != 17 {
		fmt.Fprintf(writer, "[provide showpage] malformed mac address")
		return
	}

	manager := rtls.GetManager()

	var history model.BatteryResponse
	err = manager.GetBattery(&history, mac)
	if err != nil {
		fmt.Fprintf(writer, "[provide showpage] could not get battery history - error: %s", err)
		return
	}

	json, err := json.Marshal(history)
	if err != nil {
		fmt.Fprintf(writer, "[provide showpage] could not marshal data history - error: %s", err)
		return
	}

	showPage := showPage{
		Data: string(json),
	}

	// nolint: typecheck
	template, err := template.New("show").ParseFS(static,
		"static/html/pages/show.html",
		"static/html/partials/include.html",
		"static/html/partials/chart.html",
	)
	if err != nil {
		fmt.Fprintf(writer, "[provide showpage] could not provide template - error: %s", err)
		return
	}

	writer.Header().Add("Content-Type", "text/html")
	err = template.Execute(writer, showPage)
	if err != nil {
		log.Panicf("[provide showpage] could not execute parsed template - error: %s", err)
	}
}
