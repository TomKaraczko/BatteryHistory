package model

type TagResponse struct {
	Tags []Tag `xml:"TAG"`
}

type Tag struct {
	Mac string `xml:"mac"`
}

type BatteryResponse struct {
	Batteries []Battery `xml:"BATTERY"`
}

type Battery struct {
	Load uint8  `xml:"battery"`
	Time uint64 `xml:"time"`
}
