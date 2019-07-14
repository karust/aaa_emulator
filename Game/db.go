package main

var sqliteSchema = `
`

var postgesSchema = `
DROP TABLE IF EXISTS abilities;
CREATE TABLE IF NOT EXISTS abilities (
  id BIGSERIAL NOT NULL,
  exp INT NOT NULL,
  owner BIGINT NOT NULL,
  PRIMARY KEY (id,owner)
);

DROP TABLE IF EXISTS actabilities;
CREATE TABLE IF NOT EXISTS actabilities (
  id SERIAL NOT NULL,
  point INT NOT NULL DEFAULT '0',
  step SMALLINT NOT NULL DEFAULT '0',
  owner INT NOT NULL,
  PRIMARY KEY (owner,id)
);

DROP TABLE IF EXISTS appellations;
CREATE TABLE IF NOT EXISTS appellations (
  id BIGSERIAL NOT NULL,
  active SMALLINT NOT NULL DEFAULT '0',
  owner INT NOT NULL,
  PRIMARY KEY (id,owner)
);

DROP TABLE IF EXISTS blocked;
CREATE TABLE IF NOT EXISTS blocked (
  owner INT NOT NULL,
  blocked_id INT NOT NULL,
  PRIMARY KEY (owner,blocked_id)
);

DROP TABLE IF EXISTS characters;
CREATE TABLE IF NOT EXISTS characters (
  id BIGSERIAL NOT NULL,
  account_id BIGINT NOT NULL,
  name varchar(128) NOT NULL,
  race SMALLINT NOT NULL,
  gender SMALLINT NOT NULL,
  unit_model_params BYTEA NOT NULL,
  level SMALLINT NOT NULL,
  expirience INT NOT NULL,
  recoverable_exp INT NOT NULL,
  hp INT NOT NULL,
  mp INT NOT NULL,
  labor_power SMALLINT NOT NULL,
  labor_power_modified DATE NOT NULL DEFAULT '0001-01-01 00:00:00',
  consumed_lp INT NOT NULL,
  ability1 SMALLINT NOT NULL,
  ability2 SMALLINT NOT NULL,
  ability3 SMALLINT NOT NULL,
  world_id INT NOT NULL,
  zone_id INT NOT NULL,
  x REAL NOT NULL,
  y REAL NOT NULL,
  z REAL NOT NULL,
  rotation_x SMALLINT NOT NULL,
  rotation_y SMALLINT NOT NULL,
  rotation_z SMALLINT NOT NULL,
  faction_id INT NOT NULL,
  faction_name varchar(128) NOT NULL,
  expedition_id INT NOT NULL,
  family INT NOT NULL,
  dead_count SMALLINT NOT NULL,
  dead_time DATE NOT NULL DEFAULT '0001-01-01 00:00:00',
  rez_wait_duration INT NOT NULL,
  rez_time DATE NOT NULL DEFAULT '0001-01-01 00:00:00',
  rez_penalty_duration INT NOT NULL,
  leave_time DATE NOT NULL DEFAULT '0001-01-01 00:00:00',
  money BIGINT NOT NULL,
  money2 BIGINT NOT NULL,
  honor_point INT NOT NULL,
  vocation_point INT NOT NULL,
  crime_point INT NOT NULL,
  crime_record INT NOT NULL,
  delete_request_time DATE NOT NULL DEFAULT '0001-01-01 00:00:00',
  transfer_request_time DATE NOT NULL DEFAULT '0001-01-01 00:00:00',
  delete_time DATE NOT NULL DEFAULT '0001-01-01 00:00:00',
  bm_point INT NOT NULL,
  auto_use_aapoint SMALLINT NOT NULL,
  prev_point INT NOT NULL,
  point INT NOT NULL,
  gift INT NOT NULL,
  num_inv_slot SMALLINT NOT NULL DEFAULT 50,
  num_bank_slot SMALLINT NOT NULL DEFAULT 50,
  expanded_expert SMALLINT NOT NULL,
  slots BYTEA NOT NULL,
  created_at DATE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATE NOT NULL DEFAULT '0001-01-01 00:00:00',
  PRIMARY KEY (id,account_id)
);

DROP TABLE IF EXISTS completed_quests;
CREATE TABLE IF NOT EXISTS completed_quests (
  id BIGSERIAL NOT NULL,
  data BYTEA NOT NULL,
  owner BIGINT NOT NULL,
  PRIMARY KEY (id,owner)
);

DROP TABLE IF EXISTS expeditions;
CREATE TABLE IF NOT EXISTS expeditions (
  id BIGSERIAL NOT NULL,
  owner BIGINT NOT NULL,
  owner_name varchar(128) NOT NULL,
  name varchar(128) NOT NULL,
  mother INT NOT NULL,
  created_at DATE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id,owner)
);

DROP TABLE IF EXISTS family_members;
CREATE TABLE IF NOT EXISTS family_members (
  character_id BIGSERIAL NOT NULL,
  family_id INT NOT NULL,
  name varchar(45) NOT NULL,
  role SMALLINT NOT NULL DEFAULT '0',
  title varchar(45) DEFAULT NULL,
  PRIMARY KEY (family_id,character_id)
);

DROP TABLE IF EXISTS friends;
CREATE TABLE IF NOT EXISTS friends (
  id BIGSERIAL NOT NULL,
  friend_id BIGINT NOT NULL,
  owner  BIGINT NOT NULL,
  PRIMARY KEY (id,owner)
);

DROP TABLE IF EXISTS housings;
CREATE TABLE IF NOT EXISTS housings (
  id BIGSERIAL NOT NULL,
  account_id BIGINT NOT NULL,
  owner BIGINT NOT NULL,
  template_id BIGINT NOT NULL,
  x REAL NOT NULL,
  y REAL NOT NULL,
  z REAL NOT NULL,
  rotation_z SMALLINT NOT NULL,
  current_step SMALLINT NOT NULL,
  permission SMALLINT NOT NULL,
  PRIMARY KEY (account_id,owner,id)
);

DROP TABLE IF EXISTS items;
DROP TYPE IF EXISTS ST;
CREATE TYPE ST AS ENUM ('Equipment','Inventory','Bank');
CREATE TABLE IF NOT EXISTS items (
  id BIGSERIAL NOT NULL,
  type varchar(100) NOT NULL,
  template_id BIGINT NOT NULL,
  slot_type ST NOT NULL,
  slot INT NOT NULL,
  count INT NOT NULL,
  details BYTEA NULL,
  lifespan_mins INT NOT NULL,
  made_unit_id BIGINT NOT NULL DEFAULT '0',
  unsecure_time DATE NOT NULL DEFAULT '0001-01-01 00:00:00',
  unpack_time DATE NOT NULL DEFAULT '0001-01-01 00:00:00',
  owner BIGINT NOT NULL,
  grade SMALLINT DEFAULT '0',
  created_at DATE NOT NULL DEFAULT '0001-01-01 00:00:00',
  PRIMARY KEY (id)
);

DROP TABLE IF EXISTS mates;
CREATE TABLE IF NOT EXISTS mates (
  id BIGSERIAL NOT NULL,
  item_id BIGINT NOT NULL,
  name text NOT NULL,
  xp INT NOT NULL,
  level SMALLINT NOT NULL,
  mileage INT NOT NULL,
  hp INT NOT NULL,
  mp INT NOT NULL,
  owner BIGINT NOT NULL,
  updated_at DATE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  created_at DATE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id,item_id,owner)
);

DROP TABLE IF EXISTS options;
CREATE TABLE IF NOT EXISTS options (
  key varchar(100) NOT NULL,
  value text NOT NULL,
  owner BIGINT NOT NULL,
  PRIMARY KEY (key,owner)
);

DROP TABLE IF EXISTS portal_book_coords;
CREATE TABLE IF NOT EXISTS portal_book_coords (
  id BIGSERIAL NOT NULL,
  name varchar(128) NOT NULL,
  x INT DEFAULT '0',
  y INT DEFAULT '0',
  z INT DEFAULT '0',
  zone_id INT DEFAULT '0',
  z_rot INT DEFAULT '0',
  sub_zone_id INT DEFAULT '0',
  owner INT NOT NULL,
  PRIMARY KEY (id,owner)
);

DROP TABLE IF EXISTS portal_visited_districts;
CREATE TABLE IF NOT EXISTS portal_visited_districts (
  id int NOT NULL,
  subzone int NOT NULL,
  owner int NOT NULL,
  PRIMARY KEY (id,subzone,owner)
);

DROP TABLE IF EXISTS quests;
CREATE TABLE IF NOT EXISTS quests (
  id BIGSERIAL NOT NULL,
  template_id BIGINT NOT NULL,
  data BYTEA NOT NULL,
  status SMALLINT NOT NULL,
  owner BIGINT NOT NULL,
  PRIMARY KEY (id,owner)
);

DROP TABLE IF EXISTS skills;
DROP TYPE IF EXISTS T;
CREATE TYPE T AS ENUM ('Skill','Buff');
CREATE TABLE IF NOT EXISTS skills (
  id BIGSERIAL NOT NULL,
  level SMALLINT NOT NULL,
  type T NOT NULL,
  owner BIGINT NOT NULL,
  PRIMARY KEY (id,owner)
);

DROP TABLE IF EXISTS cash_shop_items;
CREATE TABLE cash_shop_items  (
  id BIGSERIAL NOT NULL,
  uniq_id BIGINT NULL DEFAULT 0 ,
  cash_name varchar(255) NOT NULL ,
  main_tab SMALLINT CHECK (main_tab > 0) NULL DEFAULT 1 ,
  sub_tab SMALLINT CHECK (sub_tab > 0) NULL DEFAULT 1 ,
  level_min SMALLINT CHECK (level_min > 0) NULL DEFAULT 0 ,
  level_max SMALLINT CHECK (level_max > 0) NULL DEFAULT 0 ,
  item_template_id int CHECK (item_template_id > 0) NULL DEFAULT 0 ,
  is_sell SMALLINT CHECK (is_sell > 0) NULL DEFAULT 0 ,
  is_hidden SMALLINT CHECK (is_hidden > 0) NULL DEFAULT 0 ,
  limit_type SMALLINT CHECK (limit_type > 0) NULL DEFAULT 0,
  buy_count SMALLINT CHECK (buy_count > 0) NULL DEFAULT 0,
  buy_type SMALLINT CHECK (buy_type > 0) NULL DEFAULT 0,
  buy_id int CHECK (buy_id > 0) NULL DEFAULT 0,
  start_date timestamp(0) NULL DEFAULT '0001-01-01 00:00:00' ,
  end_date timestamp(0) NULL DEFAULT '0001-01-01 00:00:00' ,
  type SMALLINT CHECK (type > 0) NULL DEFAULT 0 ,
  price int CHECK (price > 0) NULL DEFAULT 0 ,
  remain int CHECK (remain > 0) NULL DEFAULT 0 ,
  bonus_type int CHECK (bonus_type > 0) NULL DEFAULT 0 ,
  bouns_count int CHECK (bouns_count > 0) NULL DEFAULT 0 ,
  cmd_ui SMALLINT CHECK (cmd_ui > 0) NULL DEFAULT 0 ,
  item_count int CHECK (item_count > 0) NULL DEFAULT 1 ,
  select_type SMALLINT CHECK (select_type > 0) NULL DEFAULT 0,
  default_flag SMALLINT CHECK (default_flag > 0) NULL DEFAULT 0,
  event_type SMALLINT CHECK (event_type > 0) NULL DEFAULT 0 ,
  event_date timestamp(0) NULL DEFAULT '0001-01-01 00:00:00' ,
  dis_price int CHECK (dis_price > 0) NULL DEFAULT 0 ,
  PRIMARY KEY (id)
); 

SET CLIENT_ENCODING TO 'UTF-8';
COMMIT;`

