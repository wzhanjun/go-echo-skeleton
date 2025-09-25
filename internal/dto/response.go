package dto

// 响应数据体
type Response struct {
	// 业务状态错误码 0 = 正常， 其他：错误
	Code int32 `json:"code"`
	// 错误信息
	Msg string `json:"msg"`
	// 业务数据
	Data interface{} `json:"data"`
}

// 分页内容
type WrapPageData struct {
	// 总记录数
	Count int32 `json:"count"`
	// 分页数据
	List interface{} `json:"list"`
}

type PageParams struct {
	// 分页页码
	PageIndex int `json:"page_index" query:"page_index"`
	// 分页大小
	PageSize int `json:"page_size" query:"page_size"`
}
