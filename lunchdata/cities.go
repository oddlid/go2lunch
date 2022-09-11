package lunchdata

type Cities []*City

func (cs Cities) Len() int {
	return len(cs)
}
