package main

// Assets represents an asset storage access
type Assets interface {
	String(name string) string
}
