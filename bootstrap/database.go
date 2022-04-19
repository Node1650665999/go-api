package bootstrap

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gin-api/application/http/model"
	"gin-api/pkg/config"
)

//setupDB 初始化数据库链接
func setupDB() *gorm.DB {
	db, err := gorm.Open(dialector(), &gorm.Config{
		Logger: model.Logger(),
	})
	if err != nil {
		panic(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	//registerCalllback(db)
	model.SetDB(db)
	return db
}

func dialector() gorm.Dialector {
	//TODO 适配驱动
	driver   := config.GetString("database.connection")
	return mysql.Open(mysqlDns(driver))
}

//dns 返回数据库的连接信息
func mysqlDns(driver string) string {
	username := config.GetString("database." + driver + ".username")
	password := config.GetString("database." + driver + ".password")
	host     := config.GetString("database." + driver + ".host")
	port     := config.GetInt("database." + driver + ".port")
	database := config.GetString("database." + driver + ".database")

	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		username,
		password,
		host,
		port,
		database,
	)
}

//registerCalllback 注册回调函数 refer : https://segmentfault.com/a/1190000039070187
func registerCalllback(db *gorm.DB)  {
	// 监听delete方法
	db.Callback().Delete().Register("gorm:deleteSql", callback)
	// 监听查询
	db.Callback().Query().Register("gorm:querySql", callback)
	// 监听update方法
	db.Callback().Update().Register("gorm:updateSql", callback)
	// 监听create方法
	db.Callback().Create().Register("gorm:insertSql", callback)
	// 监听row 方法
	db.Callback().Row().Register("gorm:row", callback)
	// 监听raw 方法
	db.Callback().Raw().Register("gorm:raw", callback)
}

func callback(db *gorm.DB) {
	//TODO
}
