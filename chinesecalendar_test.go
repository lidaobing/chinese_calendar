package chinesecalendar

import "testing"

func TestA(t *testing.T) {
	cc := fromOffset(0)
	if(cc.year != 0) {
		t.FailNow()
	}
}
