package service

type PageParam struct {
	Page     int
	PageSize int
}

func (p *PageParam) Offset() int {
	return (p.Page - 1) * p.PageSize
}

func (p *PageParam) Limit() int {
	return p.PageSize
}
