package triggers

import (
	"time"

	"github.com/google/uuid"
)

// TODO: ver si podemos integrar las conditions, y si tiene sentido. Ahora no se estan usando.
type FrequencyTrigger struct {
	ID                string    `json:"id"`
	Name              string    `json:"name"`
	IsActive          string    `json:"active"`
	Condition         Condition `json:"condition"`
	Actions           []Action  `json:"actions"`
	Frequency         time.Duration
	lastExecuteTime   time.Time
	lastEventReceived *Event
	Ticker            *time.Ticker
	isTickerActive    bool
}

func NewFrequencyTrigger(name string, frequency time.Duration) *FrequencyTrigger {
	return &FrequencyTrigger{
		ID:              uuid.New().String(),
		Name:            name,
		Frequency:       frequency,
		lastExecuteTime: time.Now(),
		isTickerActive:  false,
	}
}

func (t *FrequencyTrigger) Update(event Event) {
	t.lastEventReceived = &event
	if time.Since(t.lastExecuteTime) >= t.Frequency && !t.isTickerActive {
		t.execute(event)
	} else {
		if !t.isTickerActive {
			t.isTickerActive = true
			t.Ticker = time.NewTicker(t.Frequency)
			go func() {
				<-t.Ticker.C
				t.execute(*t.lastEventReceived)
				t.isTickerActive = false
			}()
		}
	}
}

func (t *FrequencyTrigger) execute(event Event) {
	t.lastExecuteTime = time.Now()
	for _, action := range t.Actions {
		action.Execute(event)
	}
}

func (t *FrequencyTrigger) GetID() string {
	return t.ID
}

func (t *FrequencyTrigger) SetID(id string) {
	t.ID = id
}
