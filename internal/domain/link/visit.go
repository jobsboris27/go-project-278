package link

import "time"

type LinkVisit struct {
	ID        int64
	LinkID    int64
	IP        string
	UserAgent string
	Referer   string
	Status    int
	CreatedAt time.Time
}

func NewLinkVisit(linkID int64, ip, userAgent, referer string, status int) *LinkVisit {
	return &LinkVisit{
		LinkID:    linkID,
		IP:        ip,
		UserAgent: userAgent,
		Referer:   referer,
		Status:    status,
		CreatedAt: time.Now(),
	}
}
