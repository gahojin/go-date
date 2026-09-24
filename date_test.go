package date_test

import (
	"testing"
	"time"

	"github.com/gahojin/go-date"
	"github.com/stretchr/testify/assert"
)

// TestNew は New コンストラクタが正しく Date を生成し正規化するかをテストする。
func TestNew(t *testing.T) {
	tests := []struct {
		name      string
		year      int
		month     time.Month
		day       int
		wantYear  int
		wantMonth time.Month
		wantDay   int
	}{
		{
			name:      "standard date",
			year:      2024,
			month:     time.May,
			day:       1,
			wantYear:  2024,
			wantMonth: time.May,
			wantDay:   1,
		},
		{
			name:      "leap year feb 29",
			year:      2024,
			month:     time.February,
			day:       29,
			wantYear:  2024,
			wantMonth: time.February,
			wantDay:   29,
		},
		{
			name:      "normalize Jan 32 to Feb 1",
			year:      2024,
			month:     time.January,
			day:       32,
			wantYear:  2024,
			wantMonth: time.February,
			wantDay:   1,
		},
		{
			name:      "normalize day 0 to previous month end",
			year:      2024,
			month:     time.March,
			day:       0,
			wantYear:  2024,
			wantMonth: time.February,
			wantDay:   29,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := date.New(tt.year, tt.month, tt.day)
			assert.Equal(t, tt.wantYear, d.Year())
			assert.Equal(t, tt.wantMonth, d.Month())
			assert.Equal(t, tt.wantDay, d.Day())
		})
	}
}

// TestFromTime は FromTime 関数が time.Time から Date を正しく生成するかをテストする。
func TestFromTime(t *testing.T) {
	jst := time.FixedZone("JST", 9*60*60)

	tests := []struct {
		name string
		t    time.Time
		want date.Date
	}{
		{
			name: "UTC time",
			t:    time.Date(2024, 5, 1, 10, 0, 0, 0, time.UTC),
			want: date.New(2024, 5, 1),
		},
		{
			name: "keeps the time.Time's own location (JST, no conversion)",
			t:    time.Date(2024, 5, 1, 23, 30, 0, 0, jst),
			want: date.New(2024, 5, 1),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := date.FromTime(tt.t)
			assert.True(t, tt.want.Equal(got))
		})
	}
}

// TestFromTimeIn は FromTimeIn 関数をテストする。
func TestFromTimeIn(t *testing.T) {
	jst := time.FixedZone("JST", 9*60*60)

	tests := []struct {
		name string
		tm   time.Time
		loc  *time.Location
		want date.Date
	}{
		{
			name: "converts UTC to JST, crossing the date boundary",
			tm:   time.Date(2024, 5, 1, 23, 30, 0, 0, time.UTC),
			loc:  jst,
			want: date.New(2024, 5, 2),
		},
		{
			name: "converts JST to UTC, crossing the date boundary backwards",
			tm:   time.Date(2024, 5, 1, 1, 0, 0, 0, jst),
			loc:  time.UTC,
			want: date.New(2024, 4, 30),
		},
		{
			name: "nil location defaults to time.Local",
			tm:   time.Date(2024, 5, 1, 10, 0, 0, 0, time.UTC),
			loc:  nil,
			want: date.FromTime(time.Date(2024, 5, 1, 10, 0, 0, 0, time.UTC).In(time.Local)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := date.FromTimeIn(tt.tm, tt.loc)
			assert.True(t, tt.want.Equal(got))
		})
	}
}

// TestDate_Year_Month_Day は Year, Month, Day メソッドが各値を正しく返すかをテストする。
func TestDate_Year_Month_Day(t *testing.T) {
	d := date.New(2024, time.May, 1)
	assert.Equal(t, 2024, d.Year())
	assert.Equal(t, time.May, d.Month())
	assert.Equal(t, 1, d.Day())

	var zero date.Date
	assert.Equal(t, 0, zero.Year())
	assert.Equal(t, time.Month(0), zero.Month())
	assert.Equal(t, 0, zero.Day())
}

