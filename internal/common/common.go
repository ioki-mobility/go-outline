package common

const (
	HdrKeyAuthorization string = "Authorization"
	HdrKeyContentType   string = "Content-Type"
	HdrValueContentType string = "application/json"
	HdrKeyAccept        string = "Accept"
	HdrValueAccept      string = "application/json"
)

func HdrValueAuthorization(key string) string {
	return "Bearer " + key
}

func CollectionsGetEndpoint() string {
	return "collections.info"
}

func CollectionsListEndpoint() string {
	return "collections.list"
}
