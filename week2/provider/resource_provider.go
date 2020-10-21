package provider

import (
	"golang-training/week2/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type RP interface {
	GetDB() *gorm.DB
}

type implDB struct {
	db *gorm.DB
}

func MustBuildResourceProvider(cfg *config.Config) RP {
	dsn := cfg.MySQLDSN
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// Get generic database object sql.DB to use its functions
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	// sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(10)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	// sqlDB.SetConnMaxLifetime(time.Hour)

	return &implDB{
		db: db,
	}
}

func (iDB *implDB) GetDB() *gorm.DB {
	return iDB.db
}
