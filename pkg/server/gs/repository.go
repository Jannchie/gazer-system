package gs

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/jannchie/gazer-system/pkg/server/variables"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

type Repository struct {
	db                *gorm.DB
	taskUpdateChannel chan Task
	rawConsumeChannel chan uint64
}

func NewRepository(dsn string, logLevel logger.LogLevel) *Repository {
	db := initDB(dsn, logLevel)
	r := &Repository{db: db, taskUpdateChannel: make(chan Task, 128), rawConsumeChannel: make(chan uint64, 128)}
	go r.updateTask()
	go r.recoverRaw()
	return r
}

func initDB(dsn string, logLevel logger.LogLevel) *gorm.DB {
	var myLog logger.Writer
	filePath := ""
	if filePath != "" {
		f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}
		myLog = log.New(f, "\r\n", log.LstdFlags)
	} else {
		myLog = log.New(os.Stdout, "\r\n", log.LstdFlags)
	}
	config := &gorm.Config{
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
	}
	if *variables.DB == "sqlite" {
		return initSqlite(dsn, config)
	} else {
		return initPostgres(dsn, config)
	}
}

func initPostgres(dsn string, config *gorm.Config) *gorm.DB {
	ticker := time.NewTicker(time.Second * 5)
	for range ticker.C {
		db, err := gorm.Open(postgres.Open(dsn), config)
		if err != nil {
			log.Println(err)
		}
		err = db.AutoMigrate(&Task{}, &Raw{})
		if err != nil {
			log.Println(err)
		}
		return db
	}
	return nil
}

