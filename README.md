# robin

Package `robin` exports a round-robin structure `Robin` for comparable types that supports O(1) addition and removal of values. It can either be unbounded or bounded by a maximum length in which case an optional buffer can be provided. The buffer will receive values added to a full `Robin` and then replace values that are removed. A buffer implementation `LIFOBuffer` is included in the package. This is a stack-like buffer with a fixed capacity, where pushing to a full buffer overwrites the oldest value.

New values are added in the current position and a subsequent `Next()` call will return the first of the added values. `Robin` is constrained to unique values, duplicates are ignored. Additionally, `Robin` is not thread-safe on its own, wrap function calls with a mutex or some other synchronization primitive for concurrent access.

## Usage

### Basic usage

```go
r := robin.NewUnbounded[int]()
r.Add(0, 1, 2)
for i := 0; i < 2; i++ {
    n, _ := r.Next()
    fmt.Printf("%d ", n)
}
r.Add(3)
for i := 0; i < 6; i++ {
    n, _ := r.Next()
    fmt.Printf("%d ", n)
}
fmt.Println()

// Output:
// 0 1 3 2 0 1 3 2 
```

### Buffered robin

```go
r := robin.NewBounded(
    3,
    robin.WithBuffer[int](robin.NewLIFOBuffer[int](2)),
)
r.Add(0, 1, 2, 3, 4, 5)
for i := 0; i < 6; i++ {
    n, _ := r.Next()
    fmt.Printf("%d ", n)
}
fmt.Println()

r.Remove(0, 1)
for i := 0; i < 6; i++ {
    n, _ := r.Next()
    fmt.Printf("%d ", n)
}
fmt.Println()

// Output:
// 0 1 2 0 1 2 
// 5 4 2 5 4 2 
```
