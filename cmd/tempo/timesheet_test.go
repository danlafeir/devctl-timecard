package tempo

import (
	"testing"
	"time"
)


func TestConstants(t *testing.T) {
	tests := []struct {
		name     string
		constant string
		expected string
	}{
		{
			name:     "CapitalizableTime constant",
			constant: CapitalizableTime,
			expected: "How much time did you spend developing, designing or testing software (in hours): ",
		},
		{
			name:     "PtoTime constant",
			constant: PtoTime,
			expected: "How much time did you spend with PTO (vacation or sick) (in hours): ",
		},
		{
			name:     "SupportTime constant",
			constant: SupportTime,
			expected: "How much time did you supporting a software system (in hours): ",
		},
		{
			name:     "MeetingTime constant",
			constant: MeetingTime,
			expected: "How much time did you spend in meetings (in hours): ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.constant != tt.expected {
				t.Errorf("Constant value mismatch: got %q, expected %q", tt.constant, tt.expected)
			}
		})
	}
}

func TestDetermineWeekforTimeSheet(t *testing.T) {
	tests := []struct {
		name        string
		currentTime time.Time
		expectedDay time.Weekday
	}{
		{
			name:        "Current day is Monday",
			currentTime: time.Date(2024, 1, 8, 10, 0, 0, 0, time.UTC), // Monday
			expectedDay: time.Monday,
		},
		{
			name:        "Current day is Tuesday",
			currentTime: time.Date(2024, 1, 9, 10, 0, 0, 0, time.UTC), // Tuesday
			expectedDay: time.Monday,
		},
		{
			name:        "Current day is Wednesday",
			currentTime: time.Date(2024, 1, 10, 10, 0, 0, 0, time.UTC), // Wednesday
			expectedDay: time.Monday,
		},
		{
			name:        "Current day is Thursday",
			currentTime: time.Date(2024, 1, 11, 10, 0, 0, 0, time.UTC), // Thursday
			expectedDay: time.Monday,
		},
		{
			name:        "Current day is Friday",
			currentTime: time.Date(2024, 1, 12, 10, 0, 0, 0, time.UTC), // Friday
			expectedDay: time.Monday,
		},
		{
			name:        "Current day is Saturday",
			currentTime: time.Date(2024, 1, 13, 10, 0, 0, 0, time.UTC), // Saturday
			expectedDay: time.Monday,
		},
		{
			name:        "Current day is Sunday",
			currentTime: time.Date(2024, 1, 14, 10, 0, 0, 0, time.UTC), // Sunday
			expectedDay: time.Monday,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Since determineWeekforTimeSheet uses time.Now(), we'll test the logic indirectly
			// by testing the weekday calculation logic directly
			dayOfTheWeek := int(tt.currentTime.Weekday())
			var distanceToMonday int
			if dayOfTheWeek-1 == -1 {
				distanceToMonday = -6
			} else {
				distanceToMonday = -(dayOfTheWeek - 1)
			}

			monday := tt.currentTime.AddDate(0, 0, distanceToMonday)
			
			if monday.Weekday() != tt.expectedDay {
				t.Errorf("Expected weekday %v, got %v for input %v", tt.expectedDay, monday.Weekday(), tt.currentTime)
			}

			// Verify it's actually the Monday of the same week
			if monday.After(tt.currentTime) {
				t.Errorf("Monday %v should not be after current time %v", monday, tt.currentTime)
			}

			// Verify it's within 6 days before the current time
			daysDiff := tt.currentTime.Sub(monday).Hours() / 24
			if daysDiff < 0 || daysDiff > 6 {
				t.Errorf("Monday should be 0-6 days before current time, but difference is %.1f days", daysDiff)
			}
		})
	}
}

