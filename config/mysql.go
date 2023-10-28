package config

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
	"log"
	"time"
)

type DatabaseBusiness string //业务类型
type DatabaseAccess string   //访问类型

const (
	DatabaseBusinessXCamera = DatabaseBusiness("xcamera")
)

// 获取一个Mysql实例
func MysqlGet(business DatabaseBusiness) *gorm.DB {
	return localDatabaseManager.get(business, false)
}

func MysqlGetLocal(business DatabaseBusiness) *gorm.DB {
	return localDatabaseManager.get(business, true)
}

var localDatabaseManager *DatabaseManager = &DatabaseManager{
	databases: map[DatabaseBusiness]*databaseConfig{},
}

type databaseParams struct {
	Address  string `yaml:"address"`  //IP:Port
	Name     string `yaml:"name"`     //数据库名称
	User     string `yaml:"user"`     //用户名
	Password string `yaml:"password"` //密码
}

type databaseConfig struct {
	Master  *databaseParams   `yaml:"master"` //一主库
	Slaves  []*databaseParams `yaml:"slaves"` //多从库
	db      *gorm.DB          `yaml:"-"`      //db
	dbLocal *gorm.DB          `yaml:"-"`      //db-local
}

type DatabaseManager struct {
	databases map[DatabaseBusiness]*databaseConfig //数据库表
}

func (t *DatabaseManager) add(business DatabaseBusiness, dbConfig *databaseConfig) {
	t.databases[business] = dbConfig

	source := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True",
		dbConfig.Master.User, dbConfig.Master.Password,
		dbConfig.Master.Address, dbConfig.Master.Name)
	sourceLocal := source + "&loc=Local"

	replicas := make([]string, 0)
	replicasLocal := make([]string, 0)
	for _, iVal := range dbConfig.Slaves {
		replica := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True",
			iVal.User, iVal.Password, iVal.Address, iVal.Name)
		replicas = append(replicas, replica)
		replicasLocal = append(replicasLocal, replica+"&loc=Local")
	}
	dbConfig.db = getMysqlDatabase(source, replicas)
	dbConfig.dbLocal = getMysqlDatabase(sourceLocal, replicasLocal)
}

func (t *DatabaseManager) get(business DatabaseBusiness, isLocal bool) *gorm.DB {
	dbConfig, ok := t.databases[business]
	if ok == false {
		return nil
	}
	if isLocal {
		return dbConfig.dbLocal
	}
	return dbConfig.db
}

func getMysqlDatabase(source string, replicas []string) *gorm.DB {
	logSet := logger.Default.LogMode(logger.Silent)
	if Instance.EnvParams.GORMLogLevel > int(logger.Silent) &&
		Instance.EnvParams.GORMLogLevel <= int(logger.Info) {
		logSet = logger.Default.LogMode(logger.LogLevel(Instance.EnvParams.GORMLogLevel))
	}
	db, err := gorm.Open(mysql.Open(source), &gorm.Config{
		Logger:                 logSet,
		SkipDefaultTransaction: true,
	})
	if err != nil {
		log.Fatalln(err)
		return nil
	}

	dialects := make([]gorm.Dialector, 0)
	for _, iVal := range replicas {
		dialects = append(dialects, mysql.Open(iVal))
	}
	plugin := dbresolver.Register(dbresolver.Config{
		Replicas: dialects,
		Policy:   dbresolver.RandomPolicy{},
	}).SetMaxIdleConns(100).
		SetMaxOpenConns(1000).
		SetConnMaxIdleTime(time.Minute).
		SetConnMaxLifetime(time.Hour)

	if err := db.Use(plugin); err != nil {
		log.Fatalln(err)
		return nil
	}
	if sqlDB, err := db.DB(); err != nil {
		log.Fatalln(err)
		return nil
	} else if err := sqlDB.Ping(); err != nil {
		log.Fatalf("Mysql ping %s, err: %+v", source, err)
	}

	return db
}
