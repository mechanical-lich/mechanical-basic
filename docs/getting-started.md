---
layout: default
title: Getting Started
nav_order: 2
---

# Getting Started

This guide will help you get up and running with Mechanical Basic in your Go project.

## Installation

Install Mechanical Basic using Go modules:

```bash
go get github.com/mechanical-lich/mechanical-basic
```

## Basic Usage

### Creating an Interpreter Instance

```go
package main

import (
    "github.com/mechanical-lich/mechanical-basic/pkg/basic"
)

func main() {
    // Create a new interpreter instance
    mBasic := basic.NewMechanicalBasic()
    
    // Run your BASIC code
    code := `
        let x = 10
        let y = 20
        let sum = x + y
        print "The sum is: " + sum
    `
    
    err := mBasic.Run(code)
    if err != nil {
        panic(err)
    }
}
```

### Running a Simple Script

Here's a complete example that demonstrates variables, conditionals, and loops:

```go
package main

import (
    "github.com/mechanical-lich/mechanical-basic/pkg/basic"
)

func main() {
    mBasic := basic.NewMechanicalBasic()
    
    code := `
        # Calculate factorial
        let n = 5
        let result = 1
        
        for i = 1 to n
            result = result * i
        next i
        
        print "Factorial of " + n + " is " + result
    `
    
    mBasic.Run(code)
}
```

## Multiple Interpreter Instances

Mechanical Basic is designed to support multiple interpreter instances, each with their own scope and registered functions:

```go
// Entity-specific interpreter
entityInterpreter := basic.NewMechanicalBasic()

// Game event interpreter
gameInterpreter := basic.NewMechanicalBasic()

// Each can have different registered functions
// and operate independently
```

## Next Steps

- Learn the complete [Syntax Reference](syntax-reference.md)
- Explore [Built-in Functions](built-in-functions.md)
- Register [External Functions](external-functions.md) to extend capabilities
- Check out examples in the [examples directory](https://github.com/mechanical-lich/mechanical-basic/tree/main/examples)

## Common Use Cases

### Entity Behavior Scripts

```go
type Entity struct {
    ID     int
    X, Y   float64
    Script *basic.MechanicalBasic
}

func (e *Entity) Update() {
    // Run entity's update script
    e.Script.Run(`
        # Entity behavior
        if getHealth() < 50 then
            seekHealth()
        else
            patrol()
        endif
    `)
}
```

### Game Event System

```go
gameEvents := basic.NewMechanicalBasic()

// Trigger weather changes
gameEvents.Run(`
    let weather = getWeather()
    if weather = "clear" then
        if RND() > 0.95 then
            setWeather("rain")
            print "It starts to rain..."
        endif
    endif
`)
```

## Error Handling

Always check for errors when running scripts:

```go
err := mBasic.Run(code)
if err != nil {
    log.Printf("Script error: %v", err)
    // Handle error appropriately
}
```
