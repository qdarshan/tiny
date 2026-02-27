# CSP Sieve of Eratosthenes

A concurrent prime number sieve implemented in Go using Communicating Sequential Processes (CSP) style channels and goroutines.

## How It Works

The program finds all prime numbers below 100 by chaining together goroutine filters connected via channels — a concurrent take on the classic [Sieve of Eratosthenes](https://en.wikipedia.org/wiki/Sieve_of_Eratosthenes).

1. **Generator** — `main` sends integers 3–99 into the first channel.
2. **Filter chain** — Each `Filter` goroutine owns a prime number. It reads from its input channel, discards multiples of its prime, and forwards survivors to the next filter in the chain.
3. **Dynamic growth** — The first number that passes through a filter is itself prime, so a new `Filter` goroutine is spawned on the fly to handle it. This builds a pipeline of filters, one per discovered prime.

```
3,4,5,6,7,… ──▶ Filter(2) ──▶ Filter(3) ──▶ Filter(5) ──▶ Filter(7) ──▶ …
                 drops 4,6,8…   drops 9,15…   drops 25,35…   drops 49,77…
```

A `sync.WaitGroup` ensures the program waits for all goroutines to finish before exiting.

## Running

```sh
go run main.go
```
