package model

import (
	"fmt"
	"gorm.io/gorm"
	gorm_logger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"math"
	"gin-api/pkg/logger"
	"sync"
	"time"
)

var (
	db       *gorm.DB
	once     sync.Once
	pageObj  *Paginate
)

type Model struct{}

// TimestampsField 时间戳
type TimestampsField struct {
	CreatedAt time.Time `gorm:"column:created_at;index;" json:"created_at,omitempty"`
	UpdatedAt time.Time `gorm:"column:updated_at;index;" json:"updated_at,omitempty"`
}

// Logger 返回 gorm logger
func Logger() gorm_logger.Interface {
	return gorm_logger.New(
		writer{},
		gorm_logger.Config{
			LogLevel:                  gorm_logger.Info,
			SlowThreshold:             time.Second,
			IgnoreRecordNotFoundError: true,
		})
}

//writer 实现 gorm_logger.Writer 接口
type writer struct {}

//Printf 实现 gorm_logger.Writer 接口
func (m writer) Printf(format string, v ...interface{}) {
	log := fmt.Sprintf(format, v...)
	//sql := strings.Split(log, "\n")[1]
	logger.AccessLog("sql-log", log)
}

//SetDB 设置数据库连接
func SetDB(DB *gorm.DB)  {
	db = DB
}

//GetDB 获取数据库连接
func GetDB() *gorm.DB {
	return  db.Session(&gorm.Session{NewDB: true})
}

//Create 插入数据
func Create(model interface{}) int64 {
	return GetDB().Create(model).RowsAffected
}

//CreateBatch 分配插入
func CreateBatch(model interface{}, chunkSize int) int64 {
	return GetDB().CreateInBatches(model, chunkSize).RowsAffected
}

//Column 获取一列数据
func Column(model interface{}, field string, where string) int64 {
	return GetDB().Where(where).Pluck(field, model).RowsAffected
}

//Take 获取单行数据
func Take(model interface{}, where string, order string) int64 {
	return GetDB().Where(where).Order(order).Take(model).RowsAffected
}

//Find 获取多行数据
func Find(model interface{}, where string, order string) int64 {
	return GetDB().Where(where).Order(order).Find(model).RowsAffected
}

//TableName 基于模型返回数据表名称
func TableName(model interface{}) string {
	return Schema(model).Table
}

//Schema 返回数据表 Schema
func Schema(model interface{}) *schema.Schema{
	stmt := &gorm.Statement{DB: GetDB()}
	stmt.Parse(model)
	return  stmt.Schema
}

//Updates 更新记录
func Updates(data interface{}, where string) int64 {
	return GetDB().Where(where).Updates(data).RowsAffected
}

//Query 执行sql查询
func Query(model interface{}, sql string) int64 {
	return GetDB().Raw(sql).Scan(model).RowsAffected
}

//Exec 执行sql增删改
func Exec(sql string) int64 {
	return  GetDB().Exec(sql).RowsAffected
}

//Delete 删除数据
func Delete(model interface{}, where string) int64  {
	return GetDB().Where(where).Delete(model).RowsAffected
}

//Paginate 定义一个存储分页信息的对象
type Paginate struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
	Pages    int `json:"pages"`
	Total    int `json:"total"`
	Offset   int `json:"offset"`
}

//FindPage 获取分页数据
func FindPage(model interface{}, where string, order string, page, pageSize int) int64 {
	count := int64(0)
	GetDB().Where(where).Count(&count)
	pageObj = SetPaginate(page, pageSize, int(count))
	return GetDB().Where(where).Order(order).Offset(pageObj.Offset).Limit(pageObj.PageSize).Find(model).RowsAffected
}

//SetPaginate 初始化分页存储对象
func SetPaginate(page, pageSize, count int) *Paginate {

	//总页数
	pages := 0
	if pageSize > 0 {
		pages = int(math.Ceil(float64(count) / float64(pageSize)))
	}

	//当前页
	switch {
	case page <= 0:
		page = 1
	case pages > 0 && page >= pages:
		page = pages
	}

	//每页数量
	switch {
	case pageSize > 100:
		pageSize = 100
	case pageSize <= 0:
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	pageObj = &Paginate{
		Page:     page,
		PageSize: pageSize,
		Offset:   offset,
		Total:    count,
		Pages:    pages,
	}

	return pageObj
}

//GetPaginate 获取分页存储对象
func GetPaginate() Paginate {
	return *pageObj
}

//TransStart 开启事务
func TransStart() {
	GetDB().Begin()
}

//TransRollback 回滚事务
func TransRollback() {
	GetDB().Rollback()
}

//TransCommit 提交事务
func TransCommit() {
	GetDB().Commit()
}

