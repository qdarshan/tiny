# Sequential Print with Threads

A demonstration of thread synchronization using two different approaches.

## Overview
This program prints numbers (1-26) and letters (A-Z) in alternating sequence from two separate threads, ensuring they execute in order despite running concurrently.

## Versions

**VersionOne.java** - Uses `wait()` and `notify()`
- Employs a shared lock object and boolean flag to coordinate between threads
- One thread waits until the other signals it's ready

**VersionTwo.java** - Uses Semaphores
- Uses two binary semaphores to control thread execution order
- Cleaner approach: numbers thread starts first (semaphore=1), letters thread waits (semaphore=0)

## Output
Both versions produce identical output:
```
1
A
2
B
3
C
...
26
Z
```
