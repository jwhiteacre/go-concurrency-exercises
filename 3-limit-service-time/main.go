//////////////////////////////////////////////////////////////////////
//
// Your video processing service has a freemium model. Everyone has 10
// sec of free processing time on your service. After that, the
// service will kill your process, unless you are a paid premium user.
//
// Beginner Level: 10s max per request
// Advanced Level: 10s max per user (accumulated)
//

package main

import (
	"time"
)

const MaxFreeTimeSeconds = 10

// User defines the UserModel. Use this to check whether a User is a
// Premium user or not
type User struct {
	ID        int
	IsPremium bool
	TimeUsed  int64 // in seconds
}

// HandleRequest runs the processes requested by users. Returns false
// if process had to be killed
func HandleRequest(process func(), u *User) bool {
	if u.IsPremium {
		process()
	} else {

		// channel to signal done
		done := make(chan bool)

		go func(d chan bool) {
			process()

			// let caller know we are done
			d <- true
		}(done)

		// Ticker for every second
		ticker := time.NewTicker(1 * time.Second)

		// Loop until we hit 10 seconds or are completed
		for {

			select {

			// see if go routine has completed
			case <-done:
				return true

			// check if we have hit max time
			case <-ticker.C:
				u.TimeUsed++
				if u.TimeUsed >= MaxFreeTimeSeconds {
					return false
				}
			}
		}
	}

	return true
}

func main() {
	RunMockServer()
}
