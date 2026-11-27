package strings

import "sync"

func GoLeanInspectNewReplacer(r *Replacer) ([]string, *sync.Once, bool) {
	return r.oldnew, &r.once, r.r == nil
}
