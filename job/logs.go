package job

// logs provides thread-safe operations for reading job logs
type logs struct {
	stdout *threadSafeBuffer
	stderr *threadSafeBuffer
}

func (l *logs) Output() string {
	return l.stdout.String()
}

func (l *logs) Errors() string {
	return l.stderr.String()
}
