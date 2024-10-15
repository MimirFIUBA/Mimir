package db

import "mimir/triggers"

func GetTriggers() []triggers.TriggerObserver {
	var triggerList []triggers.TriggerObserver
	for _, sensor := range SensorsData.sensors {
		for _, trigger := range sensor.GetTriggers() {
			triggerList = append(triggerList, trigger)
		}
	}
	return triggerList
}
