package hash

func String(s string) int {
	v := 0
	for _, rn := range s {
		v = 31*v + int(rn)
	}
	return v
}
