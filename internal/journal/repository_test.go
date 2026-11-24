package journal

import (
	"errors"
	"testing"
)

func TestProcedureName(t *testing.T) {
	cases := map[string]string{
		ActionLogin:  "log_authorization",
		ActionCreate: "log_create",
		ActionUpdate: "log_update",
		ActionDelete: "log_delete",
	}
	for action, want := range cases {
		got, err := procedureName(action)
		if err != nil || got != want {
			t.Fatalf("%s: got %s, %v", action, got, err)
		}
	}
	if _, err := procedureName("drop"); !errors.Is(err, ErrInvalid) {
		t.Fatalf("unknown action: %v", err)
	}
}