// TestDate_CompactString は CompactString メソッドが YYYYMMDD 形式の文字列を返すかをテストする。
func TestDate_CompactString(t *testing.T) {
	d := date.New(2024, time.May, 1)
	assert.Equal(t, "20240501", d.CompactString())

	d2 := date.New(2024, time.December, 31)
	assert.Equal(t, "20241231", d2.CompactString())
}

// TestDate_AddDate は AddDate メソッドが指定された年月日を加算したDateを返すかをテストする。
func TestDate_AddDate(t *testing.T) {
	tests := []struct {
		name   string
		d      date.Date
		years  int
		months int
		days   int
		want   date.Date
	}{
		{
			name:   "add days",
			d:      date.New(2024, time.May, 1),
			years:  0,
			months: 0,
			days:   10,
			want:   date.New(2024, time.May, 11),
		},
		{
			name:   "add months",
			d:      date.New(2024, time.May, 1),
			years:  0,
			months: 2,
			days:   0,
			want:   date.New(2024, time.July, 1),
		},
		{
			name:   "add years",
			d:      date.New(2024, time.May, 1),
			years:  3,
			months: 0,
			days:   0,
			want:   date.New(2027, time.May, 1),
		},
		{
			name:   "subtract days across month boundary",
			d:      date.New(2024, time.May, 1),
			years:  0,
			months: 0,
			days:   -1,
			want:   date.New(2024, time.April, 30),
		},
		{
			name:   "leap year feb 29 + 1 year",
			d:      date.New(2024, time.February, 29),
			years:  1,
			months: 0,
			days:   0,
			want:   date.New(2025, time.March, 1),
		},
		{
			name:   "composite add",
			d:      date.New(2024, time.January, 15),
			years:  1,
			months: 2,
			days:   5,
			want:   date.New(2025, time.March, 20),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.d.AddDate(tt.years, tt.months, tt.days)
			assert.Equal(t, tt.want, got)
		})
	}
}

// TestDate_FirstDayOfMonth は FirstDayOfMonth メソッドがその月の初日（1日）を返すかをテストする。
func TestDate_FirstDayOfMonth(t *testing.T) {
	tests := []struct {
		name string
		d    date.Date
		want date.Date
	}{
		{
			name: "normal date",
			d:    date.New(2024, time.May, 15),
			want: date.New(2024, time.May, 1),
		},
		{
			name: "already first day of month",
			d:    date.New(2024, time.May, 1),
			want: date.New(2024, time.May, 1),
		},
		{
			name: "end of month (leap year feb)",
			d:    date.New(2024, time.February, 29),
			want: date.New(2024, time.February, 1),
		},
		{
			name: "zero value",
			d:    date.Date{},
			want: date.Date{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.d.FirstDayOfMonth()
			assert.Equal(t, tt.want, got)
		})
	}
}

// TestDate_LastDayOfMonth は LastDayOfMonth メソッドがその月の末日を返すかをテストする。
func TestDate_LastDayOfMonth(t *testing.T) {
	tests := []struct {
		name string
		d    date.Date
		want date.Date
	}{
		{
			name: "normal date (30-day month)",
			d:    date.New(2024, time.April, 15),
			want: date.New(2024, time.April, 30),
		},
		{
			name: "normal date (31-day month)",
			d:    date.New(2024, time.May, 15),
			want: date.New(2024, time.May, 31),
		},
		{
			name: "already last day of month",
			d:    date.New(2024, time.May, 31),
			want: date.New(2024, time.May, 31),
		},
		{
			name: "first day of month",
			d:    date.New(2024, time.May, 1),
			want: date.New(2024, time.May, 31),
		},
		{
			name: "february in leap year",
			d:    date.New(2024, time.February, 15),
			want: date.New(2024, time.February, 29),
		},
		{
			name: "february in non-leap year",
			d:    date.New(2023, time.February, 15),
			want: date.New(2023, time.February, 28),
		},
		{
			name: "december (year boundary)",
			d:    date.New(2024, time.December, 15),
			want: date.New(2024, time.December, 31),
		},
		{
			name: "zero value",
			d:    date.Date{},
			want: date.Date{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.d.LastDayOfMonth()
			assert.Equal(t, tt.want, got)
		})
	}
}

