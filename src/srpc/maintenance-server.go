package srpc


type MaintenanceOrder struct {
	Id              int64       `json:"Id"`
		CommonUser      string      `json:"CommonUser"`
		MaintenanceUser string      `json:"MaintenanceUser"`
		Timestamp       int64       `json:"Timestamp"`
		Receiver        string      `json:"Receiver"`
		Reason          string      `json:"Reason"`
		Note            string      `json:"Note`
		Status          int `json:"column:status"`
}

func RequestOrder() (res struct{
	Status bool
	Message string
	Order MaintenanceOrder }) {

	}