package date

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestNewCompactDate は NewCompactDate コンストラクタが正しく CompactDate を生成し正規化するかをテストする。
func TestNewCompactDate(t *testing.T) {
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
			d := NewCompact(tt.year, tt.month, tt.day)
			assert.Equal(t, tt.wantYear, d.Year())
			assert.Equal(t, tt.wantMonth, d.Month())
			assert.Equal(t, tt.wantDay, d.Day())
		})
	}
}

type testCompactDateData struct {
	Value    CompactDate
	Nullable *CompactDate `json:",omitempty"`
	Empty    CompactDate  `json:",omitzero"`
}

func TestCompactDate_Marshal(t *testing.T) {
	sut := NewCompact(2025, time.October, 24)
	data := testCompactDateData{Value: sut, Nullable: nil, Empty: CompactDate{}}

	actual, err := json.Marshal(data)
	assert.NoError(t, err)
	assert.Equal(t, []byte(`{"Value":"20251024"}`), actual)
}

func TestCompactDate_Unmarshal(t *testing.T) {
	sut := NewCompact(2025, 10, 24)
	var actual testCompactDateData

	err := json.Unmarshal([]byte(`{"Value":"20251024"}`), &actual)
	assert.NoError(t, err)
	assert.Equal(t, testCompactDateData{Value: sut, Nullable: nil, Empty: CompactDate{}}, actual)
}

func TestCompactDate_Date(t *testing.T) {
	sut := NewCompact(2025, 10, 24)

	assert.Equal(t, New(2025, 10, 24), sut.Date)
}

func TestCompactDate_Equal(t *testing.T) {
	sut := NewCompact(2025, 10, 24)

	assert.True(t, sut.Equal(CompactFromTime(time.Date(2025, 10, 24, 0, 0, 0, 0, time.Local)).Date))
	assert.True(t, sut.Equal(CompactFromTime(time.Date(2025, 10, 24, 10, 20, 30, 4, time.Local)).Date))
	assert.False(t, sut.Equal(CompactFromTime(time.Date(2025, 10, 25, 0, 0, 0, 0, time.Local)).Date))
}

func TestCompactDate_String(t *testing.T) {
	sut := NewCompact(2025, 10, 24)

	assert.Equal(t, "20251024", sut.String())
}
