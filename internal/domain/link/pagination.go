package link

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	DefaultLimit = 10
	MaxLimit     = 100
)

var ErrInvalidRange = errors.New("invalid range format")

type Pagination struct {
	Offset int
	Limit  int
}

func NewPagination(offset, limit int) *Pagination {
	if limit <= 0 {
		limit = DefaultLimit
	}
	if limit > MaxLimit {
		limit = MaxLimit
	}
	if offset < 0 {
		offset = 0
	}
	return &Pagination{
		Offset: offset,
		Limit:  limit,
	}
}

func ParseRange(rangeStr string) (*Pagination, error) {
	if rangeStr == "" {
		return NewPagination(0, DefaultLimit), nil
	}

	rangeStr = strings.TrimSpace(rangeStr)
	if !strings.HasPrefix(rangeStr, "[") || !strings.HasSuffix(rangeStr, "]") {
		return nil, ErrInvalidRange
	}

	inner := strings.TrimPrefix(strings.TrimSuffix(rangeStr, "]"), "[")
	parts := strings.Split(inner, ",")
	if len(parts) != 2 {
		return nil, ErrInvalidRange
	}

	offset, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return nil, ErrInvalidRange
	}

	limit, err := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil {
		return nil, ErrInvalidRange
	}

	if offset < 0 || limit <= 0 {
		return nil, ErrInvalidRange
	}

	return NewPagination(offset, limit), nil
}

func (p *Pagination) ContentRange(total int) string {
	if total == 0 {
		return fmt.Sprintf("links 0-0/0")
	}

	last := p.Offset + p.Limit - 1
	if last >= total {
		last = total - 1
	}

	return fmt.Sprintf("links %d-%d/%d", p.Offset, last, total)
}
