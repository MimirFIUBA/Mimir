package models

import "time"

type Alert struct {
	message string
	read    bool
	time    time.Time
	data    interface{}
}
