package chinesecalendar

import "testing"

type yearInfo2YearDayTest struct {
    in int
    out int
}

var yearInfo2YearDayTests = []yearInfo2YearDayTest {
    yearInfo2YearDayTest{0, 348}, // no leap month, and every month has 29 days
    yearInfo2YearDayTest{1, 377}, // 1 leap month, and every month has 29 days.
    yearInfo2YearDayTest{(1<<12-1)*16, 360}, // no leap month, and every month has 30 days.
    yearInfo2YearDayTest{(1<<13-1)*16+1, 390}, // 1 leap month, and every month has 30 days.
    yearInfo2YearDayTest{(1<<12-1)*16+1, 389}, // 1 leap month, and every normal month has 30 days, and leap month has 29 days.
}

func TestYearInfo2YearDay(t *testing.T) {
    for _, dt := range yearInfo2YearDayTests {
        v := yearInfo2yearDay(dt.in)
        if(v != dt.out) {
            t.Errorf("yearInfo2yearDay(%d) = %d, expect %d.", dt.in, v, dt.out)
        }
    }
}

func TestA(t *testing.T) {
    cc := fromOffset(0)
    if(cc.year != 0) {
        t.FailNow()
    }
}
