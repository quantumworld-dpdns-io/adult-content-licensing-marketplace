package audit

import "testing"

func TestStreamAppendList(t *testing.T) {
	t.Parallel()
	s := NewStream()
	s.Append(Event{Type: "license.created", EntityID: "lic_1"})
	all := s.List()
	if len(all) != 1 {
		t.Fatalf("expected 1 event, got %d", len(all))
	}
	if all[0].Type != "license.created" {
		t.Fatalf("unexpected event type: %s", all[0].Type)
	}
}
