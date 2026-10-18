package db

import (
	"sync"
	"fibera/internal/config"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	dbCli *gorm.DB
	once  sync.Once
	dbErr error
)

// GetDB 获取数据库连接
func GetDB() (*gorm.DB, error) {
	once.Do(func() {
		conf, _ := config.LoadConfig()

		dbCli, dbErr = gorm.Open(mysql.Open(conf.Database.Source), &gorm.Config{
			SkipDefaultTransaction: true,
		})
		if dbErr != nil {
			zap.L().Error("failed to connect database", zap.Error(dbErr))
			return
		}

		sqlDB, err := dbCli.DB()
		if err != nil {
			zap.L().Error("failed to get sql db", zap.Error(err))
			return
		}

		// 设置连接池参数
		sqlDB.SetMaxIdleConns(conf.Database.MaxIdleConns)
		sqlDB.SetMaxOpenConns(conf.Database.MaxOpenConns)

		// 启用 SQL 日志
		dbCli = dbCli.Debug()

		// 测试连接
		if dbErr = sqlDB.Ping(); dbErr != nil {
			zap.L().Error("failed to ping database", zap.Error(dbErr))
			return
		}
	})
	return dbCli, dbErr
}
