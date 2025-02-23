package kafka

type CdcEventMsg struct {
	Source CdcSource              `json:"source"`
	Op     CdcOperation           `json:"op"`
	Before map[string]interface{} `json:"before"`
	After  map[string]interface{} `json:"after"`
}

type CdcSource struct {
	Connector string `json:"connector"`
	Name      string `json:"name"`
	Db        string `json:"db"`
	Schema    string `json:"schema"`
	Table     string `json:"table"`
}

type CdcOperation string

const (
	CdcOperationCreate CdcOperation = "c"
	CdcOperationUpdate CdcOperation = "u"
	CdcOperationRead   CdcOperation = "r"
	CdcOperationDelete CdcOperation = "d"
)
