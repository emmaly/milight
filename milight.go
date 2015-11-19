package milight

import (
	"net"
	"time"

	"github.com/lucasb-eyer/go-colorful"
)

// Zone constants are used to declare the target zone
const (
	ZoneAll int = iota
	Zone1
	Zone2
	Zone3
	Zone4
)

var (
	trailer byte = 0x55

	on    = [][]byte{{0x42, 0x00}, {0x45, 0x00}, {0x47, 0x00}, {0x49, 0x00}, {0x4B, 0x00}}
	off   = [][]byte{{0x41, 0x00}, {0x46, 0x00}, {0x48, 0x00}, {0x4A, 0x00}, {0x4C, 0x00}}
	night = [][]byte{{0xC1, 0x00}, {0xC6, 0x00}, {0xC8, 0x00}, {0xCA, 0x00}, {0xCC, 0x00}}
	white = [][]byte{{0xC2, 0x00}, {0xC5, 0x00}, {0xC7, 0x00}, {0xC9, 0x00}, {0xCB, 0x00}}

	discoModeOn     = []byte{0x4D, 0x00}
	discoModeFaster = []byte{0x44, 0x00}
	discoModeSlower = []byte{0x43, 0x00}

	brightnessPrefix   byte = 0x4E
	brightnessValueMin byte = 0x02
	brightnessValueMax byte = 0x1B

	colorPrefix byte = 0x40
)

// Send sends the command to the light controller
func Send(addr string, values ...[]byte) error {
	a, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return err
	}
	conn, err := net.DialUDP("udp", nil, a)
	if err != nil {
		return err
	}
	defer conn.Close()
	for i1, value := range values {
		for i := 0; i < len(value); i += 2 {
			if i+i1 > 0 {
				time.Sleep(200 * time.Millisecond)
			}
			_, err = conn.Write([]byte{value[i], value[i+1], trailer})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// TurnOff turns off the lights in the selected zone(s)
func TurnOff(zone int) []byte {
	return off[zone]
}

// TurnOn turns on the lights in the selected zone(s) without changing their in-bulb values
func TurnOn(zone int) []byte {
	return on[zone]
}

// SetWhite sets the lights in the selected zone(s) to full white
func SetWhite(zone int) []byte {
	return append(white[zone], white[zone]...) // apparently have to turn the bulb on first before setting WHITE
}

// SetNight sets the lights in the selected zone(s) to night mode
func SetNight(zone int) []byte {
	return append(night[zone], night[zone]...) // apparently have to turn the bulb on first before setting WHITE, so we'll do it for night too
}

// SetBrightness sets the brightness for the lights in the already selected zone(s)
func SetBrightness(brightness float64) []byte {
	value := byte(brightness*float64(brightnessValueMax-brightnessValueMin)) + brightnessValueMin
	return []byte{brightnessPrefix, value}
}

// SetColorRGB sets the color for the selected zone(s), using int 0-255 as the RGB values
func SetColorRGB(r, g, b int) []byte {
	return SetColorRGBFloat(float64(r)/255, float64(g)/255, float64(b)/255)
}

// SetColorRGBHex sets the color for the selected zone(s), using #000000 as the RGB values
func SetColorRGBHex(rgb string) []byte {
	hex, _ := colorful.Hex(rgb)
	return setColor(hex)
}

// SetColorRGBFloat sets the color for the selected zone(s), using float64 0..1 as the RGB values
func SetColorRGBFloat(r, g, b float64) []byte {
	return setColor(colorful.Color{r, g, b})
}

func setColor(color colorful.Color) []byte {
	hue, _, _ := color.Hsv()
	value := byte(int(256+176-(hue/360*255)) % 256) // math stolen from https://goo.gl/rMwABR
	return []byte{colorPrefix, value}
}
