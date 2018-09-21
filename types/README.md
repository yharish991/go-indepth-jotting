### [Type Identity](https://golang.org/ref/spec#Type_identity)

> Two types are either identical or different.

### [Function Types](https://golang.org/ref/spec#Function_types)

> A function type denotes the set of all functions with the same parameter and result types.

means that each function has its own function type
If two functions have the same signature (parameter and result types), they share one function type

**By writing `type Hello func...` we're just giving a name to a particular function type, not defining a new one**

_Two function types are identical if they have the same number of parameters and result values, corresponding parameter and result types are identical, and either both functions are variadic or neither is. Parameter and result names are not required to match._
