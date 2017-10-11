package chinese_calendar

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
	{0, 348},                // no leap month, and every month has 29 days
	{1, 377},                // 1 leap month, and every month has 29 days.
	{(1<<12 - 1) * 16, 360}, // no leap month, and every month has 30 days.
	{(1<<13-1)*16 + 1, 390}, // 1 leap month, and every month has 30 days.
	{(1<<12-1)*16 + 1, 389}, // 1 leap month, and every normal month has 30 days, and leap month has 29 days.
}

func TestYearInfo_TotalDays(t *testing.T) {
	for _, dt := range yearInfo2YearDayTests {
		v := yearInfo{1900, dt.in}.TotalDays()
		if v != dt.out {
			t.Errorf("YearInfo{info:%d}.TotalDays() = %d, expect %d.", dt.in, v, dt.out)
		}
	}
}

type fromOffsetTest struct {
	in     int
	expect ChineseCalendar
}

var fromOffsetTests = []fromOffsetTest{
	{0, ChineseCalendar{1900, 1, 1, false}},
}

func TestFromOffset(t *testing.T) {
	for _, item := range fromOffsetTests {
		out, err := fromOffset(item.in)
		assert.NoError(t, err)
		assert.Equal(t, out, item.expect)
	}
}

func TestFromSolarDate(t *testing.T) {
	t1, err := FromSolarDate(1976, 10, 1)
	assert.NoError(t, err)
	assert.Equal(t, t1, ChineseCalendar{1976, 8, 8, true})
}

func TestChineseCalendar_Validate(t *testing.T) {
	assert.Equal(t, ChineseCalendar{1899, 12, 29, false}.Validate(), ErrYearOutOfRange)
	assert.Equal(t, ChineseCalendar{1899, 12, 30, false}.Validate(), ErrYearOutOfRange)
	assert.Equal(t, ChineseCalendar{1900, 1, 1, false}.Validate(), nil)
	assert.Equal(t, ChineseCalendar{2049, 1, 1, false}.Validate(), nil)
	assert.Equal(t, ChineseCalendar{2050, 1, 1, false}.Validate(), ErrYearOutOfRange)
	assert.Equal(t, ChineseCalendar{2051, 1, 1, false}.Validate(), ErrYearOutOfRange)

	assert.Equal(t, ChineseCalendar{1900, 1, 1, false}.Validate(), nil)
	assert.Equal(t, ChineseCalendar{1900, 1, 1, true}.Validate(), ErrNotLeapMonth)
	assert.Equal(t, ChineseCalendar{1976, 8, 1, false}.Validate(), nil)
	assert.Equal(t, ChineseCalendar{1976, 8, 1, true}.Validate(), nil)

	assert.Equal(t, ChineseCalendar{1900, 1, -1, false}.Validate(), ErrDayOutOfRange)
	assert.Equal(t, ChineseCalendar{1900, 1, 0, false}.Validate(), ErrDayOutOfRange)
	assert.Equal(t, ChineseCalendar{1900, 1, 1, false}.Validate(), nil)
	assert.Equal(t, ChineseCalendar{1900, 1, 29, false}.Validate(), nil)
	assert.Equal(t, ChineseCalendar{1900, 1, 30, false}.Validate(), ErrDayOutOfRange)
	assert.Equal(t, ChineseCalendar{1900, 2, 29, false}.Validate(), nil)
	assert.Equal(t, ChineseCalendar{1900, 2, 30, false}.Validate(), nil)
	assert.Equal(t, ChineseCalendar{1900, 2, 31, false}.Validate(), ErrDayOutOfRange)

	assert.Equal(t, ChineseCalendar{1900, 0, 1, false}.Validate(), ErrMonthOutOfRange)
	assert.Equal(t, ChineseCalendar{1900, 1, 1, false}.Validate(), nil)
	assert.Equal(t, ChineseCalendar{1900, 12, 1, false}.Validate(), nil)
	assert.Equal(t, ChineseCalendar{1900, 13, 1, false}.Validate(), ErrMonthOutOfRange)
}
func TestChineseCalendar_ToTime(t *testing.T) {
	assert.Equal(t,
		ChineseCalendar{1976, 8, 8, true}.MustToTime().Format("2006-01-02"),
		"1976-10-01",
	)

	oldLocal := time.Local
	local, err := time.LoadLocation("Asia/Shanghai")
	time.Local = local
	assert.NoError(t, err)
	name, offset := ChineseCalendar{1976, 8, 8, true}.MustToTime().Zone()
	assert.Equal(t, name, "CST")
	assert.Equal(t, offset, 8*3600)
	time.Local = oldLocal

	assert.Equal(t,
		ChineseCalendar{2049, 12, 29, false}.MustToTime().Format("2006-01-02"),
		"2050-01-22",
	)
}

func TestChineseCalendar_Before(t *testing.T) {
	assert.True(t,
		ChineseCalendar{1982, 11, 20, false}.Before(ChineseCalendar{1982, 11, 21, false}))
}

func TestChineseCalendar_After(t *testing.T) {
	assert.False(t,
		ChineseCalendar{1982, 11, 20, false}.After(ChineseCalendar{1982, 11, 21, false}))
}

func TestChineseCalendar_NextDay(t *testing.T) {
	d := ChineseCalendar{1982, 11, 20, false}.NextDay()
	assert.Equal(t, d, ChineseCalendar{1982, 11, 21, false})
}

func TestChineseCalendar_PrevDay(t *testing.T) {
	d := ChineseCalendar{1982, 11, 1, false}.PrevDay()
	assert.Equal(t, d, ChineseCalendar{1982, 10, 30, false})
	d = ChineseCalendar{1976, 9, 1, false}.PrevDay()
	assert.Equal(t, d, ChineseCalendar{1976, 8, 29, true})

}

func TestFromTime(t *testing.T) {
	t1, err := FromTime(time.Date(1900, time.January, 30, 0, 0, 0, 0, time.Local))
	assert.Equal(t, err, ErrTimeOutOfRange)

	t1, err = FromTime(time.Date(1900, time.January, 30, 23, 59, 0, 0, time.Local))
	assert.Equal(t, err, ErrTimeOutOfRange)

	t1 = MustFromTime(time.Date(1900, time.January, 31, 0, 0, 0, 0, time.Local))
	assert.Equal(t, t1, ChineseCalendar{1900, 1, 1, false})

	t1 = MustFromTime(time.Date(1983, time.January, 3, 0, 0, 0, 0, time.Local))
	assert.Equal(t, t1, ChineseCalendar{1982, 11, 20, false})

	t1 = MustFromTime(time.Date(2050, time.January, 22, 0, 0, 0, 0, time.Local))
	assert.Equal(t, t1, ChineseCalendar{2049, 12, 29, false})

	t1 = MustFromTime(time.Date(2050, time.January, 22, 23, 59, 59, 0, time.Local))
	assert.Equal(t, t1, ChineseCalendar{2049, 12, 29, false})

	_, err = FromTime(time.Date(2050, time.January, 23, 0, 0, 0, 0, time.Local))
	assert.Equal(t, err, ErrTimeOutOfRange)
}

func TestToday(t *testing.T) {
	c := Today()
	assert.NotNil(t, c)
}
