func printStars(count: int) {
    if count > 0 {
        print("*")
        printStars(count - 1)
    }
}

func printSpaces(count: int) {
    if count > 0 {
        print(" ")
        printSpaces(count - 1)
    }
}

func printTriangle(rows: int, currentRow: int) {
    if currentRow <= rows {
        printSpaces(rows - currentRow) // Print leading spaces
        printStars(2 * currentRow - 1) // Print stars for the current row
        print("\n") // Move to the next line
        printTriangle(rows, currentRow + 1)
    }
}

pub func main() {
    set n = 5
    printTriangle(n, 1)
}