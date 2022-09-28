package router

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"time"

	"github.com/Plaenkler/BatteryHistory/pkg/config"
	"github.com/Plaenkler/BatteryHistory/pkg/handler"
	"github.com/Plaenkler/BatteryHistory/pkg/router/routes"
)

var (
	//go:embed routes/static
	static   embed.FS
	instance *Manager
)

type Manager struct {
	Router *http.ServeMux
	config *config.Config
}

func GetManager() *Manager {
	defer handler.HandlePanic("router")

	if instance == nil {
		instance = &Manager{
			config: config.GetConfig(),
		}
	}

	return instance
}

func (manager *Manager) Start() {
	defer handler.HandlePanic("router")

	manager.Router = http.NewServeMux()

	manager.Router.HandleFunc("/",
		routes.ProvideHomePage)

	manager.Router.HandleFunc("/show/",
		routes.ProvideShowPage)

	err := manager.provideFiles(); 
	if err != nil {
		log.Panicf("[router] could not provide files - error: %s", err)
	}
	
	server := &http.Server{
		Addr: ":"+manager.config.WebPort,
		ReadHeaderTimeout: 3* time.Second,
		Handler: manager.Router,
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Panicf("[router] could not listen and serve - error: %s", err)
	}
}

func (manager *Manager) provideFiles() error {
	defer handler.HandlePanic("router")

	css, err := fs.Sub(static, "routes/static")
	if err != nil {
		return err
	}

	manager.Router.Handle("/css/", http.FileServer(http.FS(css)))

	return nil
}
