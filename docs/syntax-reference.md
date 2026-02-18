---
layout: default
title: Syntax Reference
nav_order: 3
---

# Syntax Reference

Complete reference for Mechanical Basic syntax and language constructs.

## Comments

Use the hash/pound symbol for comments:

```basic
# This is a comment
let x = 10  # Comments can also appear at the end of lines
```

## Variables

### Declaring Variables

Use `LET` to establish a new scoped version of a variable:

```basic
LET X = 5
```

If done within a function, `LET` will create a local variable that overwrites any global variables with the same name.

### Assigning Variables

Variables can be assigned without `LET` to reuse existing scope:

```basic
X = 5
```

If the variable doesn't exist, it will be created in the current scope.

### Variable Types

Variables are dynamically typed and can hold:
- **Numbers** (integers and floats)
- **Strings**
- **Boolean values**

## Data Types and Operations

### Numeric Operations

```basic
x = 1 + 5 / 5 - 1 * 1  # Basic arithmetic
x = x + 5               # Variable arithmetic
x += 1                  # Compound assignment
x++                     # Increment
x--                     # Decrement
```

**Type Behavior:** Numbers maintain their data type. If an int and a float are added together, the result is an int.

### String Operations

When strings are involved, types are automatically converted to strings:

```basic
x = "Hello" + " " + "World"  # "Hello World"
x = "Test " + 1               # "Test 1"
x = "Value: " + 3.14          # "Value: 3.14"
```

## Conditionals

### Basic If Statement

```basic
if <condition> then
    # code block
endif
```

### If-Else Statement

```basic
if <condition> then
    # code when true
else
    # code when false
endif
```

### If-ElseIf-Else Statement

```basic
if <condition> then
    # first condition code
elseif <condition> then
    # second condition code
elseif <condition> then
    # third condition code
else
    # default code
endif
```

### Conditional Examples

```basic
let age = 25

if age < 18 then
    print "Minor"
elseif age < 65 then
    print "Adult"
else
    print "Senior"
endif
```

## Loops

### For Loop

```basic
for I = 1 to 10
    # code block
next I
```

**Loop Variable:** The loop variable (`I` in this example) is automatically incremented and can be used within the loop.

### For Loop Examples

```basic
# Simple counter
for i = 1 to 5
    print i
next i

# Using loop variable in calculations
for count = 1 to 10
    let squared = count * count
    print "Square of " + count + " is " + squared
next count
```

### Break Statement

Use `break` to exit a loop early:

```basic
for i = 1 to 100
    if i * i > 50 then
        break
    endif
    print i
next i
```

## Functions

### Defining Functions

```basic
function functionName(param1, param2):
    # function body
    return value
endfunction
```

### Function Examples

```basic
# Simple addition function
function add(x, y):
    let z = x + y
    return z
endfunction

let result = add(10, 10)
print result  # Prints 20
```

```basic
# Function with conditional logic
function max(a, b):
    if a > b then
        return a
    else
        return b
    endif
endfunction

print max(5, 10)  # Prints 10
```

```basic
# Function using local variables
function calculateArea(width, height):
    let area = width * height
    return area
endfunction

let roomArea = calculateArea(12, 15)
print "Room area: " + roomArea
```

### Function Scope

Functions create their own scope. Variables declared with `LET` inside a function are local to that function:

```basic
let globalVar = 100

function testScope():
    let localVar = 50       # Local to function
    let globalVar = 200     # Shadows global variable
    print localVar          # 50
    print globalVar         # 200
endfunction

testScope()
print globalVar             # Still 100
```

## Debug Output

Use `print` to output to the terminal console or configured logger:

```basic
print "Hello world"
print "Value: " + x
print x + y + z
```

## Operator Precedence

Operations follow standard mathematical precedence:

1. Parentheses `()`
2. Multiplication `*` and Division `/`
3. Addition `+` and Subtraction `-`

```basic
let result = 2 + 3 * 4      # 14 (not 20)
let result = (2 + 3) * 4    # 20
```

## Comparison Operators

```basic
=   # Equal to
<>  # Not equal to  
<   # Less than
>   # Greater than
<=  # Less than or equal to
>=  # Greater than or equal to
```

## Logical Operators

```basic
and  # Logical AND
or   # Logical OR
not  # Logical NOT
```

### Examples

```basic
if x > 0 and x < 10 then
    print "x is between 0 and 10"
endif

if status = "ready" or status = "active" then
    print "System operational"
endif

if not gameover then
    print "Keep playing!"
endif
```

## Complete Example

Here's a comprehensive example using multiple language features:

```basic
# Fibonacci sequence calculator
function fibonacci(n):
    if n <= 1 then
        return n
    endif
    
    let a = 0
    let b = 1
    
    for i = 2 to n
        let temp = a + b
        a = b
        b = temp
    next i
    
    return b
endfunction

# Calculate and display first 10 Fibonacci numbers
print "Fibonacci Sequence:"
for i = 0 to 9
    let fib = fibonacci(i)
    print "F(" + i + ") = " + fib
next i
```

## Next Steps

- Explore [Built-in Functions](built-in-functions.md) for math and utility functions
- Learn about [External Functions](external-functions.md) to extend the language
