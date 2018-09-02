package runtime

import (
	"unsafe"
)

const Compiler = "tgo"

func Sleep(d Duration) {
	// This function is treated specially by the compiler: when goroutines are
	// used, it is transformed into a llvm.coro.suspend() call.
	// When goroutines are not used this function behaves as normal.
	sleep(d)
}

func GOMAXPROCS(n int) int {
	// Note: setting GOMAXPROCS is ignored.
	return 1
}

func GOROOT() string {
	// TODO: don't hardcode but take the one at compile time.
	return "/usr/local/go"
}

// Copy size bytes from src to dst. The memory areas must not overlap.
func memcpy(dst, src unsafe.Pointer, size uintptr) {
	for i := uintptr(0); i < size; i++ {
		*(*uint8)(unsafe.Pointer(uintptr(dst) + i)) = *(*uint8)(unsafe.Pointer(uintptr(src) + i))
	}
}

// Set the given number of bytes to zero.
func memzero(ptr unsafe.Pointer, size uintptr) {
	for i := uintptr(0); i < size; i++ {
		*(*byte)(unsafe.Pointer(uintptr(ptr) + i)) = 0
	}
}

// Compare two same-size buffers for equality.
func memequal(x, y unsafe.Pointer, n uintptr) bool {
	for i := uintptr(0); i < n; i++ {
		cx := *(*uint8)(unsafe.Pointer(uintptr(x) + i))
		cy := *(*uint8)(unsafe.Pointer(uintptr(y) + i))
		if cx != cy {
			return false
		}
	}
	return true
}

func _panic(message interface{}) {
	printstring("panic: ")
	printitf(message)
	printnl()
	abort()
}

// Check for bounds in *ssa.IndexAddr and *ssa.Lookup.
func lookupBoundsCheck(length, index int) {
	if index < 0 || index >= length {
		// printstring() here is safe as this function is excluded from bounds
		// checking.
		printstring("panic: runtime error: index out of range\n")
		abort()
	}
}

// Check for bounds in *ssa.Slice
func sliceBoundsCheck(length, low, high uint) {
	if !(0 <= low && low <= high && high <= length) {
		printstring("panic: runtime error: slice out of range\n")
		abort()
	}
}