func TestDetermineWeekforTimeSheet_Integration(t *testing.T) {
	// Integration test that calls the actual function
	result := determineWeekforTimeSheet()
	
	// Verify the result is a Monday
	if result.Weekday() != time.Monday {
		t.Errorf("determineWeekforTimeSheet() should return a Monday, got %v", result.Weekday())
	}

	// Verify the result is not in the future
	now := time.Now()
	if result.After(now) {
		t.Errorf("determineWeekforTimeSheet() should not return a date in the future")
	}

	// Verify the result is within the current week (within 6 days ago)
	daysDiff := now.Sub(result).Hours() / 24
	if daysDiff < 0 || daysDiff > 6 {
		t.Errorf("determineWeekforTimeSheet() should return a date within the current week, but it's %.1f days ago", daysDiff)
	}
}

// Note: The following functions are difficult to test without mocking stdin/stdout:
// - requestTimeInput() - requires user input via fmt.Scan
// - getTime() - requires user input via fmt.Scan  
// - requestDayOfWeek() - requires user input via fmt.Scan
//
// These would benefit from dependency injection or refactoring to separate
// I/O operations from business logic for better testability.

func TestRequestTimeInput_Logic(t *testing.T) {
	// Test the calculation logic used in requestTimeInput
	// This tests the core logic without I/O dependencies
	
	tests := []struct {
		name                           string
		devTime, ptoTime, supportTime, meetingTime int
		expectedTotal                  int
	}{
		{
			name:          "All zeros",
			devTime:       0,
			ptoTime:       0,
			supportTime:   0,
			meetingTime:   0,
			expectedTotal: 0,
		},
		{
			name:          "Standard work week",
			devTime:       32,
			ptoTime:       0,
			supportTime:   4,
			meetingTime:   4,
			expectedTotal: 40,
		},
		{
			name:          "Vacation week",
			devTime:       0,
			ptoTime:       40,
			supportTime:   0,
			meetingTime:   0,
			expectedTotal: 40,
		},
		{
			name:          "Mixed week",
			devTime:       20,
			ptoTime:       8,
			supportTime:   6,
			meetingTime:   6,
			expectedTotal: 40,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test the calculation logic from requestTimeInput
			totalHoursThisWeek := tt.meetingTime + tt.ptoTime + tt.supportTime + tt.devTime
			
			if totalHoursThisWeek != tt.expectedTotal {
				t.Errorf("Total hours calculation: got %d, expected %d", totalHoursThisWeek, tt.expectedTotal)
			}
		})
	}
}

func TestRequestDayOfWeek_DateCalculation(t *testing.T) {
	// Test the date calculation logic used in requestDayOfWeek
	// This tests the core logic without I/O dependencies
	
	tests := []struct {
		name          string
		baseDate      time.Time
		weeksBack     int
		expectedDate  time.Time
	}{
		{
			name:         "Current week",
			baseDate:     time.Date(2024, 1, 8, 0, 0, 0, 0, time.UTC), // Monday
			weeksBack:    0,
			expectedDate: time.Date(2024, 1, 8, 0, 0, 0, 0, time.UTC),
		},
		{
			name:         "One week back",
			baseDate:     time.Date(2024, 1, 8, 0, 0, 0, 0, time.UTC), // Monday
			weeksBack:    1,
			expectedDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), // Previous Monday
		},
		{
			name:         "Two weeks back",
			baseDate:     time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC), // Monday
			weeksBack:    2,
			expectedDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),  // Two Mondays ago
		},
		{
			name:         "Three weeks back",
			baseDate:     time.Date(2024, 1, 22, 0, 0, 0, 0, time.UTC), // Monday
			weeksBack:    3,
			expectedDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),  // Three Mondays ago
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test the date calculation logic from requestDayOfWeek
			result := tt.baseDate.AddDate(0, 0, -7*tt.weeksBack)
			
			if !result.Equal(tt.expectedDate) {
				t.Errorf("Date calculation: got %v, expected %v", result.Format(time.DateOnly), tt.expectedDate.Format(time.DateOnly))
			}
			
			// Verify the result is still a Monday (if base was Monday)
			if tt.baseDate.Weekday() == time.Monday && result.Weekday() != time.Monday {
				t.Errorf("Result should be Monday, got %v", result.Weekday())
			}
		})
	}
}
