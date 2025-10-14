package system

type SystemHealthResponse struct {
	ServerStatus string `json:"server_status"`
	DBStatus     string `json:"db_status"`
	Timestamp    string `json:"timestamp"`
}

type SystemLogFileResponse struct {
	FileName   string `json:"file_name"`
	Size       string `json:"size"`
	ModifiedAt string `json:"modified_at"`
}
