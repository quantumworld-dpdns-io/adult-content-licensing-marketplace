package audit

import "time"

type Event struct {
	Type      string         `json:"type"`
	ActorSub  string         `json:"actor_sub"`
	ActorRole string         `json:"actor_role"`
	EntityID  string         `json:"entity_id"`
	Meta      map[string]any `json:"meta"`
	At        time.Time      `json:"at"`
}
