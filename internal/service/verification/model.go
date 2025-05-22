package verification

import "time"

const (
	verificationTable = "verification"
)

type Verification struct {
	ID        int64
	UserID    int64
	Token     string
	IsUsed    bool
	CreatedAt time.Time
	ExpiredAt time.Time
}
