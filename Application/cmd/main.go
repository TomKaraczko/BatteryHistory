package main

import (
	"plaenkler.com/avanis/rtls-battery-monitor/pkg/router"
)

var (
	rtManager *router.Manager
)

func main() {
	rtManager = router.GetManager()
	rtManager.Start()
}
