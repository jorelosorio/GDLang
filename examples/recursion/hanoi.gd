func moveDisk(from: string, to: string) {
    print("Move disk from " + from + " to " + to + "\n")
}

func towersOfHanoi(n: int, from: string, to: string, aux: string) {
    if n == 1 {
        moveDisk(from, to)
    } else {
        towersOfHanoi(n - 1, from, aux, to)
        moveDisk(from, to)
        towersOfHanoi(n - 1, aux, to, from)
    }
}

pub func main() {
    set numDisks = 3
    towersOfHanoi(numDisks, "A", "C", "B")
}