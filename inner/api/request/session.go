package request

// ---------------------------------- ReadKey ----------------------------------//

type ReadKeyRequest struct {
	Key string
}

type ReadKeyResponse struct {
	Key    string
	Value  string
	Exists bool
}

// ---------------------------------- ReadPrefixKey ----------------------------------//

type ReadPrefixKeyRequest struct {
	PrefixKey string
}

type ReadPrefixKeyResponse struct {
	PrefixKey string
	Value     map[string]string
}

// ---------------------------------- PutKey ----------------------------------//

type PutKeyRequest struct {
	Key   string
	Value string
}

type PutKeyResponse struct {
	Key   string
	Value string
}

// ---------------------------------- DeleteKey ----------------------------------//

type DeleteKeyRequest struct {
	Key string
}

type DeleteKeyResponse struct {
	Key     string
	Deleted bool
}

// ---------------------------------- DeletePrefixKey ----------------------------------//

type DeletePrefixKeyRequest struct {
	PrefixKey string
}

type DeletePrefixKeyResponse struct {
	PrefixKey  string
	Deleted    bool
	DeletedKey []string
}