// TestDate_IsZero は IsZero メソッドがゼロ値を正しく判定できるかをテストする。
func TestDate_IsZero(t *testing.T) {
	tests := []struct {
		name string
		d    date.Date
		want bool
	}{
		{
			name: "zero value",
			d:    date.Date{},
			want: true,
		},
		{
			name: "non-zero date",
			d:    date.New(2024, time.May, 1),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.d.IsZero())
		})
	}
}

// TestDate_Equal は Equal メソッドが二つの Date が等しいかを正しく判定できるかをテストする。
func TestDate_Equal(t *testing.T) {
	tests := []struct {
		name  string
		d     date.Date
		other date.Date
		want  bool
	}{
		{
			name:  "equal dates",
			d:     date.New(2024, time.May, 1),
			other: date.New(2024, time.May, 1),
			want:  true,
		},
		{
			name:  "different year",
			d:     date.New(2024, time.May, 1),
			other: date.New(2025, time.May, 1),
			want:  false,
		},
		{
			name:  "different month",
			d:     date.New(2024, time.May, 1),
			other: date.New(2024, time.June, 1),
			want:  false,
		},
		{
			name:  "different day",
			d:     date.New(2024, time.May, 1),
			other: date.New(2024, time.May, 2),
			want:  false,
		},
		{
			name:  "both zero values",
			d:     date.Date{},
			other: date.Date{},
			want:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.d.Equal(tt.other))
		})
	}
}

// TestDate_EqualMonth は EqualMonth メソッドが年と月のみを比較して判定できるかをテストする。
func TestDate_EqualMonth(t *testing.T) {
	tests := []struct {
		name  string
		d     date.Date
		other date.Date
		want  bool
	}{
		{
			name:  "same year and month, same day",
			d:     date.New(2024, time.May, 1),
			other: date.New(2024, time.May, 1),
			want:  true,
		},
		{
			name:  "same year and month, different day",
			d:     date.New(2024, time.May, 1),
			other: date.New(2024, time.May, 20),
			want:  true,
		},
		{
			name:  "different year, same month and day",
			d:     date.New(2024, time.May, 1),
			other: date.New(2025, time.May, 1),
			want:  false,
		},
		{
			name:  "same year, different month, same day",
			d:     date.New(2024, time.May, 1),
			other: date.New(2024, time.June, 1),
			want:  false,
		},
		{
			name:  "both zero values",
			d:     date.Date{},
			other: date.Date{},
			want:  true,
		},
		{
			name:  "one zero value",
			d:     date.New(2024, time.May, 1),
			other: date.Date{},
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.d.EqualMonth(tt.other))
		})
	}
}

// TestDate_Compare は Compare メソッドが二つの Date の大小・前後関係を正しく比較できるかをテストする。
func TestDate_Compare(t *testing.T) {
	tests := []struct {
		name  string
		d     date.Date
		other date.Date
		want  int
	}{
		{
			name:  "equal",
			d:     date.New(2024, time.May, 10),
			other: date.New(2024, time.May, 10),
			want:  0,
		},
		{
			name:  "earlier year",
			d:     date.New(2023, time.May, 10),
			other: date.New(2024, time.May, 10),
			want:  -1,
		},
		{
			name:  "later year",
			d:     date.New(2025, time.May, 10),
			other: date.New(2024, time.May, 10),
			want:  1,
		},
		{
			name:  "earlier month",
			d:     date.New(2024, time.April, 10),
			other: date.New(2024, time.May, 10),
			want:  -1,
		},
		{
			name:  "later month",
			d:     date.New(2024, time.June, 10),
			other: date.New(2024, time.May, 10),
			want:  1,
		},
		{
			name:  "earlier day",
			d:     date.New(2024, time.May, 9),
			other: date.New(2024, time.May, 10),
			want:  -1,
		},
		{
			name:  "later day",
			d:     date.New(2024, time.May, 11),
			other: date.New(2024, time.May, 10),
			want:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.d.Compare(tt.other))
		})
	}
}

