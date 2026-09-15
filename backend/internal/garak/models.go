package garak

type GarakEval struct {
	EntryType string `json:"entry_type"` // e.g. "eval"
	Eval      map[string]map[string]interface{} `json:"eval"`
}
