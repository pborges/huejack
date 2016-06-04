package main

import (
	"github.com/pborges/huejack"
	"fmt"
	"os"
)

func main() {
	huejack.SetLogger(os.Stdout)
	huejack.Handle("test", func(req huejack.Request, res *huejack.Response) {
		fmt.Println("im handling test from", req.RemoteAddr, req.RequestedOnState)
		res.OnState = req.RequestedOnState
		// res.ErrorState = true //set ErrorState to true to have the echo respond with "unable to reach device"
		return
	})

	// it is very important to use a full IP here or the UPNP does not work correctly.
	// one day ill fix this
	panic(huejack.ListenAndServe("192.168.2.103:5000"))
}
