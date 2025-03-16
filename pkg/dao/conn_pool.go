package dao

import (
	"github.com/MoWan-inc/aqua/pkg/config"
	aquadao "github.com/MoWan-inc/aqua/pkg/dao/gorm"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"reflect"
	"time"
)

func newDB(config *mysql.Config, option *config.ConnectionOption, customLog logger.Interface) (*gorm.DB, error) {
	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name
	// dsn := "user:password@tcp(localhost:5555)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	var newLogger logger.Interface
	if customLog != nil && !reflect.ValueOf(customLog).IsZero() {
		newLogger = customLog
	} else {
		newLogger = logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold: time.Second, // 慢sql阈值
				LogLevel:      logger.Info, // 日志级别
				Colorful:      false,       // 是否彩色打印
			},
		)
	}

	db, err := gorm.Open(mysql.New(*config), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	if option != nil {
		// 设置连接池连接最大数量
		sqlDB.SetMaxIdleConns(option.MaxIdleConns)
		// 设置数据库最大打开连接数
		sqlDB.SetMaxOpenConns(option.MaxOpenConns)
		// 设置连接最大可复用时间
		sqlDB.SetConnMaxLifetime(option.ConnMaxLifeTime)
	}
	return db, nil
}

func NewDAO(cfg config.MysqlConfig) aquadao.DAO {
	return NewBaseDAO(cfg)
}

func NewBaseDAO(cfg config.MysqlConfig) *aquadao.BaseDAO {
	// todo: 从配置文件中读取配置初始化logger
	db, err := newDB(&mysql.Config{DSN: cfg.DSN}, cfg.ConnOption, nil)
	if err != nil {
		panic(err)
	}
	return aquadao.NewBaseDAO(db)
}
