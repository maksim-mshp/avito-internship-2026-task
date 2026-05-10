package domain

const (
	StatusPending = "pending"
	StatusSuccess = "success"
	StatusFailed  = "failed"
)

func IsValidStatus(status string) bool {
	return status == StatusPending || status == StatusSuccess || status == StatusFailed
}
