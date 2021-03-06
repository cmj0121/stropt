package stropt

import (
	"net"
	"testing"
	"time"
)

func TestParseTime(t *testing.T) {
	foo := struct {
		D1 time.Time
		D2 *time.Time
	}{
		D1: time.Now(),
		D2: nil,
	}

	time_str := "2006-01-02T15:04:05Z"
	timestamp, _ := time.Parse(time.RFC3339, time_str)
	parser := MustNew(&foo)
	if _, err := parser.Parse("--d1", time_str); err != nil {
		// parse the flag failure
		t.Fatalf("cannot parse flag: %v", err)
	} else if foo.D1 != timestamp {
		// parse but get wrong value
		t.Errorf("expect --d1 %v: %v", time_str, foo.D1)
	}

	if _, err := parser.Parse(time_str); err != nil {
		// parse the flag failure
		t.Fatalf("cannot parse flag: %v", err)
	} else if *foo.D2 != timestamp {
		// parse but get wrong value
		t.Errorf("expect --d2 %v: %v", time_str, *foo.D2)
	}
}

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
		t.Errorf("expect --d2 %v: %v", duration, foo.D2)
	}
}

func TestParseIP(t *testing.T) {
	foo := struct {
		IP1 net.IP
		IP2 *net.IP
	}{
		IP1: nil,
		IP2: nil,
	}

	parser := MustNew(&foo)
	if _, err := parser.Parse("--ip1", "127.0.0.2"); err != nil {
		// parse the flag failure
		t.Fatalf("cannot parse flag: %v", err)
	} else if foo.IP1.String() != "127.0.0.2" {
		// parse but get wrong value
		t.Errorf("expect --ip1 127.0.0.2: %v", foo.IP1)
	}

	if _, err := parser.Parse("--ip1", "localhost"); err != nil {
		// parse the flag failure
		t.Fatalf("cannot parse flag: %v", err)
	} else if !foo.IP1.IsLoopback() {
		// parse but get wrong value
		t.Errorf("expect --ip1 localhost: %v", foo.IP1)
	}

	if _, err := parser.Parse("0.0.0.0"); err != nil {
		// parse the flag failure
		t.Fatalf("cannot parse flag: %v", err)
	} else if foo.IP2.String() != "0.0.0.0" {
		// parse but get wrong value
		t.Errorf("expect --ip1 0.0.0.0: %v", foo.IP2)
	}
}

func TestParseCIDR(t *testing.T) {
	foo := struct {
		IP1 *net.IPNet `attr:"flag"`
		IP2 *net.IPNet
	}{
		IP1: nil,
		IP2: nil,
	}

	parser := MustNew(&foo)
	if _, err := parser.Parse("--ip1", "127.0.0.2/16"); err != nil {
		// parse the flag failure
		t.Fatalf("cannot parse flag: %v", err)
	} else if foo.IP1.String() != "127.0.0.0/16" {
		// parse but get wrong value
		t.Errorf("expect --ip1 127.0.0.2/16: %v", foo.IP1)
	}

	if _, err := parser.Parse("0.0.0.0/0"); err != nil {
		// parse the flag failure
		t.Fatalf("cannot parse flag: %v", err)
	} else if foo.IP2.String() != "0.0.0.0/0" {
		// parse but get wrong value
		t.Errorf("expect --ip1 0.0.0.0/0: %v", foo.IP2)
	}
}
