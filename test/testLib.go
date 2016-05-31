package main

import (
	"../lib"
	"fmt"
)

func TestLib() {
	var trackingData lib.TrackData
	trackingData.Kingdom = "test"
	trackingData.Phylum = "success"
	go lib.TrackCounter(trackingData, 3)
	info, _ := lib.TrackCounter(trackingData, 3)
	fmt.Println("Change info:", info)
}
