package googlecloud

import "context"

type StoreRequest struct {
	Filename string
	// The raw data of the object to be stored
	Data     []byte
	MetaData map[string]string
}

type StorageObject struct {
	Filename string
	Url      string
}

type DeleteRequest struct {
	Filename string
}

type DeleteResponse struct {
	Filename string
}

type Storage interface {
	PublicURL(filename string) string
	Store(ctx context.Context, filename string, data []byte, metadata map[string]string) error
	Delete(ctx context.Context, filename string) error
}
