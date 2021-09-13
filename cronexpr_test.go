/*!
 * Copyright 2013 Raymond Hill
 *
 * Project: github.com/gorhill/cronexpr
 * File: cronexpr_test.go
 * Version: 1.0
 * License: pick the one which suits you best:
 *   GPL v3 see <https://www.gnu.org/licenses/gpl.html>
 *   APL v2 see <http://www.apache.org/licenses/LICENSE-2.0>
 *
 */

package cronexpr

/******************************************************************************/

import (
	"fmt"
	"testing"
	"time"
)

/******************************************************************************/

type crontimes struct {
	from string
	next string
}

type dstLocationDates struct {
	when  string
	where string
}

type locationTestCase struct {
	when     string
	expected string // Expectation in UTC
	where    string
}

type crontest struct {
	expr   string
	layout string
	times  []crontimes
}

var crontests = []crontest{
	// Seconds
	{
		"* * * * * * *",
		"2006-01-02 15:04:05",
		[]crontimes{
			{"2013-01-01 00:00:00", "2013-01-01 00:00:01"},
			{"2013-01-01 00:00:59", "2013-01-01 00:01:00"},
			{"2013-01-01 00:59:59", "2013-01-01 01:00:00"},
			{"2013-01-01 23:59:59", "2013-01-02 00:00:00"},
			{"2013-02-28 23:59:59", "2013-03-01 00:00:00"},
			{"2016-02-28 23:59:59", "2016-02-29 00:00:00"},
			{"2012-12-31 23:59:59", "2013-01-01 00:00:00"},
		},
	},

	// every 5 Second
	{
		"*/5 * * * * * *",
		"2006-01-02 15:04:05",
		[]crontimes{
			{"2013-01-01 00:00:00", "2013-01-01 00:00:05"},
			{"2013-01-01 00:00:59", "2013-01-01 00:01:00"},
			{"2013-01-01 00:59:59", "2013-01-01 01:00:00"},
			{"2013-01-01 23:59:59", "2013-01-02 00:00:00"},
			{"2013-02-28 23:59:59", "2013-03-01 00:00:00"},
			{"2016-02-28 23:59:59", "2016-02-29 00:00:00"},
			{"2012-12-31 23:59:59", "2013-01-01 00:00:00"},
		},
	},

	// Minutes
	{
		"* * * * *",
		"2006-01-02 15:04:05",
		[]crontimes{
			{"2013-01-01 00:00:00", "2013-01-01 00:01:00"},
			{"2013-01-01 00:00:59", "2013-01-01 00:01:00"},
			{"2013-01-01 00:59:00", "2013-01-01 01:00:00"},
			{"2013-01-01 23:59:00", "2013-01-02 00:00:00"},
			{"2013-02-28 23:59:00", "2013-03-01 00:00:00"},
			{"2016-02-28 23:59:00", "2016-02-29 00:00:00"},
			{"2012-12-31 23:59:00", "2013-01-01 00:00:00"},
		},
	},

	// Minutes with interval
	{
		"17-43/5 * * * *",
		"2006-01-02 15:04:05",
		[]crontimes{
			{"2013-01-01 00:00:00", "2013-01-01 00:17:00"},
			{"2013-01-01 00:16:59", "2013-01-01 00:17:00"},
			{"2013-01-01 00:30:00", "2013-01-01 00:32:00"},
			{"2013-01-01 00:50:00", "2013-01-01 01:17:00"},
			{"2013-01-01 23:50:00", "2013-01-02 00:17:00"},
			{"2013-02-28 23:50:00", "2013-03-01 00:17:00"},
			{"2016-02-28 23:50:00", "2016-02-29 00:17:00"},
			{"2012-12-31 23:50:00", "2013-01-01 00:17:00"},
		},
	},

	// Minutes interval, list
	{
		"15-30/4,55 * * * *",
		"2006-01-02 15:04:05",
		[]crontimes{
			{"2013-01-01 00:00:00", "2013-01-01 00:15:00"},
			{"2013-01-01 00:16:00", "2013-01-01 00:19:00"},
			{"2013-01-01 00:30:00", "2013-01-01 00:55:00"},
			{"2013-01-01 00:55:00", "2013-01-01 01:15:00"},
			{"2013-01-01 23:55:00", "2013-01-02 00:15:00"},
			{"2013-02-28 23:55:00", "2013-03-01 00:15:00"},
			{"2016-02-28 23:55:00", "2016-02-29 00:15:00"},
			{"2012-12-31 23:54:00", "2012-12-31 23:55:00"},
			{"2012-12-31 23:55:00", "2013-01-01 00:15:00"},
		},
	},

	// Hours
	{
		"0 0 13-15 ? * * *",
		"2006-01-02 15:04:05",
		[]crontimes{
			{"2013-04-04 12:00:00", "2013-04-04 13:00:00"},
			{"2013-04-04 12:53:00", "2013-04-04 13:00:00"},
			{"2013-04-04 23:00:00", "2013-04-05 13:00:00"},
			{"2013-04-04 13:00:00", "2013-04-04 14:00:00"},
			{"2013-04-04 13:40:00", "2013-04-04 14:00:00"},
			{"2013-04-04 13:40:00", "2013-04-04 14:00:00"},
			{"2013-04-04 14:40:00", "2013-04-04 15:00:00"},
			{"2013-04-04 15:00:00", "2013-04-05 13:00:00"},
		},
	},

	// Days of week
	{
		"0 0 * * MON",
		"Mon 2006-01-02 15:04",
		[]crontimes{
			{"2013-01-01 00:00:00", "Mon 2013-01-07 00:00"},
			{"2013-01-28 00:00:00", "Mon 2013-02-04 00:00"},
			{"2013-12-30 00:30:00", "Mon 2014-01-06 00:00"},
		},
	},
	{
		"0 0 * * friday",
		"Mon 2006-01-02 15:04",
		[]crontimes{
			{"2013-01-01 00:00:00", "Fri 2013-01-04 00:00"},
			{"2013-01-28 00:00:00", "Fri 2013-02-01 00:00"},
			{"2013-12-30 00:30:00", "Fri 2014-01-03 00:00"},
		},
	},
	{
		"0 0 * * 6,7",
		"Mon 2006-01-02 15:04",
		[]crontimes{
			{"2013-01-01 00:00:00", "Sat 2013-01-05 00:00"},
			{"2013-01-28 00:00:00", "Sat 2013-02-02 00:00"},
			{"2013-12-30 00:30:00", "Sat 2014-01-04 00:00"},
		},
	},

	// Specific days of week
	{
		"0 0 * * 6#5",
		"Mon 2006-01-02 15:04",
		[]crontimes{
			{"2013-09-02 00:00:00", "Sat 2013-11-30 00:00"},
		},
	},

	// Work day of month
	{
		"0 0 14W * *",
		"Mon 2006-01-02 15:04",
		[]crontimes{
			{"2013-03-31 00:00:00", "Mon 2013-04-15 00:00"},
			{"2013-08-31 00:00:00", "Fri 2013-09-13 00:00"},
		},
	},

	// Work day of month -- end of month
	{
		"0 0 30W * *",
		"Mon 2006-01-02 15:04",
		[]crontimes{
			{"2013-03-02 00:00:00", "Fri 2013-03-29 00:00"},
			{"2013-06-02 00:00:00", "Fri 2013-06-28 00:00"},
			{"2013-09-02 00:00:00", "Mon 2013-09-30 00:00"},
			{"2013-11-02 00:00:00", "Fri 2013-11-29 00:00"},
		},
	},

	// Last day of month
	{
		"0 0 L * *",
		"Mon 2006-01-02 15:04",
		[]crontimes{
			{"2013-09-02 00:00:00", "Mon 2013-09-30 00:00"},
			{"2014-01-01 00:00:00", "Fri 2014-01-31 00:00"},
			{"2014-02-01 00:00:00", "Fri 2014-02-28 00:00"},
			{"2016-02-15 00:00:00", "Mon 2016-02-29 00:00"},
		},
	},

	// Zero padded months
	{
		"0 0 * 04 * *",
		"Mon 2006-01-02 15:04",
		[]crontimes{
			{"2013-09-02 00:00:00", "Tue 2014-04-01 00:00"},
			{"2014-04-03 03:00:00", "Fri 2014-04-04 00:00"},
			{"2014-08-15 00:00:00", "Wed 2015-04-01 00:00"},
		},
	},

	{ // Zero leading values
		"00 01 03 07 *",
		"Mon 2006-01-02 15:04",
		[]crontimes{
			{"2013-01-01 00:00:00", "Wed 2013-07-03 01:00"},
			{"2014-01-28 00:00:00", "Thu 2014-07-03 01:00"},
			{"2013-12-30 00:30:00", "Thu 2014-07-03 01:00"},
			{"2015-07-03 02:01:00", "Sun 2016-07-03 01:00"},
		},
	},

	// Last work day of month
	{
		"0 0 LW * *",
		"Mon 2006-01-02 15:04",
		[]crontimes{
			{"2013-09-02 00:00:00", "Mon 2013-09-30 00:00"},
			{"2013-11-02 00:00:00", "Fri 2013-11-29 00:00"},
			{"2014-08-15 00:00:00", "Fri 2014-08-29 00:00"},
		},
	},

	// TODO: more tests
}

