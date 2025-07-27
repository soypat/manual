package manual

import (
	"errors"
	"math"
	"unsafe"
)

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

// Malloc provides generic slice allocation similar to the `make` built in.
func Malloc[T any](alloc Allocator, n int) []T {
	var z T
	ptr := alloc.Malloc(n * int(unsafe.Sizeof(z)))
	if ptr == nil {
		return nil
	}
	return unsafe.Slice((*T)(ptr), n)
}

// Free calls [Allocator.Free] on the data portion of the buffer. Keep in mind the buffer
// must be the same one returned my Malloc. If the slice start pointer has been resliced Free will fail.
func Free[T any](alloc Allocator, b []T) error {
	if b == nil {
		panic("nil pointer argument to Free")
	}
	return alloc.Free(unsafe.Pointer(unsafe.SliceData(b)))
}

// TestAllocator is a simple implementation of an [Allocator].
// It has the added complexity of being able to reuse freed memory later on to potentially detect
// use-after-free and double-free bugs.
type TestAllocator struct {
	maxmem    int
	currlive  int
	live      [][]byte
	currfree  int
	maxfree   int
	free      [][]byte
	afterFree func(b []byte)
}

// SetMaxMemory sets the maximum amount of memory the allocator can have allocated, combined between live and freed.
// Set to zero to disable memory limits.
func (a *TestAllocator) SetMaxMemory(maxmemInBytes int) {
	a.maxmem = maxmemInBytes
}

// SetMaxFree sets the maximum amount of freed memory to have allocated and ready for reuse.
// Set to zero for limitless memory. Set to -1 to disable freed memory reuse.
func (a *TestAllocator) SetMaxFree(maxFreeMemInBytes int) {
	a.maxfree = maxFreeMemInBytes
}

// SetOnFreeCallback is called on [TestAllocator.Free] argument buffer after it successfully finds the memory to free.
func (a *TestAllocator) SetOnFreeCallback(onFree func(b []byte)) {
	a.afterFree = onFree
}

// Malloc implements [Allocator.Malloc].
func (a *TestAllocator) Malloc(n int) unsafe.Pointer {
	minFreeMatch := math.MaxInt
	freeMatch := -1
	for i := range a.free {
		l := len(a.free[i])
		if l >= n && l < minFreeMatch {
			freeMatch = i
			minFreeMatch = l
		}
	}
	if freeMatch >= 0 {
		nowlive := a.free[freeMatch]
		a.free[freeMatch] = a.free[len(a.free)-1]
		a.free = a.free[:len(a.free)-1]
		a.live = append(a.live, nowlive)
		a.currfree -= minFreeMatch
		a.currlive += minFreeMatch
		return unsafe.Pointer(unsafe.SliceData(nowlive))
	}
	currlive := a.currlive + n
	if a.maxmem != 0 && currlive+a.currfree > a.maxmem {
		return nil
	}
	b := make([]byte, n)
	a.currlive = currlive
	a.live = append(a.live, b)
	return unsafe.Pointer(unsafe.SliceData(b))
}

// Free implements [Allocator.Free].
func (a *TestAllocator) Free(p unsafe.Pointer) error {
	for i := range a.live {
		memstart := unsafe.Pointer(unsafe.SliceData(a.live[i]))
		if memstart == p {
			freed := a.live[i]
			a.live[i] = a.live[len(a.live)-1]
			a.live = a.live[:len(a.live)-1] // pop.
			currfree := a.currfree + len(freed)
			if a.maxfree == 0 || currfree < a.maxfree {
				a.free = append(a.free, freed)
				a.currfree = currfree
			}
			if a.afterFree != nil {
				a.afterFree(freed)
			}
			return nil
		}
	}
	return errors.New("pointer not found in allocations")
}
