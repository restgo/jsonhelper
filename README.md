
Go JSON Utilities

JsonHelper [![GoDoc](https://godoc.org/github.com/fxding/jsonhelper?status.svg)](https://godoc.org/github.com/fxding/jsonhelper)
===========================================================================================================================

A Go Json Helper, focused on Type coercion, and json path query. You can use it for your request body, or pass json config to you function. 

*This lib is forked from [araddon/gou](https://github.com/araddon/gou), and remove other utilities. json helper is remained.

## Create a JsonHelper from json string

```go
    jh := NewJsonHelper([]byte(`{
        "name":"aaron",
        "nullstring":null,
        "ints":[1,2,3,4],
        "int":1,
        "bool": false,
        "intstr":"1",
        "int64":1234567890,
        "MaxSize" : 1048576,
        "strings":["string1"],
        "stringscsv":"string1,string2",
        "nested":{
            "nest":"string2",
            "strings":["string1"],
            "int":2,
            "list":["value"],
            "nest2":{
                "test":"good"
            }
        },
        "nested2":[
            {"sub":2}
        ],
        "period.name":"value"
    }`)
```


## Get Values

```go
    // String method
    jh.String("name", "defaultValue")    // "aaron"
    jh.String("noname", "defaultValue")  // "defaultValue"
    
    
    // Int Method
    jh.Int("int", 0)    //1
    jh.Int("noint", 20) //20
    
    // Bool Method
    jh.Bool("bool", true)   // fale
    jh.Bool("nobool", true) // true

```

## Get Value by Path

```go
    jh.Int("nested.int", 0)             // 2
    jh.String("nested.nest", "")        // "string2"
    jh.String("nested.nest2.test", "")  // "good"
    jh.String("nested.list[0]", "")     // "value"
    jh.Int("nested2[0].sub", 10)        // 2
```

License
===============
MIT License