func TestExpressions(t *testing.T) {
	for _, test := range crontests {
		for _, times := range test.times {
			from, _ := time.Parse("2006-01-02 15:04:05", times.from)
			expr, err := Parse(test.expr)
			if err != nil {
				t.Errorf(`Parse("%s") returned "%s"`, test.expr, err.Error())
			}
			next := expr.Next(from)
			nextstr := next.Format(test.layout)
			if nextstr != times.next {
				t.Errorf(`("%s").Next("%s") = "%s", got "%s"`, test.expr, times.from, times.next, nextstr)
			}
		}
	}
}

/******************************************************************************/
var dstLocationTests = []dstLocationDates{
	{"2019-03-31 00:00:00", "Europe/London"},
	{"2019-10-27 00:00:00", "Europe/London"},
	{"2019-03-31 00:00:00", "Europe/Paris"},
	{"2019-10-27 00:00:00", "Europe/Paris"},
	{"2019-03-10 00:00:00", "America/New_York"},
	{"2019-11-03 00:00:00", "America/New_York"},
	{"2019-03-10 00:00:00", "America/Los_Angeles"},
	{"2019-11-03 00:00:00", "America/Los_Angeles"},
	{"2019-03-10 00:00:00", "US/Pacific"},
	{"2019-11-03 00:00:00", "US/Pacific"},
	{"2019-03-31 00:00:00", "Asia/Shanghai"},  // Non-DST
	{"2019-10-02 00:00:00", "Asia/Shanghai"},  // Non-DST
	{"2019-03-31 00:00:00", "Asia/Hong_Kong"}, // Non-DST
	{"2019-10-02 00:00:00", "Asia/Hong_Kong"}, // Non-DST
	{"2019-04-05 00:00:00", "Australia/Melbourne"},
	{"2019-10-04 00:00:00", "Australia/Melbourne"},
	{"2019-03-31 00:00:00", "Australia/Perth"}, // Non-DST
	{"2019-10-02 00:00:00", "Australia/Perth"}, // Non-DST
	{"2019-03-22 00:00:00", "America/Asuncion"},
	{"2019-10-04 00:00:00", "America/Asuncion"},
}

