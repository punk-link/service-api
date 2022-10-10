package deezer

type UpcResponse struct {
	Error DeezerError `json:"error"`
	Upc   string      `json:"upc"`
	Url   string      `json:"link"`
}
