```s
Type 'help' for list of commands.
(dlv) break main.main
Breakpoint 1 set at 0x45473f for main.main() ./main.go:7
(dlv) ls
> _rt0_amd64_linux() /usr/local/go/src/runtime/rt0_linux_amd64.s:8 (PC: 0x44ff30)
Warning: debugging optimized function
     3:	// license that can be found in the LICENSE file.
     4:
     5:	#include "textflag.h"
     6:
     7:	TEXT _rt0_amd64_linux(SB),NOSPLIT,$-8
=>   8:		JMP	_rt0_amd64(SB)
     9:
    10:	TEXT _rt0_amd64_linux_lib(SB),NOSPLIT,$0
    11:		JMP	_rt0_amd64_lib(SB)
(dlv) n
> main.main() ./main.go:7 (hits goroutine(1):1 total:1) (PC: 0x45473f)
     2:
     3:	func goroutineA(c2 chan int) {
     4:		c2 <- 2
     5:	}
     6:
=>   7:	func main() {
     8:		c2 := make(chan int)
     9:		go goroutineA(c2)
    10:
    11:		for {
    12:		}
(dlv) n
> main.main() ./main.go:8 (PC: 0x45474d)
     3:	func goroutineA(c2 chan int) {
     4:		c2 <- 2
     5:	}
     6:
     7:	func main() {
=>   8:		c2 := make(chan int)
     9:		go goroutineA(c2)
    10:
    11:		for {
    12:		}
    13:	}
(dlv) p c2
chan int {
	qcount: 0,
	dataqsiz: 0,
	buf: *[0]int nil,
	elemsize: 0,
	closed: 0,
	elemtype: *runtime._type nil,
	sendx: 0,
	recvx: 0,
	recvq: waitq<int> {
		first: *sudog<int> nil,
		last: *sudog<int> nil,},
	sendq: waitq<int> {
		first: *sudog<int> nil,
		last: *sudog<int> nil,},
	lock: runtime.mutex {key: 0},}
(dlv) n
> main.main() ./main.go:9 (PC: 0x454770)
     4:		c2 <- 2
     5:	}
     6:
     7:	func main() {
     8:		c2 := make(chan int)
=>   9:		go goroutineA(c2)
    10:
    11:		for {
    12:		}
    13:	}
(dlv) p c2
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
		alg: *(*runtime.typeAlg)(0x4bff90),
		gcdata: *1,
		str: 775,
		ptrToThis: 28320,},
	sendx: 0,
	recvx: 0,
	recvq: waitq<int> {
		first: *sudog<int> nil,
		last: *sudog<int> nil,},
	sendq: waitq<int> {
		first: *sudog<int> nil,
		last: *sudog<int> nil,},
	lock: runtime.mutex {key: 0},}
(dlv) c2.buf
Command failed: command not available
(dlv) p c2.buf
*[0]int []
(dlv) n
> main.main() ./main.go:11 (PC: 0x45478d)
     6:
     7:	func main() {
     8:		c2 := make(chan int)
     9:		go goroutineA(c2)
    10:
=>  11:		for {
    12:		}
    13:	}
(dlv) p c2
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
		alg: *(*runtime.typeAlg)(0x4bff90),
		gcdata: *1,
		str: 775,
		ptrToThis: 28320,},
	sendx: 0,
	recvx: 0,
	recvq: waitq<int> {
		first: *sudog<int> nil,
		last: *sudog<int> nil,},
	sendq: waitq<int> {
		first: *(*sudog<int>)(0xc000074000),
		last: *(*sudog<int>)(0xc000074000),},
	lock: runtime.mutex {key: 0},}
(dlv) p c2.buf
*[0]int []
(dlv) p c2.sendq
waitq<int> {
	first: *sudog<int> {
		g: *(*runtime.g)(0xc000001080),
		isSelect: false,
		next: *runtime.sudog nil,
		prev: *runtime.sudog nil,
		elem: *2,
		acquiretime: 0,
		releasetime: 0,
		ticket: 0,
		parent: *runtime.sudog nil,
		waitlink: *runtime.sudog nil,
		waittail: *runtime.sudog nil,
		c: *(*runtime.hchan)(0xc00001e120),},
	last: *sudog<int> {
		g: *(*runtime.g)(0xc000001080),
		isSelect: false,
		next: *runtime.sudog nil,
		prev: *runtime.sudog nil,
		elem: *2,
		acquiretime: 0,
		releasetime: 0,
		ticket: 0,
		parent: *runtime.sudog nil,
		waitlink: *runtime.sudog nil,
		waittail: *runtime.sudog nil,
		c: *(*runtime.hchan)(0xc00001e120),},}
(dlv) p c2.sendq.first.elem
*2
(dlv) p c2
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
		alg: *(*runtime.typeAlg)(0x4bff90),
		gcdata: *1,
		str: 775,
		ptrToThis: 28320,},
	sendx: 0,
	recvx: 0,
	recvq: waitq<int> {
		first: *sudog<int> nil,
		last: *sudog<int> nil,},
	sendq: waitq<int> {
		first: *(*sudog<int>)(0xc000074000),
		last: *(*sudog<int>)(0xc000074000),},
	lock: runtime.mutex {key: 0},}
(dlv)
```

In channel , If the buffer is full (`unbuffered channel is always full if no receiver is there`) then the element to be written is saved in the structure of the currently executing goroutine i.e `sudog`. That's why during the demo @Gaurav Agarwal was unable to access the element for `pch.buf`
