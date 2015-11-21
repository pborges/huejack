package main
import (
	"github.com/pborges/huejack"
	"fmt"
	"os"
)


func main() {
	huejack.SetLogger(os.Stdin)
	huejack.Handle("test", func(req huejack.Request) (state bool) {
		fmt.Println("im handling test from", req.EchoId, req.RequestedOnState)
		state = !req.RequestedOnState
		return req.RequestedOnState
	})

	// it is very important to use a full IP here or the UPNP does not work correctly.
	huejack.ListenAndServe("192.168.2.102:5000")
}
