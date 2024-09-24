package test

import "testing"

func TestRecursive(t *testing.T) {
	RunTests(t, []Test{
		{`func factorial(n: int) => int {
			if n == 0 {
				return 1
			} else {
				return n * factorial(n - 1)
			}
		}

		pub func main() {
			set num = 5, result = factorial(num)
			print("Factorial of " + num + " is: " + result)
		}`, "Factorial of 5 is: 120", ""},
	})
}
