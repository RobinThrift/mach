package mach

func Shell(expr string) *Runner {
	return NewRunner(cmd{
		cmd:  "sh",
		args: append([]string{"-c"}, expr),
	})
}
