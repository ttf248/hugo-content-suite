package operations

import "testing"

type testStatus string

func (s testStatus) GetStatus() string { return string(s) }

func TestFilterByModeExcludesSkippedItems(t *testing.T) {
	items := []testStatus{ModeCreate, ModeUpdate, "skip"}
	if got := filterByMode(items, ModeAll); len(got) != 2 {
		t.Fatalf("all 模式不应包含 skip: %#v", got)
	}
}
