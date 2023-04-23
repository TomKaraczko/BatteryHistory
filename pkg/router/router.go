package router

import (
	"embed"
	"fmt"
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
		Addr:              fmt.Sprintf(":%v", config.GetConfig().Port),
		ReadTimeout:       15 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       120 * time.Second,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Read X-Forwarded-For-Header for proxying
			remoteIP := r.Header.Get("X-Forwarded-For")
			remotePort := r.Header.Get("X-Forwarded-Port")
			if remotePort == "" {
				remotePort = "80"
			}
			r.URL.Scheme = "http"
			r.URL.Host = fmt.Sprintf("%s:%s", remoteIP, remotePort)
			manager.Router.ServeHTTP(w, r)
		}),
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

	return nil
}
