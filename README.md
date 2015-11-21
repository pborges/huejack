# huejack
## A bare bones emulator library for the Phillips Hue bulb system to be used with the Amazon Echo (Alexa)

### Acknowledgements
Thanks to armzilla for doing all the real work ([https://github.com/armzilla/amazon-echo-ha-bridge](https://github.com/armzilla/amazon-echo-ha-bridge))

I feel this implementation is a little heavy.

I wanted a library rather then a server for easy extensibility.

I wanted to remove alot of code that I felt was _voodoo_ magic.

The only difference between this library and [this one](https://github.com/pborges/huemulator) is that in this library I attempted to emulate the golang http package and remove the need for a database of lights
### Example
see ```examples/huejack.go```

### Tips

Returning true for error in a huejack.Handler will make the echo reply with "Sorry the device is not responding"