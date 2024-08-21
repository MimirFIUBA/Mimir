package trigger

import "time"

type Event struct {
	Name      string
	Timestamp time.Time
	Data      interface{}
}
