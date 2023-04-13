package contracts

const maxPageSize = 100

type Pagination struct {
	Page      uint  `form:"page" json:"page"`
	PageSize  uint  `form:"page_size" json:"page_size"`
	TotalRows int64 `json:"total_rows"`
}

func (p *Pagination) Validate() {
	if p.PageSize > maxPageSize {
		p.PageSize = maxPageSize
	}
}

func (p *Pagination) ToOffset() int {
	return int((p.Page - 1) * p.PageSize)
}

func (p *Pagination) ToLimit() int {
	if p.PageSize == 0 {
		return -1 //this is done so gorm cancels the Limit() operation
	}
	return int(p.PageSize)
}
