# json
Go json delete with field

Getting Started
===============

Installing
----------

To start using DJSON, install Go and run `go get`:

```sh
$ go get -u github.com/relunctance/djson
```

This will retrieve the library.

Set a value
-----------
Set sets the value for the specified path. 
A path is in dot syntax, such as "name.last" or "age". 
This function expects that the json is well-formed and validated. 
Invalid json will not panic, but it may return back unexpected results.
Invalid paths may return an error.


```go
package main

import "github.com/relunctance/djson"

const json = `{"name":{"first":"Janet","last":"Prichard"},"age":47}`

func main() {
    value, _ := djson.Delete(json, "name.last", "Anderson")
        println(value)
}
```

This will print:

```json
{"name":{"first":"Janet","last":"Anderson"},"age":47}
```

Path syntax
-----------

A path is a series of keys separated by a dot.
The dot and colon characters can be escaped with ``\``.

```json
{
    "name": {"first": "Tom", "last": "Anderson"},
        "age":37,
        "children": ["Sara","Alex","Jack"],
        "fav.movie": "Deer Hunter",
        "friends": [
        {"first": "James", "last": "Murphy"},
        {"first": "Roger", "last": "Craig"}
        ]
}
```
```
"name.last"          >> "Anderson"
"age"                >> 37
"children.1"         >> "Alex"
"friends.1.last"     >> "Craig"
```

The `-1` key can be used to append a value to an existing array:

```
"children.-1"  >> appends a new value to the end of the children array
```

Normally number keys are used to modify arrays, but it's possible to force a numeric object key by using the colon character:

```json
{
    "users":{
        "2313":{"name":"Sara"},
        "7839":{"name":"Andy"}
    }
}
```

A colon path would look like:

```
"users.:2313.name"    >> "Sara"
```

Supported types
---------------

Pretty much any type is supported:

```go
djson.Set(`{"key":true}`, "key", nil)
djson.Set(`{"key":true}`, "key", false)
djson.Set(`{"key":true}`, "key", 1)
djson.Set(`{"key":true}`, "key", 10.5)
djson.Set(`{"key":true}`, "key", "hello")
djson.Set(`{"key":true}`, "key", map[string]interface{}{"hello":"world"})
```

When a type is not recognized, SJSON will fallback to the `encoding/json` Marshaller.


Examples
--------

Set a value from empty document:
```go
    value, _ := djson.Set("", "name", "Tom")
println(value)

    // Output:
    // {"name":"Tom"}
    ```

    Set a nested value from empty document:
    ```go
    value, _ := djson.Set("", "name.last", "Anderson")
println(value)

    // Output:
    // {"name":{"last":"Anderson"}}
    ```

    Set a new value:
    ```go
    value, _ := djson.Set(`{"name":{"last":"Anderson"}}`, "name.first", "Sara")
println(value)

    // Output:
    // {"name":{"first":"Sara","last":"Anderson"}}
    ```

    Update an existing value:
    ```go
    value, _ := djson.Set(`{"name":{"last":"Anderson"}}`, "name.last", "Smith")
println(value)

    // Output:
    // {"name":{"last":"Smith"}}
    ```

    Set a new array value:
    ```go
    value, _ := djson.Set(`{"friends":["Andy","Carol"]}`, "friends.2", "Sara")
println(value)

    // Output:
    // {"friends":["Andy","Carol","Sara"]
    ```

    Append an array value by using the `-1` key in a path:
    ```go
    value, _ := djson.Set(`{"friends":["Andy","Carol"]}`, "friends.-1", "Sara")
println(value)

    // Output:
    // {"friends":["Andy","Carol","Sara"]
    ```

    Append an array value that is past the end:
    ```go
    value, _ := djson.Set(`{"friends":["Andy","Carol"]}`, "friends.4", "Sara")
println(value)

    // Output:
    // {"friends":["Andy","Carol",null,null,"Sara"]
    ```

    Delete a value:
    ```go
    value, _ := djson.Delete(`{"name":{"first":"Sara","last":"Anderson"}}`, "name.first")
println(value)

    // Output:
    // {"name":{"last":"Anderson"}}
    ```

    Delete an array value:
    ```go
    value, _ := djson.Delete(`{"friends":["Andy","Carol"]}`, "friends.1")
println(value)

    // Output:
    // {"friends":["Andy"]}
    ```

    Delete the last array value:
    ```go
    value, _ := djson.Delete(`{"friends":["Andy","Carol"]}`, "friends.-1")
println(value)

    // Output:
    // {"friends":["Andy"]}
    ```

