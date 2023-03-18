package lunchdata

// The intention of this struct is to use it to wrap any other struct from this
// package and have a field for GTag besides it, instead of storing a GTag field
// in all other structs and try to keep that in sync.
type TagContainer struct {
	Content any
	GTag    string
}
