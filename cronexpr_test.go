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
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
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

	{
		"* * 16 L-3 2/3 ? *",
		"Mon 2006-01-02 15:04",
		[]crontimes{
			{"2024-05-15 00:00:00", "Tue 2024-05-28 16:00"},
			{"2024-05-30 00:00:00", "Wed 2024-08-28 16:00"},
			{"2024-08-30 00:00:00", "Wed 2024-11-27 16:00"},
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

func TestExpressionWithNonTrivialTimeZone(t *testing.T) {
	location, err := time.LoadLocation("Asia/Calcutta")
	cronExp := "0 59 23 1 9 * 2023" // At 11:59 PM, on day 1 of the month, only in September, only in 2023

	startTimeDST := time.Date(2023, time.September, 1, 23, 58, 0, 0, location)

	parsedExpression, err := Parse(cronExp)
	if err != nil {
		t.Errorf(`Failed to parse cron expression, error: %s`, err.Error())
	}

	nextTime := parsedExpression.Next(startTimeDST)
	expectedNextTime := startTimeDST.Add(1 * time.Minute)
	if nextTime.Unix() != expectedNextTime.Unix() {
		t.Errorf(`Incorrect next time. Want: %s, got: %s`, expectedNextTime, nextTime)
	}
}

var specificLocationTests = []locationTestCase{
	{"2019-03-31 02:30:00", "2019-04-01T02:30:00+01:00", "Europe/London"},
	{"2019-10-27 02:30:00", "2019-10-28T02:30:00Z", "Europe/London"},
	{"2019-03-31 02:30:00", "2019-04-01T02:30:00+02:00", "Europe/Paris"},
	{"2019-10-27 02:30:00", "2019-10-28T02:30:00+01:00", "Europe/Paris"},
	{"2019-02-09 02:30:00", "2019-02-10T02:30:00-05:00", "America/New_York"},
	{"2019-03-08 02:30:00", "2019-03-09T02:30:00-05:00", "America/New_York"},
	{"2019-03-09 02:30:00", "2019-03-11T02:30:00-04:00", "America/New_York"}, // Note that 02:30am on the 10th doesn't exist due to DST
	{"2019-03-10 03:30:00", "2019-03-11T02:30:00-04:00", "America/New_York"},
	{"2019-11-03 02:30:00", "2019-11-04T02:30:00-05:00", "America/New_York"},
	{"2019-11-04 02:30:00", "2019-11-05T02:30:00-05:00", "America/New_York"},
	{"2019-03-09 02:30:00", "2019-03-11T02:30:00-07:00", "America/Los_Angeles"}, // Note that 02:30am on the 10th doesn't exist due to DST
	{"2019-11-03 02:30:00", "2019-11-04T02:30:00-08:00", "America/Los_Angeles"},
	{"2019-03-09 02:30:00", "2019-03-11T03:30:00-06:00", "US/Pacific"}, // Note that 02:30am on the 10th doesn't exist due to DST
	{"2019-11-03 02:30:00", "2019-11-04T02:30:00-08:00", "US/Pacific"},
	{"2019-04-05 02:30:00", "2019-04-06T02:30:00+11:00", "Australia/Melbourne"},
	{"2019-10-04 02:30:00", "2019-10-05T02:30:00+10:00", "Australia/Melbourne"},
	{"2019-10-04 02:30:00", "2019-10-05T02:30:00-04:00", "America/Asuncion"},
	{"2019-03-09 02:30:00", "2019-03-10T02:30:00-03:00", "America/Sao_Paulo"},
	{"2019-11-03 02:30:00", "2019-11-04T02:30:00-03:00", "America/Sao_Paulo"},
	{"2018-02-15 02:30:00", "2018-02-16T02:30:00-02:00", "America/Sao_Paulo"},
	{"2018-02-16 02:30:00", "2018-02-17T02:30:00-02:00", "America/Sao_Paulo"},
	{"2018-02-17 02:30:00", "2018-02-18T02:30:00-03:00", "America/Sao_Paulo"}, // Note that the DST for Brazil until 2018 was at Feb 18 12:00am (clock would be turned back 1 hr)
}

func TestDailyExpressionsWithTimeZones(t *testing.T) {
	for _, test := range specificLocationTests {
		location, err := time.LoadLocation(test.where)
		require.NoError(t, err)

		from, err := time.ParseInLocation("2006-01-02 15:04:05", test.when, location)
		require.NoError(t, err)

		to, _ := time.Parse(time.RFC3339, test.expected)
		require.NoError(t, err)

		expr, err := Parse("0 30 2 * * * *") // At 02:30 AM
		require.NoError(t, err)

		expectedTm := to.UTC()
		resultTm := expr.Next(from).UTC()

		require.Equal(t, expectedTm, resultTm, test)
	}
}

/******************************************************************************/

func TestZero(t *testing.T) {
	from, _ := time.Parse("2006-01-02", "2013-08-31")
	next := MustParse("* * * * * 1980").Next(from)
	if next.IsZero() == true {
		t.Error(`("* * * * * 1980").Next("2013-08-31").IsZero() returned 'true', expected 'false'`)
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

func TestPeriodicConfig_DSTChange_Transitions(t *testing.T) {
	locName := "America/Los_Angeles"
	loc, err := time.LoadLocation(locName)
	require.NoError(t, err)

	cases := []struct {
		name     string
		pattern  string
		initTime time.Time
		expected []time.Time
	}{
		{
			"normal time",
			"0 2 * * * 2019",
			time.Date(2019, time.February, 7, 1, 0, 0, 0, loc),
			[]time.Time{
				time.Date(2019, time.February, 7, 2, 0, 0, 0, loc),
				time.Date(2019, time.February, 8, 2, 0, 0, 0, loc),
				time.Date(2019, time.February, 9, 2, 0, 0, 0, loc),
			},
		},
		{
			"Spring forward but not in switch time",
			"0 4 * * * 2019",
			time.Date(2019, time.March, 9, 1, 0, 0, 0, loc),
			[]time.Time{
				time.Date(2019, time.March, 9, 4, 0, 0, 0, loc),
				time.Date(2019, time.March, 10, 4, 0, 0, 0, loc),
				time.Date(2019, time.March, 11, 4, 0, 0, 0, loc),
			},
		},
		{
			"Spring forward at a skipped time odd",
			"2 2 * * * 2019",
			time.Date(2019, time.March, 9, 1, 0, 0, 0, loc),
			[]time.Time{
				time.Date(2019, time.March, 9, 2, 2, 0, 0, loc),
				// no time in March 10!
				time.Date(2019, time.March, 11, 2, 2, 0, 0, loc),
				time.Date(2019, time.March, 12, 2, 2, 0, 0, loc),
			},
		},
		{
			"Spring forward at a skipped time",
			"1 2 * * * 2019",
			time.Date(2019, time.March, 9, 1, 0, 0, 0, loc),
			[]time.Time{
				time.Date(2019, time.March, 9, 2, 1, 0, 0, loc),
				// no time in March 8!
				time.Date(2019, time.March, 11, 2, 1, 0, 0, loc),
				time.Date(2019, time.March, 12, 2, 1, 0, 0, loc),
			},
		},
		{
			"Spring forward at a skipped time boundary",
			"0 2 * * * 2019",
			time.Date(2019, time.March, 9, 1, 0, 0, 0, loc),
			[]time.Time{
				time.Date(2019, time.March, 9, 2, 0, 0, 0, loc),
				// no time in March 8!
				time.Date(2019, time.March, 11, 2, 0, 0, 0, loc),
				time.Date(2019, time.March, 12, 2, 0, 0, 0, loc),
			},
		},
		{
			"Spring forward at a boundary of repeating time",
			"0 1 * * * 2019",
			time.Date(2019, time.March, 9, 0, 0, 0, 0, loc),
			[]time.Time{
				time.Date(2019, time.March, 9, 1, 0, 0, 0, loc),
				time.Date(2019, time.March, 10, 0, 0, 0, 0, loc).Add(1 * time.Hour),
				time.Date(2019, time.March, 11, 1, 0, 0, 0, loc),
				time.Date(2019, time.March, 12, 1, 0, 0, 0, loc),
			},
		},
		{
			"Fall back: before transition",
			"30 0 * * * 2019",
			time.Date(2019, time.November, 3, 0, 0, 0, 0, loc),
			[]time.Time{
				time.Date(2019, time.November, 3, 0, 30, 0, 0, loc),
				time.Date(2019, time.November, 4, 0, 30, 0, 0, loc),
				time.Date(2019, time.November, 5, 0, 30, 0, 0, loc),
				time.Date(2019, time.November, 6, 0, 30, 0, 0, loc),
			},
		},
		{
			"Fall back: after transition",
			"30 3 * * * 2019",
			time.Date(2019, time.November, 3, 0, 0, 0, 0, loc),
			[]time.Time{
				time.Date(2019, time.November, 3, 3, 30, 0, 0, loc),
				time.Date(2019, time.November, 4, 3, 30, 0, 0, loc),
				time.Date(2019, time.November, 5, 3, 30, 0, 0, loc),
				time.Date(2019, time.November, 6, 3, 30, 0, 0, loc),
			},
		},
		{
			"Fall back: after transition starting in repeated span before",
			"30 3 * * * 2019",
			time.Date(2019, time.November, 3, 0, 10, 0, 0, loc).Add(1 * time.Hour),
			[]time.Time{
				time.Date(2019, time.November, 3, 3, 30, 0, 0, loc),
				time.Date(2019, time.November, 4, 3, 30, 0, 0, loc),
				time.Date(2019, time.November, 5, 3, 30, 0, 0, loc),
				time.Date(2019, time.November, 6, 3, 30, 0, 0, loc),
			},
		},
		{
			"Fall back: after transition starting in repeated span after",
			"30 3 * * * 2019",
			time.Date(2019, time.November, 3, 0, 10, 0, 0, loc).Add(2 * time.Hour),
			[]time.Time{
				time.Date(2019, time.November, 3, 3, 30, 0, 0, loc),
				time.Date(2019, time.November, 4, 3, 30, 0, 0, loc),
				time.Date(2019, time.November, 5, 3, 30, 0, 0, loc),
				time.Date(2019, time.November, 6, 3, 30, 0, 0, loc),
			},
		},
		{
			"Fall back: in repeated region",
			"30 1 * * * 2019",
			time.Date(2019, time.November, 3, 0, 0, 0, 0, loc),
			[]time.Time{
				time.Date(2019, time.November, 3, 0, 30, 0, 0, loc).Add(1 * time.Hour),
				time.Date(2019, time.November, 3, 0, 30, 0, 0, loc).Add(2 * time.Hour),
				time.Date(2019, time.November, 4, 1, 30, 0, 0, loc),
				time.Date(2019, time.November, 5, 1, 30, 0, 0, loc),
				time.Date(2019, time.November, 6, 1, 30, 0, 0, loc),
			},
		},
		{
			"Fall back: in repeated region boundary",
			"0 1 * * * 2019",
			time.Date(2019, time.November, 3, 0, 0, 0, 0, loc),
			[]time.Time{
				time.Date(2019, time.November, 3, 0, 0, 0, 0, loc).Add(1 * time.Hour),
				time.Date(2019, time.November, 3, 0, 0, 0, 0, loc).Add(2 * time.Hour),
				time.Date(2019, time.November, 4, 1, 0, 0, 0, loc),
				time.Date(2019, time.November, 5, 1, 0, 0, 0, loc),
				time.Date(2019, time.November, 6, 1, 0, 0, 0, loc),
			},
		},
		{
			"Fall back: in repeated region boundary 2",
			"0 2 * * * 2019",
			time.Date(2019, time.November, 3, 0, 0, 0, 0, loc),
			[]time.Time{
				time.Date(2019, time.November, 3, 0, 0, 0, 0, loc).Add(3 * time.Hour),
				time.Date(2019, time.November, 4, 2, 0, 0, 0, loc),
				time.Date(2019, time.November, 5, 2, 0, 0, 0, loc),
				time.Date(2019, time.November, 6, 2, 0, 0, 0, loc),
			},
		},
		{
			"Fall back: in repeated region, starting from within region",
			"30 1 * * * 2019",
			time.Date(2019, time.November, 3, 0, 40, 0, 0, loc).Add(1 * time.Hour),
			[]time.Time{
				time.Date(2019, time.November, 3, 0, 30, 0, 0, loc).Add(2 * time.Hour),
				time.Date(2019, time.November, 4, 1, 30, 0, 0, loc),
				time.Date(2019, time.November, 5, 1, 30, 0, 0, loc),
				time.Date(2019, time.November, 6, 1, 30, 0, 0, loc),
			},
		},
		{
			"Fall back: in repeated region, starting from within region 2",
			"30 1 * * * 2019",
			time.Date(2019, time.November, 3, 0, 40, 0, 0, loc).Add(2 * time.Hour),
			[]time.Time{
				time.Date(2019, time.November, 4, 1, 30, 0, 0, loc),
				time.Date(2019, time.November, 5, 1, 30, 0, 0, loc),
				time.Date(2019, time.November, 6, 1, 30, 0, 0, loc),
			},
		},
		{
			"Fall back: wildcard",
			"30 * * * * 2019",
			time.Date(2019, time.November, 3, 0, 0, 0, 0, loc),
			[]time.Time{
				time.Date(2019, time.November, 3, 0, 30, 0, 0, loc),
				time.Date(2019, time.November, 3, 0, 30, 0, 0, loc).Add(1 * time.Hour),
				time.Date(2019, time.November, 3, 0, 30, 0, 0, loc).Add(2 * time.Hour),
				time.Date(2019, time.November, 3, 2, 30, 0, 0, loc),
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			expr := MustParse(c.pattern)

			starting := c.initTime
			for _, next := range c.expected {
				n := expr.Next(starting)
				if next != n {
					t.Fatalf("next(%v) = %v not %v", starting, next, n)
				}

				starting = next
			}
		})
	}
}

func TestPeriodicConfig_DSTChange_Transitions_LordHowe(t *testing.T) {
	locName := "Australia/Lord_Howe"
	loc, err := time.LoadLocation(locName)
	require.NoError(t, err)

	cases := []struct {
		name     string
		pattern  string
		initTime time.Time
		expected []time.Time
	}{
		{
			"normal time",
			"0 2 * * * 2019",
			time.Date(2019, time.February, 7, 1, 0, 0, 0, loc),
			[]time.Time{
				time.Date(2019, time.February, 7, 2, 0, 0, 0, loc),
				time.Date(2019, time.February, 8, 2, 0, 0, 0, loc),
				time.Date(2019, time.February, 9, 2, 0, 0, 0, loc),
			},
		},
		{
			"backward: non repeated portion of the hour",
			"3 1 * * * 2019",
			time.Date(2019, time.April, 6, 0, 0, 0, 0, loc),
			[]time.Time{
				time.Date(2019, time.April, 6, 1, 3, 0, 0, loc),
				time.Date(2019, time.April, 7, 1, 3, 0, 0, loc),
				time.Date(2019, time.April, 8, 1, 3, 0, 0, loc),
				time.Date(2019, time.April, 9, 1, 3, 0, 0, loc),
			},
		},
		{
			"backward: repeated portion of the hour",
			"31 1 * * * 2019",
			time.Date(2019, time.April, 6, 0, 0, 0, 0, loc),
			[]time.Time{
				time.Date(2019, time.April, 6, 1, 31, 0, 0, loc),
				time.Date(2019, time.April, 7, 0, 31, 0, 0, loc).Add(60 * time.Minute),
				time.Date(2019, time.April, 7, 0, 31, 0, 0, loc).Add(90 * time.Minute),
				time.Date(2019, time.April, 8, 1, 31, 0, 0, loc),
				time.Date(2019, time.April, 9, 1, 31, 0, 0, loc),
			},
		},
		{
			"forward: skipped portion of the hour",
			"3 2 * * * 2019",
			time.Date(2019, time.October, 5, 0, 0, 0, 0, loc),
			[]time.Time{
				time.Date(2019, time.October, 5, 2, 3, 0, 0, loc),
				// no Oct 6
				time.Date(2019, time.October, 7, 2, 3, 0, 0, loc),
				time.Date(2019, time.October, 8, 2, 3, 0, 0, loc),
				time.Date(2019, time.October, 9, 2, 3, 0, 0, loc),
			},
		},
		{
			"forward: non-skipped portion of the hour",
			"31 2 * * * 2019",
			time.Date(2019, time.October, 5, 0, 0, 0, 0, loc),
			[]time.Time{
				time.Date(2019, time.October, 5, 2, 31, 0, 0, loc),
				// no Oct 6
				time.Date(2019, time.October, 7, 2, 31, 0, 0, loc),
				time.Date(2019, time.October, 8, 2, 31, 0, 0, loc),
				time.Date(2019, time.October, 9, 2, 31, 0, 0, loc),
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			expr := MustParse(c.pattern)

			starting := c.initTime
			for _, next := range c.expected {
				n := expr.Next(starting)
				if next != n {
					t.Fatalf("next(%v) = %v not %v", starting, next, n)
				}

				starting = next
			}
		})
	}
}

func TestNext_DaylightSaving_Property(t *testing.T) {
	locName := "America/Los_Angeles"
	loc, err := time.LoadLocation(locName)
	if err != nil {
		t.Fatalf("failed to get location: %v", err)
	}

	cronExprs := []string{
		"* * * * *",
		"0 2 * * *",
		"* 1 * * *",
	}

	times := []time.Time{
		// spring forward
		time.Date(2019, time.March, 11, 0, 0, 0, 0, loc),
		time.Date(2019, time.March, 10, 0, 0, 0, 0, loc),
		time.Date(2019, time.March, 11, 0, 0, 0, 0, loc),

		// leap backwards
		time.Date(2019, time.November, 4, 0, 0, 0, 0, loc),
		time.Date(2019, time.November, 5, 0, 0, 0, 0, loc),
		time.Date(2019, time.November, 6, 0, 0, 0, 0, loc),
	}

	testSpan := 4 * time.Hour

	testCase := func(t *testing.T, cronExpr string, init time.Time) {
		cron := MustParse(cronExpr)

		prevNext := init
		for start := init; start.Before(init.Add(testSpan)); start = start.Add(1 * time.Minute) {
			next := cron.Next(start)
			if !next.After(start) {
				t.Fatalf("next(%v) = %v is not after start time", start, next)
			}

			if next.Before(prevNext) {
				t.Fatalf("next(%v) = %v reverted back in time from %v", start, next, prevNext)
			}

			if strings.HasPrefix(cronExpr, "* * ") {
				if next.Sub(start) != time.Minute {
					t.Fatalf("next(%v) = %v should be the next minute", start, next)
				}
			}

			prevNext = next
		}
	}

	for _, cron := range cronExprs {
		for _, startTime := range times {
			t.Run(fmt.Sprintf("%v: %v", cron, startTime), func(t *testing.T) {
				testCase(t, cron, startTime)
			})
		}
	}
}

func TestNext_DaylightSaving_Property_LordHowe(t *testing.T) {
	// Lord Howe, Australia is at GMT+1100 April-October and GMT+1030 otherwise.
	//
	// On April 7, 2019, at when clock approches 2am, the clock
	// transitions to 1.30am.
	//
	// On October 6, when the clock approaches 2am, the clock transitions
	// to 2.30am.
	locName := "Australia/Lord_Howe"
	loc, err := time.LoadLocation(locName)
	if err != nil {
		t.Fatalf("failed to get location: %v", err)
	}

	cronExprs := []string{
		"* * * * *",
		"0 2 * * *",
		"* 1 * * *",
		"35 1 * * *",
		"5 2 * * *",
	}

	times := []time.Time{
		// spring forward
		time.Date(2019, time.April, 6, 0, 0, 0, 0, loc),
		time.Date(2019, time.April, 7, 0, 0, 0, 0, loc),
		time.Date(2019, time.April, 8, 0, 0, 0, 0, loc),

		// leap backwards
		time.Date(2019, time.October, 5, 0, 0, 0, 0, loc),
		time.Date(2019, time.October, 6, 0, 0, 0, 0, loc),
		time.Date(2019, time.October, 7, 0, 0, 0, 0, loc),
	}

	testSpan := 4 * time.Hour

	testCase := func(t *testing.T, cronExpr string, init time.Time) {
		cron := MustParse(cronExpr)

		prevNext := init
		for start := init; start.Before(init.Add(testSpan)); start = start.Add(1 * time.Minute) {
			next := cron.Next(start)
			if !next.After(start) {
				t.Fatalf("next(%v) = %v is not after start time", start, next)
			}

			if next.Before(prevNext) {
				t.Fatalf("next(%v) = %v reverted back in time from %v", start, next, prevNext)
			}

			if strings.HasPrefix(cronExpr, "* * ") {
				if next.Sub(start) != time.Minute {
					t.Fatalf("next(%v) = %v should be the next minute", start, next)
				}
			}

			prevNext = next
		}
	}

	for _, cron := range cronExprs {
		for _, startTime := range times {
			t.Run(fmt.Sprintf("%v: %v", cron, startTime), func(t *testing.T) {
				testCase(t, cron, startTime)
			})
		}
	}
}

func TestNext_DaylightSaving_Property_Brazil(t *testing.T) {
	// Until 2018, Brazil/Sao Paulo and some South American countries used
	// to transition for daylight savings at midnight.
	//
	// When the clock approaches 2018-11-04 midnight, the clock transitions to 1am.
	locName := "America/Sao_Paulo"
	loc, err := time.LoadLocation(locName)
	if err != nil {
		t.Fatalf("failed to get location: %v", err)
	}

	cronExprs := []string{
		"* * * * *",
		"0 2 * * *",
		"* 1 * * *",
		"5 1 * * *",
		"5 23 * * *",
	}

	times := []time.Time{
		// spring forward
		time.Date(2018, time.February, 16, 22, 0, 0, 0, loc),
		time.Date(2018, time.February, 17, 22, 0, 0, 0, loc),
		time.Date(2018, time.February, 18, 22, 0, 0, 0, loc),

		// leap backwards
		time.Date(2018, time.November, 3, 23, 0, 0, 0, loc),
		time.Date(2018, time.November, 3, 23, 0, 0, 0, loc),
		time.Date(2018, time.November, 3, 23, 0, 0, 0, loc),
	}

	testSpan := 4 * time.Hour

	testCase := func(t *testing.T, cronExpr string, init time.Time) {
		cron := MustParse(cronExpr)

		prevNext := init
		for start := init; start.Before(init.Add(testSpan)); start = start.Add(1 * time.Minute) {
			next := cron.Next(start)
			if !next.After(start) {
				t.Fatalf("next(%v) = %v is not after start time", start, next)
			}

			if next.Before(prevNext) {
				t.Fatalf("next(%v) = %v reverted back in time from %v", start, next, prevNext)
			}

			if strings.HasPrefix(cronExpr, "* * ") {
				if next.Sub(start) != time.Minute {
					t.Fatalf("next(%v) = %v should be the next minute", start, next)
				}
			}

			prevNext = next
		}
	}

	for _, cron := range cronExprs {
		for _, startTime := range times {
			t.Run(fmt.Sprintf("%v: %v", cron, startTime), func(t *testing.T) {
				testCase(t, cron, startTime)
			})
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
