package image

// ImageData represents the structure of the prebaked JSON data
type Data struct {
	Size     int      `json:"size"`
	Manifest string   `json:"manifest"`
	Blobs    []string `json:"blob"`
}
