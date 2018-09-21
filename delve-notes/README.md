```
$ ▶ dlv debug main.go
Type 'help' for list of commands.
(dlv)
```

**list (alias: ls | l) -------- Show source code.**

```
(dlv) list main.main
Showing /home/ankur/Documents/Others/go-http-indepth/delve-notes/hello-world/main.go:7 (PC: 0x49c01f)
   2:
   3:	import (
   4:		"fmt"
   5:	)
   6:
   7:	func main() {
   8:		fmt.Println("Hello World")
   9:	}
(dlv)
```

**set break point**

```
(dlv) break main.main
Breakpoint 1 set at 0x49c01f for main.main() ./main.go:7
(dlv)
```

**continue at the break point**

```
(dlv) c
> main.main() ./main.go:7 (hits goroutine(1):1 total:1) (PC: 0x49c01f)
     2:
     3:	import (
     4:		"fmt"
     5:	)
     6:
=>   7:	func main() {
     8:		fmt.Println("Hello World")
     9:	}
(dlv)
```

> What is the state of the goroutine at this breakpoint ?

**goroutine shows current goroutine**

**goroutines List program goroutines**

```
(dlv) goroutine
Thread 5607 at ./main.go:7
Goroutine 1:
	Runtime: ./main.go:7 main.main (0x49c01f)
	User: ./main.go:7 main.main (0x49c01f)
	Go: /usr/local/go/src/runtime/asm_amd64.s:201 runtime.rt0_go (0x454cd9)
	Start: /usr/local/go/src/runtime/proc.go:110 runtime.main (0x42c990)
(dlv) goroutines
[4 goroutines]
* Goroutine 1 - User: ./main.go:7 main.main (0x49c01f) (thread 5607)
  Goroutine 2 - User: /usr/local/go/src/runtime/proc.go:303 runtime.gopark (0x42cf24)
  Goroutine 3 - User: /usr/local/go/src/runtime/proc.go:303 runtime.gopark (0x42cf24)
  Goroutine 4 - User: /usr/local/go/src/runtime/proc.go:303 runtime.gopark (0x42cf24)
(dlv)
```

state of threads

```
(dlv) threads
* Thread 5607 at 0x49c01f ./main.go:7 main.main
  Thread 5973 at 0x4585fd /usr/local/go/src/runtime/sys_linux_amd64.s:131 runtime.usleep
  Thread 5974 at 0x458b83 /usr/local/go/src/runtime/sys_linux_amd64.s:532 runtime.futex
  Thread 5975 at 0x458b83 /usr/local/go/src/runtime/sys_linux_amd64.s:532 runtime.futex
  Thread 5976 at 0x458b83 /usr/local/go/src/runtime/sys_linux_amd64.s:532 runtime.futex
```

**4 goroutines and 5 threads ?** why

```Shell
  ankur ▶ ~ ▶ $ ▶ lscpu | grep "CPU(s)"
CPU(s):                4
```

Also, in addition to threads for the number of logical CPUs, it seems that one extra thread is used by golang itself for control(need verification).

**backtrace**

```
(dlv) goroutine 1
Switched from 1 to 1 (thread 5607)
(dlv) bt
0  0x000000000049c01f in main.main
   at ./main.go:7
1  0x000000000042cb55 in runtime.main
   at /usr/local/go/src/runtime/proc.go:201
2  0x0000000000456d31 in runtime.goexit
   at /usr/local/go/src/runtime/asm_amd64.s:1333
(dlv)
```