// TestDate_Before は Before メソッドが指定した日付より前かを正しく判定できるかをテストする。
func TestDate_Before(t *testing.T) {
	d1 := date.New(2024, time.May, 1)
	d2 := date.New(2024, time.May, 2)

	assert.True(t, d1.Before(d2))
	assert.False(t, d2.Before(d1))
	assert.False(t, d1.Before(d1))
}

// TestDate_After は After メソッドが指定した日付より後かを正しく判定できるかをテストする。
func TestDate_After(t *testing.T) {
	d1 := date.New(2024, time.May, 2)
	d2 := date.New(2024, time.May, 1)

	assert.True(t, d1.After(d2))
	assert.False(t, d2.After(d1))
	assert.False(t, d1.After(d1))
}

// TestOrder は Order 関数が2つの日付を昇順に並べ替えて返すかをテストする。
func TestOrder(t *testing.T) {
	d1 := date.New(2024, time.January, 1)
	d2 := date.New(2024, time.May, 1)

	tests := []struct {
		name      string
		a         date.Date
		b         date.Date
		wantFirst date.Date
		wantSec   date.Date
	}{
		{
			name:      "already in order",
			a:         d1,
			b:         d2,
			wantFirst: d1,
			wantSec:   d2,
		},
		{
			name:      "reverse order (swapped)",
			a:         d2,
			b:         d1,
			wantFirst: d1,
			wantSec:   d2,
		},
		{
			name:      "equal dates",
			a:         d1,
			b:         d1,
			wantFirst: d1,
			wantSec:   d1,
		},
		{
			name:      "both zero values",
			a:         date.Date{},
			b:         date.Date{},
			wantFirst: date.Date{},
			wantSec:   date.Date{},
		},
		{
			name:      "zero value and non-zero value (zero is earlier)",
			a:         date.Date{},
			b:         d1,
			wantFirst: date.Date{},
			wantSec:   d1,
		},
		{
			name:      "non-zero value and zero value (swapped)",
			a:         d1,
			b:         date.Date{},
			wantFirst: date.Date{},
			wantSec:   d1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFirst, gotSec := date.Order(tt.a, tt.b)
			assert.Equal(t, tt.wantFirst, gotFirst)
			assert.Equal(t, tt.wantSec, gotSec)

			// メソッド呼び出しの挙動も同時に検証
			gotMethodFirst, gotMethodSec := tt.a.Order(tt.b)
			assert.Equal(t, tt.wantFirst, gotMethodFirst)
			assert.Equal(t, tt.wantSec, gotMethodSec)
		})
	}
}

// TestDaysIn は DaysIn 関数およびメソッドが指定年月の末日（日数）を正しく返すかをテストする。
func TestDaysIn(t *testing.T) {
	tests := []struct {
		name  string
		year  int
		month time.Month
		want  int
	}{
		{
			name:  "January (31 days)",
			year:  2024,
			month: time.January,
			want:  31,
		},
		{
			name:  "February in leap year (29 days)",
			year:  2024,
			month: time.February,
			want:  29,
		},
		{
			name:  "February in common year (28 days)",
			year:  2023,
			month: time.February,
			want:  28,
		},
		{
			name:  "April (30 days)",
			year:  2024,
			month: time.April,
			want:  30,
		},
		{
			name:  "December (31 days)",
			year:  2024,
			month: time.December,
			want:  31,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, date.DaysIn(tt.year, tt.month))

			d := date.New(tt.year, tt.month, 1)
			assert.Equal(t, tt.want, d.DaysIn())
		})
	}
}

// TestDate_Date は Date メソッドが年・月・日を正しく分解して返すかをテストする。
func TestDate_Date(t *testing.T) {
	d := date.New(2024, time.December, 25)
	year, month, day := d.Date()

	assert.Equal(t, 2024, year)
	assert.Equal(t, time.December, month)
	assert.Equal(t, 25, day)
}

