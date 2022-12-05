package kraken

type ResizeStrategy string

const (
	Auto      ResizeStrategy = "auto"
	Exact                    = "exact"
	Portrait                 = "portrait"
	Landscape                = "landscape"
	Fit                      = "fit"
	Crop                     = "crop"
)

type Resize struct {
	Width    int            `json:"width"`
	Height   int            `json:"height"`
	Strategy ResizeStrategy `json:"strategy"`
}

type Auth struct {
	ApiKey    string `json:"api_key"`
	ApiSecret string `json:"api_secret"`
}

type Request struct {
	Auth   Auth   `json:"auth"`
	Url    string `json:"url"`
	Wait   bool   `json:"wait"`
	Resize Resize `json:"resize"`
}

type Response struct {
	FileName       string `json:"file_name"`
	OriginalSize   int    `json:"original_size"`
	KrakedSize     int    `json:"kraked_size"`
	SavedBytes     int    `json:"saved_bytes"`
	KrakedUrl      string `json:"kraked_url"`
	OriginalWidth  int    `json:"original_width"`
	OriginalHeight int    `json:"original_height"`
	KrakedWidth    int    `json:"kraked_width"`
	KrakedHeight   int    `json:"kraked_height"`
	Success        bool   `json:"success"`
}

type ErrorResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}
