```s
Type 'help' for list of commands.
(dlv) c
received SIGINT, stopping process (will not forward signal)
> main.main() ./main.go:24 (PC: 0x49c320)
    19:		c2 := make(chan int)
    20:
    21:		go goroutineA(c1, c2)
    22:		go goroutineB(c1, c2)
    23:
=>  24:		for {
    25:		}
    26:	}
(dlv)
```

**no of goroutines**

```s
(dlv) goroutines
[6 goroutines]
* Goroutine 1 - User: ./main.go:24 main.main (0x49c320) (thread 12823)
  Goroutine 2 - User: /usr/local/go/src/runtime/proc.go:303 runtime.gopark (0x42cfa4)
  Goroutine 17 - User: /usr/local/go/src/runtime/proc.go:303 runtime.gopark (0x42cfa4)
  Goroutine 18 - User: /usr/local/go/src/runtime/proc.go:303 runtime.gopark (0x42cfa4)
  Goroutine 19 - User: ./main.go:6 main.goroutineA (0x49c0ce)
  Goroutine 20 - User: ./main.go:12 main.goroutineB (0x49c153)
(dlv)
```

```s
(dlv) threads
* Thread 12823 at 0x49c320 ./main.go:24 main.main
  Thread 12831 at 0x45867d /usr/local/go/src/runtime/sys_linux_amd64.s:131 runtime.usleep
  Thread 12832 at 0x458c03 /usr/local/go/src/runtime/sys_linux_amd64.s:532 runtime.futex
  Thread 12833 at 0x458c03 /usr/local/go/src/runtime/sys_linux_amd64.s:532 runtime.futex
  Thread 12834 at 0x458c03 /usr/local/go/src/runtime/sys_linux_amd64.s:532 runtime.futex
(dlv)
```

**Goroutine that is about to receive the channels**

```s
(dlv) goroutine 19
Switched from 1 to 19 (thread 12823)
(dlv) bt
0  0x000000000042cfa4 in runtime.gopark
   at /usr/local/go/src/runtime/proc.go:303
1  0x000000000042d063 in runtime.goparkunlock
   at /usr/local/go/src/runtime/proc.go:308
2  0x0000000000405784 in runtime.chanrecv
   at /usr/local/go/src/runtime/chan.go:520
3  0x000000000040558b in runtime.chanrecv1
   at /usr/local/go/src/runtime/chan.go:402
4  0x000000000049c0ce in main.goroutineA
   at ./main.go:6
5  0x0000000000456db1 in runtime.goexit
   at /usr/local/go/src/runtime/asm_amd64.s:1333
(dlv)
```

**content of channel c1 of goroutineA by specifying frame**

```s
(dlv) frame 4 p c1
chan int {
	qcount: 0,
	dataqsiz: 0,
	buf: *[0]int [],
	elemsize: 8,
	closed: 0,
	elemtype: *runtime._type {
		size: 8,
		ptrdata: 0,
		hash: 4149441018,
		tflag: tflagUncommon|tflagExtraStar|tflagNamed,
		align: 8,
		fieldalign: 8,
		kind: 130,
		alg: *(*runtime.typeAlg)(0x568eb0),
		gcdata: *1,
		str: 1015,
		ptrToThis: 45376,},
	sendx: 0,
	recvx: 0,
	recvq: waitq<int> {
		first: *(*sudog<int>)(0xc00009a000),
		last: *(*sudog<int>)(0xc00009a000),},
	sendq: waitq<int> {
		first: *sudog<int> nil,
		last: *sudog<int> nil,},
	lock: runtime.mutex {key: 0},}
```

**What is at the top of the stack ?**

```s
(dlv) stack
0  0x000000000042cfa4 in runtime.gopark
   at /usr/local/go/src/runtime/proc.go:303
1  0x000000000042d063 in runtime.goparkunlock
   at /usr/local/go/src/runtime/proc.go:308
2  0x0000000000405784 in runtime.chanrecv
   at /usr/local/go/src/runtime/chan.go:520
3  0x000000000040558b in runtime.chanrecv1
   at /usr/local/go/src/runtime/chan.go:402
4  0x000000000049c0ce in main.goroutineA
   at ./main.go:6
5  0x0000000000456db1 in runtime.goexit
   at /usr/local/go/src/runtime/asm_amd64.s:1333
```

