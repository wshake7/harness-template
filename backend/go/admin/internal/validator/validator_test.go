package validator

import "testing"

func TestNotblankValidation(t *testing.T) {
	type req struct {
		Name     string  `binding:"notblank"`
		Nickname *string `binding:"omitempty,notblank"`
	}

	nickname := "  "
	if err := Struct(req{Name: "alice"}); err != nil {
		t.Fatalf("expected non-blank string with nil optional pointer to pass: %v", err)
	}
	if err := Struct(req{Name: "  "}); err == nil {
		t.Fatal("expected blank string to fail")
	}
	if err := Struct(req{Name: "alice", Nickname: &nickname}); err == nil {
		t.Fatal("expected blank string pointer to fail")
	}
}