// Ability ... Ability table
type Ability struct {
	ID    uint
	Exp   uint
	Owner uint
}

// Actabilitiy ... Actabilitiy table
type Actabilitiy struct {
	ID    uint
	Point uint
	Step  uint16
	Owner uint
}

// Appellation ... Appellation table
type Appellation struct {
	ID     uint
	Active uint8
	Owner  uint
}

// Blocked ... Blocked table
type Blocked struct {
	Owner     int
	BlockedID int `db:"blocked_id"`
}

// Character ... Character table
type Character struct {
	ID                  uint
	AccountID           uint `db:"account_id"`
	Name                string
	Race                uint8
	Gender              uint8
	UnitModelParams     []byte `db:"unit_model_params"`
	Level               uint8
	Expirience          uint
	RecoverableEXP      uint `db:"recoverable_exp"`
	HP                  int
	MP                  int
	LaborPower          uint8  `db:"labor_power"`
	LaborPowerModifed   string `db:"labor_power_modified"`
	ConsumedLP          uint   `db:"consumed_lp"`
	Ability1            uint8
	Ability2            uint8
	Ability3            uint8
	WorldID             int `db:"world_id"`
	ZoneID              int `db:"zone_id"`
	X                   float32
	Y                   float32
	Z                   float32
	RotationX           uint8  `db:"rotation_x"`
	RotationY           uint8  `db:"rotation_y"`
	RotationZ           uint8  `db:"rotation_z"`
	FactionID           uint   `db:"faction_id"`
	FactionName         string `db:"faction_name"`
	ExpeditionID        uint   `db:"expedition_id"`
	Family              uint
	DeadCount           uint8  `db:"dead_count"`
	DeadTime            string `db:"dead_time"`
	RezWaitDuration     uint   `db:"rez_wait_duration"`
	RezTime             string `db:"rez_time"`
	RezPenaltyDuration  uint   `db:"rez_penalty_duration"`
	LeaveTime           string `db:"leave_time"`
	Money               uint
	Money2              uint
	HonorPoint          uint   `db:"honor_point"`
	VocationPoint       uint   `db:"vocation_point"`
	CrimePoint          uint   `db:"crime_point"`
	CrimeRecord         uint   `db:"crime_record"`
	DeleteRequesTime    string `db:"delete_request_time"`
	TransferRequestTime string `db:"transfer_request_time"`
	DeleteTime          string `db:"delete_time"`
	BmPoint             uint   `db:"bm_point"`
	AutoUseAApoint      uint8  `db:"auto_use_aapoint"`
	PrevPoint           int    `db:"prev_point"`
	Point               int
	Gift                int
	NumInvSlot          int8 `db:"num_inv_slot"`
	NumBankSlot         int8 `db:"num_bank_slot"`
	ExpandedExpert      int8 `db:"expanded_expert"`
	Slots               []byte
	CreatedAt           string `db:"created_at"`
	UpdatedAt           string `db:"updated_at"`
}

