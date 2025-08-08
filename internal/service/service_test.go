package service

import "testing"

func Test(t *testing.T) {
	m := [5]int{1, 2, 3, 4, 5}
	s := m[:]

	s[0] = 10
	t.Log(s) // Output: 10
	t.Log(m)
}
