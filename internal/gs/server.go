package gs

import (
	"context"
	"errors"
	"fmt"
	"github.com/jannchie/gazer-system/api"
	"github.com/jannchie/gazer-system/internal/variables"
	"github.com/jannchie/speedo"
	"github.com/wybiral/torgo"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"io"
	"log"
	"net"
	"net/http"
	"time"
)

type Server struct {
	repository    *Repository
	collector     *Collector
	collectSpeedo *speedo.Speedometer
	consumeSpeedo *speedo.Speedometer
	receiveSpeedo *speedo.Speedometer
}

func (s *Server) ConsumeRaws(_ context.Context, req *api.ConsumeRawsReq) (*api.OperationResp, error) {
	s.repository.ConsumeRaws(req.IdList)
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
	if req.Limit > 20 || req.Limit == 0 {
		req.Limit = 20
	}
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
	count, err := s.repository.AddTasks(ctx, tasks)
	if err != nil {
		return nil, err
	}
	s.receiveSpeedo.AddCount(count)
	return &api.OperationResp{Code: 1, Msg: "ok"}, nil
}

type Config struct {
	Debug             bool
	DSN               string
	TorSock5Host      string
	TorControllerHost string
	CollectHandle     CollectHandle
}

func getDefaultConfig() *Config {
	variables.Init()
	return &Config{
		Debug:             true,
		DSN:               *variables.DSN,
		TorSock5Host:      *variables.TorAddr,
		TorControllerHost: *variables.TorCtlAddr,
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
}

func NewDefaultServer() *Server {
	return NewServer(getDefaultConfig())
}

func NewServer(cfg *Config) *Server {
	variables.Init()
	logLevel := getLogLevel(cfg)

	return &Server{
		repository:    NewRepository(cfg.DSN, logLevel),
		collector:     NewCollector(cfg.TorSock5Host, cfg.TorControllerHost, cfg.CollectHandle),
		collectSpeedo: speedo.NewSpeedometer(speedo.Config{Log: false, Name: "Collect"}),
		consumeSpeedo: speedo.NewSpeedometer(speedo.Config{Log: false, Name: "Consume"}),
		receiveSpeedo: speedo.NewSpeedometer(speedo.Config{Log: false, Name: "Receive"}),
	}
}

func getLogLevel(cfg *Config) logger.LogLevel {
	var logLevel logger.LogLevel
	if cfg.Debug {
		logLevel = logger.Warn
	} else {
		logLevel = logger.Silent
	}
	return logLevel
}

func (s *Server) Run() {
	taskChan := make(chan Task)
	go s.addToChan(taskChan)
	go s.consumeTask(taskChan)
	s.serve()
}

func (s *Server) serve() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *variables.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	opts = append(opts, grpc.MaxSendMsgSize(10000000))
	opts = append(opts, grpc.MaxRecvMsgSize(10000000))
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
	lastRefresh     time.Time
	concurrency     uint8
	semChannel      chan struct{}
}

func NewCollector(proxySock5Host, proxyControllerHost string, collect CollectHandle) *Collector {
	for {
		client, err := torgo.NewClient(proxySock5Host)
		if err != nil {
			log.Println(err)
			log.Println("[PROXY] Cannot connect to proxy host")
			time.Sleep(time.Second)
			continue
		}

		proxyController, err := torgo.NewController(proxyControllerHost)
		if err != nil {
			log.Println(err)
			log.Println("[PROXY] Cannot connect to proxy controller")
			time.Sleep(time.Second)
			continue
		}

		err = proxyController.AuthenticatePassword(*variables.TorPassword)
		if err != nil {
			log.Println(err)
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
			lastRefresh:     time.Now().UTC(),
		}
	}
}

func (c *Collector) RefreshClient() error {
	duration := time.Since(c.lastRefresh)
	if duration > 10*time.Second {
		now := time.Now().UTC()
		c.lastRefresh = now
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
		if duration < time.Hour {
			log.Println("Refreshed Client! Duration: ", duration)
		}
	}
	return nil
}

func (c *Collector) collect(collector *Collector, url string) ([]byte, error) {
	return c.collectHandle(collector, url)
}
