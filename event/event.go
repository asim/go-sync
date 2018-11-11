// Package event provides a distributed log interface
package event

// Event provides a distributed log interface
type Event interface {
	Log(name string) (Log, error)
}

type Log interface {
	// Close the log handle
	Close() error
	// Read will read the next record
	Read() (*Record, error)
	// Go to an offset
	Seek(offset int64) error
	// Write an event to the log
	Write(*Record) error
}

type Record struct {
	Metadata map[string]interface{}
	Data     []byte
}
