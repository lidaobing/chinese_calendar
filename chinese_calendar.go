// The chinesecalendar package provide support conversion between chinese calendar and time.Time
//
// for more information on Chinese Calendar: http://en.wikipedia.org/wiki/Chinese_calendar
//
// https://github.com/lidaobing/gochinesecalendar
package chinese_calendar

import (
	"errors"
	"fmt"
	"time"
)

var (
	ErrYearOutOfRange  = errors.New("year out of range [1900, 2050)")
	ErrMonthOutOfRange = errors.New("month out of range [1, 12]")
	ErrNotLeapMonth    = errors.New("this month is not a leap month")
	ErrDayOutOfRange   = errors.New("day out of range")
	ErrTimeOutOfRange  = errors.New("time out of range [1900-01-31, 2050-01-22]")
)

// ChineseCalendar
type ChineseCalendar struct {
	Year        int
	Month       int
	Day         int
	IsLeapMonth bool
}

func (c ChineseCalendar) MustToTime() time.Time {
	res, err := c.ToTime()
	if err != nil {
		panic(err)
	}
	return res
}

func (c ChineseCalendar) ToTime() (res time.Time, err error) {
	offset := 0
	if c.Year < 1900 || c.Year >= 2050 {
		return res, ErrYearOutOfRange
	}
	yearIdx := c.Year - 1900
	for i := 0; i < yearIdx; i++ {
		offset += yearDays[i]
	}

	var days int
	days, err = calcDays(yearInfos[yearIdx], c.Month, c.Day, c.IsLeapMonth)
	if err != nil {
		return
	}
	offset += days
	res = startDate(time.Local).AddDate(0, 0, offset)
	return
}

func (c ChineseCalendar) Validate() (err error) {
	_, err = c.ToTime()
	return
}

func (c ChineseCalendar) IsValid() bool {
	return c.Validate() == nil
}

func (c ChineseCalendar) Before(u ChineseCalendar) bool {
	return c.MustToTime().Before(u.MustToTime())
}

func (c ChineseCalendar) After(u ChineseCalendar) bool {
	return c.MustToTime().After(u.MustToTime())
}

func (c ChineseCalendar) NextDay() (res ChineseCalendar) {
	return MustFromTime(c.MustToTime().AddDate(0, 0, 1))
}

func (c ChineseCalendar) PrevDay() (res ChineseCalendar) {
	return MustFromTime(c.MustToTime().AddDate(0, 0, -1))
}

type yearInfoItem struct {
	Month       int
	Days        int
	IsLeapMonth bool
}

func FromSolarDate(year, month, day int) (res ChineseCalendar, err error) {
	var t time.Time
	t, err = time.Parse("2006-01-02", fmt.Sprintf("%04d-%02d-%02d", year, month, day))
	if err != nil {
		return
	}
	return FromTime(t)
}

func FromTime(t time.Time) (res ChineseCalendar, err error) {
	// reset to begin of day, a workaround for the leap second problem
	t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
	floatDays := t.Sub(startDate(t.Location())).Hours() / 24
	if floatDays < 0 {
		err = ErrTimeOutOfRange
		return
	}
	days := int(floatDays + 0.5)
	if days > 54778 {
		err = ErrTimeOutOfRange
		return
	}
	return fromOffset(days)
}

func MustFromTime(t time.Time) (res ChineseCalendar) {
	res, err := FromTime(t)
	if err != nil {
		panic(err)
	}
	return
}

func Today() (res ChineseCalendar) {
	res, err := FromTime(time.Now())
	if err != nil {
		panic(err)
	}
	return
}

// the information of one year
type yearInfo struct {
	Year int
	info int
}

// return the total days of the year
func (self yearInfo) TotalDays() int {
	var res = 29 * 12
	var leap = 0
	var yearInfo = self.info
	if yearInfo%16 != 0 {
		leap = 1
		res += 29
	}
	yearInfo = yearInfo / 16
	for i := 0; i < 12+leap; i++ {
		if yearInfo%2 == 1 {
			res += 1
		}
		yearInfo = yearInfo / 2
	}
	return res
}

// return the idx-th day in this year
// if idx < 0 or idx > TotalDays(), return nil
func (self yearInfo) Day(idx int) *ChineseCalendar {
	if idx < 0 || idx > self.TotalDays() {
		return nil
	}
	var month = 1
	var day = 1
	var isLeapMonth = false
	return &ChineseCalendar{Year: self.Year, Month: month, Day: day, IsLeapMonth: isLeapMonth}
}

