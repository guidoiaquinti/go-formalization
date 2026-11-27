package sync

func GoLeanOnceDoSlow(o *Once, f func()) { o.doSlow(f) }
func GoLeanOnceDone(o *Once) bool        { return o.done.Load() }
func GoLeanOnceLocked(o *Once) bool {
	if o.m.TryLock() {
		o.m.Unlock()
		return false
	}
	return true
}
