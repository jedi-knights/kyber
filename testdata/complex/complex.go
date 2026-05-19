// Package complex is a kyber test fixture: a deliberately branchy function
// whose McCabe cyclomatic complexity is hand-computed below.
//
// Branches counted (in addition to the base 1):
//   - 3 if statements (nested or sequential, all count)               +3
//   - 1 for loop                                                       +1
//   - 1 range loop                                                     +1
//   - switch with 4 non-default cases (the default does NOT count)     +4
//   - 1 `&&` in the condition `x > 0 && y > 0`                         +1
//   - 1 `||` in the condition `x < 0 || y < 0`                         +1
//
// Total: 1 + 3 + 1 + 1 + 4 + 1 + 1 = 12.
package complex

import "fmt"

// Branchy is intentionally complex. Expected cyclomatic complexity: 12.
func Branchy(x, y int, items []string) string {
	if x > 0 && y > 0 {
		fmt.Println("both positive")
	}
	if x < 0 || y < 0 {
		fmt.Println("at least one negative")
	}
	if x == y {
		fmt.Println("equal")
	}

	for i := 0; i < x; i++ {
		fmt.Println(i)
	}
	for _, item := range items {
		fmt.Println(item)
	}

	switch x {
	case 1:
		return "one"
	case 2:
		return "two"
	case 3:
		return "three"
	case 4:
		return "four"
	default:
		return "many"
	}
}
