# Mechanical Basic

A BASIC interpreter built to add scriptable customizability to other projects.    


## Features
- BASIC-like syntax
- Ability to register functions to extend the syntax
    - You can use these functions to grab information from an entity/environemnt for example to provide realtime information to the script or to allow the script to update entities/environment


## Thoughts
The intent is that the MechBasic interpreter is light enough to support multiple instances of it dedicated for specific purposes.   IE, each entity that has events has a MechBasic interpeter instnace that has been registered with functions to update that particular entity.   Or there can be a central MechBasic instance that controls game events like weather, story triggers, etc.    


## Registering External Functions
Go:
```
mBasic := basic.NewMechanicalBasic()

getX := func(args ...interface{}) (interface{}, error) {
		if len(args) > 0 {
			return nil, nil
		}
		return 0, nil
	}

pow := func(args ...interface{}) (interface{}, error) {
    if len(args) != 2 {
        return nil, nil
    }
    a, aOk := args[0].(float64)
    b, bOk := args[1].(float64)
    if !aOk || !bOk {
        return nil, nil
    }
    math.Pow(a, b)
    return nil, nil
}

// RegisterFunc(name, arguementCount, function pointer)
mBasic.RegisterFunc("getX", getX)
mBasic.RegisterFunc("pow", pow)

code := "let x = getX()\nprint x\nx = pow(2,2)\nprint x"
mBasic.Run(code)
```

## Syntax Overview
### Conditionals:
```
if <condition> then
    ...
endif

if <condition> then
    ...
else
    ...
endif

if <condition> then
    ...
elseif <condition> then
    ...
else
    ...
endif
```

### Looping:
```
for I=1 to 10
    ...
next I
```
`break` - Special command to exit a for loop early.

### Variables:
Use `LET` to establish a new scoped version of a variable.  If done in a function this will overwrite any global variables that may be floating in.
```
LET X = 5
```
Will reuse scope if available:
```
X = 5
```

### Operations
Supported operations
Numeric:
```
x = 1 + 5 / 5 - 1 * 1 # Basic math
x = x + 5 # Variable math
x += 1 # Operate on x
x++ # Increment x
x-- # Decrement x
```
Numbers stay their data type.   If an int and a float are added together the result is an int.  

When strings are involved, types are always converted to a string.
```
x = "Hello" + " " + "World" # "Hello World"
x = "Test " + 1 # "Test 1"
```


### Functions
```
function add(x, y):
    let z = x + y
    return z
endfunction

let result = add(10,10)
print result # Prints 20
```


### Debug
Output to terminal console or supported logger if configured.
`print "Hello world"`

Comments:
```
# Use a pound or hash as a comment 
```