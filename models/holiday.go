package models

// Holiday 代表一个节假日
type Holiday struct {
	Name  string `json:"name"`  // 节假日名称
	Start string `json:"start"` // 开始日期
	End   string `json:"end"`   // 结束日期
}

// Workday 代表一个调休的工作日
type Workday struct {
	Name string `json:"name"` // 调休名称
	Date string `json:"date"` // 调休日期
}

// YearData 包含一年的所有节假日和调休工作日
type YearData struct {
	Holidays []Holiday `json:"holidays"` // 节假日数组
	Work     []Workday `json:"work"`     // 调休工作日数组
}
