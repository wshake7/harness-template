package middleware

import (
	"strings"
	"testing"
)

func TestSanitizeApiLogPayloadRedactsJSON(t *testing.T) {
	payload := `{"authorization":"Bearer abc","nested":{"password":"secret","ok":"value"},"token":"abc"}`

	result := sanitizeApiLogPayload(payload)

	if strings.Contains(result, "Bearer abc") || strings.Contains(result, "secret") || strings.Contains(result, `"token":"abc"`) {
		t.Fatalf("expected sensitive values to be redacted, got %s", result)
	}
	if !strings.Contains(result, `"ok":"value"`) {
		t.Fatalf("expected non-sensitive value to remain, got %s", result)
	}
}

func TestSanitizeApiLogPayloadRedactsText(t *testing.T) {
	payload := `authorization=Bearer abc password=secret normal=value`

	result := sanitizeApiLogPayload(payload)

	if strings.Contains(result, "Bearer abc") || strings.Contains(result, "secret") {
		t.Fatalf("expected sensitive text values to be redacted, got %s", result)
	}
	if !strings.Contains(result, "normal=value") {
		t.Fatalf("expected normal text value to remain, got %s", result)
	}
}

func TestSanitizeApiLogPayloadTruncates(t *testing.T) {
	payload := strings.Repeat("a", maxApiLogFieldLength+1)

	result := sanitizeApiLogPayload(payload)

	if len(result) <= maxApiLogFieldLength {
		t.Fatalf("expected sanitized value to include truncation marker, got length %d", len(result))
	}
	if !strings.HasSuffix(result, "...[truncated]") {
		t.Fatalf("expected truncation marker, got suffix %q", result[len(result)-20:])
	}
}

func TestDiffChangeFormatsPointerValues(t *testing.T) {
	type changeReq struct {
		Nickname  *string `change:"昵称"`
		IsEnabled *bool   `change:"启用状态"`
	}

	oldNickname := "旧昵称"
	newNickname := "新昵称"
	enabled := true
	disabled := false
	before := &changeReq{
		Nickname:  &oldNickname,
		IsEnabled: &enabled,
	}
	after := &changeReq{
		Nickname:  &newNickname,
		IsEnabled: &disabled,
	}

	got := DiffChange(before, after)
	want := "昵称：旧昵称->新昵称, 启用状态：true->false"
	if got != want {
		t.Fatalf("DiffChange() = %q, want %q", got, want)
	}
}

func TestDiffChangeFormatsNilPointerValue(t *testing.T) {
	type changeReq struct {
		Remark *string `change:"备注"`
	}

	remark := "补充说明"
	got := DiffChange(&changeReq{Remark: nil}, &changeReq{Remark: &remark})
	want := "备注：空->补充说明"
	if got != want {
		t.Fatalf("DiffChange() = %q, want %q", got, want)
	}
}
