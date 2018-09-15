## tags in Go.

tag the structure properties

```Go
type People struct {
    Name string `json:"name"`
    Age  int8   `json:"age"`
}
```

> Tag can be any string, the following will be a common form, don't think about it.

Generally, **Tag** is in key:"value" form of such a key-value pairs.

```Go
type User struct {
    Name string `json:"name" xml:"name"`
}
```

The **key** generally refers to the package name to be used. For example, json here means that the Name field will be used and processed by the **encoding/json** package.

If there are multiple value information to be passed, it is usually comma ,delimited.

```Go
type User struct {
    Name string `json:"name,omitempty" xml:"name"`
}
```

**omittempty** indicates that this field is ignored if the value of this field is empty during conversion (Defined as false, 0, a nil pointer, a nil interface value, and any empty array, slice, map, or string). Another common one is **-** , which means that this field is ignored directly.
