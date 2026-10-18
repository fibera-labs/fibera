package db

import (
	"go.uber.org/zap"
)

// SyncDatabase 同步数据库，遍历已注册的数据表模型
func SyncDatabase() error {
	db, err := GetDB()
	if err != nil {
		zap.L().Error("failed to get db", zap.Error(err))
		return err
	}

	for _, m := range Tables() {
		err = db.AutoMigrate(m)
		if err != nil {
			zap.L().Error("failed to sync table", zap.Error(err), zap.String("table", m.TableName()))
			return err
		}
		zap.L().Info("table sync successfully", zap.String("table", m.TableName()))
	}
	return nil
}
