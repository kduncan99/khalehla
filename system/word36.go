package system

type Word36 struct {
	Value uint64
}

func (w *Word36) IsNegative() bool {
	if w.Value&0400000000000 == 0400000000000 {
		return true
	} else {
		return false
	}
}

func (w *Word36) IsZero() bool {
	if (w.Value == 0) || (w.Value == 0777777777777) {
		return true
	} else {
		return false
	}
}

func (w *Word36) Get() uint64 {
	return w.Value
}

func (w *Word36) GetH1() uint64 {
	return w.Value >> 18
}

func (w *Word36) GetH2() uint64 {
	return w.Value & 0777777
}

func (w *Word36) GetQ1() uint64 {
	return w.Value >> 27
}

func (w *Word36) GetQ2() uint64 {
	return (w.Value >> 18) & 0777
}

func (w *Word36) GetQ3() uint64 {
	return (w.Value >> 9) & 0777
}

func (w *Word36) GetQ4() uint64 {
	return w.Value & 0777
}

func (w *Word36) GetS1() uint64 {
	return w.Value >> 30
}

func (w *Word36) GetS2() uint64 {
	return (w.Value >> 24) & 077
}

func (w *Word36) GetS3() uint64 {
	return (w.Value >> 18) & 077
}

func (w *Word36) GetS4() uint64 {
	return (w.Value >> 12) & 077
}

func (w *Word36) GetS5() uint64 {
	return (w.Value >> 6) & 077
}

func (w *Word36) GetS6() uint64 {
	return w.Value & 077
}

func (w *Word36) And(op *Word36) {
	w.Value = w.Value & op.Value
}

func (w *Word36) Not() {
	w.Value = w.Value ^ 0777777777777
}

func (w *Word36) Or(op *Word36) {
	w.Value = w.Value | op.Value
}

func New(value uint64) *Word36 {
	return &Word36{Value: value & 0777777777777}
}

//	TODO Add
//	TODO Subtract
//	TODO other fun stuff