func TestHourlyExpressionsWithTimeZones(t *testing.T) {
	for _, test := range dstLocationTests {
		location, err := time.LoadLocation(test.where)
		if err != nil {
			t.Errorf("Invalid Test Location:%s", test.where)
		}

		from, _ := time.ParseInLocation("2006-01-02 15:04:05", test.when, location)
		// Every 30 mins
		expr, err := Parse("0 0/30 * * * * *")
		if err != nil {
			t.Errorf(err.Error())
		}
		for j := 0; j < 6; j++ {
			next := expr.Next(from)
			diff := next.Sub(from)
			if diff != time.Duration(30)*time.Minute {
				t.Errorf(`DST changed failed for %s at %s. From:%s, To %s From:%s, To %s, diff %s`,
					test.where, test.when, from, next, from.UTC(), next.UTC(), next.UTC().Sub(from.UTC()))
			}
			from = next
		}
	}
}

func TestSingleExpressionsWithTimeZones(t *testing.T) {
	for _, test := range []dstLocationDates{
		{"2019-03-31 01:30:00", "Europe/London"},
		{"2019-10-27 01:30:00", "Europe/London"},
		{"2019-03-31 01:30:00", "Europe/Paris"},
		{"2019-10-02 01:30:00", "Europe/Paris"},
		{"2019-03-10 01:30:00", "America/New_York"},
		{"2019-11-03 01:30:00", "America/New_York"},
		{"2019-03-10 01:30:00", "America/Los_Angeles"},
		{"2019-11-03 01:30:00", "America/Los_Angeles"},
		{"2019-03-10 01:30:00", "US/Pacific"},
		{"2019-11-03 01:30:00", "US/Pacific"},
	} {
		location, err := time.LoadLocation(test.where)
		if err != nil {
			t.Errorf("Invalid Test Location:%s", test.where)
		}
		start, _ := time.ParseInLocation("2006-01-02 15:04:05", test.when, location)
		expected := start.UTC().Add(30 * time.Minute)
		// Every Hour
		expr, err := Parse("0 0 * * * * *")
		if err != nil {
			t.Errorf(err.Error())
		}
		next := expr.Next(start)
		if !next.Equal(expected) {
			t.Errorf(`Unexpected convertion at :%s, expected :%s, actual: %s from: %s`,
				test.where, expected, next.UTC(), start.UTC())
		}

	}
}

