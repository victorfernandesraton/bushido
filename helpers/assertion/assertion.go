package assertion

import (
	"fmt"
	"reflect"
	"testing"
)

func Assert(t *testing.T, ok bool, msg any) {
	t.Helper()
	if !ok {
		t.Fatalf("assert fail: %v", msg)
	}
}
func AssertEqual(t *testing.T, want, got any) {
	t.Helper()
	Assert(t, want == got, fmt.Sprintf("want: '%v', got: '%v'", want, got))
}

func AssertDeepEqual(t *testing.T, want, got any) {
	t.Helper()
	Assert(t, reflect.DeepEqual(want, got), fmt.Sprintf("want:\n%#v\n\ngot:\n%#v", want, got))
}
