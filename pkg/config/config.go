package config

import (
	"time"
)

type Config struct {
	Options  Options
	Bulbs    []Bulb
	Canvases []Canvas
	Presets  []Preset
}

type Options struct {
	Debug bool
}

type Bulb struct {
	Host string
	Name string
	Zone Zone
}

type Zone struct {
	Min   *uint8
	Max   *uint8
	Limit int
}

type Canvas struct {
	Name  string
	Bulbs []Bulb
}

type Preset struct {
	Name   string
	Canvas string
	Delay  time.Duration
	Effect Effect
}

type Effect struct {
	Name   string
	Config map[string]interface{}
}
