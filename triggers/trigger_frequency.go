package triggers

import (
	"time"

	"github.com/google/uuid"
)

// TODO: ver si podemos integrar las conditions, y si tiene sentido. Ahora no se estan usando.
type FrequencyTrigger struct {
	ID                string
	Name              string
	IsActive          bool
	Condition         Condition
	Actions           []Action
	Frequency         time.Duration
	lastExecuteTime   time.Time
	lastEventReceived *Event
	Ticker            *time.Ticker
	isTickerActive    bool
	observedSubjects  []Subject
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

func (t *FrequencyTrigger) UpdateCondition(newCondition string) error {
	condition, err := BuildConditionFromString(newCondition)
	if err != nil {
		return err
	}
	t.Condition = condition
	return nil
}

func (t *FrequencyTrigger) UpdateActions(actions []Action, _ TriggerOptions) error {
	t.Actions = actions
	return nil
}

func (t *FrequencyTrigger) AddAction(a Action, _ TriggerOptions) {
	t.Actions = append(t.Actions, a)
}

func (t *FrequencyTrigger) AddSubject(subject Subject) {
	t.observedSubjects = append(t.observedSubjects, subject)
}

func (t *FrequencyTrigger) StopWatching() {
	for _, subject := range t.observedSubjects {
		subject.Deregister(t)
	}
}

func (t *FrequencyTrigger) Activate() {
	t.IsActive = true
}

func (t *FrequencyTrigger) Deactivate() {
	t.IsActive = false
}

func (t *FrequencyTrigger) GetType() TriggerType {
	return FREQUENCY_TRIGGER
}

func (t *FrequencyTrigger) SetStatus(active bool) {
	if t.IsActive != active {
		if active {
			t.Activate()
		} else {
			t.Deactivate()
		}
	}
}

func (t *FrequencyTrigger) SetScheduled(scheduled bool) {
	panic("cannot schedule frequency trigger")
}

func (t *FrequencyTrigger) GetName() string {
	return t.Name
}
