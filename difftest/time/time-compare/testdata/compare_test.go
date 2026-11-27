package time_test

import (
	"bufio"
	"cmp"
	"fmt"
	"io"
	"os"
	"testing"
	"time"
)

type instant struct {
	seconds int64
	nanos   int32
	reading int64
	hasMono bool
}

type observations struct {
	nsec                       int32
	sec                        int64
	unixSec                    int64
	mono                       int64
	nano                       int
	unix                       int64
	zero, before, after, equal bool
	compare                    int
}

func decode(wall uint64, ext int64) instant {
	const nanosPerField = uint64(1 << 30)
	const monoFlag = uint64(1 << 63)
	v := instant{seconds: ext, nanos: int32(wall % nanosPerField)}
	if wall >= monoFlag {
		v.seconds = 59453308800 + int64((wall-monoFlag)/nanosPerField)
		v.reading, v.hasMono = ext, true
	}
	return v
}

func reference(tw uint64, te int64, uw uint64, ue int64) observations {
	t, u := decode(tw, te), decode(uw, ue)
	var order int
	if t.hasMono && u.hasMono {
		order = cmp.Compare(t.reading, u.reading)
	} else {
		order = cmp.Compare(t.seconds, u.seconds)
		if order == 0 {
			order = cmp.Compare(t.nanos, u.nanos)
		}
	}
	unix := t.seconds - 62135596800
	return observations{t.nanos, t.seconds, unix, t.reading, int(t.nanos), unix,
		tw == 0 && te == 0, order < 0, order > 0, order == 0, order}
}

func TestGoLeanClockCompare(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_CLOCK_COMPARE"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	west := time.FixedZone("west", -18000)
	east := time.FixedZone("east", 19800)
	for row := 0; ; row++ {
		var tw, uw uint64
		var te, ue int64
		var want observations
		n, err := fmt.Fscan(reader, &tw, &te, &uw, &ue,
			&want.nsec, &want.sec, &want.unixSec, &want.mono, &want.nano, &want.unix,
			&want.zero, &want.before, &want.after, &want.equal, &want.compare)
		if n == 0 && err == io.EOF && row > 0 {
			t.Logf("%d pinned observations per declaration", row)
			return
		}
		if err != nil {
			t.Fatal(err)
		}
		a, b := time.GoLeanTimeFromWords(tw, te, west), time.GoLeanTimeFromWords(uw, ue, east)
		ns, sec, us, mono := time.GoLeanTimeWords(a)
		got := observations{ns, sec, us, mono, a.Nanosecond(), a.Unix(),
			a.IsZero(), a.Before(b), a.After(b), a.Equal(b), a.Compare(b)}
		oracle := reference(tw, te, uw, ue)
		if got != oracle || got != want {
			t.Fatalf("row %d: Go %+v, Lean %+v, oracle %+v", row, got, want, oracle)
		}
	}
}
