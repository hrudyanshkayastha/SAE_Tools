package cortex


type JobReport struct {
	ID           string                 `json:"id"`
	Status       string                 `json:"status"`
	Date         int64                  `json:"date"`
	AnalyzerId   string                 `json:"analyzerId"`
	AnalyzerName string                 `json:"analyzerName"`
	Report       map[string]interface{} `json:"report"`
}