// CompletedQuests ... CompletedQuests table
type CompletedQuests struct {
	ID    uint
	Data  []byte
	Owner uint
}

// Expedition  ... Expedition table
type Expedition struct {
	ID        uint
	Owner     uint
	OwnerName string `db:"owner_name"`
	Name      string
	Mother    int
	CreatedAt string `db:"created_at"`
}

// FamilyMember  ... FamilyMember table
type FamilyMember struct {
	CharacterID uint `db:"character_id"`
	FamilyID    uint `db:"family_id"`
	Name        string
	Role        uint8
	Title       string
}

// Friend  ... Friend table
type Friend struct {
	CharacterID uint `db:"character_id"`
	ID          uint
	FriendID    uint `db:"friend_id"`
	Owner       uint
}

// Hounsing  ... Hounsing table
type Hounsing struct {
	ID          uint
	AccountID   uint `db:"account_id"`
	Owner       uint
	TemplateID  uint `db:"template_id"`
	X           float32
	Y           float32
	Z           float32
	RotationZ   uint8 `db:"rotation_z"`
	CurrentStep uint8 `db:"current_step"`
	Permission  uint8
}

// Item  ... Item table
type Item struct {
	ID           uint
	Type         string
	TemplateID   uint   `db:"template_id"`
	SlotType     string `db:"slot_type"`
	Slot         int
	Count        int
	Details      []byte
	LifespanMins int    `db:"lifespan_mins"`
	MadeUnitID   uint   `db:"made_unit_id"`
	UnsecureTime string `db:"unsecure_time"`
	UnpackTime   string `db:"unpack_time"`
	Owner        uint
	Grade        uint8
	CreatedAt    string `db:"created_at"`
}

