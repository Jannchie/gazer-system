package gs

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type Repository struct {
	db                *gorm.DB
	taskUpdateChannel chan Task
	rawConsumeChannel chan uint64
}

func NewRepository(db *gorm.DB) *Repository {
	r := &Repository{db: db, taskUpdateChannel: make(chan Task, 128)}
	go r.updateTask()
	return r
}
func (r *Repository) updateTask() {
	ticker := time.NewTicker(time.Second * 1)
	deleteList := make([]uint64, 0, 100)
	updateList := make([]uint64, 0, 100)
	consumeList := make([]uint64, 0, 100)
	for {
		select {
		case <-ticker.C:
			if len(deleteList) != 0 {
				r.db.Unscoped().Where("id in ?", deleteList).Delete(&Task{})
				deleteList = deleteList[:0]
			}
			if len(updateList) != 0 {
				r.db.Where("id in ?", updateList).Update("next", time.Now().UTC().Unix())
				updateList = updateList[:0]
			}
			if len(consumeList) != 0 {
				r.db.Unscoped().Where("id in ?", consumeList).Delete(&Raw{})
				consumeList = consumeList[:0]
			}
		case task := <-r.taskUpdateChannel:
			if task.IntervalMS == 0 {
				deleteList = append(deleteList, task.ID)
			} else {
				updateList = append(updateList, task.ID)
			}
		case raw := <-r.rawConsumeChannel:
			consumeList = append(consumeList, raw)
		}
	}
}
func (r *Repository) AddTasks(ctx context.Context, tasks []Task) error {
	return r.db.WithContext(ctx).Create(&tasks).Error
}
func (r *Repository) ListRaws(ctx context.Context, tag string, limit uint32) ([]Raw, error) {
	var raws []Raw
	err := r.db.WithContext(ctx).Where("tag = ?", tag).Limit(int(limit)).Find(&raws).Error
	if err != nil {
		return nil, err
	}
	return raws, nil
}
func (r *Repository) ConsumePendingTasks(ctx context.Context, limit uint32) ([]Task, error) {
	if limit == 0 || limit > 100 {
		limit = 10
	}
	var tasks []Task
	if err := r.db.WithContext(ctx).Where("next < ?", time.Now().UTC().Unix()).Limit(int(limit)).Find(&tasks).Error; err != nil {
		return nil, err
	} else {
		if len(tasks) != 0 {
			next := uint64(time.Now().UTC().Add(time.Second * 30).Unix())
			idList := make([]uint64, len(tasks))
			for i := range tasks {
				idList[i] = tasks[i].ID
			}
			r.db.WithContext(ctx).Model(&Task{}).Where("id in ?", idList).Update("next", next)
		}
		return tasks, nil
	}
}

func (r *Repository) SaveRaw(tag string, url string, data []byte) error {
	now := time.Now().UTC()
	return r.db.Create(&Raw{
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
