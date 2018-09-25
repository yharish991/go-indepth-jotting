## Chan struct

> Chan contains 4 types of operations: make, read, write, and close

`c := make(chan int,2)`
compiler will eventually make the make statement

```Go
func reflect_makechan(t *chantype, size int) *hchan {
	return makechan(t, size)
}
```

```Go
func makechan(t *chantype, size int) *hchan {
	elem := t.elem

	// compiler checks this but be safe.
	if elem.size >= 1<<16 {
		throw("makechan: invalid channel element type")
	}
	if hchanSize%maxAlign != 0 || elem.align > maxAlign {
		throw("makechan: bad alignment")
	}

	if size < 0 || uintptr(size) > maxSliceCap(elem.size) || uintptr(size)*elem.size > maxAlloc-hchanSize {
		panic(plainError("makechan: size out of range"))
	}

	// Hchan does not contain pointers interesting for GC when elements stored in buf do not contain pointers.
	// buf points into the same allocation, elemtype is persistent.
	// SudoG's are referenced from their owning thread so they can't be collected.
	// TODO(dvyukov,rlh): Rethink when collector can move allocated objects.
	var c *hchan
	switch {
	case size == 0 || elem.size == 0:
		// Queue or element size is zero.
		c = (*hchan)(mallocgc(hchanSize, nil, true))
		// Race detector uses this location for synchronization.
		c.buf = unsafe.Pointer(c)
	case elem.kind&kindNoPointers != 0:
		// Elements do not contain pointers.
		// Allocate hchan and buf in one call.
		c = (*hchan)(mallocgc(hchanSize+uintptr(size)*elem.size, nil, true))
		c.buf = add(unsafe.Pointer(c), hchanSize)
	default:
		// Elements contain pointers.
		c = new(hchan)
		c.buf = mallocgc(uintptr(size)*elem.size, elem, true)
	}

	c.elemsize = uint16(elem.size)
	c.elemtype = elem
	c.dataqsiz = uint(size)

	if debugChan {
		print("makechan: chan=", c, "; elemsize=", elem.size, "; elemalg=", elem.alg, "; dataqsiz=", size, "\n")
	}
	return c
}
```

hchan structs

```Go
type hchan struct {
	qcount   uint           // total data in the queue
	dataqsiz uint           // size of the circular queue
	buf      unsafe.Pointer // points to an array of dataqsiz elements
	elemsize uint16
	closed   uint32
	elemtype *_type // element type
	sendx    uint   // send index
	recvx    uint   // receive index
	recvq    waitq  // list of recv waiters
	sendq    waitq  // list of send waiters

	// lock protects all fields in hchan, as well as several
	// fields in sudogs blocked on this channel.
	//
	// Do not change another G's status while holding this lock
	// (in particular, do not ready a G), as this can deadlock
	// with stack shrinking.
	lock mutex
}

type waitq struct {
	first *sudog
	last  *sudog
}
```

**structure at the runtime for non-buffered**

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

**structure at the runtime for buffered**

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

`dataqsize` Is the size of the buffer mentioned above, that is make(chan T, N), the N.

`elemsize` Is the size of a channel corresponding to a single element.

`closed` Indicates whether the current channel is in the closed state. After make is created, this field is set to 0, that is, the channel is open; by calling close to set it to 1, the channel is closed.

`sendx` and `recvx` is state field of a ring buffer, which indicates the current sending position and receiving position.

`recvq` and `sendq` waiting queues, which are used to store the blocked goroutines during the sending and receiving on channel.

`lock` as sending and receiving must be mutually exclusive operations.

### sudog

```Go
// sudog represents a g in a wait list, such as for sending/receiving
// on a channel.
type sudog struct {

	g *g

	// isSelect indicates g is participating in a select, so
	// g.selectDone must be CAS'd to win the wake-up race.
	isSelect bool
	next     *sudog
	prev     *sudog
	elem     unsafe.Pointer // data element (may point to stack)

	...
	c           *hchan // channel
}
```

`g` is currently blocked goroutines.

## Send

