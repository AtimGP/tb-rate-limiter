package storage

import (
	"testing"
	"time"
)

func TestBucket(t *testing.T) {
	timeStart := time.Now()

	tests := []struct{
		name 			string
		now 			time.Time
		limit			float64
		rateRefill		float64
		lastRefill		time.Time
		tokens			float64

		wantAllowed		bool
		wantRemaining	int
		wantResetAfter	time.Duration
	}{
		{
			name: 		"Доступ разрешен тест",
			now: 		timeStart,
			limit: 		10.0,
			rateRefill: 1.0,
			lastRefill: timeStart,
			tokens: 	10.0,

			wantAllowed:	 true,
			wantRemaining:	 9,
			wantResetAfter:	 1 * time.Second,
		},
		{
			name: 		"Доступ запрещен тест",
			now: 		timeStart,
			limit: 		10.0,
			rateRefill: 1.0,
			lastRefill: timeStart,
			tokens: 	0.5,

			wantAllowed:	 false,
			wantRemaining:	 0,
			wantResetAfter:	 500 * time.Millisecond,
		},
		{
			name: 		"Доступ разрешен после восстановления тест",
			now: 		timeStart,
			limit: 		10.0,
			rateRefill: 1.0,
			lastRefill: timeStart.Add(-2 * time.Second),
			tokens: 	0.0,

			wantAllowed:	 true,
			wantRemaining:	 1,
			wantResetAfter:	 9 * time.Second,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			b := &bucket{
				tokens: test.tokens,
				lastRefill: test.lastRefill,
			}

			allowed, remaining, resetAfter := b.Take(test.now, test.limit, test.rateRefill)

			if allowed != test.wantAllowed {
				t.Errorf("FAIL %s: allowed = %v; ожидалось %v", test.name, allowed, test.wantAllowed)
			}

			if remaining != test.wantRemaining {
				t.Errorf("FAIL %s: remaining= %v; ожидалось %v", test.name, remaining, test.wantRemaining)
			}

			if resetAfter != test.wantResetAfter {
				t.Errorf("FAIL %s: resetAfter = %v; ожидалось %v", test.name, resetAfter, test.wantResetAfter)
			}

		})
	}
}