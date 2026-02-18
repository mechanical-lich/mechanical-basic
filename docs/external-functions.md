---
layout: default
title: External Functions
nav_order: 5
---

# External Functions

One of the most powerful features of Mechanical Basic is the ability to register Go functions that can be called from your BASIC scripts. This allows scripts to interact with your application's runtime environment, query game state, modify entities, and more.

## Overview

External functions enable:
- **Real-time Data Access** - Scripts can query entity positions, health, game state, etc.
- **Environment Modification** - Scripts can update entities, trigger events, or change game world
- **Custom Operations** - Extend the language with domain-specific functionality
- **Integration** - Seamlessly bridge between BASIC scripts and Go code

## Registering Functions

### Basic Registration

Use `RegisterFunc` to register a Go function with your interpreter instance:

```go
mBasic := basic.NewMechanicalBasic()

// Register a function
mBasic.RegisterFunc("functionName", functionPointer)
```

### Function Signature

External functions must follow this signature:

```go
func(args ...interface{}) (interface{}, error)
```

- **Parameters**: Variable number of arguments as `interface{}`
- **Return**: Single value (any type) and an error

## Simple Examples

### Zero-Argument Function

```go
mBasic := basic.NewMechanicalBasic()

getX := func(args ...interface{}) (interface{}, error) {
    if len(args) > 0 {
        return nil, fmt.Errorf("getX takes no arguments")
    }
    return 42, nil
}

mBasic.RegisterFunc("getX", getX)
```

**Usage in BASIC:**
```basic
let x = getX()
print x  # Prints 42
```

### Single-Argument Function

```go
square := func(args ...interface{}) (interface{}, error) {
    if len(args) != 1 {
        return nil, fmt.Errorf("square requires exactly 1 argument")
    }
    
    num, ok := args[0].(float64)
    if !ok {
        return nil, fmt.Errorf("square requires a numeric argument")
    }
    
    return num * num, nil
}

mBasic.RegisterFunc("square", square)
```

**Usage in BASIC:**
```basic
let result = square(5)
print result  # Prints 25
```

### Multi-Argument Function

```go
pow := func(args ...interface{}) (interface{}, error) {
    if len(args) != 2 {
        return nil, fmt.Errorf("pow requires exactly 2 arguments")
    }
    
    base, baseOk := args[0].(float64)
    exp, expOk := args[1].(float64)
    
    if !baseOk || !expOk {
        return nil, fmt.Errorf("pow requires numeric arguments")
    }
    
    return math.Pow(base, exp), nil
}

mBasic.RegisterFunc("pow", pow)
```

**Usage in BASIC:**
```basic
let result = pow(2, 8)
print result  # Prints 256
```

## Practical Examples

### Entity Query Functions

```go
type Entity struct {
    ID     int
    X, Y   float64
    Health int
}

var currentEntity *Entity

func setupEntityFunctions(mBasic *basic.MechanicalBasic, entity *Entity) {
    currentEntity = entity
    
    // Get entity X position
    getEntityX := func(args ...interface{}) (interface{}, error) {
        return currentEntity.X, nil
    }
    
    // Get entity Y position
    getEntityY := func(args ...interface{}) (interface{}, error) {
        return currentEntity.Y, nil
    }
    
    // Get entity health
    getHealth := func(args ...interface{}) (interface{}, error) {
        return float64(currentEntity.Health), nil
    }
    
    mBasic.RegisterFunc("getX", getEntityX)
    mBasic.RegisterFunc("getY", getEntityY)
    mBasic.RegisterFunc("getHealth", getHealth)
}
```

**Usage in BASIC:**
```basic
let x = getX()
let y = getY()
let health = getHealth()

print "Entity at (" + x + ", " + y + ")"
print "Health: " + health
```

### Entity Modification Functions

```go
func setupEntityModifiers(mBasic *basic.MechanicalBasic, entity *Entity) {
    // Move entity relative to current position
    moveEntity := func(args ...interface{}) (interface{}, error) {
        if len(args) != 2 {
            return nil, fmt.Errorf("moveEntity requires dx, dy")
        }
        
        dx, _ := args[0].(float64)
        dy, _ := args[1].(float64)
        
        entity.X += dx
        entity.Y += dy
        return nil, nil
    }
    
    // Set entity position absolutely
    setPosition := func(args ...interface{}) (interface{}, error) {
        if len(args) != 2 {
            return nil, fmt.Errorf("setPosition requires x, y")
        }
        
        x, _ := args[0].(float64)
        y, _ := args[1].(float64)
        
        entity.X = x
        entity.Y = y
        return nil, nil
    }
    
    // Damage entity
    takeDamage := func(args ...interface{}) (interface{}, error) {
        if len(args) != 1 {
            return nil, fmt.Errorf("takeDamage requires damage amount")
        }
        
        damage, _ := args[0].(float64)
        entity.Health -= int(damage)
        
        if entity.Health < 0 {
            entity.Health = 0
        }
        
        return float64(entity.Health), nil
    }
    
    mBasic.RegisterFunc("moveEntity", moveEntity)
    mBasic.RegisterFunc("setPosition", setPosition)
    mBasic.RegisterFunc("takeDamage", takeDamage)
}
```

**Usage in BASIC:**
```basic
# Move entity based on health
if getHealth() < 50 then
    # Retreat when low health
    moveEntity(-10, 0)
else
    # Advance when healthy
    moveEntity(10, 0)
endif

# Take environmental damage
if getY() < 0 then
    takeDamage(5)
    print "Fell into hazard!"
endif
```

### Game State Functions

