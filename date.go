package date

import (
	"cmp"
	"time"
)

type Date struct {
	Year  int
	Month time.Month
	Day   int
}

// IsZero はゼロ値かを返す
func (d Date) IsZero() bool {
	return d.Year == 0 && d.Month == 0 && d.Day == 0
}

// Equal は値が一致するかを返す
func (d Date) Equal(other Date) bool {
	return d.Year == other.Year && d.Month == other.Month && d.Day == other.Day
}

func (d Date) Before(other Date) bool {
	c := d.Compare(other)
	return c < 0
}

func (d Date) After(other Date) bool {
	c := other.Compare(d)
	return c < 0
}

// Compare は他のDateと比較する
func (d Date) Compare(other Date) int {
	c := cmp.Compare(d.Year, other.Year)
	if c == 0 {
		c = cmp.Compare(d.Month, other.Month)
	}
	if c == 0 {
		c = cmp.Compare(d.Day, other.Day)
	}
	return c
}

func (d Date) Date() (int, time.Month, int) {
	return d.Year, d.Month, d.Day
}

// ToTime はtime.Timeに変換して返す
func (d Date) ToTime() time.Time {
	return time.Date(d.Year, d.Month, d.Day, 0, 0, 0, 0, time.Local)
}

func (d Date) Format(layout string) string {
	return d.ToTime().Format(layout)
}

func (d Date) String() string {
	return d.Format(time.DateOnly)
}
