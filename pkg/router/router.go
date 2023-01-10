package router

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/Plaenkler/BatteryHistory/pkg/config"
	"github.com/Plaenkler/BatteryHistory/pkg/handler"
	"github.com/Plaenkler/BatteryHistory/pkg/router/routes"
)

var (
	//go:embed routes/static
	static   embed.FS
	instance *Manager
	once     sync.Once
)

type Manager struct {
	Router *http.ServeMux
}

func GetManager() *Manager {
	defer handler.HandlePanic("router")

	once.Do(func() {
		instance = &Manager{
			Router: http.NewServeMux(),
		}
	})

	return instance
}

func (manager *Manager) Start() {
	defer handler.HandlePanic("router")

	manager.Router.HandleFunc("/",
		routes.ProvideHomePage)

	manager.Router.HandleFunc("/show",
		routes.ProvideShowPage)

	err := manager.provideFiles()
	if err != nil {
		log.Panicf("[router] could not provide files - error: %s", err)
	}

	server := &http.Server{
		Addr:              ":" + config.GetConfig().Port,
		ReadHeaderTimeout: 3 * time.Second,
		Handler:           manager.Router,
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Panicf("[router] could not listen and serve - error: %s", err)
	}
}

func (manager *Manager) provideFiles() error {
	defer handler.HandlePanic("router")

	fs, err := fs.Sub(static, "routes/static")
	if err != nil {
		return err
	}

	manager.Router.Handle("/css/", http.FileServer(http.FS(fs)))
	manager.Router.Handle("/img/", http.FileServer(http.FS(fs)))
	manager.Router.Handle("/js/", http.FileServer(http.FS(fs)))

	return nil
}
