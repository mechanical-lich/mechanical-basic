---
layout: default
title: Built-in Functions
nav_order: 4
---

# Built-in Functions

Mechanical Basic includes a set of built-in mathematical functions that are available without any registration.

## Mathematical Functions

### ABS - Absolute Value

Returns the absolute value of a number.

**Syntax:**
```basic
result = ABS(number)
```

**Parameters:**
- `number` - Any numeric value

**Returns:** The absolute value (always positive or zero)

**Examples:**
```basic
let x = ABS(-5)      # Returns 5
let y = ABS(3.14)    # Returns 3.14
let z = ABS(-10.5)   # Returns 10.5

print ABS(-42)       # Prints 42
```

---

### SQR - Square Root

Returns the square root of a number.

**Syntax:**
```basic
result = SQR(number)
```

**Parameters:**
- `number` - Any non-negative numeric value

**Returns:** The square root of the number

**Examples:**
```basic
let x = SQR(16)      # Returns 4
let y = SQR(2)       # Returns 1.414...
let z = SQR(100)     # Returns 10

# Calculate distance formula
let dx = 3
let dy = 4
let distance = SQR(dx * dx + dy * dy)  # Returns 5
```

---

### INT - Integer Part

Returns the integer part of a number (typically floor function).

**Syntax:**
```basic
result = INT(number)
```

**Parameters:**
- `number` - Any numeric value

**Returns:** The integer portion of the number

**Examples:**
```basic
let x = INT(3.7)     # Returns 3
let y = INT(5.2)     # Returns 5
let z = INT(-2.8)    # Returns -2 (floor behavior)

# Round to nearest integer
let rounded = INT(value + 0.5)
```

---

### RND - Random Number

Generates a random number.

**Syntax:**
```basic
result = RND()
```

**Parameters:** None

**Returns:** A random floating-point number (typically between 0 and 1)

**Examples:**
```basic
# Basic random number
let r = RND()
print r  # Prints value like 0.742891

# Random integer from 1 to 6 (dice roll)
let dice = INT(RND() * 6) + 1
print dice

# Random integer in range [min, max]
let min = 10
let max = 20
let random = INT(RND() * (max - min + 1)) + min

# Random chance (50% probability)
if RND() > 0.5 then
    print "Heads!"
else
    print "Tails!"
endif
```

---

## Trigonometric Functions

All trigonometric functions use radians (not degrees).

### SIN - Sine

Returns the sine of an angle in radians.

**Syntax:**
```basic
result = SIN(radians)
```

**Parameters:**
- `radians` - Angle in radians

**Returns:** Sine value (range: -1 to 1)

**Examples:**
```basic
let x = SIN(0)       # Returns 0
let y = SIN(1.5708)  # Returns ~1 (π/2 radians = 90°)

# Convert degrees to radians: radians = degrees * π / 180
let pi = 3.14159265
let degrees = 45
let radians = degrees * pi / 180
let result = SIN(radians)  # Returns ~0.707
```

---

### COS - Cosine

Returns the cosine of an angle in radians.

**Syntax:**
```basic
result = COS(radians)
```

**Parameters:**
- `radians` - Angle in radians

**Returns:** Cosine value (range: -1 to 1)

**Examples:**
```basic
let x = COS(0)       # Returns 1
let y = COS(3.14159) # Returns ~-1 (π radians = 180°)

# Calculate position on circle
let angle = 1.0
let radius = 10
let x = radius * COS(angle)
let y = radius * SIN(angle)
```

---

### TAN - Tangent

Returns the tangent of an angle in radians.

**Syntax:**
```basic
result = TAN(radians)
```

**Parameters:**
- `radians` - Angle in radians

**Returns:** Tangent value

**Examples:**
```basic
let x = TAN(0)       # Returns 0
let y = TAN(0.7854)  # Returns ~1 (π/4 radians = 45°)

# Calculate slope
let angle = 0.5
let slope = TAN(angle)
```

---

### ATN - Arctangent

Returns the arctangent (inverse tangent) of a number in radians.

**Syntax:**
```basic
result = ATN(number)
```

**Parameters:**
- `number` - Any numeric value

**Returns:** Angle in radians (range: -π/2 to π/2)

**Examples:**
```basic
let x = ATN(0)       # Returns 0
let y = ATN(1)       # Returns ~0.7854 (π/4 radians = 45°)

# Calculate angle to target
let dx = targetX - currentX
let dy = targetY - currentY
let angle = ATN(dy / dx)
```

---

## Exponential and Logarithmic Functions

### EXP - Exponential

Returns e raised to the power of a number (e^x).

**Syntax:**
```basic
result = EXP(number)
```

**Parameters:**
- `number` - The exponent

**Returns:** e^number where e ≈ 2.71828

