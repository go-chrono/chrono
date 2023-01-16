package chrono_test

import "github.com/go-chrono/chrono"

var (
	formatDurationCases = []struct {
		name      string
		of        chrono.Extent
		exclusive []chrono.Designator
		expected  string
	}{
		{
			name:      "default HMS",
			of:        1*chrono.Hour + 15*chrono.Minute + 30*chrono.Second + 500*chrono.Millisecond,
			exclusive: []chrono.Designator{},
			expected:  "PT1H15M30.5S",
		},
		{
			name:      "default HM",
			of:        1*chrono.Hour + 15*chrono.Minute,
			exclusive: []chrono.Designator{},
			expected:  "PT1H15M",
		},
		{
			name:      "default HS",
			of:        12*chrono.Hour + 30*chrono.Second + 500*chrono.Millisecond,
			exclusive: []chrono.Designator{},
			expected:  "PT12H0M30.5S",
		},
		{
			name:      "default H",
			of:        1 * chrono.Hour,
			exclusive: []chrono.Designator{},
			expected:  "PT1H",
		},
		{
			name:      "default MS",
			of:        15*chrono.Minute + 30*chrono.Second + 500*chrono.Millisecond,
			exclusive: []chrono.Designator{},
			expected:  "PT15M30.5S",
		},
		{
			name:      "default M",
			of:        15 * chrono.Minute,
			exclusive: []chrono.Designator{},
			expected:  "PT15M",
		},
		{
			name:      "default S",
			of:        500 * chrono.Millisecond,
			exclusive: []chrono.Designator{},
			expected:  "PT0.5S",
		},
		{
			name:      "default zero value",
			of:        0,
			exclusive: []chrono.Designator{},
			expected:  "PT0S",
		},
		{
			name:      "exclusive HMS",
			of:        1*chrono.Hour + 15*chrono.Minute + 30*chrono.Second + 500*chrono.Millisecond,
			exclusive: []chrono.Designator{chrono.Hours, chrono.Minutes, chrono.Seconds},
			expected:  "PT1H15M30.5S",
		},
		{
			name:      "exclusive HMS zero value",
			of:        0,
			exclusive: []chrono.Designator{chrono.Hours, chrono.Minutes, chrono.Seconds},
			expected:  "PT0H0M0S",
		},
		{
			name:      "exclusive HM",
			of:        1*chrono.Hour + 15*chrono.Minute + 30*chrono.Second + 600*chrono.Millisecond,
			exclusive: []chrono.Designator{chrono.Hours, chrono.Minutes},
			expected:  "PT1H15.51M",
		},
		{
			name:      "exclusive HM zero value",
			of:        0,
			exclusive: []chrono.Designator{chrono.Hours, chrono.Minutes},
			expected:  "PT0H0M",
		},
		{
			name:      "exclusive HS",
			of:        12*chrono.Hour + 1*chrono.Minute + 30*chrono.Second + 500*chrono.Millisecond,
			exclusive: []chrono.Designator{chrono.Hours, chrono.Seconds},
			expected:  "PT12H90.5S",
		},
		{
			name:      "exclusive HS zero value",
			of:        0,
			exclusive: []chrono.Designator{chrono.Hours, chrono.Seconds},
			expected:  "PT0H0S",
		},
		{
			name:      "exclusive H",
			of:        1*chrono.Hour + 30*chrono.Minute + 36*chrono.Second + 36*chrono.Millisecond,
			exclusive: []chrono.Designator{chrono.Hours},
			expected:  "PT1.51001H",
		},
		{
			name:      "exclusive H zero value",
			of:        0,
			exclusive: []chrono.Designator{chrono.Hours},
			expected:  "PT0H",
		},
		{
			name:      "exclusive MS",
			of:        1*chrono.Hour + 15*chrono.Minute + 30*chrono.Second + 500*chrono.Millisecond,
			exclusive: []chrono.Designator{chrono.Minutes, chrono.Seconds},
			expected:  "PT75M30.5S",
		},
		{
			name:      "exclusive MS zero value",
			of:        0,
			exclusive: []chrono.Designator{chrono.Minutes, chrono.Seconds},
			expected:  "PT0M0S",
		},
		{
			name:      "exclusive M",
			of:        1*chrono.Hour + 15*chrono.Minute + 30*chrono.Second + 600*chrono.Millisecond,
			exclusive: []chrono.Designator{chrono.Minutes},
			expected:  "PT75.51M",
		},
		{
			name:      "exclusive M zero value",
			of:        0,
			exclusive: []chrono.Designator{chrono.Minutes},
			expected:  "PT0M",
		},
		{
			name:      "exclusive S",
			of:        1*chrono.Hour + 15*chrono.Minute + 30*chrono.Second + 500*chrono.Millisecond,
			exclusive: []chrono.Designator{chrono.Seconds},
			expected:  "PT4530.5S",
		},
		{
			name:      "exclusive S zero value",
			of:        0,
			exclusive: []chrono.Designator{chrono.Seconds},
			expected:  "PT0S",
		},
	}

	parseDurationCases = []struct {
		name     string
		input    string
		expected chrono.Extent
	}{
		{
			name:     "valid HMS integers",
			input:    "PT5H3M1S",
			expected: 5*chrono.Hour + 3*chrono.Minute + 1*chrono.Second,
		},
		{
			name:     "valid HMS floats",
			input:    "PT4.5H3.25M1.1S",
			expected: chrono.Extent(4.5*float64(chrono.Hour) + 3.25*float64(chrono.Minute) + 1.1*float64(chrono.Second)),
		},
		{
			name:     "valid HM integers",
			input:    "PT5H3M",
			expected: 5*chrono.Hour + 3*chrono.Minute,
		},
		{
			name:     "valid HM floats",
			input:    "PT4.5H3.25M",
			expected: chrono.Extent(4.5*float64(chrono.Hour) + 3.25*float64(chrono.Minute)),
		},
		{
			name:     "valid HS integers",
			input:    "PT5H1S",
			expected: 5*chrono.Hour + 1*chrono.Second,
		},
		{
			name:     "valid HS floats",
			input:    "PT4.5H1.1S",
			expected: chrono.Extent(4.5*float64(chrono.Hour) + 1.1*float64(chrono.Second)),
		},
		{
			name:     "valid H integer",
			input:    "PT5H",
			expected: 5 * chrono.Hour,
		},
		{
			name:     "valid H float",
			input:    "PT4.5H",
			expected: chrono.Extent(4.5 * float64(chrono.Hour)),
		},
		{
			name:     "valid MS integers",
			input:    "PT3M1S",
			expected: 3*chrono.Minute + 1*chrono.Second,
		},
		{
			name:     "valid MS floats",
			input:    "PT3.25M1.1S",
			expected: chrono.Extent(3.25*float64(chrono.Minute) + 1.1*float64(chrono.Second)),
		},
		{
			name:     "valid M integer",
			input:    "PT3M",
			expected: 3 * chrono.Minute,
		},
		{
			name:     "valid M float",
			input:    "PT3.25M",
			expected: chrono.Extent(3.25 * float64(chrono.Minute)),
		},
		{
			name:     "valid S integer",
			input:    "PT1S",
			expected: 1 * chrono.Second,
		},
		{
			name:     "valid S float",
			input:    "PT1.1S",
			expected: chrono.Extent(1.1 * float64(chrono.Second)),
		},
	}
)
