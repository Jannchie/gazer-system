package gs

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Scheduler struct {
	db *gorm.DB
}

type Task struct {
	ID         uint64    `gorm:"primarykey"`
	CreatedAt  time.Time `json:"created_at" gorm:"index"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"index"`
	URL        string    `json:"url" gorm:"index:unique"`
	Tag        string    `json:"tag" gorm:"index"`
	Next       uint64    `json:"next" gorm:"index"`
	IntervalMS uint64    `json:"interval_ms"`
}

func NewScheduler(dsn string, logLevel logger.LogLevel) *Scheduler {
	return &Scheduler{
		db: initDB(dsn, logLevel, ""),
	}
}

func initDB(dsn string, logLevel logger.LogLevel, filePath string) *gorm.DB {
	var myLog logger.Writer
	if filePath != "" {
		f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}
		myLog = log.New(f, "\r\n", log.LstdFlags)
	} else {
		myLog = log.New(os.Stdout, "\r\n", log.LstdFlags)
	}
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		Logger: logger.New(
			myLog, // io writer
			logger.Config{
				SlowThreshold:             time.Second, // Slow SQL threshold
				LogLevel:                  logLevel,    // Log level
				IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
				Colorful:                  false,       // Disable color
			},
		),
	})
	if err != nil {
		panic(err)
	}

	db.Exec("PRAGMA journal_mode = WAL;")
	db.Exec("PRAGMA synchronous = NORMAL;")
	db.Exec("PRAGMA temp_store = MEMORY;")
	// db.Exec("PRAGMA foreign_keys = ON;")
	db.Exec("PRAGMA busy_timeout = 60000;")
	db.Exec("PRAGMA journal_size_limit = 134217728;") // 128M

	err = db.AutoMigrate(&Task{})
	if err != nil {
		panic(err)
	}

	return db
}
