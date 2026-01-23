package dto

// ListResponse 列表分页响应
type ListResponse struct {
	Data  interface{} `json:"data"`  // 列表数据
	Page  int         `json:"page"`  // 当前页码
	Size  int         `json:"size"`  // 每页大小
	Total int64       `json:"total"` // 总记录数
}