func TestExpressionsWithTimeZones(t *testing.T) {
	for nHourly := 1; nHourly < 12; nHourly++ {
		for _, test := range dstLocationTests {
			location, err := time.LoadLocation(test.where)
			if err != nil {
				t.Errorf("Invalid Test Location:%s", test.where)
			}

			from, _ := time.ParseInLocation("2006-01-02 15:04:05", test.when, location)

			// At every n Hour
			expr, err := Parse(fmt.Sprintf("0 0 */%v * * * *", nHourly))
			if err != nil {
				t.Errorf(err.Error())
			}

			next := expr.Next(from)
			if !next.After(from) {
				t.Errorf(`Interval %v DST changed failed for %s at %s. From:%s, To %s From:%s, To %s, diff %s`,
					nHourly, test.where, test.when, from, next, from.UTC(), next.UTC(), next.UTC().Sub(from.UTC()))
			}
		}
	}
}

var specificLocationTests = []locationTestCase{
	{"2019-03-31 02:30:00", "2019-04-01T02:30:00+01:00", "Europe/London"},
	{"2019-10-27 02:30:00", "2019-10-28T02:30:00Z", "Europe/London"},
	{"2019-03-31 02:30:00", "2019-04-01T02:30:00+02:00", "Europe/Paris"},
	{"2019-10-27 02:30:00", "2019-10-28T02:30:00+01:00", "Europe/Paris"},
	{"2019-03-09 02:30:00", "2019-03-10T02:30:00-05:00", "America/New_York"},
	{"2019-03-10 03:30:00", "2019-03-11T02:30:00-04:00", "America/New_York"},
	{"2019-11-03 02:30:00", "2019-11-04T02:30:00-05:00", "America/New_York"},
	{"2019-11-04 02:30:00", "2019-11-05T02:30:00-05:00", "America/New_York"},
	{"2019-03-09 02:30:00", "2019-03-10T02:30:00-08:00", "America/Los_Angeles"},
	{"2019-11-03 02:30:00", "2019-11-04T02:30:00-08:00", "America/Los_Angeles"},
	{"2019-03-09 02:30:00", "2019-03-10T03:30:00-07:00", "US/Pacific"},
	{"2019-11-03 02:30:00", "2019-11-04T02:30:00-08:00", "US/Pacific"},
	{"2019-04-05 02:30:00", "2019-04-06T02:30:00+11:00", "Australia/Melbourne"},
	{"2019-10-04 02:30:00", "2019-10-05T02:30:00+10:00", "Australia/Melbourne"},
	{"2019-10-04 02:30:00", "2019-10-05T02:30:00-04:00", "America/Asuncion"},
}

func TestDailyExpressionsWithTimeZones(t *testing.T) {
	for _, test := range specificLocationTests {
		location, err := time.LoadLocation(test.where)
		if err != nil {
			t.Errorf("Invalid Test Location:%s", test.where)
		}

		from, _ := time.ParseInLocation("2006-01-02 15:04:05", test.when, location)
		to, _ := time.Parse(time.RFC3339, test.expected)
		// Every 30 mins
		expr, err := Parse("0 30 2 * * * *")
		if err != nil {
			t.Errorf(err.Error())
		}
		if expr.Next(from).UTC() != to.UTC() {
			t.Errorf("Expected Next time to be %s from %s in %s, when actually %s; UTC FROM %s, UTC TO:%s ", to, from, test.where, expr.Next(from), from.UTC(), expr.Next(from).UTC())
		}
	}
}

/******************************************************************************/

func TestZero(t *testing.T) {
	from, _ := time.Parse("2006-01-02", "2013-08-31")
	next := MustParse("* * * * * 1980").Next(from)
	if next.IsZero() == false {
		t.Error(`("* * * * * 1980").Next("2013-08-31").IsZero() returned 'false', expected 'true'`)
	}

	next = MustParse("* * * * * 2050").Next(from)
	if next.IsZero() == true {
		t.Error(`("* * * * * 2050").Next("2013-08-31").IsZero() returned 'true', expected 'false'`)
	}

	next = MustParse("* * * * * 2099").Next(time.Time{})
	if next.IsZero() == false {
		t.Error(`("* * * * * 2014").Next(time.Time{}).IsZero() returned 'true', expected 'false'`)
	}
}

