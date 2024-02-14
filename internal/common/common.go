package common

var Version string = "dev"

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

func BaseURL(server string) string {
	return server + "/api/"
}

func CollectionsStructureEndpoint() string {
	return "collections.documents"
}

func CollectionsGetEndpoint() string {
	return "collections.info"
}

func CollectionsListEndpoint() string {
	return "collections.list"
}

func CollectionsCreateEndpoint() string {
	return "collections.create"
}

func CollectionsUpdateEndpoint() string {
	return "collections.update"
}

func DocumentsGetEndpoint() string {
	return "documents.info"
}

func DocumentsCreateEndpoint() string {
	return "documents.create"
}

func DocumentsUpdateEndpoint() string {
	return "documents.update"
}

func AttachmentsCreateEndpoint() string {
	return "attachments.create"
}