// TestDate_WithTime は WithTime メソッドが Date と Time とタイムゾーンから time.Time を正しく生成できるかをテストする。
func TestDate_WithTime(t *testing.T) {
	d := date.New(2024, time.May, 1)
	tm := date.NewTime(15, 30, 45, 123456789)
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)

	tests := []struct {
		name    string
		d       date.Date
		tm      date.Time
		loc     *time.Location
		wantLoc *time.Location
	}{
		{
			name:    "nil location (default to Local)",
			d:       d,
			tm:      tm,
			loc:     nil,
			wantLoc: time.Local,
		},
		{
			name:    "UTC location",
			d:       d,
			tm:      tm,
			loc:     time.UTC,
			wantLoc: time.UTC,
		},
		{
			name:    "custom location",
			d:       d,
			tm:      tm,
			loc:     jst,
			wantLoc: jst,
		},
		{
			name:    "zero date and zero time",
			d:       date.Date{},
			tm:      date.Time{},
			loc:     time.UTC,
			wantLoc: time.UTC,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.d.WithTime(tt.tm, tt.loc)
			want := time.Date(tt.d.Year(), tt.d.Month(), tt.d.Day(), tt.tm.Hour(), tt.tm.Minute(), tt.tm.Second(), tt.tm.Nanosecond(), tt.wantLoc)

			assert.True(t, got.Equal(want))
			assert.Equal(t, tt.wantLoc, got.Location())
			assert.Equal(t, tt.tm.Hour(), got.Hour())
			assert.Equal(t, tt.tm.Minute(), got.Minute())
			assert.Equal(t, tt.tm.Second(), got.Second())
			assert.Equal(t, tt.tm.Nanosecond(), got.Nanosecond())
		})
	}
}

// TestDate_ToTime は ToTime メソッドが time.Time に正しく変換できるかをテストする。
func TestDate_ToTime(t *testing.T) {
	d := date.New(2024, time.May, 1)
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)

	tests := []struct {
		name    string
		loc     *time.Location
		wantLoc *time.Location
	}{
		{
			name:    "nil location (default to Local)",
			loc:     nil,
			wantLoc: time.Local,
		},
		{
			name:    "UTC location",
			loc:     time.UTC,
			wantLoc: time.UTC,
		},
		{
			name:    "custom location",
			loc:     jst,
			wantLoc: jst,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := d.ToTime(tt.loc)
			want := time.Date(2024, time.May, 1, 0, 0, 0, 0, tt.wantLoc)

			assert.True(t, got.Equal(want))
			assert.Equal(t, tt.wantLoc, got.Location())
			assert.Equal(t, 0, got.Hour())
			assert.Equal(t, 0, got.Minute())
			assert.Equal(t, 0, got.Second())
			assert.Equal(t, 0, got.Nanosecond())
		})
	}
}

// TestDate_Format は Format メソッドが指定フォーマットに従って文字列化できるかをテストする。
func TestDate_Format(t *testing.T) {
	d := date.New(2024, time.May, 1)
	tests := []struct {
		layout string
		want   string
	}{
		{
			layout: "2006/01/02",
			want:   "2024/05/01",
		},
		{
			layout: "2006-01-02",
			want:   "2024-05-01",
		},
		{
			layout: "January 2, 2006",
			want:   "May 1, 2024",
		},
	}

	for _, tt := range tests {
		t.Run(tt.layout, func(t *testing.T) {
			assert.Equal(t, tt.want, d.Format(tt.layout))
		})
	}
}

// TestDate_YearMonth は YearMonth メソッドが YYYYMM 形式の文字列を返すかをテストする。
func TestDate_YearMonth(t *testing.T) {
	d := date.New(2024, time.May, 1)
	assert.Equal(t, "202405", d.YearMonth())
}

// TestDate_String は String メソッドが YYYY-MM-DD 形式の日付文字列を返すかをテストする。
func TestDate_String(t *testing.T) {
	d := date.New(2024, time.May, 1)
	assert.Equal(t, "2024-05-01", d.String())
}

