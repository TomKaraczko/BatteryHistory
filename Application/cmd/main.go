package main

import (
	"github.com/Plaenkler/BatteryHistory/pkg/router"
)

var (
	rtManager *router.Manager
)

func main() {
	rtManager = router.GetManager()
	rtManager.Start()
}
