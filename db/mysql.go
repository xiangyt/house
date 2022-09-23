package db

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/xiangyt/house/config"
	"github.com/xiangyt/house/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"sync"
	"time"
)

type Mysql struct {
	db *gorm.DB
}

var (
	instance *Mysql
	once     sync.Once
)

func GetDB() *gorm.DB {
	return GetInstance().db
}

func GetInstance() *Mysql {
	once.Do(func() {
		instance = &Mysql{}
	})
	return instance
}

func (m *Mysql) Init() error {
	log := logger.Default
	if gin.IsDebugging() {
		log = logger.Default.LogMode(logger.Silent)
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/house?charset=utf8mb4&parseTime=True&loc=Local",
		config.Get().Mysql.User,
		config.Get().Mysql.Password,
		config.Get().Mysql.Host,
		config.Get().Mysql.Port,
	)
	var datetimePrecision = 2
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,                // data source name, refer https://github.com/go-sql-driver/mysql#dsn-data-source-name
		DefaultStringSize:         256,                // add default size for string fields, by default, will use db type `longtext` for fields without size, not a primary key, no index defined and don't have default values
		DisableDatetimePrecision:  true,               // disable datetime precision support, which not supported before MySQL 5.6
		DefaultDatetimePrecision:  &datetimePrecision, // default datetime precision
		DontSupportRenameIndex:    true,               // drop & create index when rename index, rename index not supported before MySQL 5.7, MariaDB
		DontSupportRenameColumn:   true,               // use change when rename column, rename rename not supported before MySQL 8, MariaDB
		SkipInitializeWithVersion: false,              // smart configure based on used version
	}), &gorm.Config{
		Logger: log,
	})
	if err != nil {
		return errors.New("init mysql failed")
	}
	m.db = db
	logrus.Info("init mysql success!")

	err = m.db.AutoMigrate(&model.House{})
	if err != nil {
		return errors.New(fmt.Sprintf("AutoMigrate failed, err:%s", err))
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Hour)
	return nil
}