// TestDate_SubMonths は SubMonths メソッドが2つのDateの間の経過月数を正しく計算できるかをテストする。
func TestDate_SubMonths(t *testing.T) {
	tests := []struct {
		name  string
		d     date.Date
		other date.Date
		want  int
	}{
		{
			name:  "same date",
			d:     date.New(2024, time.May, 15),
			other: date.New(2024, time.May, 15),
			want:  0,
		},
		{
			name:  "both zero values",
			d:     date.Date{},
			other: date.Date{},
			want:  0,
		},
		{
			name:  "same month earlier day",
			d:     date.New(2024, time.May, 10),
			other: date.New(2024, time.May, 20),
			want:  0,
		},
		{
			name:  "same month later day",
			d:     date.New(2024, time.May, 20),
			other: date.New(2024, time.May, 10),
			want:  0,
		},
		{
			name:  "1 month later same day",
			d:     date.New(2024, time.June, 15),
			other: date.New(2024, time.May, 15),
			want:  1,
		},
		{
			name:  "1 month earlier same day (order agnostic)",
			d:     date.New(2024, time.May, 15),
			other: date.New(2024, time.June, 15),
			want:  1,
		},
		{
			name:  "1 month later but day not reached (1/15 to 2/14)",
			d:     date.New(2024, time.January, 15),
			other: date.New(2024, time.February, 14),
			want:  0,
		},
		{
			name:  "1 month later and day reached (1/15 to 2/15)",
			d:     date.New(2024, time.January, 15),
			other: date.New(2024, time.February, 15),
			want:  1,
		},
		{
			name:  "1 month earlier and day not reached (reverse)",
			d:     date.New(2024, time.February, 14),
			other: date.New(2024, time.January, 15),
			want:  0,
		},
		{
			name:  "1 month earlier and day reached (reverse)",
			d:     date.New(2024, time.February, 15),
			other: date.New(2024, time.January, 15),
			want:  1,
		},
		{
			name:  "across year boundary (11/15 to next year 2/14 is 2 months)",
			d:     date.New(2024, time.November, 15),
			other: date.New(2025, time.February, 14),
			want:  2,
		},
		{
			name:  "across year boundary (11/15 to next year 2/15 is 3 months)",
			d:     date.New(2024, time.November, 15),
			other: date.New(2025, time.February, 15),
			want:  3,
		},
		{
			name:  "across year boundary (2025/2/14 to 2024/11/15 is 2 months)",
			d:     date.New(2025, time.February, 14),
			other: date.New(2024, time.November, 15),
			want:  2,
		},
		{
			name:  "across year boundary (2025/2/15 to 2024/11/15 is 3 months)",
			d:     date.New(2025, time.February, 15),
			other: date.New(2024, time.November, 15),
			want:  3,
		},
		{
			name:  "multi-year difference (2024/1/15 to 2026/7/14 is 29 months)",
			d:     date.New(2024, time.January, 15),
			other: date.New(2026, time.July, 14),
			want:  29,
		},
		{
			name:  "multi-year difference (2024/1/15 to 2026/7/15 is 30 months)",
			d:     date.New(2024, time.January, 15),
			other: date.New(2026, time.July, 15),
			want:  30,
		},
		{
			name:  "multi-year difference reverse order",
			d:     date.New(2026, time.July, 15),
			other: date.New(2024, time.January, 15),
			want:  30,
		},
		{
			name:  "month-end: Jan 31 to Feb 29 in leap year (full month)",
			d:     date.New(2024, time.February, 29),
			other: date.New(2024, time.January, 31),
			want:  1,
		},
		{
			name:  "month-end: Jan 31 to Feb 28 in leap year (not full month)",
			d:     date.New(2024, time.February, 28),
			other: date.New(2024, time.January, 31),
			want:  0,
		},
		{
			name:  "month-end: Jan 31 to Feb 28 in non-leap year (full month)",
			d:     date.New(2023, time.February, 28),
			other: date.New(2023, time.January, 31),
			want:  1,
		},
		{
			name:  "month-end: Jan 31 to Feb 27 in non-leap year (not full month)",
			d:     date.New(2023, time.February, 27),
			other: date.New(2023, time.January, 31),
			want:  0,
		},
		{
			name:  "month-end: Jan 30 to Feb 28 in non-leap year (full month)",
			d:     date.New(2023, time.February, 28),
			other: date.New(2023, time.January, 30),
			want:  1,
		},
		{
			name:  "month-end: Aug 31 to Sep 30 (30-day month end)",
			d:     date.New(2024, time.September, 30),
			other: date.New(2024, time.August, 31),
			want:  1,
		},
		{
			name:  "month-end: Aug 31 to Sep 29 (day before 30-day month end)",
			d:     date.New(2024, time.September, 29),
			other: date.New(2024, time.August, 31),
			want:  0,
		},
		{
			name:  "month-end: Feb 29 leap year to Feb 28 non-leap year (12 months)",
			d:     date.New(2021, time.February, 28),
			other: date.New(2020, time.February, 29),
			want:  12,
		},
		{
			name:  "month-end: Feb 29 leap year to Feb 27 non-leap year (11 months)",
			d:     date.New(2021, time.February, 27),
			other: date.New(2020, time.February, 29),
			want:  11,
		},
		{
			name:  "month-end: Feb 29 leap year to Feb 29 next leap year (48 months)",
			d:     date.New(2024, time.February, 29),
			other: date.New(2020, time.February, 29),
			want:  48,
		},
		{
			name:  "month-end: Feb 29 leap year to Feb 28 next leap year (47 months)",
			d:     date.New(2024, time.February, 28),
			other: date.New(2020, time.February, 29),
			want:  47,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.d.SubMonths(tt.other))
		})
	}
}

