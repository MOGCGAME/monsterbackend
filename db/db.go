package db

import (
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	
	"monster/db/model"

	"xorm.io/core"
	"xorm.io/xorm"
)

var (
	db     *xorm.Engine
	logger *log.Entry
)

type DBModule struct {
	connString   string
	maxIdleConns int
	maxOpenConns int
	showSQL      bool
}

func ExportDb() *xorm.Engine {
	return db
}

func NewDBModule() *DBModule {
	logger = log.WithField("source", "db")
	db := &DBModule{}

	if err := db.configuration(); err != nil {
		return nil
	}
	return db
}

func (d *DBModule) configuration() error {
	d.connString = viper.GetString("database.connect")
	d.maxIdleConns = viper.GetInt("database.max_idle_conns")
	d.maxOpenConns = viper.GetInt("database.max_open_conns")
	d.showSQL = viper.GetBool("database.showsql")
	return nil
}

func (d *DBModule) Init() error {
	logger.Debugf("mysql:%s", d.connString)
	database, err := xorm.NewEngine("mysql", d.connString)
	if err != nil {
		return err
	}
	db = database
	db.SetMapper(core.GonicMapper{})
	db.SetMaxIdleConns(d.maxIdleConns)
	db.SetMaxOpenConns(d.maxOpenConns)
	db.ShowSQL(d.showSQL)

	err = d.syncSchema()

	return err
}

func (d *DBModule) AfterInit() {
}

func (d *DBModule) BeforeShutdown() {
}

func (d *DBModule) Shutdown() error {
	db.Close()
	return nil
}

func (d *DBModule) syncSchema() error {
	err := db.StoreEngine("InnoDB").Sync2(
		new(model.User),
		new(model.UserEmbattle),
		new(model.UserProp),
		new(model.UserItem),
		new(model.UserEquipment),
		new(model.UserRanking),
		new(model.UserFriendList),
		new(model.UserBanList),
		new(model.UserChatList),
		new(model.RequestFriendList),
		new(model.StageBot),
		new(model.BotEmbattle),
		new(model.BotMonsterInfo),
		new(model.MonsterInfo),
		new(model.MonsterSkill),
		new(model.Skill),
		new(model.PropInfo),
		new(model.EquipmentInfo),
		new(model.StageAward),
	)

	logger.Println("re populating!!")
	return err
}
