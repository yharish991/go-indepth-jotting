### Stop the program before executing any instruction

### Start it again.

### Need a program (debugger) that can control other process (our debugged program)

`ptrace`

> The ptrace() system call provides a means by which one process (the
> "tracer") may observe and control the execution of another process
> (the "tracee"), and examine and change the tracee's memory and registers. It is primarily used to implement breakpoint debugging
> and system call tracing.

`ptrace() is a single system call`

```
PtraceCont() tells the target process to restart execution
PtraceSingleStep() only allows it to run the next machine code instruction
```
