package enums

type Status string

const (
	NotStarted Status = "not_started"
	InProgress Status = "in_progress"
	Testing    Status = "testing"
	Completed  Status = "completed"
	Backlog    Status = "backlog"
)
