package feed

import "context"

type Loader interface {
	// Load a reference in byte slice representation and return all content associated with the reference.
	Load(context.Context, []byte) ([]byte, error)
}

type Saver interface {
	// Save an arbitrary byte slice and return the reference byte slice representation.
	Save(context.Context, []byte) ([]byte, error)
}

type LoadSaver interface {
	Loader
	Saver
}