// Mate  ... Mate table
type Mate struct {
	ID        uint
	ItemID    uint `db:"item_id"`
	Name      string
	XP        int
	Level     uint8
	Mileage   int
	HP        int
	MP        int
	Owner     uint
	UpdatedAt string `db:"updated_at"`
	CreatedAt string `db:"created_at"`
}

// Option  ... Option table
type Option struct {
	Key   string
	Value string
	Owner uint
}

// PortalBookCoord   ... PortalBookCoord table
type PortalBookCoord struct {
	ID        uint
	Name      string
	X         int
	Y         int
	Z         int
	ZoneID    int `db:"zone_id"`
	ZRot      int `db:"z_rot"`
	SubZoneID int `db:"sub_zone_id"`
	Owner     int
}

// PortalVisitedDistrict   ... PortalVisitedDistrict table
type PortalVisitedDistrict struct {
	ID      uint
	Subzone int
	Owner   uint
}

// Quest  ... Quest table
type Quest struct {
	ID         uint
	TemplateID uint `db:"template_id"`
	Data       []byte
	Status     uint8
	Owner      uint
}

// Skill  ... Skill table
type Skill struct {
	ID    uint
	Level uint8
	Type  string
	Owner uint
}

// CashShopItem  ... CashShopItem table
type CashShopItem struct {
	ID             uint
	UniqID         uint   `db:"uniq_id"`
	CashName       string `db:"cash_name"`
	MainTab        uint8  `db:"main_tab"`
	SubTab         uint8  `db:"sub_tab"`
	LevelMin       uint8  `db:"level_min"`
	LevelMax       uint8  `db:"level_max"`
	ItemTemplateID int    `db:"item_template_id"`
	IsSell         uint8  `db:"uniq_id"`
	IsHidden       uint8  `db:"is_hidden"`
	LimitType      uint8  `db:"limit_type"`
	BuyCount       uint8  `db:"buy_count"`
	BuyType        uint8  `db:"buy_type"`
	BuyID          int    `db:"buy_id"`
	StartDate      uint   `db:"start_date"`
	EndDate        uint   `db:"end_date"`
	Type           uint8
	Price          int
	Remain         int
	BonusType      int   `db:"bonus_type"`
	BounsCount     int   `db:"bouns_count"`
	CmdUI          uint8 `db:"cmd_ui"`
	ItemCount      int   `db:"item_count"`
	SelectType     uint8 `db:"select_type"`
	DefaultFlag    uint8 `db:"default_flag"`
	EventType      uint8 `db:"event_type"`
	EventDate      uint  `db:"event_date"`
	DisPrice       int   `db:"dis_price"`
}
