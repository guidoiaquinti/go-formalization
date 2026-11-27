package strings

import (
	"fmt"
	"runtime"
)

func GoLeanBuildOnceControl() (failure string) {
	for _, completed := range []bool{false, true} {
		r := NewReplacer("ab", "X")
		if completed {
			r.once.Do(func() {})
		}
		r.buildOnce()
		if r.oldnew != nil || r.r.Replace("abab") != "XX" {
			return fmt.Sprintf("completed=%t: direct callback did not consume arguments and install replacement", completed)
		}
		want := "abab"
		if completed {
			want = "XX"
		}
		if got := r.Replace("abab"); got != want {
			return fmt.Sprintf("completed=%t: public Replace=%q want %q", completed, got, want)
		}
		r = NewReplacer("ab", "X")
		if completed {
			r.once.Do(func() {})
		}
		r.buildOnce()
		r.buildOnce()
		if _, ok := r.r.(*byteReplacer); !ok || r.r.Replace("abab") != "abab" || r.oldnew != nil {
			return fmt.Sprintf("completed=%t: second callback did not rebuild identity", completed)
		}
	}
	r := NewReplacer("\xff", "x", "\xff", "y")
	r.buildOnce()
	if got := r.r.Replace("\xffz\xff"); got != "xzx" {
		return fmt.Sprintf("duplicate byte precedence: got %q", got)
	}
	var nilReceiver *Replacer
	defer func() {
		fault, ok := recover().(runtime.Error)
		if !ok || fault.Error() != "runtime error: invalid memory address or nil pointer dereference" {
			failure = "nil receiver did not produce the expected runtime panic"
		}
	}()
	nilReceiver.buildOnce()
	return ""
}
