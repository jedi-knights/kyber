// Package unreadable is a kyber test fixture for the readability metric.
// Single function with deep nesting, short identifiers, no comments, and
// >40 lines — every readability sub-signal should be unfavorable.
package unreadable

func Tangled(a, b, c, d, e int) int {
	x := 0
	if a > 0 {
		if b > 0 {
			if c > 0 {
				if d > 0 {
					if e > 0 {
						x = a + b + c + d + e
					} else {
						x = a + b + c + d
					}
				} else {
					x = a + b + c
				}
			} else {
				x = a + b
			}
		} else {
			x = a
		}
	} else {
		x = 0
	}
	if a < 0 {
		if b < 0 {
			if c < 0 {
				if d < 0 {
					if e < 0 {
						x = -a - b - c - d - e
					}
				}
			}
		}
	}
	if x > 0 {
		if x < 100 {
			x = x * 2
		}
	}
	return x
}
