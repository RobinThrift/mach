package main

import (
	. "github.com/RobinThrift/mach"
)

func main() {
	Run("ls examples").Pipe(Run("wc -l")).Go()
}
