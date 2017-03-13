package listener

import (
	"encoding/json"
	"time"
)

type TransitionMessage struct {
	Wtst  int64   // Time of waypoint creation
	Lat   float32 // Latitude
	Long  float32 // Longitude
	Tst   int64   // Timestamp of transition
	Acc   uint32  // Accuracy of Lat/Long
	Tid   string  // Tracker ID
	Event string  // Enter or Leave
	Desc  string  // Description
}

func NewTransitionMessage(payload []byte) (*TransitionMessage, error) {
	tm := TransitionMessage{}
	err := json.Unmarshal([]byte(payload), &tm)
	return &tm, err
}

func (tm *TransitionMessage) Timestamp() time.Time {
	return time.Unix(tm.Tst, 0)
}
