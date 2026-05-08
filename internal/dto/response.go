package dto

// 响应数据体
type Response struct {
	Code int32       `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// 分页数据
type PageData struct {
	Count int64      `json:"count"`
	List  interface{} `json:"list"`
}

func NewPageData(count int64, list interface{}) PageData {
	return PageData{Count: count, List: list}
}

type PageParams struct {
	PageIndex int `json:"page_index" query:"page_index"`
	PageSize  int `json:"page_size" query:"page_size"`
}

const defaultPageSize = 20

func (p *PageParams) Normalize() {
	if p.PageIndex <= 0 {
		p.PageIndex = 1
	}
	if p.PageSize <= 0 {
		p.PageSize = defaultPageSize
	}
}
