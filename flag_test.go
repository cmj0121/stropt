package stropt

import (
	"testing"
	"time"
)

func TestParseTimeDuration(t *testing.T) {
	foo := struct {
		D1 time.Duration
		D2 *time.Duration
	}{
		D1: 0,
		D2: nil,
	}

	parser := MustNew(&foo)
	if _, err := parser.Parse("--d1", "10"); err != nil {
		// parse the flag failure
		t.Fatalf("cannot parse flag: %v", err)
	} else if foo.D1 != 10 {
		// parse but get wrong value
		t.Errorf("expect --d1 10: %v", foo.D1)
	}

	duration, _ := time.ParseDuration("10m20s")
	if _, err := parser.Parse("--d1", "10m20s"); err != nil {
		// parse the flag failure
		t.Fatalf("cannot parse flag: %v", err)
	} else if foo.D1 != duration {
		// parse but get wrong value
		t.Errorf("expect --d1 %v: %v", duration, foo.D1)
	}

	duration, _ = time.ParseDuration("2h3m10s")
	if _, err := parser.Parse("2h3m10s"); err != nil {
		// parse the flag failure
		t.Fatalf("cannot parse flag: %v", err)
	} else if *foo.D2 != duration {
		// parse but get wrong value
		t.Errorf("expect --d1 %v: %v", duration, foo.D2)
	}
}
