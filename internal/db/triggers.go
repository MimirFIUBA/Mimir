package db

import "mimir/triggers"

func GetTriggers() []triggers.TriggerObserver {
	var triggerList []triggers.TriggerObserver
	for _, sensor := range SensorsData.sensors {
		triggerList = append(triggerList, sensor.GetTriggers()...)
	}
	return triggerList
}
