package utils

// retryableFn is the function signature in a retry.
type retryableFn func(attempt int) (retry bool, err error)

// Do runs fn and retries if fn fails. Do exits when either error returned by fn is nil,
// retry returned by fn is false, or maxRetries exceeded.
func Do(fn retryableFn) error {
	var err error
	var cont bool
	attempt := 1
	for {
		cont, err = fn(attempt)
		if !cont || err == nil {
			break
		}
		attempt++
	}
	return err
}
