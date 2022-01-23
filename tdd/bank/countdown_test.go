type SpySleeper struct {
	Calls int
}
func (s *SpySleeper) Sleep() {
	s.Calls++
}

func TestCountdown(t *testing.T) {
	buf := &bytes.Buffer{}
	sleeper := new(SpySleeper)

	Countdown(buf, sleeper)

	got := buf.String()
	wanted := `3
2
1
Go!`

	if got != wanted {
		t.Errorf("got %q wanted %q", got, wanted)
	}
	if sleeper.Calls != 4 {
		t.Errorf("not enough calls to sleeper, want 4 got %d", sleeper.Calls)
	}
}

func ExampleCountdown() {
	sleeper := &DefaultSleeper{}
	Countdown(os.Stdout, sleeper)
}

type SpyCountdownOp struct {
	Calls []string
}

func (s *SpyCountdownOp) Sleep() {
	s.Calls = append(s.Calls, sleepOp)
}

func (s *SpyCountdownOp) Write(p []byte) (n int, err error) {
	s.Calls = append(s.Calls, writeOp)
	return
}

const writeOp = "write"
const sleepOp = "sleep"

func TestCountdown2(t *testing.T) {
	t.Run("prints 3 to Go!", func(t *testing.T) {
		buffer := &bytes.Buffer{}
		Countdown(buffer, &SpyCountdownOp{})

		got := buffer.String()
		wanted := `3
2
1
Go!`

		if got != wanted {
			t.Errorf("got %q wanted %q", got, wanted)
		}
	})

	t.Run("sleep before every print", func(t *testing.T) {
		ops := &SpyCountdownOp{}
		Countdown(ops, ops)

		wanted := []string{
			sleepOp,
			writeOp,
			sleepOp,
			writeOp,
			sleepOp,
			writeOp,
			sleepOp,
			writeOp,
		}

		if !reflect.DeepEqual(wanted, ops.Calls) {
			t.Errorf("wanted calls %v got %v", wanted, ops.Calls)
		}
	})
}