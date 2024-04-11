package mapx

import (
	"encoding/json"
	"fmt"
	"testing"
)

func assertEq[T any](subject string, a T, b T, f func(string)) {
	if f == nil {
		return
	}

	ab, _ := json.Marshal(a)
	bb, _ := json.Marshal(b)

	as := string(ab)
	bs := string(bb)

	if as == bs {
		return
	}

	s := "unexpected '%v':\n\twant: '%v'\n\thave: '%v'\n"
	f(fmt.Sprintf(s, subject, as, bs))
}

func TestNewIdeal(t *testing.T) {
	p1 := Pair[int, int]{1, 2}
	p2 := Pair[int, int]{2, 3}

	assertEq("r", map[int]int{1: 2, 2: 3}, New(p1, p2), func(s string) { t.Fatal(s) })
}
