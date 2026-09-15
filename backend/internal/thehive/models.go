package thehive

type Case struct {
	ID          string   `json:"id"`
	Number      int      `json:"number"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Severity    int      `json:"severity"` // 1=Low, 2=Medium, 3=High, 4=Critical
	StartDate   int64    `json:"startDate"`
	Status      string   `json:"status"`
	Tags        []string `json:"tags"`
	Assignee    string   `json:"assignee"`
}

type WebhookEvent struct {
	ObjectType string `json:"objectType"`
	Operation  string `json:"operation"`
	Object     Case   `json:"object"` // Note: normally this could be Alert, Task, etc., but we're focusing on Cases for this mock adapter.
}
