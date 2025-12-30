package custom_type

type Chain[E any] struct {
	chain [][]E
	cur   int
	i     int
}

func (c *Chain[E]) Add(s ...[]E) {
	c.chain = append(c.chain, s...)
}

func (c *Chain[E]) HasNext() bool {
	if c.cur >= len(c.chain) {
		return false
	}
	if c.i < len(c.chain[c.cur]) {
		return true
	}
	c.cur++
	c.i = 0
	return c.HasNext()
}

func (c *Chain[E]) Next() E {
	defer func() { c.i++ }()
	return c.chain[c.cur][c.i]
}
