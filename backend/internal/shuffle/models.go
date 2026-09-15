package shuffle

type ActionResult struct {
	Action      Action `json:"action"`
	ExecutionId string `json:"execution_id"`
	Result      string `json:"result"`
	StartedAt   int64  `json:"started_at"`
	CompletedAt int64  `json:"completed_at"`
	Status      string `json:"status"`
}

type Action struct {
	AppName string `json:"app_name"`
	Name    string `json:"name"`
	Label   string `json:"label"`
}
