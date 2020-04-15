package main

func SafeClose(ch chan struct{}) (justClosed bool) {
	defer func() {
		if recover() != nil {
			justClosed = false
		}
	}()

	close(ch)
	return true
}

func SafeSend(ch chan struct{}, value struct{}) (closed bool) {
	defer func() {
		if recover() != nil {
			closed = true
		}
	}()
	ch <- value
	return false
}