**Examples:**
```basic
let x = EXP(0)       # Returns 1
let y = EXP(1)       # Returns ~2.71828 (e)
let z = EXP(2)       # Returns ~7.389

# Exponential growth
let time = 5
let growthRate = 0.1
let population = 1000 * EXP(growthRate * time)
```

---

### LOG - Natural Logarithm

Returns the natural logarithm (base e) of a number.

**Syntax:**
```basic
result = LOG(number)
```

**Parameters:**
- `number` - Any positive numeric value

**Returns:** Natural logarithm of the number

**Examples:**
```basic
let x = LOG(1)       # Returns 0
let y = LOG(2.71828) # Returns ~1 (log of e)
let z = LOG(10)      # Returns ~2.303

# LOG is the inverse of EXP
let original = 5
let exp_val = EXP(original)
let back = LOG(exp_val)  # Returns 5
```

---

## Practical Examples

### Distance Calculation

```basic
# Calculate distance between two points
function distance(x1, y1, x2, y2):
    let dx = x2 - x1
    let dy = y2 - y1
    let dist = SQR(dx * dx + dy * dy)
    return dist
endfunction

let d = distance(0, 0, 3, 4)
print "Distance: " + d  # Prints 5
```

### Angle Calculation

```basic
# Calculate angle from current position to target
function angleToTarget(currentX, currentY, targetX, targetY):
    let dx = targetX - currentX
    let dy = targetY - currentY
    let angle = ATN(dy / dx)
    return angle
endfunction

let angle = angleToTarget(10, 10, 20, 30)
print "Angle in radians: " + angle
```

### Circular Motion

```basic
# Move entity in a circle
let pi = 3.14159265
let radius = 50
let speed = 0.1

for step = 0 to 100
    let angle = step * speed
    let x = radius * COS(angle)
    let y = radius * SIN(angle)
    print "Position: (" + x + ", " + y + ")"
next step
```

### Random Events

```basic
# Spawn random entities
for i = 1 to 10
    let x = INT(RND() * 100)
    let y = INT(RND() * 100)
    let entityType = INT(RND() * 3)  # 0, 1, or 2
    
    if entityType = 0 then
        print "Spawn enemy at (" + x + ", " + y + ")"
    elseif entityType = 1 then
        print "Spawn item at (" + x + ", " + y + ")"
    else
        print "Spawn obstacle at (" + x + ", " + y + ")"
    endif
next i
```

### Probability Distribution

```basic
# Simulate weighted random selection
function weightedRandom():
    let roll = RND()
    
    if roll < 0.5 then
        return "common"       # 50% chance
    elseif roll < 0.85 then
        return "uncommon"     # 35% chance
    elseif roll < 0.97 then
        return "rare"         # 12% chance
    else
        return "legendary"    # 3% chance
    endif
endfunction

# Test distribution
for i = 1 to 10
    let rarity = weightedRandom()
    print "Drop: " + rarity
next i
```

### Exponential Decay

```basic
# Calculate health regeneration with exponential approach
function regenerate(currentHealth, targetHealth, rate):
    let difference = targetHealth - currentHealth
    let regen = difference * (1 - EXP(-rate))
    return currentHealth + regen
endfunction

let health = 50
let maxHealth = 100
let regenRate = 0.1

# Simulate 10 ticks of regeneration
for tick = 1 to 10
    health = regenerate(health, maxHealth, regenRate)
    print "Health: " + INT(health)
next tick
```

## Function Quick Reference

| Function | Purpose | Example |
|----------|---------|---------|
| `ABS(x)` | Absolute value | `ABS(-5)` → 5 |
| `SQR(x)` | Square root | `SQR(16)` → 4 |
| `INT(x)` | Integer part | `INT(3.7)` → 3 |
| `RND()` | Random number | `RND()` → 0.742... |
| `SIN(x)` | Sine (radians) | `SIN(0)` → 0 |
| `COS(x)` | Cosine (radians) | `COS(0)` → 1 |
| `TAN(x)` | Tangent (radians) | `TAN(0)` → 0 |
| `ATN(x)` | Arctangent | `ATN(1)` → 0.7854 |
| `EXP(x)` | e raised to x | `EXP(1)` → 2.718 |
| `LOG(x)` | Natural log | `LOG(2.718)` → 1 |

## Constants

While not built-in, you can define useful mathematical constants:

```basic
let PI = 3.14159265358979
let E = 2.71828182845905
let TAU = 6.28318530717959  # 2 * PI

# Degree/Radian conversion factors
let DEG_TO_RAD = PI / 180
let RAD_TO_DEG = 180 / PI
```

## Next Steps

- Learn how to create [External Functions](external-functions.md) to extend functionality
- Review [Syntax Reference](syntax-reference.md) for language features
- See [Getting Started](getting-started.md) for setup instructions
