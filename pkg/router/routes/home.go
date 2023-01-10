package routes

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/Plaenkler/BatteryHistory/pkg/handler"
)

var (
	//go:embed static
	static embed.FS
)

func ProvideHomePage(writer http.ResponseWriter, request *http.Request) {
	defer handler.HandlePanic("routes")

	log.Printf("[provide homepage] request on %s - by address: %s", request.RequestURI, request.RemoteAddr)

	if request.URL.Path != "/" {
		http.Redirect(writer, request, "/", http.StatusSeeOther)
		return
	}

	if request.Method != http.MethodGet {
		writer.WriteHeader(http.StatusForbidden)
		_, err := writer.Write([]byte("405 - Method not allowed"))
		if err != nil {
			log.Printf("[provide homepage] could not write http reply - error: %s", err)
		}
		return
	}

	// nolint: typecheck
	template, err := template.New("home").ParseFS(static,
		"static/html/pages/home.html",
		"static/html/partials/include.html",
	)
	if err != nil {
		fmt.Fprintf(writer, "[provide homepage] could not provide template - error: %s", err)
		return
	}

	writer.Header().Add("Content-Type", "text/html")
	err = template.Execute(writer, nil)
	if err != nil {
		log.Panicf("[provide homepage] could not execute parsed template - error: %s", err)
	}
}
