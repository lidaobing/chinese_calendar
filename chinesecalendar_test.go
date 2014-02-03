package chinesecalendar

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type yearInfo2YearDayTest struct {
	in  int
	out int
}

var yearInfo2YearDayTests = []yearInfo2YearDayTest{
	yearInfo2YearDayTest{0, 348},                // no leap month, and every month has 29 days
	yearInfo2YearDayTest{1, 377},                // 1 leap month, and every month has 29 days.
	yearInfo2YearDayTest{(1<<12 - 1) * 16, 360}, // no leap month, and every month has 30 days.
	yearInfo2YearDayTest{(1<<13-1)*16 + 1, 390}, // 1 leap month, and every month has 30 days.
	yearInfo2YearDayTest{(1<<12-1)*16 + 1, 389}, // 1 leap month, and every normal month has 30 days, and leap month has 29 days.
}

func TestYearInfo_TotalDays(t *testing.T) {
	for _, dt := range yearInfo2YearDayTests {
		v := YearInfo{1900, dt.in}.TotalDays()
		if v != dt.out {
			t.Errorf("YearInfo{info:%d}.TotalDays() = %d, expect %d.", dt.in, v, dt.out)
		}
	}
}

type fromOffsetTest struct {
	in     int
	expect *ChineseCalendar
}

var fromOffsetTests = []fromOffsetTest{
	fromOffsetTest{0, &ChineseCalendar{1900, 1, 1, false}},
}

func TestFromOffset(t *testing.T) {
	for _, item := range fromOffsetTests {
		out, err := fromOffset(item.in)
		assert.NoError(t, err)
		if !out.Equal(*item.expect) {
			t.Errorf("fromOffset(%d) = %s, expect %s", item.in, out, item.expect)
		}
	}
}

func TestFromSolarDate(t *testing.T) {
	t1, err := FromSolarDate(1976, 10, 1)
	assert.NoError(t, err)
	assert.Equal(t, t1, ChineseCalendar{1976, 8, 8, true})
}

func TestChineseCalendar_ToSolarDate(t *testing.T) {
	assert.Equal(t,
		ChineseCalendar{1976, 8, 8, true}.ToSolarDate().Format("2006-01-02"),
		"1976-10-01",
	)
}

func TestFromTime(t *testing.T) {
	t1 := time.Date(2014, time.February, 14, 0, 0, 0, 0, time.Local)

	t2, err := FromTime(t1)
	assert.NoError(t, err)
	assert.Equal(t, t2, ChineseCalendar{2014, 1, 15, false})
}

func TestToday(t *testing.T) {
	c := Today()
	assert.NotNil(t, c)
}
