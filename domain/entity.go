package domain

import "time"

type Measurement struct {
	Ts        time.Time `json:"ts"`
	ActorName string    `json:"actor_name"`
	Type      string    `json:"type"`
	Data      string    `json:"data"`
}