```s
(dlv) list /usr/local/go/src/runtime/proc.go:303
Showing /usr/local/go/src/runtime/proc.go:303 (PC: 0x42cfa4)
 298:		mp.waittraceev = traceEv
 299:		mp.waittraceskip = traceskip
 300:		releasem(mp)
 301:		// can't do anything that might move the G between Ms here.
 302:		mcall(park_m)
 303:	}
 304:
 305:	// Puts the current goroutine into a waiting state and unlocks the lock.
 306:	// The goroutine can be made runnable again by calling goready(gp).
 307:	func goparkunlock(lock *mutex, reason waitReason, traceEv byte, traceskip int) {
 308:		gopark(parkunlock_c, unsafe.Pointer(lock), reason, traceEv, traceskip)
(dlv)
```

Where our goroutine started ?

```s
(dlv) frame 4 list
Goroutine 19 frame 4 at /home/ankur/Documents/Others/go-http-indepth/channels/learn-channel-debug/main.go:6 (PC: 0x49c0ce)
     1:	package main
     2:
     3:	import "fmt"
     4:
     5:	func goroutineA(c1 chan int, c2 chan int) {
=>   6:		_ = <-c2
     7:		c1 <- 1
     8:		return
     9:	}
    10:
    11:	func goroutineB(c1 chan int, c2 chan int) {
(dlv)
```

So we have a receive `chan`.

```s
Goroutine 19 frame 3 at /usr/local/go/src/runtime/chan.go:402 (PC: 0x40558b)
   397:	}
   398:
   399:	// entry points for <- c from compiled code
   400:	//go:nosplit
   401:	func chanrecv1(c *hchan, elem unsafe.Pointer) {
=> 402:		chanrecv(c, elem, true)
   403:	}
   404:
   405:	//go:nosplit
   406:	func chanrecv2(c *hchan, elem unsafe.Pointer) (received bool) {
   407:		_, received = chanrecv(c, elem, true)
(dlv)
```

```s
(dlv) frame 2 list
Goroutine 19 frame 2 at /usr/local/go/src/runtime/chan.go:520 (PC: 0x405784)
   515:		mysg.g = gp
   516:		mysg.isSelect = false
   517:		mysg.c = c
   518:		gp.param = nil
   519:		c.recvq.enqueue(mysg)
=> 520:		goparkunlock(&c.lock, waitReasonChanReceive, traceEvGoBlockRecv, 3)
   521:
   522:		// someone woke us up
   523:		if mysg != gp.waiting {
   524:			throw("G waiting list is corrupted")
   525:		}
(dlv)
```

**put the golang in the park state**

```s
(dlv) frame 1 list
Goroutine 19 frame 1 at /usr/local/go/src/runtime/proc.go:308 (PC: 0x42d063)
   303:	}
   304:
   305:	// Puts the current goroutine into a waiting state and unlocks the lock.
   306:	// The goroutine can be made runnable again by calling goready(gp).
   307:	func goparkunlock(lock *mutex, reason waitReason, traceEv byte, traceskip int) {
=> 308:		gopark(parkunlock_c, unsafe.Pointer(lock), reason, traceEv, traceskip)
   309:	}
   310:
   311:	func goready(gp *g, traceskip int) {
   312:		systemstack(func() {
   313:			ready(gp, traceskip, true)
(dlv)
```

**A Buffered Channel at runtime**

```s
(dlv) frame 4 p c1
chan int {
	qcount: 0,
	dataqsiz: 3,
	buf: *[3]int [0,0,0],
	elemsize: 8,
	closed: 0,
	elemtype: *runtime._type {
		size: 8,
		ptrdata: 0,
		hash: 4149441018,
		tflag: tflagUncommon|tflagExtraStar|tflagNamed,
		align: 8,
		fieldalign: 8,
		kind: 130,
		alg: *(*runtime.typeAlg)(0x568eb0),
		gcdata: *1,
		str: 1015,
		ptrToThis: 45376,},
	sendx: 0,
	recvx: 0,
	recvq: waitq<int> {
		first: *(*sudog<int>)(0xc000080060),
		last: *(*sudog<int>)(0xc000080060),},
	sendq: waitq<int> {
		first: *sudog<int> nil,
		last: *sudog<int> nil,},
	lock: runtime.mutex {key: 0},}
(dlv)
```

**Channel structure after putting the value**

```s
(dlv) print c1
chan int {
	qcount: 1,
	dataqsiz: 3,
	buf: *[3]int [2,0,0],
	elemsize: 8,
	closed: 0,
	elemtype: *runtime._type {
		size: 8,
		ptrdata: 0,
		hash: 4149441018,
		tflag: tflagUncommon|tflagExtraStar|tflagNamed,
		align: 8,
		fieldalign: 8,
		kind: 130,
		alg: *(*runtime.typeAlg)(0x568eb0),
		gcdata: *1,
		str: 1015,
		ptrToThis: 45376,},
	sendx: 1,
	recvx: 0,
	recvq: waitq<int> {
		first: *sudog<int> nil,
		last: *sudog<int> nil,},
	sendq: waitq<int> {
		first: *sudog<int> nil,
		last: *sudog<int> nil,},
	lock: runtime.mutex {key: 0},}
(dlv) step
```