func calcYearInfos() []int {
	return []int{
		/* encoding:
		              b bbbbbbbbbbbb bbbb
		      bit#    1 111111000000 0000
		              6 543210987654 3210
		              . ............ ....
		      month#    000000000111
		              M 123456789012   L

		   b_j = 1 for long month, b_j = 0 for short month
		   L is the leap month of the year if 1<=L<=12; NO leap month if L = 0.
		   The leap month (if exists) is long one iff M = 1.
		*/
		0x04bd8,                                     /* 1900 */
		0x04ae0, 0x0a570, 0x054d5, 0x0d260, 0x0d950, /* 1905 */
		0x16554, 0x056a0, 0x09ad0, 0x055d2, 0x04ae0, /* 1910 */
		0x0a5b6, 0x0a4d0, 0x0d250, 0x1d255, 0x0b540, /* 1915 */
		0x0d6a0, 0x0ada2, 0x095b0, 0x14977, 0x04970, /* 1920 */
		0x0a4b0, 0x0b4b5, 0x06a50, 0x06d40, 0x1ab54, /* 1925 */
		0x02b60, 0x09570, 0x052f2, 0x04970, 0x06566, /* 1930 */
		0x0d4a0, 0x0ea50, 0x06e95, 0x05ad0, 0x02b60, /* 1935 */
		0x186e3, 0x092e0, 0x1c8d7, 0x0c950, 0x0d4a0, /* 1940 */
		0x1d8a6, 0x0b550, 0x056a0, 0x1a5b4, 0x025d0, /* 1945 */
		0x092d0, 0x0d2b2, 0x0a950, 0x0b557, 0x06ca0, /* 1950 */
		0x0b550, 0x15355, 0x04da0, 0x0a5d0, 0x14573, /* 1955 */
		0x052d0, 0x0a9a8, 0x0e950, 0x06aa0, 0x0aea6, /* 1960 */
		0x0ab50, 0x04b60, 0x0aae4, 0x0a570, 0x05260, /* 1965 */
		0x0f263, 0x0d950, 0x05b57, 0x056a0, 0x096d0, /* 1970 */
		0x04dd5, 0x04ad0, 0x0a4d0, 0x0d4d4, 0x0d250, /* 1975 */
		0x0d558, 0x0b540, 0x0b5a0, 0x195a6, 0x095b0, /* 1980 */
		0x049b0, 0x0a974, 0x0a4b0, 0x0b27a, 0x06a50, /* 1985 */
		0x06d40, 0x0af46, 0x0ab60, 0x09570, 0x04af5, /* 1990 */
		0x04970, 0x064b0, 0x074a3, 0x0ea50, 0x06b58, /* 1995 */
		0x05ac0, 0x0ab60, 0x096d5, 0x092e0, 0x0c960, /* 2000 */
		0x0d954, 0x0d4a0, 0x0da50, 0x07552, 0x056a0, /* 2005 */
		0x0abb7, 0x025d0, 0x092d0, 0x0cab5, 0x0a950, /* 2010 */
		0x0b4a0, 0x0baa4, 0x0ad50, 0x055d9, 0x04ba0, /* 2015 */
		0x0a5b0, 0x15176, 0x052b0, 0x0a930, 0x07954, /* 2020 */
		0x06aa0, 0x0ad50, 0x05b52, 0x04b60, 0x0a6e6, /* 2025 */
		0x0a4e0, 0x0d260, 0x0ea65, 0x0d530, 0x05aa0, /* 2030 */
		0x076a3, 0x096d0, 0x04afb, 0x04ad0, 0x0a4d0, /* 2035 */
		0x1d0b6, 0x0d250, 0x0d520, 0x0dd45, 0x0b5a0, /* 2040 */
		0x056d0, 0x055b2, 0x049b0, 0x0a577, 0x0a4b0, /* 2045 */
		0x0aa50, 0x1b255, 0x06d20, 0x0ada0} /* 2049 */
}

func calcYearDays() (res []int) {
	for _, yi := range yearInfos {
		res = append(res, yearInfo{info: yi}.TotalDays())
	}
	return res
}

func fromOffset(offset int) (res ChineseCalendar, err error) {
	if offset < 0 {
		err = fmt.Errorf("offset must >= 0")
		return
	}
	var yearInfo int
	for idx, yearDay := range yearDays {
		if offset < yearDay {
			res.Year = 1900 + idx
			yearInfo = yearInfos[idx]
			break
		}
		offset -= yearDay
	}
	if res.Year == 0 {
		err = fmt.Errorf("offset too large")
		return
	}
	month, day, isLeapMonth := calcMonthDay(yearInfo, offset)
	res.Month = month
	res.Day = day
	res.IsLeapMonth = isLeapMonth
	return
}

func calcMonthDay(yearInfo, offset int) (month, day int, isLeapMonth bool) {
	for _, yii := range enumMonth(yearInfo) {
		if offset < yii.Days {
			month = yii.Month
			day = offset + 1
			isLeapMonth = yii.IsLeapMonth
			return
		}
		offset -= yii.Days
	}
	panic("offset too large for the yearInfo")
	return
}

func calcDays(yearInfo, month, day int, isLeapMonth bool) (offset int, err error) {
	if month < 1 || month > 12 {
		return 0, ErrMonthOutOfRange
	}
	for _, yii := range enumMonth(yearInfo) {
		if month == yii.Month && isLeapMonth == yii.IsLeapMonth {
			if day <= 0 || day > yii.Days {
				return 0, ErrDayOutOfRange
			}
			offset += day - 1
			return
		}
		offset += yii.Days
	}
	return 0, ErrNotLeapMonth
}

func enumMonth(yearInfo int) (res []yearInfoItem) {
	leapMonth := yearInfo % 16
	if leapMonth > 12 {
		panic(fmt.Sprintf("invalid yearInfo: %d", yearInfo))
	}
	for i := 1; i < 13; i++ {
		yii := yearInfoItem{
			Month:       i,
			Days:        (yearInfo>>uint(16-i))%2 + 29,
			IsLeapMonth: false,
		}
		res = append(res, yii)
		if leapMonth == i {
			yii := yearInfoItem{
				Month:       i,
				Days:        (yearInfo>>16)%2 + 29,
				IsLeapMonth: true,
			}
			res = append(res, yii)
		}
	}
	return
}

func startDate(l *time.Location) time.Time {
	return time.Date(1900, time.January, 31, 0, 0, 0, 0, l)
}

var (
	yearInfos = calcYearInfos()
	yearDays  = calcYearDays()
)
