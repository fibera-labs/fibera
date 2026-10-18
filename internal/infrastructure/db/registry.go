package db

// TableModel 定义参与数据库同步的数据表模型
type TableModel interface {
	TableName() string
}

var tables []TableModel

// Register 注册数据表模型，由各模型在 init 阶段调用
func Register(m TableModel) {
	tables = append(tables, m)
}

// Tables 返回全部已注册的数据表模型
func Tables() []TableModel {
	return tables
}
