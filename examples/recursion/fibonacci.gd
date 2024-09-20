func fibonacci(n: int) => int {
    if n == 0 {
        return 0
    } else if n == 1 {
        return 1
    } else {
        return fibonacci(n - 1) + fibonacci(n - 2)
    }
}

pub func main() {
    set n = 15, result = fibonacci(n)
    print("Fibonacci number at position " + n + " is: " + result)
}