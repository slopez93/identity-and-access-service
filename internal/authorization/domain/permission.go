package domain

type ResourceType string

type Action string

const (
	ResourceTransaction ResourceType = "TRANSACTIONS"
	ResourcePayment     ResourceType = "PAYMENTS"
	ResourceAccount     ResourceType = "ACCOUNTS"

	ActionRead    Action = "READ"
	ActionExecute Action = "EXECUTE"
	ActionCreate  Action = "CREATE"
	ActionUpdate  Action = "UPDATE"
)

type Permission struct {
	ID                string
	Resource          ResourceType
	Action            Action
	RequiresOwnership bool
}
