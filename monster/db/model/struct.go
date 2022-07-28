package model

import "time"

type User struct { // new
	Id          int `xorm:"not null pk autoincr"`
	Uid         int `xorm:"unique"`
	NickName    string
	HeadIcon    int `xorm:"not null DEFAULT 1"`
	Frame       int `xorm:"not null DEFAULT 1"`
	GameCoin    int `xorm:"not null DEFAULT 0"`
	Strength    int
	Rank1       int `xorm:"not null DEFAULT 1000"`
	Rank2       int `xorm:"not null DEFAULT 1000"`
	EnergyLimit int `xorm:"not null DEFAULT 10"`
	CheckPoint  int `xorm:"not null DEFAULT 1"`
	Stage       int `xorm:"not null DEFAULT 1"`
	Award       int `xorm:"not null DEFAULT 0"`
	Matching    int `xorm:"not null DEFAULT 0"`
}

type UserEmbattle struct {
	Id             int `xorm:"not null pk autoincr"`
	UserId         int
	TeamId         int
	UserMonsterUid int
	UserMonsterId  int
	SequenceId     int
	Length         int
	Current        int `xorm:"not null DEFAULT 0"`
}

type UserProp struct {
	Id     int `xorm:"not null pk autoincr"`
	UserId int
	PropId int
	Amount int
}

type UserItem struct {
	Id        int `xorm:"not null pk autoincr"`
	UserId    int
	ItemId    int
	ItemName  string
	Type      string
	Rarity    int
	Introduce string
	Used      int `xorm:"not null DEFAULT 0"`
}

type UserEquipment struct {
	Id            int `xorm:"not null pk autoincr"`
	UserId        int
	EquipmentUid  int `xorm:"unique"`
	EquipmentId   int
	MainStatus    string
	ViceStatus    string
	UserMonsterId int `xorm:"not null DEFAULT 0"`
}

type UserRanking struct {
	Id           int `xorm:"not null pk autoincr"`
	UserId       int `xorm:"unique"`
	NickName     string
	HeadIcon     int
	MonsterIds   string
	MonsterLevel string
	Strength     int
	CreateAt     time.Time
}

type UserFriendList struct {
	Id       int `xorm:"not null pk autoincr"`
	UserId   int
	FriendId int
}

type UserBanList struct {
	Id     int `xorm:"not null pk autoincr"`
	UserId int
	BanId  int
}

type UserChatList struct { // new
	Id         int `xorm:"not null pk autoincr"`
	SenderId   int
	ReceiverId int
	Msg        string
	Read       int `xorm:"not null DEFAULT 0"`
	SendTime   string
}

type RequestFriendList struct {
	Id       int `xorm:"not null pk autoincr"`
	UserId   int
	FriendId int
}

type StageBot struct {
	Id         int `xorm:"not null pk autoincr"`
	BotId      int
	CheckPoint int
	Stage      int
}

type BotMonsterInfo struct {
	Id        int `xorm:"not null pk autoincr"`
	UserId    int
	MonsterId int
	Uid       int
	Name      string
	Rarity    int
	Element   int
	Hp        int
	Attack    int
	Defend    int
	Speed     int
	Lv        int `xorm:"not null DEFAULT 0"`
	Exp       int `xorm:"not null DEFAULT 0"`
}

type BotEmbattle struct {
	Id             int `xorm:"not null pk autoincr"`
	BotId          int
	UserMonsterUid int
	UserMonsterId  int
	SequenceId     int
}

type MonsterInfo struct {
	Id        int `xorm:"not null pk autoincr"`
	UserId    int
	MonsterId int
	Uid       int `xorm:"unique"`
	Name      string
	Rarity    int
	Element   int
	Hp        int
	Attack    int
	Defend    int
	Speed     int
	Lv        int     `xorm:"not null default 0"`
	Exp       float64 `xorm:"DECIMAL(16,4) not null default 0.0000"`
	Energy    int     `xorm:"not null DEFAULT 20"`
}

type MonsterSkill struct {
	Id        int `xorm:"not null pk autoincr"`
	MonsterId int
	Skill     int
}

type Skill struct {
	Id        int `xorm:"not null pk autoincr"`
	Skill     int
	Trigger   int
	Introduce string
}

type PropInfo struct {
	Id        int `xorm:"not null pk autoincr"`
	Name      string
	Rarity    int
	Classify  string
	Introduce string
}

type EquipmentInfo struct {
	Id        int `xorm:"not null pk autoincr"`
	UserId    int
	Name      int
	Rarity    int
	Introduce string
}

type StageAward struct {
	Id         int `xorm:"not null pk autoincr"`
	CheckPoint int
	Stage      int
	Award      int
	Exp        int
}
