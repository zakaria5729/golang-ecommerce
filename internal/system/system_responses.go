package system

type SystemHealthResponse struct {
	ServerStatus string `json:"server_status"`
	DBStatus     string `json:"db_status"`
	Timestamp    string `json:"timestamp"`
}
