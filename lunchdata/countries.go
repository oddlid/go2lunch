package lunchdata

type Countries []*Country

func (cs *Countries) Add(c *Country) {
	*cs = append(*cs, c)
}

func (cs *Countries) Len() int {
	return len(*cs)
}