// TestDate_SpanMonths は SpanMonths メソッドが2つのDateが含まれるカレンダー上の月数を正しく計算できるかをテストする。
func TestDate_SpanMonths(t *testing.T) {
	tests := []struct {
		name  string
		d     date.Date
		other date.Date
		want  int
	}{
		{
			name:  "both zero values",
			d:     date.Date{},
			other: date.Date{},
			want:  0,
		},
		{
			name:  "same date",
			d:     date.New(2024, time.January, 15),
			other: date.New(2024, time.January, 15),
			want:  1,
		},
		{
			name:  "same month different day",
			d:     date.New(2024, time.January, 15),
			other: date.New(2024, time.January, 20),
			want:  1,
		},
		{
			name:  "1/15 to 2/14 is 2 months",
			d:     date.New(2024, time.January, 15),
			other: date.New(2024, time.February, 14),
			want:  2,
		},
		{
			name:  "2/14 to 1/15 (reverse order) is 2 months",
			d:     date.New(2024, time.February, 14),
			other: date.New(2024, time.January, 15),
			want:  2,
		},
		{
			name:  "1/15 to 2/15 is 2 months",
			d:     date.New(2024, time.January, 15),
			other: date.New(2024, time.February, 15),
			want:  2,
		},
		{
			name:  "1/1 to 12/31 is 12 months",
			d:     date.New(2024, time.January, 1),
			other: date.New(2024, time.December, 31),
			want:  12,
		},
		{
			name:  "across year boundary (Jan 15 to Jan 14 next year is 13 months)",
			d:     date.New(2024, time.January, 15),
			other: date.New(2025, time.January, 14),
			want:  13,
		},
		{
			name:  "multi-year difference (2024/1/15 to 2026/7/15 is 31 months)",
			d:     date.New(2024, time.January, 15),
			other: date.New(2026, time.July, 15),
			want:  31,
		},
		{
			name:  "multi-year difference reverse order",
			d:     date.New(2026, time.July, 15),
			other: date.New(2024, time.January, 15),
			want:  31,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.d.SpanMonths(tt.other))
		})
	}
}

// TestParse は Parse 関数が指定フォーマットに従って文字列をDateにパースできるかをテストする。
func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		layout  string
		value   string
		want    date.Date
		wantErr bool
	}{
		{
			name:    "valid standard date format",
			layout:  time.DateOnly,
			value:   "2024-05-01",
			want:    date.New(2024, time.May, 1),
			wantErr: false,
		},
		{
			name:    "valid slash format",
			layout:  "2006/01/02",
			value:   "2024/12/31",
			want:    date.New(2024, time.December, 31),
			wantErr: false,
		},
		{
			name:    "invalid format",
			layout:  time.DateOnly,
			value:   "2024/05/01",
			want:    date.Date{},
			wantErr: true,
		},
		{
			name:    "invalid date value",
			layout:  time.DateOnly,
			value:   "invalid-date",
			want:    date.Date{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := date.Parse(tt.layout, tt.value)
			if tt.wantErr {
				assert.Error(t, err)
				assert.True(t, got.IsZero())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
