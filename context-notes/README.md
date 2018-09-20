### Context package in simple language

WaitGroup Example
My business need two goroutine to run and print there goroutine number on exit.

```Go
func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		time.Sleep(2 * time.Second)
		fmt.Println("1st goroutine done")
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		time.Sleep(2 * time.Second)
		fmt.Println("2nd goroutine done")
		wg.Done()
	}()
	// wait for all the unfinished goroutine
	wg.Wait()
	fmt.Println("All goroutine executed")
}
```

We wait for all the goroutine to finish.

Now suppose we have a scenario(business need): "We want the 2nd goroutine to wait for 10sec and may want to notify it at certain earlier point of time to end and also send info after every 2 sec".

Now Since Go Routine runs in the background how can we achieve such scenario.

**Once goroutine starts, we can't control them. Most of the time it is waiting for it to end itself**.
Also we can have goroutine running in the background that won't end itself.

### chan stop way.

make a global variable, to send a notification to the background `goroutine` to end.

```Go
func main() {
	stop := make(chan bool)
	go func() {
		for {
			select {
			case <-stop:
				fmt.Println("Received Signal to stop")
				return
			case <-time.After(10 * time.Second):
				// timeout
				fmt.Println("timing out")
			default:
				fmt.Println("Printing in 2 second")
				time.Sleep(2 * time.Second)
			}
		}
	}()
	time.Sleep(8 * time.Second)
	stop <- true
	time.Sleep(3 * time.Second)
	fmt.Println("program stopped")
}
```

`chan+select` method is a more elegant way to end a goroutine, but this method also has limitations. If there are many goroutines that need control to end, what should I do? What if these goroutines have spawned more goroutines? What if there is an endless goroutine?

**Network Request**, each Request needs to open a goroutine to do something, these goroutine may open other goroutine. So we need a way to track goroutine, in order to achieve their purpose of control.

**Context Of Goroutine**

Rewrite the above example with `Go Context` such that Context is used to track the goroutine for control.

```Go
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("Context has been cancelled")
				return
			default:
				fmt.Println("go-routine is executing with proper time")
				time.Sleep(2 * time.Second)
			}
		}
	}(ctx)
	time.Sleep(8 * time.Second)
	fmt.Println("Cancelling the Context")
	cancel()
	time.Sleep(2 * time.Second)
	fmt.Println("program exit")
}
```

`context.Background()` Returns an empty Context, this empty Context is generally used for the root node of the entire Context tree.

`context.WithCancel(parent)` function to create a cancelable sub-Context, which was then passed to the goroutine as a parameter, so that we can use this sub-Context to track the goroutine.
**WithCancel returns a copy of parent with a new Done channel.**

The returned context's Done channel is closed when the returned cancel function is called or when the parent context's Done channel is closed, whichever happens first.

Canceling this context releases resources associated with it, so code should call cancel as soon as the operations running in this Context complete.

```Go
// multiple go Routine context
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go watch(ctx, "[1]")
	go watch(ctx, "[2]")
	go watch(ctx, "[3]")

	time.Sleep(10 * time.Second)
	fmt.Println("Cacelling all goroutines")
	cancel()
	time.Sleep(5 * time.Second)
}

func watch(ctx context.Context, name string) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println(name, "Cancelled")
			return
		default:
			fmt.Println(name, "goroutine running")
			go watchInner(ctx, name)
			time.Sleep(2 * time.Second)
		}
	}
}

func watchInner(ctx context.Context, name string) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println(name, "Cancelled Inner")
			return
		default:
			fmt.Println(name, "Inner Goroutine running")
			time.Sleep(1 * time.Second)
		}
	}
}
```

When we use cancel function notification, all three goroutine will be ended. This is the control ability of the Context. It is like a controller. After pressing the switch, all the sub-Contexts based on this Context or derived will receive a notification.
