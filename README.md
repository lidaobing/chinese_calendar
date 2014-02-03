# A Chinese Calendar Library in Golang

https://github.com/lidaobing/chinese_calendar

[![Build Status](https://travis-ci.org/lidaobing/chinese_calendar.png?branch=master)](https://travis-ci.org/lidaobing/chinese_calendar)
[![GitHub version](https://badge.fury.io/gh/lidaobing%2Fchinese_calendar.png)](http://badge.fury.io/gh/lidaobing%2Fchinese_calendar)


## Install

```go
go get -u github.com/lidaobing/chinese_calendar
```

## Usage

```go
package main

import (
  "fmt"
  "time"

  cc "github.com/lidaobing/chinese_calendar"
)

func main() {
  today := cc.Today()
  fmt.Printf("today:\t%#v\n", today)

  today, _ = cc.FromTime(time.Now())
  fmt.Printf("today(from time.Now()):\t%#v\n", today)

  now, _ := today.ToTime()
  fmt.Printf("today in time: %s\n", now)

  fmt.Printf("Today: year: %d: month: %d, day: %d, isLeapMonth: %v\n",
    today.Year, today.Month, today.Day, today.IsLeapMonth)

  tomorrow := today.NextDay()
  fmt.Printf("Tomorrow: year: %d: month: %d, day: %d, isLeapMonth: %v\n",
    tomorrow.Year, tomorrow.Month, tomorrow.Day, tomorrow.IsLeapMonth)

  yesterday := today.PrevDay()
  fmt.Printf("Yesterday: year: %d: month: %d, day: %d, isLeapMonth: %v\n",
    yesterday.Year, yesterday.Month, yesterday.Day, yesterday.IsLeapMonth)

  fmt.Printf("yesterday is before today? %v\n", yesterday.Before(today))
  fmt.Printf("today is after tomorrow? %v\n", today.After(tomorrow))
}
```

sample output:

```
today:  chinesecalendar.ChineseCalendar{Year:2014, Month:1, Day:4, IsLeapMonth:false}
today(from time.Now()): chinesecalendar.ChineseCalendar{Year:2014, Month:1, Day:4, IsLeapMonth:false}
today in time: 2014-02-03 00:00:00 +0800 CST
Today: year: 2014: month: 1, day: 4, isLeapMonth: false
Tomorrow: year: 2014: month: 1, day: 5, isLeapMonth: false
Yesterday: year: 2014: month: 1, day: 3, isLeapMonth: false
yesterday is before today? true
today is after tomorrow? false
```

## LICENSE

Permission is hereby granted, free of charge, to any person obtaining
a copy of this software and associated documentation files (the
'Software'), to deal in the Software without restriction, including
without limitation the rights to use, copy, modify, merge, publish,
distribute, sublicense, and/or sell copies of the Software, and to
permit persons to whom the Software is furnished to do so, subject to
the following conditions:

The above copyright notice and this permission notice shall be
included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED 'AS IS', WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY
CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT,
TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE
SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