func TestNextN(t *testing.T) {
	expected := []string{
		"Sat, 30 Nov 2013 00:00:00",
		"Sat, 29 Mar 2014 00:00:00",
		"Sat, 31 May 2014 00:00:00",
		"Sat, 30 Aug 2014 00:00:00",
		"Sat, 29 Nov 2014 00:00:00",
	}
	from, _ := time.Parse("2006-01-02 15:04:05", "2013-09-02 08:44:30")
	result := MustParse("0 0 * * 6#5").NextN(from, uint(len(expected)))
	if len(result) != len(expected) {
		t.Errorf(`MustParse("0 0 * * 6#5").NextN("2013-09-02 08:44:30", 5):\n"`)
		t.Errorf(`  Expected %d returned time values but got %d instead`, len(expected), len(result))
	}
	for i, next := range result {
		nextStr := next.Format("Mon, 2 Jan 2006 15:04:15")
		if nextStr != expected[i] {
			t.Errorf(`MustParse("0 0 * * 6#5").NextN("2013-09-02 08:44:30", 5):\n"`)
			t.Errorf(`  result[%d]: expected "%s" but got "%s"`, i, expected[i], nextStr)
		}
	}
}

func TestNextN_every5min(t *testing.T) {
	expected := []string{
		"Mon, 2 Sep 2013 08:45:00",
		"Mon, 2 Sep 2013 08:50:00",
		"Mon, 2 Sep 2013 08:55:00",
		"Mon, 2 Sep 2013 09:00:00",
		"Mon, 2 Sep 2013 09:05:00",
	}
	from, _ := time.Parse("2006-01-02 15:04:05", "2013-09-02 08:44:32")
	result := MustParse("*/5 * * * *").NextN(from, uint(len(expected)))
	if len(result) != len(expected) {
		t.Errorf(`MustParse("*/5 * * * *").NextN("2013-09-02 08:44:30", 5):\n"`)
		t.Errorf(`  Expected %d returned time values but got %d instead`, len(expected), len(result))
	}
	for i, next := range result {
		nextStr := next.Format("Mon, 2 Jan 2006 15:04:05")
		if nextStr != expected[i] {
			t.Errorf(`MustParse("*/5 * * * *").NextN("2013-09-02 08:44:30", 5):\n"`)
			t.Errorf(`  result[%d]: expected "%s" but got "%s"`, i, expected[i], nextStr)
		}
	}
}

// Issue: https://github.com/gorhill/cronexpr/issues/16
func TestInterval_Interval60Issue(t *testing.T) {
	_, err := Parse("*/60 * * * * *")
	if err == nil {
		t.Errorf("parsing with interval 60 should return err")
	}

	_, err = Parse("*/61 * * * * *")
	if err == nil {
		t.Errorf("parsing with interval 61 should return err")
	}

	_, err = Parse("2/60 * * * * *")
	if err == nil {
		t.Errorf("parsing with interval 60 should return err")
	}

	_, err = Parse("2-20/61 * * * * *")
	if err == nil {
		t.Errorf("parsing with interval 60 should return err")
	}
}

/******************************************************************************/

var benchmarkExpressions = []string{
	"* * * * *",
	"@hourly",
	"@weekly",
	"@yearly",
	"30 3 15W 3/3 *",
	"30 0 0 1-31/5 Oct-Dec * 2000,2006,2008,2013-2015",
	"0 0 0 * Feb-Nov/2 thu#3 2000-2050",
}
var benchmarkExpressionsLen = len(benchmarkExpressions)

func BenchmarkParse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = MustParse(benchmarkExpressions[i%benchmarkExpressionsLen])
	}
}

func BenchmarkNext(b *testing.B) {
	exprs := make([]*Expression, benchmarkExpressionsLen)
	for i := 0; i < benchmarkExpressionsLen; i++ {
		exprs[i] = MustParse(benchmarkExpressions[i])
	}
	from := time.Now()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		expr := exprs[i%benchmarkExpressionsLen]
		next := expr.Next(from)
		next = expr.Next(next)
		next = expr.Next(next)
		next = expr.Next(next)
		next = expr.Next(next)
	}
}
