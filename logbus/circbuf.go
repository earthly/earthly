package logbus

type circBuf struct {
	data    []byte
	offset  int
	maxSize int
}

func newCircBuf(maxSize int) *circBuf {
	return &circBuf{
		maxSize: maxSize,
	}
}

func (c *circBuf) Write(buf []byte) (int, error) {
	l := len(c.data)
	n := len(buf)

	if n > c.maxSize {
		buf = buf[n-c.maxSize:]
	}

	if l < c.maxSize {
		r := c.maxSize - l
		if n > r {
			c.data = append(c.data, buf[:r]...)
			buf = buf[r:]
			c.offset = 0
		} else {
			c.data = append(c.data, buf...)
			c.offset = n
			return n, nil
		}
	}

	remain := c.maxSize - c.offset
	copy(c.data[c.offset:], buf)
	if len(buf) > remain {
		copy(c.data, buf[remain:])
	}

	c.offset = (c.offset + len(buf)) % c.maxSize

	return n, nil
}

func (c *circBuf) bytes() []byte {
	if len(c.data) < c.maxSize {
		return c.data
	}
	ret := make([]byte, c.maxSize)
	copy(ret, c.data[c.offset:])
	copy(ret[c.maxSize-c.offset:], c.data[:c.offset])
	return ret
}
