package main

import (
	"math/rand"
	"time"

	"github.com/dustywilson/milight"
)

func main() {
	addr := "10.10.10.104:8899"

	milight.Send(addr,
		milight.TurnOn(milight.ZoneAll),
		milight.SetBrightness(0.1),
	)

	for {
		red := rand.Intn(256)
		green := rand.Intn(256)
		blue := rand.Intn(256)
		milight.Send(addr,
			milight.SetColorRGB(red, green, blue),
		)
		time.Sleep(250 * time.Millisecond)
	}
}
