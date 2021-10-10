package gs

import (
	"context"
	"errors"
	"fmt"
	"github.com/jannchie/gs-server/api"
	"github.com/jannchie/speedo"
	"github.com/wybiral/torgo"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

type Server struct {
	repository    *Repository
	collector     *Collector
	collectSpeedo *speedo.Speedometer
	consumeSpeedo *speedo.Speedometer
}

func (s *Server) ConsumeRaws(ctx context.Context, req *api.ConsumeRawsReq) (*api.OperationResp, error) {
	err := s.repository.ConsumeRaws(ctx, req.IdList)
	if err != nil {
		return nil, err
	}
	s.consumeSpeedo.AddCount(uint64(len(req.IdList)))
	return &api.OperationResp{Msg: "ok", Code: 1}, nil
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

type Raw struct {
	ID        uint64         `gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at" gorm:"index"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	URL       string         `json:"url" gorm:"index"`
	Tag       string         `json:"tag" gorm:"index"`
	Data      []byte         `json:"data"`
}

func (s *Server) ListRaws(ctx context.Context, req *api.ListRawsReq) (*api.RawsResp, error) {
	raws, err := s.repository.ListRaws(ctx, req.Tag, req.Limit)
	if err != nil {
		return nil, err
	}
	resRaws := make([]*api.Raw, len(raws))
	for i, raw := range raws {
		resRaws[i] = &api.Raw{
			Url:       raw.URL,
			Tag:       raw.Tag,
			Id:        raw.ID,
			Timestamp: timestamppb.New(raw.CreatedAt),
			Data:      raw.Data,
		}
	}
	resp := api.RawsResp{
		Raws: resRaws,
	}
	return &resp, nil
}

func (s *Server) AddTasks(ctx context.Context, req *api.AddTasksReq) (*api.OperationResp, error) {
	length := len(req.Tasks)
	tasks := make([]Task, length)
	for i := range req.Tasks {
		tasks[i].URL = req.Tasks[i].Url
		tasks[i].Tag = req.Tasks[i].Tag
		tasks[i].IntervalMS = req.Tasks[i].IntervalMS
	}
	err := s.repository.AddTasks(ctx, tasks)
	if err != nil {
		return nil, err
	}
	return &api.OperationResp{Code: 1, Msg: "ok"}, nil
}

type Repository struct {
	db                *gorm.DB
	taskUpdateChannel chan Task
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
		default:
			task := <-r.taskUpdateChannel
			if task.IntervalMS == 0 {
				deleteList = append(deleteList, task.ID)
			} else {
				updateList = append(updateList, task.ID)
			}
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

func (r *Repository) ConsumeRaws(ctx context.Context, list []uint64) error {
	return r.db.Unscoped().WithContext(ctx).Where("id IN ?", list).Delete(&Raw{}).Error
}

type Config struct {
	Debug             bool
	DSN               string
	TorSock5Host      string
	TorControllerHost string
	CollectHandle     CollectHandle
}
type temporary interface {
	Temporary() bool // IsTemporary returns true if err is temporary.
}

func IsTemporary(err error) bool {
	te, ok := err.(temporary)
	return ok && te.Temporary()
}

type TemporaryError struct {
	error
}

func (t *TemporaryError) Temporary() bool {
	return true
}

var DefaultConfig = Config{
	Debug:             true,
	DSN:               "./data.db",
	TorSock5Host:      "localhost:9050",
	TorControllerHost: "localhost:9051",
	CollectHandle: func(c *Collector, targetURL string) ([]byte, error) {
		resp, err := c.client.Get(targetURL)
		if err != nil {
			return nil, err
		}
		defer func(Body io.ReadCloser) {
			_ = Body.Close()
		}(resp.Body)
		if resp.StatusCode >= 400 && resp.StatusCode < 500 {
			// if is 4XX error, should wait for proxy refresh.
			err := c.RefreshClient()
			if err != nil {
				return nil, TemporaryError{err}
			}
			return nil, TemporaryError{fmt.Errorf("status code error: %d", resp.StatusCode)}
		}
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		//log.Printf("Download succeed: %s\n", targetURL)
		return data, nil
	},
}

func NewDefaultServer() *Server {
	return NewServer(&DefaultConfig)
}

func NewServer(cfg *Config) *Server {
	logLevel := getLogLevel(cfg)
	db, err := gorm.Open(sqlite.Open(cfg.DSN), &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold:             time.Second, // Slow SQL threshold
				LogLevel:                  logLevel,    // Log level
				IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
				Colorful:                  true,        // Disable color
			},
		),
	})
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&Task{}, &Raw{})
	if err != nil {
		panic(err)
	}
	return &Server{
		repository:    NewRepository(db),
		collector:     NewCollector(cfg.TorSock5Host, cfg.TorControllerHost, cfg.CollectHandle),
		collectSpeedo: speedo.NewSpeedometer(speedo.Config{Log: cfg.Debug, Name: "Collect"}),
		consumeSpeedo: speedo.NewSpeedometer(speedo.Config{Log: cfg.Debug, Name: "Consume"}),
	}
}

func getLogLevel(cfg *Config) logger.LogLevel {
	var logLevel logger.LogLevel
	if cfg.Debug {
		logLevel = logger.Info
	} else {
		logLevel = logger.Silent
	}
	return logLevel
}

func (s *Server) Run(port uint16) {
	taskChan := make(chan Task)
	go s.addToChan(taskChan)
	go s.consumeTask(taskChan)
	s.serve(port)
}

func (s *Server) serve(port uint16) {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	api.RegisterGazerSystemServer(grpcServer, s)
	log.Println("* * * * * * * * * * * *")
	log.Println("* GAZER SYSTEM SERVER *")
	log.Println("* * * * * * * * * * * *")
	log.Printf("%+v\n", grpcServer.GetServiceInfo())
	err = grpcServer.Serve(lis)
	if err != nil {
		panic(err)
	}
}

func (s *Server) addToChan(taskChan chan<- Task) {
	for {
		tasks, err := s.repository.ConsumePendingTasks(context.Background(), 100)
		if err != nil {
			log.Println(err)
			continue
		}
		if len(tasks) == 0 {
			time.Sleep(time.Second)
			continue
		}
		for _, task := range tasks {
			taskChan <- task
		}
	}
}

func (s *Server) consumeTask(taskChan <-chan Task) {
	for {
		s.collector.semChannel <- struct{}{}
		select {
		case task := <-taskChan:
			go func() {
				s.consumeOneTask(task)
				s.collectSpeedo.AddCount(1)
			}()
		}
	}
}

func (s *Server) consumeOneTask(task Task) {
	targetURL := task.URL
	data, err := s.collector.collect(s.collector, targetURL)

	<-s.collector.semChannel

	if err != nil {
		if !IsTemporary(err) {
			log.Println(err)
			go s.repository.DeleteTask(context.Background(), task.ID)
		}
	} else {
		err = s.repository.SaveRaw(task.Tag, task.URL, data)
		if err != nil {
			log.Println(err)
		} else {
			go s.repository.AddToUpdateChannel(task)
		}
	}
}

type CollectHandle func(collector *Collector, targetURL string) ([]byte, error)
type Collector struct {
	collectHandle   CollectHandle
	proxyController *torgo.Controller
	client          *http.Client
	proxySock5Host  string
	lastRefresh     *time.Time
	concurrency     uint8
	semChannel      chan struct{}
}

func NewCollector(proxySock5Host, proxyControllerHost string, collect CollectHandle) *Collector {
	for {
		proxyController, err := torgo.NewController(proxyControllerHost)
		if err != nil {
			log.Println("[PROXY] Cannot connect to proxy controller")
			time.Sleep(time.Second)
			continue
		}
		client, err := torgo.NewClient(proxySock5Host)
		if err != nil {
			log.Println("[PROXY] Cannot connect to proxy host")
			time.Sleep(time.Second)
			continue
		}
		concurrency := uint8(128)
		return &Collector{
			proxyController: proxyController,
			client:          client,
			proxySock5Host:  proxySock5Host,
			collectHandle:   collect,
			concurrency:     concurrency,
			semChannel:      make(chan struct{}, concurrency),
		}
	}
}

func (c *Collector) RefreshClient() error {
	duration := time.Since(*c.lastRefresh)
	if duration > 10*time.Second {
		if duration < time.Hour {
			log.Println("Refreshed Client! Duration: ", duration)
		}
		now := time.Now()
		c.lastRefresh = &now
		if c.proxyController == nil {
			return errors.New("proxy not ready")
		}
		err := c.proxyController.Signal("NEWNYM")
		if err != nil {
			log.Println(err)
			return err
		}
		client, err := torgo.NewClient(c.proxySock5Host)
		if err != nil {
			log.Println(err)
			return err
		}
		c.client = client
		c.client.Timeout = time.Second * 10
	}
	return nil
}

func (c *Collector) collect(collector *Collector, url string) ([]byte, error) {
	return c.collectHandle(collector, url)
}
