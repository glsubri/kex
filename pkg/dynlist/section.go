package dynlist

type section struct {
	start int // included
	end   int // excluded
}

func (s section) scrollDown(n int, maxHeight int) section {
	if s.end+n > maxHeight {
		n = max(0, maxHeight-s.end)
	}

	return section{
		start: s.start + n,
		end:   s.end + n,
	}
}

func (s section) scrollUp(n int) section {
	if s.start-n < 0 {
		n = s.start
	}

	return section{
		start: s.start - n,
		end:   s.end - n,
	}
}

func (s section) height() int {
	return s.end - s.start
}

func (s section) center() int {
	return (s.start + s.end) / 2
}
