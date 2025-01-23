Golang Cron expression parser
=============================
Given a cron expression and a timestamp, you can get the next timestamp which satisfies the cron expression.

This standalone library is only to parse and apply timestamps to cron expressions.

There is also a companion command-line utility to evaluate cron time expressions: <https://github.com/gorhill/cronexpr/tree/master/cronexpr>, which uses this library.

Implementation
--------------
The reference documentation for this implementation is given below:

    Field name     Mandatory?   Allowed values    Allowed special characters
    ----------     ----------   --------------    --------------------------
    Seconds        No           0-59              * / , -
    Minutes        Yes          0-59              * / , -
    Hours          Yes          0-23              * / , -
    Day of month   Yes          1-31              * / , - L W
    Month          Yes          1-12 or JAN-DEC   * / , -
    Day of week    Yes          0-6 or SUN-SAT    * / , - L #
    Year           No           1970–2099         * / , -

#### Asterisk ( * )
The asterisk indicates that the cron expression matches for all values of the field. E.g., using an asterisk in the 4th field (month) indicates every month.

#### Slash ( / )
Slashes describe increments of ranges. For example `3-59/15` in the minute field indicate the third minute of the hour and every 15 minutes thereafter. The form `*/...` is equivalent to the form "first-last/...", that is, an increment over the largest possible range of the field.

#### Comma ( , )
Commas are used to separate items of a list. For example, using `MON,WED,FRI` in the 5th field (day of week) means Mondays, Wednesdays and Fridays.

#### Hyphen ( - )
Hyphens define ranges. For example, 2000-2010 indicates every year between 2000 and 2010 AD, inclusive.

#### L
`L` stands for "last". When used in the day-of-week field, it allows you to specify constructs such as "the last Friday" (`5L`) of a given month. In the day-of-month field, it specifies the last day of the month.

#### W
The `W` character is allowed for the day-of-month field. This character is used to specify the business day (Monday-Friday) nearest the given day. As an example, if you were to specify `15W` as the value for the day-of-month field, the meaning is: "the nearest business day to the 15th of the month."

So, if the 15th is a Saturday, the trigger fires on Friday the 14th. If the 15th is a Sunday, the trigger fires on Monday the 16th. If the 15th is a Tuesday, then it fires on Tuesday the 15th. However if you specify `1W` as the value for day-of-month, and the 1st is a Saturday, the trigger fires on Monday the 3rd, as it does not 'jump' over the boundary of a month's days.

The `W` character can be specified only when the day-of-month is a single day, not a range or list of days.

The `W` character can also be combined with `L`, i.e. `LW` to mean "the last business day of the month."

#### Hash ( # )
`#` is allowed for the day-of-week field, and must be followed by a number between one and five. It allows you to specify constructs such as "the second Friday" of a given month.

Pre-defined cron expressions
---------------------------

    Entry       Description                                                             Equivalent to
    @annually   Run once a year at midnight in the morning of January 1                 0 0 0 1 1 * *
    @yearly     Run once a year at midnight in the morning of January 1                 0 0 0 1 1 * *
    @monthly    Run once a month at midnight in the morning of the first of the month   0 0 0 1 * * *
    @weekly     Run once a week at midnight in the morning of Sunday                    0 0 0 * * 0 *
    @daily      Run once a day at midnight                                              0 0 0 * * * *
    @hourly     Run once an hour at the beginning of the hour                           0 0 * * * * *
    @reboot     Not supported

Other details
-------------
* If only six fields are present, a `0` second field is prepended, that is, `* * * * * 2013` internally become `0 * * * * * 2013`.
* If only five fields are present, a `0` second field is prepended and a wildcard year field is appended, that is, `* * * * Mon` internally become `0 * * * * Mon *`.
* Domain for day-of-week field is [0-7] instead of [0-6], 7 being Sunday (like 0). This to comply with http://linux.die.net/man/5/crontab#.
* As of now, the behaviour of the code is undetermined if a malformed cron expression is supplied

Install
-------
    go get github.com/thought-machine/cronexpr

Usage
-----
Import the library:

    import "github.com/thought-machine/cronexpr"
    import "time"

Simplest way:

    nextTime := cronexpr.MustParse("0 0 29 2 *").Next(time.Now())

Assuming `time.Now()` is "2013-08-29 09:28:00", then `nextTime` will be "2016-02-29 00:00:00".

You can keep the returned Expression pointer around if you want to reuse it:

    expr := cronexpr.MustParse("0 0 29 2 *")
    nextTime := expr.Next(time.Now())
    ...
    nextTime = expr.Next(nextTime)

Use `time.IsZero()` to find out whether a valid time was returned. For example,

    cronexpr.MustParse("* * * * * 1980").Next(time.Now()).IsZero()

will return `true`, whereas

    cronexpr.MustParse("* * * * * 2050").Next(time.Now()).IsZero()

will return `false` (as of 2013-08-29...)

You may also query for `n` next time stamps:

    cronexpr.MustParse("0 0 29 2 *").NextN(time.Now(), 5)

which returns a slice of time.Time objects, containing the following time stamps (as of 2013-08-30):

    2016-02-29 00:00:00
    2020-02-29 00:00:00
    2024-02-29 00:00:00
    2028-02-29 00:00:00
    2032-02-29 00:00:00

The time zone of time values returned by `Next` and `NextN` is always the
time zone of the time value passed as argument, unless a zero time value is
returned.

API
---
<http://godoc.org/github.com/thought-machine/cronexpr>

License
-------

License: pick the one which suits you best:

- GPL v3 see <https://www.gnu.org/licenses/gpl.html>
- APL v2 see <http://www.apache.org/licenses/LICENSE-2.0>