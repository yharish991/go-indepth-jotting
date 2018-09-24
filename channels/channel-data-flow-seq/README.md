### recvq is a linked list, observed goroutine registered first gets the data

```Go
func goRoutineA(a <-chan int) {
	val := <-a
	fmt.Println("goRoutineA received the data", val)
}

func goRoutineB(b <-chan int) {
	val := <-b
	fmt.Println("goRoutineB received the data", val)
}

func main() {
	ch := make(chan int)
	go goRoutineA(ch)
	go goRoutineB(ch)
	ch <- 3
	time.Sleep(time.Second * 1)
}
```

Channel structure

```s
*runtime.hchan {
	qcount: 0,
	dataqsiz: 0,
	buf: unsafe.Pointer(0xc0000720c0),
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
		alg: *(*runtime.typeAlg)(0x569eb0),
		gcdata: *1,
		str: 1015,
		ptrToThis: 45408,},
	sendx: 0,
	recvx: 0,
	recvq: runtime.waitq {
		first: *(*runtime.sudog)(0xc00009a000),
		last: *(*runtime.sudog)(0xc000054060),},
	sendq: runtime.waitq {
		first: *runtime.sudog nil,
		last: *runtime.sudog nil,},
	lock: runtime.mutex {key: 0},}
```

waitq structure at the time

```s
704:	func (q *waitq) dequeue() *sudog {
   705:		for {
=> 706:			sgp := q.first
   707:			if sgp == nil {
   708:				return nil
   709:			}
   710:			y := sgp.next
   711:			if y == nil {
(dlv) p q
*runtime.waitq {
	first: *runtime.sudog {
		g: *(*runtime.g)(0xc00008e180),
		isSelect: false,
		next: *(*runtime.sudog)(0xc000054060),
		prev: *runtime.sudog nil,
		elem: unsafe.Pointer(0xc000032770),
		acquiretime: 0,
		releasetime: 0,
		ticket: 0,
		parent: *runtime.sudog nil,
		waitlink: *runtime.sudog nil,
		waittail: *runtime.sudog nil,
		c: *(*runtime.hchan)(0xc0000720c0),},
	last: *runtime.sudog {
		g: *(*runtime.g)(0xc00008e300),
		isSelect: false,
		next: *runtime.sudog nil,
		prev: *(*runtime.sudog)(0xc00009a000),
		elem: unsafe.Pointer(0xc000032f70),
		acquiretime: 0,
		releasetime: 0,
		ticket: 0,
		parent: *runtime.sudog nil,
		waitlink: *runtime.sudog nil,
		waittail: *runtime.sudog nil,
		c: *(*runtime.hchan)(0xc0000720c0),},}
(dlv)
```

gp structure

```s
*runtime.g {
	stack: runtime.stack {lo: 824633925632, hi: 824633927680},
	stackguard0: 824633926512,
	stackguard1: 18446744073709551615,
	_panic: *runtime._panic nil,
	_defer: *runtime._defer nil,
	m: *runtime.m nil,
	sched: runtime.gobuf {sp: 824633927184, pc: 4378532, g: 824634302848, ctxt: unsafe.Pointer(0x0), ret: 0, lr: 0, bp: 824633927216},
	syscallsp: 0,
	syscallpc: 0,
	stktopsp: 824633927640,
	param: unsafe.Pointer(0xc00009a000),
	atomicstatus: 4,
	stackLock: 0,
	goid: 18,
	schedlink: 0,
	waitsince: 0,
	waitreason: waitReasonChanReceive,
	preempt: false,
	paniconfault: false,
	preemptscan: false,
	gcscandone: false,
	gcscanvalid: false,
	throwsplit: false,
	raceignore: 0,
	sysblocktraced: false,
	sysexitticks: 0,
	traceseq: 0,
	tracelastp: 0,
	lockedm: 0,
	sig: 0,
	writebuf: []uint8 len: 0, cap: 0, nil,
	sigcode0: 0,
	sigcode1: 0,
	sigpc: 0,
	gopc: 4836289,
	ancestors: *[]runtime.ancestorInfo nil,
	startpc: 4835520,
	racectx: 0,
	waiting: *runtime.sudog {
		g: *(*runtime.g)(0xc00008e180),
		isSelect: false,
		next: *runtime.sudog nil,
		prev: *runtime.sudog nil,
		elem: unsafe.Pointer(0x0),
		acquiretime: 0,
		releasetime: 0,
		ticket: 0,
		parent: *runtime.sudog nil,
		waitlink: *runtime.sudog nil,
		waittail: *runtime.sudog nil,
		c: *(*runtime.hchan)(0xc0000720c0),},
	cgoCtxt: []uintptr len: 0, cap: 0, nil,
	labels: unsafe.Pointer(0x0),
	timer: *runtime.timer nil,
	selectDone: 0,
	gcAssistBytes: 0,}
```

