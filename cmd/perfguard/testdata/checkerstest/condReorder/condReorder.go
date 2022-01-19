package checkerstest

func Warn1(simple bool) {
	if predicate() && simple { // want `predicate() && simple => simple && predicate()`
	}
	if predicate() || simple { // want `predicate() || simple => simple || predicate()`
	}
}

func Fixed1(simple bool) {
	if simple && predicate() {
	}
	if simple || predicate() {
	}
}

func Warn2(simple bool) {
	for predicate() && simple { // want `predicate() && simple => simple && predicate()`
	}
	for predicate() || simple { // want `predicate() || simple => simple || predicate()`
	}

	for i := 0; predicate() && simple; i++ { // want `predicate() && simple => simple && predicate()`
	}
	for i := 0; predicate() || simple; i++ { // want `predicate() || simple => simple || predicate()`
	}
}

func Fixed2(simple bool) {
	for simple && predicate() {
	}
	for simple || predicate() {
	}

	for i := 0; simple && predicate(); i++ {
	}
	for i := 0; simple || predicate(); i++ {
	}
}

func Warn3(simple bool) {
	{
		_ = predicate() && simple // want `predicate() && simple => simple && predicate()`
	}
}

func Fixed3(simple bool) {
	{
		_ = simple && predicate()
	}
}

func Ignore1(xs []bool) {
	// In the conditions below, we can't change the order
	// even if `xs[0]` is "simpler" than `len(xs) > 0`.
	if len(xs) > 0 && xs[0] {
	}
	if len(xs) >= +0 && xs[0] {
	}
	if len(xs) != 0 && xs[0] {
	}

	type object struct {
		x int
		b bool
	}
	o := new(object)
	if o != nil && o.b {
	}
	if o == nil || o.b {
	}
}

func Ignore2(cond bool, xs []int) {
	// Conditions below are +/- identical.
	// There is no point in moving them around.
	if (cond) && cond {
	}
	if !cond && cond {
	}
	if !!cond && cond {
	}
	if !!!cond && cond {
	}
	if cond && len(xs) > 0 {
	}
	if len(xs) > 0 && cond {
	}
}

func Ignore3(cond bool) {
	// Never suggest to put a function call before a simple expr.
	if cond && predicate() {
	}
	if !cond && predicate() {
	}
	if !!cond && predicate() {
	}
	if !!!cond && predicate() {
	}
}

func Ignore4(isMtime, fitsOctal, needsNano bool) {
	if !isMtime || !fitsOctal || needsNano {
	}
}

func Ignore5() {
	{
		// Can't reorder here: o.Check may alter the o state, making the o.cond different.
		o := &trickyObject{}
		if o.Check() && o.cond {
		}
	}
}

type trickyObject struct {
	cond bool
}

func (o *trickyObject) Check() bool { return false }

func predicate() bool { return false }