```go
type GameState struct {
    Weather    string
    TimeOfDay  int
    Score      int
    GameOver   bool
}

var gameState *GameState

func setupGameFunctions(mBasic *basic.MechanicalBasic, gs *GameState) {
    gameState = gs
    
    getWeather := func(args ...interface{}) (interface{}, error) {
        return gameState.Weather, nil
    }
    
    setWeather := func(args ...interface{}) (interface{}, error) {
        if len(args) != 1 {
            return nil, fmt.Errorf("setWeather requires weather type")
        }
        weather, _ := args[0].(string)
        gameState.Weather = weather
        return nil, nil
    }
    
    getScore := func(args ...interface{}) (interface{}, error) {
        return float64(gameState.Score), nil
    }
    
    addScore := func(args ...interface{}) (interface{}, error) {
        if len(args) != 1 {
            return nil, fmt.Errorf("addScore requires points")
        }
        points, _ := args[0].(float64)
        gameState.Score += int(points)
        return float64(gameState.Score), nil
    }
    
    mBasic.RegisterFunc("getWeather", getWeather)
    mBasic.RegisterFunc("setWeather", setWeather)
    mBasic.RegisterFunc("getScore", getScore)
    mBasic.RegisterFunc("addScore", addScore)
}
```

**Usage in BASIC:**
```basic
# Weather-based behavior
let weather = getWeather()
if weather = "rain" then
    print "Seeking shelter..."
    moveEntity(0, -5)
endif

# Score milestone events
let score = getScore()
if score > 1000 then
    print "Achievement unlocked!"
    addScore(500)
endif
```

## Best Practices

### 1. Validate Arguments

Always check argument count and types:

```go
myFunc := func(args ...interface{}) (interface{}, error) {
    if len(args) != 2 {
        return nil, fmt.Errorf("expected 2 arguments, got %d", len(args))
    }
    
    arg1, ok := args[0].(float64)
    if !ok {
        return nil, fmt.Errorf("argument 1 must be a number")
    }
    
    // ... rest of function
}
```

### 2. Return Meaningful Errors

Provide clear error messages that help debug script issues:

```go
return nil, fmt.Errorf("moveEntity: invalid target position (%v, %v)", x, y)
```

### 3. Keep Functions Focused

Register many small, specific functions rather than few large ones:

```go
// Good: Separate functions
mBasic.RegisterFunc("getX", getX)
mBasic.RegisterFunc("getY", getY)

// Less ideal: One function with mode parameter
mBasic.RegisterFunc("getPosition", getPosition)  // Requires "x" or "y" arg
```

### 4. Document Expected Arguments

```go
// getDistanceTo(targetX, targetY) - Calculate distance to target position
getDistanceTo := func(args ...interface{}) (interface{}, error) {
    // implementation
}
```

### 5. Use Closures for Context

Capture necessary context in the closure:

```go
func CreateEntityScriptInstance(entity *Entity, world *World) *basic.MechanicalBasic {
    mBasic := basic.NewMechanicalBasic()
    
    // Functions capture entity and world
    getX := func(args ...interface{}) (interface{}, error) {
        return entity.X, nil
    }
    
    getNearbyEntities := func(args ...interface{}) (interface{}, error) {
        // Can access both entity and world
        return float64(len(world.GetEntitiesNear(entity.X, entity.Y, 50))), nil
    }
    
    mBasic.RegisterFunc("getX", getX)
    mBasic.RegisterFunc("getNearbyEntities", getNearbyEntities)
    
    return mBasic
}
```

## Complete Example

Here's a complete example showing entity scripting with external functions:

```go
package main

import (
    "fmt"
    "github.com/mechanical-lich/mechanical-basic/pkg/basic"
)

type Entity struct {
    Name   string
    X, Y   float64
    Health int
}

func main() {
    entity := &Entity{
        Name:   "Player",
        X:      50,
        Y:      50,
        Health: 100,
    }
    
    // Create interpreter with entity functions
    mBasic := basic.NewMechanicalBasic()
    
    // Register entity functions
    mBasic.RegisterFunc("getX", func(args ...interface{}) (interface{}, error) {
        return entity.X, nil
    })
    
    mBasic.RegisterFunc("getY", func(args ...interface{}) (interface{}, error) {
        return entity.Y, nil
    })
    
    mBasic.RegisterFunc("getHealth", func(args ...interface{}) (interface{}, error) {
        return float64(entity.Health), nil
    })
    
    mBasic.RegisterFunc("moveBy", func(args ...interface{}) (interface{}, error) {
        if len(args) != 2 {
            return nil, fmt.Errorf("moveBy requires dx, dy")
        }
        dx, _ := args[0].(float64)
        dy, _ := args[1].(float64)
        entity.X += dx
        entity.Y += dy
        return nil, nil
    })
    
    // Run entity behavior script
    script := `
        print "Entity: " + getName()
        print "Position: (" + getX() + ", " + getY() + ")"
        print "Health: " + getHealth()
        
        # Behavior logic
        if getHealth() < 50 then
            print "Low health! Retreating..."
            moveBy(-10, 0)
        else
            print "Advancing..."
            moveBy(10, 0)
        endif
        
        print "New position: (" + getX() + ", " + getY() + ")"
    `
    
    err := mBasic.Run(script)
    if err != nil {
        fmt.Printf("Script error: %v\n", err)
    }
}
```

## Next Steps

- Review [Built-in Functions](built-in-functions.md) for available math functions
- See [examples directory](https://github.com/mechanical-lich/mechanical-basic/tree/main/examples) for more complete examples
- Read [Syntax Reference](syntax-reference.md) for language features
