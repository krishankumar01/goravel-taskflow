package enums

type Status string

const (
	StatusPending    Status = "pending"
	StatusInProgress Status = "in-progress"
	StatusDone       Status = "done"
)
