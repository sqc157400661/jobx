package types

type Pagination struct {
	PageIndex int `form:"pageIndex" json:"pageIndex" query:"pageIndex"`
	PageSize  int `form:"pageSize"  json:"pageSize"  query:"pageIndex"`
}

func (m *Pagination) GetPageIndex() int {
	if m.PageIndex <= 0 {
		m.PageIndex = 1
	}
	return m.PageIndex
}

func (m *Pagination) GetPageSize() int {
	if m.PageSize <= 0 {
		m.PageSize = 10
	}
	return m.PageSize
}
