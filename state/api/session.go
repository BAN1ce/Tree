package api

type ReadKeyRequest struct {
	Key string
}

type ReadKeyResponse struct {
	Key    string
	Value  string
	Exists bool
}

type ReadPrefixKeyRequest struct {
	PrefixKey string
}

type ReadPrefixKeyResponse struct {
	PrefixKey string
	Value     map[string]string
}
