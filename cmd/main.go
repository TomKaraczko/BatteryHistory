package main

import (
	"github.com/Plaenkler/BatteryHistory/pkg/router"
)

func main() {
	router.GetManager().Start()
}
