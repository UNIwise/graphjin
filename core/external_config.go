package core

import (
	"errors"
	"fmt"
	"github.com/dosco/graphjin/core/internal/allow"
	"github.com/dosco/graphjin/core/internal/qcode"
	"github.com/robfig/cron/v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strconv"
	"sync/atomic"
	"time"
)

type ExternalConfigRequest struct {
	ServiceName string                   `json:"serviceName"`
	Queries     map[string]InternalQuery `json:"queries"`
	Scripts     map[string]string        `json:"scripts"`
	Fragments   map[string]string        `json:"fragments"`
}

type InternalQuery struct {
	Query string `json:"query"`
	Vars  string `json:"vars"`
}

type DBConfig struct {
	Host        string
	Port        uint16
	User        string
	Password    string
	Name        string
	CronPattern string
}

type ExternalConfig struct {
	connString   string
	dbConfig     *DBConfig
	gj           *graphjin
	lastUpdated  int64
	dbConnection *gorm.DB
}

type GJQuery struct {
	gorm.Model
	Service string `gorm:"uniqueIndex:svc_name;index;not null,size:255"`
	Name    string `gorm:"uniqueIndex:svc_name;index;not null,size:255"`
	Query   string `binding:"required"`
	Vars    string `binding:"required"`
}

func (GJQuery) TableName() string {
	return "query"
}

type GJScript struct {
	gorm.Model
	Service string `gorm:"uniqueIndex:svc_name;index;not null,size:255"`
	Name    string `gorm:"uniqueIndex:svc_name;index;not null,size:255"`
	Script  string `binding:"required"`
}

func (GJScript) TableName() string {
	return "script"
}

type GJFragment struct {
	gorm.Model
	Service  string `gorm:"uniqueIndex:svc_name;index;not null,size:255"`
	Name     string `gorm:"uniqueIndex:svc_name;index;not null,size:255"`
	Fragment string `binding:"required"`
}

func (GJFragment) TableName() string {
	return "fragment"
}

func (ec *ExternalConfig) getConnString() string {
	c := ec.dbConfig
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True", c.User, c.Password, c.Host, c.Port, c.Name)
}

func (ec *ExternalConfig) getDBConnection() (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(ec.getConnString()), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (ec *ExternalConfig) Load() error {
	now := time.Now().Unix()
	diff := now - atomic.LoadInt64(&ec.lastUpdated)
	if diff < 10 {
		return errors.New("aborted update due to throttling")
	}

	atomic.StoreInt64(&ec.lastUpdated, now)
	db := ec.dbConnection

	gj := ec.gj
	var scripts []GJScript
	if err := db.Find(&scripts).Error; err != nil {
		return err
	}
	for _, s := range scripts {
		err := gj.loadScriptFromString(s.Name, s.Service, s.Script)
		if err != nil {
			gj.log.Printf("External Config: Failed to add script. Error: %s", err)
			continue
		}
	}

	if gj.conf.DisableAllowList {
		return nil
	}

	var queries []GJQuery
	if err := db.Find(&queries).Error; err != nil {
		return err
	}
	for _, q := range queries {
		if q.Query == "" {
			continue
		}

		qt, _ := qcode.GetQType(q.Query)
		item := allow.Item{
			Name:     q.Name,
			Comment:  "",
			Query:    q.Query,
			Vars:     q.Vars,
			Metadata: allow.Metadata{},
		}
		qk := gj.getQueryKeys(item)

		for _, v := range qk {
			qc := &queryComp{qr: queryReq{
				op:    qt,
				name:  item.Name,
				query: []byte(item.Query),
				vars:  []byte(item.Vars),
			}}

			if item.Metadata.Order.Var != "" {
				qc.qr.order = [2]string{item.Metadata.Order.Var, strconv.Quote(v.val)}
			}
			gj.queries[v.key] = qc
		}

		op, _ := qcode.GetQType(item.Query)
		gj.apq.Set(item.Name, apqInfo{op: op, name: item.Name})
	}

	return nil
}

func (ec *ExternalConfig) LoadFragment(name string) (string, error) {
	db := ec.dbConnection
	var fragment GJFragment
	err := db.Order("created_at DESC").First(&fragment, "name = ?", name).Error
	if err != nil {
		return "", err
	}

	return fragment.Fragment, nil
}

func (ec *ExternalConfig) Store(c ExternalConfigRequest) error {
	db := ec.dbConnection
	var queries []GJQuery
	for queryName, q := range c.Queries {
		queries = append(queries, GJQuery{
			Service: c.ServiceName,
			Name:    queryName,
			Query:   q.Query,
			Vars:    q.Vars,
		})
	}
	if len(queries) > 0 {
		db.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(&queries)
	}

	var scripts []GJScript
	for scriptName, s := range c.Scripts {
		scripts = append(scripts, GJScript{
			Service: c.ServiceName,
			Name:    scriptName,
			Script:  s,
		})
	}
	if len(scripts) > 0 {
		db.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(scripts)
	}

	var fragments []GJFragment
	for fragmentName, f := range c.Fragments {
		fragments = append(fragments, GJFragment{
			Service:  c.ServiceName,
			Name:     fragmentName,
			Fragment: f,
		})
	}
	if len(fragments) > 0 {
		db.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(fragments)
	}

	return nil
}

func (gj *graphjin) initExternalConfig(c *DBConfig) error {
	ec := ExternalConfig{
		dbConfig: c,
		gj:       gj,
	}

	db, err := ec.getDBConnection()
	if err != nil {
		gj.log.Printf("External Config: Failed to connect to database. Error: %s", err)
		return err
	}

	ec.dbConnection = db

	err = db.AutoMigrate(&GJQuery{}, &GJScript{}, &GJFragment{})
	if err != nil {
		gj.log.Printf("External Config: Failed to auto migrate. Error: %s", err)
		return err
	}

	gj.externalConfig = &ec

	if c.CronPattern != "" {
		cronScheduler := cron.New()
		_, err = cronScheduler.AddFunc(c.CronPattern, func() {
			err := ec.Load()
			if err != nil {
				gj.log.Printf("External Config: Failed to load config using cron. Error: %s", err)
				return
			}
		})
		if err != nil {
			gj.log.Printf("External Config: Failed to add cron job. Error: %s", err)
			return err
		}

		// Start the cron scheduler
		cronScheduler.Start()
	}

	gj.log.Println("External Config initialized")

	return nil
}