**channel c1 after value has been read.**

```s
(dlv) frame 0 print c1
chan int {
	qcount: 0,
	dataqsiz: 3,
	buf: *[3]int [0,0,0],
	elemsize: 8,
	closed: 0,
	elemtype: *runtime._type {
		size: 8,
		ptrdata: 0,
		hash: 4149441018,
		tflag: tflagUncommon|tflagExtraStar|tflagNamed,
		align: 8,
		fieldalign: 8,
		kind: 130,
		alg: *(*runtime.typeAlg)(0x568eb0),
		gcdata: *1,
		str: 1015,
		ptrToThis: 45376,},
	sendx: 1,
	recvx: 1,
	recvq: waitq<int> {
		first: *sudog<int> nil,
		last: *sudog<int> nil,},
	sendq: waitq<int> {
		first: *sudog<int> nil,
		last: *sudog<int> nil,},
	lock: runtime.mutex {key: 0},}
(dlv)
```

**\*Data when put on the blocking channel**

```s
(dlv) p c
*runtime.hchan {
	qcount: 0,
	dataqsiz: 0,
	buf: unsafe.Pointer(0xc000082060),
	elemsize: 8,
	closed: 0,
	elemtype: *runtime._type {
		size: 8,
		ptrdata: 0,
		hash: 4149441018,
		tflag: tflagUncommon|tflagExtraStar|tflagNamed,
		align: 8,
		fieldalign: 8,
		kind: 130,
		alg: *(*runtime.typeAlg)(0x568eb0),
		gcdata: *1,
		str: 1015,
		ptrToThis: 45376,},
	sendx: 0,
	recvx: 0,
	recvq: runtime.waitq {
		first: *runtime.sudog nil,
		last: *runtime.sudog nil,},
	sendq: runtime.waitq {
		first: *runtime.sudog nil,
		last: *runtime.sudog nil,},
	lock: runtime.mutex {key: 0},}
(dlv)
```

```Go
func main() {
	c1 := make(chan int, 3)
	c2 := make(chan int)
	//c1 <- 2
	c2 <- 2
```

```s
(dlv) p c
*runtime.hchan {
	qcount: 0,
	dataqsiz: 0,
	buf: unsafe.Pointer(0xc000082060),
	elemsize: 8,
	closed: 0,
	elemtype: *runtime._type {
		size: 8,
		ptrdata: 0,
		hash: 4149441018,
		tflag: tflagUncommon|tflagExtraStar|tflagNamed,
		align: 8,
		fieldalign: 8,
		kind: 130,
		alg: *(*runtime.typeAlg)(0x568eb0),
		gcdata: *1,
		str: 1015,
		ptrToThis: 45376,},
	sendx: 0,
	recvx: 0,
	recvq: runtime.waitq {
		first: *runtime.sudog nil,
		last: *runtime.sudog nil,},
	sendq: runtime.waitq {
		first: *(*runtime.sudog)(0xc000080000),
		last: *(*runtime.sudog)(0xc000080000),},
	lock: runtime.mutex {key: 1},}
(dlv)
```

```Go
func main() {
	c1 := make(chan int, 3)
	c2 := make(chan int)
	//c1 <- 2
	// c2 <- 2
	go goroutineA(c1, c2)
	c2 <- 2
	go goroutineB(c1, c2)

	for {
	}
}
```

```s
(dlv) p c
*runtime.hchan {
	qcount: 0,
	dataqsiz: 0,
	buf: unsafe.Pointer(0xc000072060),
	elemsize: 8,
	closed: 0,
	elemtype: *runtime._type {
		size: 8,
		ptrdata: 0,
		hash: 4149441018,
		tflag: tflagUncommon|tflagExtraStar|tflagNamed,
		align: 8,
		fieldalign: 8,
		kind: 130,
		alg: *(*runtime.typeAlg)(0x568eb0),
		gcdata: *1,
		str: 1015,
		ptrToThis: 45376,},
	sendx: 0,
	recvx: 0,
	recvq: runtime.waitq {
		first: *(*runtime.sudog)(0xc000088000),
		last: *(*runtime.sudog)(0xc000088000),},
	sendq: runtime.waitq {
		first: *runtime.sudog nil,
		last: *runtime.sudog nil,},
	lock: runtime.mutex {key: 0},}
(dlv)
```
