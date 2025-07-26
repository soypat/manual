# manual
[![go.dev reference](https://pkg.go.dev/badge/github.com/soypat/manual)](https://pkg.go.dev/github.com/soypat/manual)
[![Go Report Card](https://goreportcard.com/badge/github.com/soypat/manual)](https://goreportcard.com/report/github.com/soypat/manual)
[![codecov](https://codecov.io/gh/soypat/manual/branch/main/graph/badge.svg)](https://codecov.io/gh/soypat/manual)
[![Go](https://github.com/soypat/manual/actions/workflows/go.yml/badge.svg)](https://github.com/soypat/manual/actions/workflows/go.yml)
[![sourcegraph](https://sourcegraph.com/github.com/soypat/manual/-/badge.svg)](https://sourcegraph.com/github.com/soypat/manual?badge)
[![License: BSD3](https://img.shields.io/badge/License-BSD3-yellow.svg)](https://opensource.org/license/bsd-3-clause)

Manual provides abstractions and implementations to work with manual memory management.

This repo is intended to be used to demonstrate manual memory management principles to students in the Go programming language.

## The Allocator interface

```go
// Allocator is the interface of a manual memory allocator.
type Allocator interface {
	// Malloc allocates a slab of memory of argument number of bytes.
	// The pointer returned points to the start address of the slab.
	// If the memory fails to allocate then nil is returned.
	Malloc(sizeInBytes int) unsafe.Pointer
	// Free receives a point previously allocated by Malloc and frees it.
	// After the memory is freed the pointer should be discarded and no other operation done with it.
	Free(toBeFreed unsafe.Pointer) error
}
```

## Example
```go
func newAllocator() manual.Allocator {
    var ta manual.TestAllocator // Ready to use as zero value.
    ta.SetMaxMemory(1024)
    return &ta
}

func doWork(alloc manual.Allocator) error {
    slice := manual.Malloc[int](alloc, 20)
    if slice == nil {
        return errors.New("allocation failed")
    }
    // do work with slice.
    err := manual.Free(alloc, slice)
    return err
}
```

## Installation

How to install package with newer versions of Go (+1.16):
```sh
go mod download github.com/soypat/manual@latest
```
