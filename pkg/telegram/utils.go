package telegram

import "io"

// It's bit `a python way` to pick non-blank value of a two.
func orInt(first int, second int) int {
	if first == 0 {
		return second
	}
	return first
}

func orString(first string, second string) string {
	if first == "" {
		return second
	}
	return first
}


// Close invokes close method without handling error
func Close(closer io.Closer) {
	_ = closer.Close()
}

func safeClose(values []interface{}) func() {
	var closers []io.Closer
	for _, value := range values {
		if in, ok := value.(io.Closer); ok {
			closers = append(closers, in)
		}
	}
	if len(closers) == 0 {
		return noop
	}
	return func(){
		for _, closer := range closers {
			_ = closer.Close()
		}
	}
}

func noop() {}
