---
layout: default
title: Home
nav_order: 1
---

# Mechanical Basic

A BASIC interpreter built to add scriptable customizability to other projects.

## Overview

Mechanical Basic is a lightweight BASIC interpreter designed to enable scripting capabilities in Go applications. It provides a familiar BASIC-like syntax while allowing seamless integration with Go code through external function registration.

## Key Features

- **BASIC-like Syntax** - Easy to learn and use familiar programming constructs
- **Extensible Function System** - Register Go functions to extend the language capabilities
- **Lightweight Design** - Support multiple interpreter instances for different purposes
- **Real-time Integration** - Scripts can interact with and modify runtime entities/environments
- **Entity-Specific Scripting** - Each entity can have its own interpreter instance with custom functions

## Use Cases

The Mechanical Basic interpreter is designed to be light enough to support multiple instances dedicated for specific purposes:

- **Entity Scripting** - Each entity with events has its own interpreter instance registered with functions to update that particular entity
- **Game Events** - A central instance controls game events like weather, story triggers, etc.
- **Environment Control** - Scripts can query and modify the game environment in real-time
- **Custom Behaviors** - Define complex behaviors without recompiling your application

## Quick Example

```basic
# Simple script with external function call
let x = getEntityX()
print "Entity X position: " + x

if x > 100 then
    moveEntity(-10, 0)
    print "Moving entity left"
endif
```

## Documentation

- [Getting Started](getting-started.md) - Installation and basic setup
- [Syntax Reference](syntax-reference.md) - Complete language syntax guide
- [Built-in Functions](built-in-functions.md) - Math and utility functions
- [External Functions](external-functions.md) - Registering Go functions

## Installation

```bash
go get github.com/mechanical-lich/mechanical-basic
```

## License

See [LICENSE](https://github.com/mechanical-lich/mechanical-basic/blob/main/LICENSE) for details.
