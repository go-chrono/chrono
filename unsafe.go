package chrono

import _ "unsafe" // for go:linkname

//go:linkname monotime runtime.nanotime
func monotime() int64

//go:linkname walltime runtime.walltime
func walltime() (sec int64, nsec int32)