```Go
// entry point for c <- x from compiled code
//go:nosplit
func chansend1(c *hchan, elem unsafe.Pointer) {
	chansend(c, elem, true, getcallerpc())
}

/*
 * generic single channel send/recv
 * If block is not nil,
 * then the protocol will not
 * sleep but return if it could
 * not complete.
 *
 * sleep can wake up with g.param == nil
 * when a channel involved in the sleep has
 * been closed.  it is easiest to loop and re-run
 * the operation; we'll see that it's now closed.
 */
func chansend(c *hchan, ep unsafe.Pointer, block bool, callerpc uintptr) bool {
	if c == nil {
		if !block {
			return false
		}
		gopark(nil, nil, waitReasonChanSendNilChan, traceEvGoStop, 2)
		throw("unreachable")
	}

	if debugChan {
		print("chansend: chan=", c, "\n")
	}

	if raceenabled {
		racereadpc(c.raceaddr(), callerpc, funcPC(chansend))
	}

	// Fast path: check for failed non-blocking operation without acquiring the lock.
	//
	// After observing that the channel is not closed, we observe that the channel is
	// not ready for sending. Each of these observations is a single word-sized read
	// (first c.closed and second c.recvq.first or c.qcount depending on kind of channel).
	// Because a closed channel cannot transition from 'ready for sending' to
	// 'not ready for sending', even if the channel is closed between the two observations,
	// they imply a moment between the two when the channel was both not yet closed
	// and not ready for sending. We behave as if we observed the channel at that moment,
	// and report that the send cannot proceed.
	//
	// It is okay if the reads are reordered here: if we observe that the channel is not
	// ready for sending and then observe that it is not closed, that implies that the
	// channel wasn't closed during the first observation.
	if !block && c.closed == 0 && ((c.dataqsiz == 0 && c.recvq.first == nil) ||
		(c.dataqsiz > 0 && c.qcount == c.dataqsiz)) {
		return false
	}

	var t0 int64
	if blockprofilerate > 0 {
		t0 = cputicks()
	}

	lock(&c.lock)

	if c.closed != 0 {
		unlock(&c.lock)
		panic(plainError("send on closed channel"))
	}

	if sg := c.recvq.dequeue(); sg != nil {
		// Found a waiting receiver. We pass the value we want to send
		// directly to the receiver, bypassing the channel buffer (if any).
		send(c, sg, ep, func() { unlock(&c.lock) }, 3)
		return true
	}

	if c.qcount < c.dataqsiz {
		// Space is available in the channel buffer. Enqueue the element to send.
		qp := chanbuf(c, c.sendx)
		if raceenabled {
			raceacquire(qp)
			racerelease(qp)
		}
		typedmemmove(c.elemtype, qp, ep)
		c.sendx++
		if c.sendx == c.dataqsiz {
			c.sendx = 0
		}
		c.qcount++
		unlock(&c.lock)
		return true
	}

	if !block {
		unlock(&c.lock)
		return false
	}

	// Block on the channel. Some receiver will complete our operation for us.
	gp := getg()
	mysg := acquireSudog()
	mysg.releasetime = 0
	if t0 != 0 {
		mysg.releasetime = -1
	}
	// No stack splits between assigning elem and enqueuing mysg
	// on gp.waiting where copystack can find it.
	mysg.elem = ep
	mysg.waitlink = nil
	mysg.g = gp
	mysg.isSelect = false
	mysg.c = c
	gp.waiting = mysg
	gp.param = nil
	c.sendq.enqueue(mysg)
	goparkunlock(&c.lock, waitReasonChanSend, traceEvGoBlockSend, 3)

	// someone woke us up.
	if mysg != gp.waiting {
		throw("G waiting list is corrupted")
	}
	gp.waiting = nil
	if gp.param == nil {
		if c.closed == 0 {
			throw("chansend: spurious wakeup")
		}
		panic(plainError("send on closed channel"))
	}
	gp.param = nil
	if mysg.releasetime > 0 {
		blockevent(mysg.releasetime-t0, 2)
	}
	mysg.c = nil
	releaseSudog(mysg)
	return true
}
```

### sending on nil channel.

```Go
if c == nil {
		...
		gopark(nil, nil, waitReasonChanSendNilChan, traceEvGoStop, 2)
		throw("unreachable")
	}
```

**gopark** - `runtime/proc.go`
Puts the current goroutine into a waiting state with reason.

Number of reason to park the goroutines can be seen in file `runtime/runtime2.go`.

**So goroutine that sends data to the nil channel will be removed from the runnable queue.**

### sending on the closed channel.

```Go
if c.closed != 0 {
		unlock(&c.lock)
		panic(plainError("send on closed channel"))
	}
```

**Send data to the close channel, directly panic.**

### Sending Data Cases.

1. A goroutine is blocked on the channel: the data is sent directly to the goroutine

   ````Go
   if sg := c.recvq.dequeue(); sg != nil {
   			// Found a waiting receiver. We pass the value we want to send
   			// directly to the receiver, bypassing the channel buffer (if any).
   			send(c, sg, ep, func() { unlock(&c.lock) }, 3)
   			return true
   		}
   		```
   ````

   Take the waiting goroutine from the `recvq` queue of the current channel and then call send

   ```Go
   func send(c *hchan, sg *sudog, ep unsafe.Pointer, unlockf func(), skip int) {
   	...
   	if sg.elem != nil {
   		sendDirect(c.elemtype, sg, ep)
   		sg.elem = nil
   	}
   	gp := sg.g
   	unlockf()
   	gp.param = unsafe.Pointer(sg)
   	if sg.releasetime != 0 {
   		sg.releasetime = cputicks()
   	}
   	goready(gp, skip+1)
   }
   ```

   **The goroutine can be made runnable again by calling goready(gp)**
