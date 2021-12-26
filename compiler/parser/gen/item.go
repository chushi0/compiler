package gen

import (
	"github.com/chushi0/compiler/utils/hash"
	"github.com/chushi0/compiler/utils/set"
)

func (item *Item) HashCode() int {
	v := item.State
	for _, i := range item.Production {
		v = 31*v + hash.String(i)
	}
	v = 31*v + len(item.Lookahead)
	return v
}

func (item *Item) Equals(sk set.SetKey) bool {
	o, ok := sk.(*Item)
	if !ok {
		return false
	}
	if item == o {
		return true
	}
	if o.State != item.State {
		return false
	}
	if len(item.Production) != len(o.Production) {
		return false
	}
	for i, p := range item.Production {
		if o.Production[i] != p {
			return false
		}
	}
	return item.Lookahead.Equals(o.Lookahead)
}
