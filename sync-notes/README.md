## Golang sync package

> Go recommends communicating and synchronizing via channels.

**Values containing the types defined in this package should not be copied.**

### `Locker` type

A Locker represents an object that can be locked and unlocked.

```Go
type Locker interface {
    Lock()
    Unlock()
}
```

#### Need of Lock

```Go
func main() {
	a := 0
	for i := 0; i < 10000; i++ {
		go func(id int) {
			a++
			fmt.Println("Loop", id, "a", a)
		}(i)
	}

	time.Sleep(time.Second)
	fmt.Println("Last Value of a", a)
}
```

Last value of a will be undeterminstic: why ?

Goroutine doesn't execute in sequence. How does addition happens ?

1. read the value from the register.
2. do the addition
3. write the result to the register.

It's Likely that the two `goroutine` can fetch the value from register at the same time i.e same value of `a` is taken out. The write value to the register will end up being the same.

### `Mutex` type

A Mutex is a mutual exclusion lock. The zero value for a Mutex is an unlocked mutex.

**A Mutex must not be copied after first use.**

```Go
type Mutex struct {
// contains filtered or unexported fields
}
```

It is a concrete implementation of `Locker`, there are two ways:

```Go
func (m *Mutex) Lock()
func (m *Mutex) Unlock()
```

**Unlocking an unlocked mutex will result in a runtime error**
**A locked Mutex is not associated with a particular goroutine. It is allowed for one goroutine to lock a Mutex and then arrange for another goroutine to unlock it.**

```Go
func main() {
	a := 0
	var mu sync.Mutex
	for i := 0; i < 10000; i++ {
		go func(id int) {
			mu.Lock()
			defer mu.Unlock()
			a++
			fmt.Println("Loop", id, "a", a)
		}(i)
	}

	time.Sleep(time.Second)
	fmt.Println("Last Value of a", a)
}
```

`mutex.Lock` lock the mutex, **not a piece of code**

> During the execution when the code is reached at the place where there is a Lock, if the `mutex.Lock` cannot be acquired, it will be blocked.

In the example above, goroutine can execute in any order, but will block if it cannot acquire the `mutex.Lock` achieving the purpose of controlled synchronization.

### `RWMutex` type

A RWMutex is a reader/writer mutual exclusion lock. The lock can be held by an arbitrary number of readers or a single writer. The zero value for a RWMutex is an unlocked mutex.

```Go
func (rw *RWMutex) Lock()
func (rw *RWMutex) Unlock()

func (rw *RWMutex) RLock()
func (rw *RWMutex) RUnlock()
```

- Read lock `RLock` to lock the read operation
- Read unlock `RUnlock` to unlock the read lock
- Write lock `Lock` to lock the write operation
- Write `Unlock` to unlock the write lock

```Go
func main() {
	var mu sync.RWMutex
	a := 0

	for i := 0; i < 10; i++ {
		go func(id int) {

			mu.RLock()
			fmt.Println("Acqured Read Lock in", id)
			defer mu.RUnlock()
			// try to acquire the write Lock at i = 50000
			if id == 5 {
				go func(id int) {
					mu.Lock()
					fmt.Println("Acqured Write Lock in", id)
					a++
					fmt.Println(a)
                    time.Sleep(time.Second * 5)
                    mu.Unlock()
				}(id)
			}
			fmt.Println(a)
		}(i)
	}

	time.Sleep(time.Second * 1)
	fmt.Println(a)
}
```

Above example will terminate blcoking the rest goroutine after `write` lock is acquired.

why ?

1. Only one goroutine can get a write lock at the same time.
2. At the same time, you can have any number of gorouinte to get the read lock.
3. There can only be write locks or read locks (read and write mutex).