calls `goready(gp)`
**The goroutine can be made runnable again by calling goready(gp)**

Enqueue and Dequeue Options

```s
(dlv) p ch
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
		alg: *(*runtime.typeAlg)(0x569eb0),
		gcdata: *1,
		str: 1015,
		ptrToThis: 45408,},
	sendx: 0,
	recvx: 0,
	recvq: waitq<int> {
		first: *(*sudog<int>)(0xc000088000),
		last: *(*sudog<int>)(0xc000088000),},
	sendq: waitq<int> {
		first: *sudog<int> nil,
		last: *sudog<int> nil,},
	lock: runtime.mutex {key: 0},}
(dlv) p ch.recvq
waitq<int> {
	first: *sudog<int> {
		g: *(*runtime.g)(0xc000001200),
		isSelect: false,
		next: *runtime.sudog nil,
		prev: *runtime.sudog nil,
		elem: *0,
		acquiretime: 0,
		releasetime: 0,
		ticket: 0,
		parent: *runtime.sudog nil,
		waitlink: *runtime.sudog nil,
		waittail: *runtime.sudog nil,
		c: *(*runtime.hchan)(0xc0000720c0),},
	last: *sudog<int> {
		g: *(*runtime.g)(0xc000001200),
		isSelect: false,
		next: *runtime.sudog nil,
		prev: *runtime.sudog nil,
		elem: *0,
		acquiretime: 0,
		releasetime: 0,
		ticket: 0,
		parent: *runtime.sudog nil,
		waitlink: *runtime.sudog nil,
		waittail: *runtime.sudog nil,
		c: *(*runtime.hchan)(0xc0000720c0),},}
(dlv)
```

After two operation

```s
(dlv) p ch
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
		alg: *(*runtime.typeAlg)(0x569eb0),
		gcdata: *1,
		str: 1015,
		ptrToThis: 45408,},
	sendx: 0,
	recvx: 0,
	recvq: waitq<int> {
		first: *(*sudog<int>)(0xc000088000),
		last: *(*sudog<int>)(0xc00008a000),},
	sendq: waitq<int> {
		first: *sudog<int> nil,
		last: *sudog<int> nil,},
	lock: runtime.mutex {key: 0},}
(dlv) p ch.recvq
waitq<int> {
	first: *sudog<int> {
		g: *(*runtime.g)(0xc000001200),
		isSelect: false,
		next: *(*runtime.sudog)(0xc00008a000),
		prev: *runtime.sudog nil,
		elem: *0,
		acquiretime: 0,
		releasetime: 0,
		ticket: 0,
		parent: *runtime.sudog nil,
		waitlink: *runtime.sudog nil,
		waittail: *runtime.sudog nil,
		c: *(*runtime.hchan)(0xc0000720c0),},
	last: *sudog<int> {
		g: *(*runtime.g)(0xc000001380),
		isSelect: false,
		next: *runtime.sudog nil,
		prev: *(*runtime.sudog)(0xc000088000),
		elem: *0,
		acquiretime: 0,
		releasetime: 0,
		ticket: 0,
		parent: *runtime.sudog nil,
		waitlink: *runtime.sudog nil,
		waittail: *runtime.sudog nil,
		c: *(*runtime.hchan)(0xc0000720c0),},}
```

After first dequeue

```s
(dlv) p ch
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
		alg: *(*runtime.typeAlg)(0x569eb0),
		gcdata: *1,
		str: 1015,
		ptrToThis: 45408,},
	sendx: 0,
	recvx: 0,
	recvq: waitq<int> {
		first: *(*sudog<int>)(0xc00008a000),
		last: *(*sudog<int>)(0xc00008a000),},
	sendq: waitq<int> {
		first: *sudog<int> nil,
		last: *sudog<int> nil,},
	lock: runtime.mutex {key: 0},}
(dlv) p ch.recvq
waitq<int> {
	first: *sudog<int> {
		g: *(*runtime.g)(0xc000001380),
		isSelect: false,
		next: *runtime.sudog nil,
		prev: *runtime.sudog nil,
		elem: *0,
		acquiretime: 0,
		releasetime: 0,
		ticket: 0,
		parent: *runtime.sudog nil,
		waitlink: *runtime.sudog nil,
		waittail: *runtime.sudog nil,
		c: *(*runtime.hchan)(0xc0000720c0),},
	last: *sudog<int> {
		g: *(*runtime.g)(0xc000001380),
		isSelect: false,
		next: *runtime.sudog nil,
		prev: *runtime.sudog nil,
		elem: *0,
		acquiretime: 0,
		releasetime: 0,
		ticket: 0,
		parent: *runtime.sudog nil,
		waitlink: *runtime.sudog nil,
		waittail: *runtime.sudog nil,
		c: *(*runtime.hchan)(0xc0000720c0),},}
```
