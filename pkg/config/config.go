package config

import (
	"time"
)

type Light struct {
	Host  string
	Zones int
}

type Scene struct {
	Effect   string
	Duration time.Duration
}

type Effect struct {

}

type Config struct {
	Lights []Light
}
