package main

// resolutionType holds the enums meant to explain why a subscriber was denied a particular parcel
const (
	delivered resolutionType = 1 + iota
	expired
	forbidden
)


type parcel struct {
	hashIndex string
	available bool
	resolution string
}

var subscribers = map[string]bool
