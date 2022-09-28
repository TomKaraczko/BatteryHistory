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

type homePage struct {
	Title string
}

func ProvideHomePage(writer http.ResponseWriter, request *http.Request) {
	defer handler.HandlePanic("routes")

	log.Printf("[providehomepage] request on /")

	template, err := template.New("home").ParseFS(static, "static/html/home.html")
	if err != nil {
		fmt.Fprintf(writer, "[providehomepage] could not provide template - error: %s", err)
	}

	writer.Header().Add("Content-Type", "text/html")
	err = template.Execute(writer, homePage{
		Title: "Homepage",
	})
	if err != nil {
		log.Panicf("[providehomepage] could not execute parsed template - error: %s", err)
	}
}
