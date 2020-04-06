package cmd

// requestVolume represents the number of requests
// processed over the time interval (in seconds)
// and err will be propagated if an error happened
// upstream
type requestVolume struct {
	numRequests int
	interval    int
	err         error
}