func initSqlite(dsn string, config *gorm.Config) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(dsn), config)
	if err != nil {
		panic(err)
	}
	db.Exec("PRAGMA journal_mode = WAL;")
	db.Exec("PRAGMA synchronous = NORMAL;")
	db.Exec("PRAGMA temp_store = MEMORY;")
	db.Exec("PRAGMA busy_timeout = 60000;")
	db.Exec("PRAGMA auto_vacuum = incremental;")
	db.Exec("PRAGMA journal_size_limit = 134217728;")
	db.Exec("PRAGMA mmap_size = 134217728;")

	err = db.AutoMigrate(&Task{}, &Raw{})
	if err != nil {
		panic(err)
	}
	return db
}
func (r *Repository) updateTask() {
	ticker := time.NewTicker(time.Second * 5)
	deleteList := make([]uint64, 0, 100)
	updateList := make([]Task, 0, 100)
	consumeList := make([]uint64, 0, 100)
	for {
		select {
		case <-ticker.C:
			go func() {
				if len(deleteList) != 0 {
					ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
					r.db.WithContext(ctx).Unscoped().Where("id IN ?", deleteList).Delete(&Task{})
					deleteList = deleteList[:0]
					cancel()
				}
			}()
			go func() {
				if len(updateList) != 0 {
					for _, task := range updateList {
						ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
						err := r.db.WithContext(ctx).Model(&Task{}).Where("id = ?", task.ID).Update("next", time.Now().UTC().Add(time.Millisecond*time.Duration(task.IntervalMS)).Unix()).Error
						if err != nil {
							log.Println(err)
							time.Sleep(time.Second)
						}
						cancel()
					}
					updateList = updateList[:0]
				}
			}()
			go func() {
				if len(consumeList) != 0 {
					ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
					if err := r.db.WithContext(ctx).Unscoped().Where("id IN ?", consumeList).Delete(&Raw{}).Error; err != nil {
						log.Println(err)
						time.Sleep(time.Second)
					}
					cancel()
					consumeList = consumeList[:0]
				}
			}()
		case task := <-r.taskUpdateChannel:
			if task.IntervalMS == 0 {
				deleteList = append(deleteList, task.ID)
			} else {
				updateList = append(updateList, task)
			}
		case raw := <-r.rawConsumeChannel:
			consumeList = append(consumeList, raw)
		}
	}
}
func (r *Repository) AddTasks(ctx context.Context, tasks []Task) (uint64, error) {
	var err error
	var count uint64
	nowTS := uint64(time.Now().Unix())
	for _, newTask := range tasks {
		var oriTask Task
		if err = r.db.WithContext(ctx).Find(&oriTask, "url = ?", newTask.URL).Error; err != nil {
			return 0, err
		} else {
			if oriTask.ID == 0 {
				// No previous
				if res := r.db.WithContext(ctx).Clauses(clause.OnConflict{DoNothing: true}).Create(&newTask); res.Error != nil {
					log.Println(res.Error)
					time.Sleep(time.Second)
				} else {
					count += uint64(res.RowsAffected)
				}
			} else if newTask.IntervalMS != 0 && oriTask.IntervalMS != 0 {
				// 有记录，且记录是非立即执行请求，且新的记录是非立即执行请求，则更新请求间隔
				if res := r.db.WithContext(ctx).Model(&newTask).Where("id = ?", oriTask.ID).Update("interval_ms", newTask.IntervalMS); res.Error != nil {
					log.Println(err)
					time.Sleep(time.Second)
				}
			} else if newTask.IntervalMS == 0 && oriTask.IntervalMS != 0 && (oriTask.Next > nowTS) {
				// 有定时请求记录，且新的记录是立即执行请求，且计划执行时间较迟，则更新计划执行时间为现在
				if res := r.db.WithContext(ctx).Model(&newTask).Where("id = ?", oriTask.ID).Update("next", nowTS); res.Error != nil {
					log.Println(err)
					time.Sleep(time.Second)
				} else {
					count += uint64(res.RowsAffected)
				}
			}
		}
	}
	if err != nil {
		return 0, err
	} else {
		return count, nil
	}
}
func (r *Repository) ListRaws(ctx context.Context, tag string, limit uint32) ([]Raw, error) {
	var raws []Raw
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Where("tag = ?", tag).Limit(int(limit)).Find(&raws).Error
		if err != nil {
			return err
		}
		if len(raws) != 0 {
			err = tx.Delete(&raws).Error
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return raws, nil
}
func (r *Repository) ConsumePendingTasks(ctx context.Context, limit uint32) ([]Task, error) {
	if limit == 0 || limit > 16 {
		limit = 16
	}
	var tasks []Task
	var err error
	if err := r.db.Where("next < ?", time.Now().UTC().Unix()).Order("next ASC").Limit(int(limit)).Find(&tasks).Error; err != nil {
		return nil, err
	}
	if len(tasks) != 0 {
		next := uint64(time.Now().UTC().Add(time.Second * 30).Unix())
		idList := make([]uint64, len(tasks))
		for i := range tasks {
			idList[i] = tasks[i].ID
		}
		err = r.db.Model(&tasks).Where("id IN ?", idList).Update("next", next).Error
		if err != nil {
			return nil, err
		}
	}
	return tasks, nil
}

func (r *Repository) SaveRaw(tag string, url string, data []byte) error {
	now := time.Now().UTC()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	return r.db.WithContext(ctx).Create(&Raw{
		URL:       url,
		Tag:       tag,
		Data:      data,
		CreatedAt: now,
	}).Error
}

func (r *Repository) AddToUpdateChannel(task Task) {
	r.taskUpdateChannel <- task
}

func (r *Repository) DeleteTask(ctx context.Context, id uint64) {
	r.db.WithContext(ctx).Unscoped().Delete(&Task{}, "id = ?", id)
}

func (r *Repository) ConsumeRaws(list []uint64) {
	for _, id := range list {
		r.rawConsumeChannel <- id
	}
}

func (r *Repository) ConsumeRaw(id uint64) {
	r.rawConsumeChannel <- id
}

func (r *Repository) recoverRaw() {
	ticker := time.NewTicker(time.Second * 10)
	for range ticker.C {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		res := r.db.WithContext(ctx).Unscoped().Model(&Raw{}).Where("deleted_at < ?", time.Now().Add(-time.Second*30)).Update("deleted_at", nil)
		if res.Error != nil {
			log.Println(res.Error)
			time.Sleep(time.Second)
		}
		cancel()
	}
}
