package gs

import (
	"context"
	"log"
	"net/url"
	"time"

	"github.com/jannchie/gazer-system/api"

	"github.com/jannchie/speedo"
)

type WorkerOptions func(*WorkerConfig)

type Mode string

var (
	ModeDefault    = Mode("default")
	ModeOnlyParser = Mode("only-parser")
	ModeOnlySender = Mode("only-sender")
)

type WorkerConfig struct {
	Debug             bool
	SpeedometerServer string
	Mode              Mode
}

func WithMode(mode Mode) WorkerOptions {
	return func(c *WorkerConfig) {
		c.Mode = mode
	}
}

func WithDebug(debug bool) WorkerOptions {
	return func(c *WorkerConfig) {
		c.Debug = debug
	}
}

func WithSpeedometerServer(server string) WorkerOptions {
	return func(c *WorkerConfig) {
		c.SpeedometerServer = server
	}
}

type SenderUnit interface {
	WorkUnit
	Sender(chan<- string)
}

type ParserUnit interface {
	WorkUnit
	Parser(raw *api.Raw, client Client) error
	GetTargetURL() *url.URL
}

type WorkUnit interface {
	Run(ctx context.Context)
	GetName() string
	SetConfig(cfg *WorkerConfig)
}

type WorkerGroup struct {
	WorkerList []WorkUnit
	Server     string
	Client     *Client
	cfg        *WorkerConfig
}

func NewWorkerGroup(servers []string, options ...WorkerOptions) *WorkerGroup {
	return &WorkerGroup{
		Client: NewClientWithLB(servers...),
		cfg:    getConfig(options),
	}
}

func (g *WorkerGroup) AddWorker(targetURL string, sender func(chan<- *api.Task), parser func(*api.Raw, *Client) error, options ...WorkerOptions) {
	bw := NewBothWorker(g.Client, targetURL, sender, parser, options...)
	if len(options) == 0 {
		bw.cfg = g.cfg
		bw.SenderWorker.cfg = g.cfg
		bw.ParserWorker.cfg = g.cfg
	}
	g.WorkerList = append(g.WorkerList, bw)
}

func (g *WorkerGroup) AddByWorkUnit(w WorkUnit) {
	g.WorkerList = append(g.WorkerList, w)
	w.SetConfig(g.cfg)
}

func (g *WorkerGroup) Run(ctx context.Context) {
	for i := range g.WorkerList {
		go g.WorkerList[i].Run(ctx)
	}
	<-ctx.Done()
}

type Worker interface {
	Run(ctx context.Context)
}

type BothWorker struct {
	Client       *Client
	SenderWorker *SenderWorker
	ParserWorker *ParserWorker
	name         string
	cfg          *WorkerConfig
}

func (w *BothWorker) GetName() string {
	return w.name
}

func (w *BothWorker) SetConfig(cfg *WorkerConfig) {
	w.cfg = cfg
}
func (p *ParserWorker) SetConfig(cfg *WorkerConfig) {
	p.cfg = cfg
}
func (s *SenderWorker) SetConfig(cfg *WorkerConfig) {
	s.cfg = cfg
}

type SenderWorker struct {
	name        string
	cfg         *WorkerConfig
	Client      *Client
	TaskChannel chan *api.Task
	Sender      func(chan<- *api.Task)
}

type ParserWorker struct {
	cfg        *WorkerConfig
	Client     *Client
	Tag        string
	rawChannel chan *api.Raw
	Parser     func(*api.Raw, *Client) error
}

func NewSenderWorker(client *Client, name string, sender func(chan<- *api.Task), options ...WorkerOptions) *SenderWorker {
	return &SenderWorker{
		Client:      client,
		Sender:      sender,
		TaskChannel: make(chan *api.Task),
		name:        name,
		cfg:         getConfig(options),
	}
}

func NewParserWorker(gazerSystemClient *Client, tag string, parser func(*api.Raw, *Client) error, options ...WorkerOptions) *ParserWorker {
	return &ParserWorker{
		Client:     gazerSystemClient,
		Parser:     parser,
		rawChannel: make(chan *api.Raw),
		cfg:        getConfig(options),
		Tag:        tag,
	}
}
func (p *ParserWorker) GetName() string {
	return p.Tag
}

func (s *SenderWorker) GetName() string {
	return s.name
}

func NewBothWorker(client *Client, tag string, sender func(chan<- *api.Task), parser func(*api.Raw, *Client) error, options ...WorkerOptions) *BothWorker {
	return &BothWorker{
		Client:       client,
		cfg:          getConfig(options),
		SenderWorker: NewSenderWorker(client, tag, sender, options...),
		ParserWorker: NewParserWorker(client, tag, parser, options...),
	}
}

func (p *ParserWorker) Run(ctx context.Context) {
	if p.Tag != "" {
		go func() {
			// 填 Raw 通道
			for {
				timeoutCtx, cancel := context.WithTimeout(ctx, time.Second*10)
				dataList, err := p.Client.ListRaws(timeoutCtx, &api.ListRawsReq{Tag: p.Tag})
				cancel()
				if err != nil {
					log.Println(err)
					time.Sleep(time.Second)
					continue
				}
				raws := dataList.GetRaws()
				if len(raws) == 0 {
					time.Sleep(time.Second)
				} else {
					for _, data := range raws {
						p.rawChannel <- data
					}
				}
			}
		}()
		go func() {
			// 消费 RAW 通道
			s := speedo.NewSpeedometer(speedo.Config{Name: p.GetName() + " Parser", Log: p.cfg.Debug})
			for data := range p.rawChannel {
				err := p.Parser(data, p.Client)
				if err != nil {
					log.Println(data.Url)
					log.Println(string(data.Data))
					log.Println(err)
					continue
				} else {
					if err != nil {
						log.Println(err)
						continue
					}
					timeoutCtx, cancel := context.WithTimeout(ctx, time.Second*10)
					_, err = p.Client.ConsumeRaws(timeoutCtx, &api.ConsumeRawsReq{IdList: []uint64{data.GetId()}})
					cancel()
					if err != nil {
						log.Println(err)
						continue
					}
					s.AddCount(1)
				}
			}
		}()
	}
	<-ctx.Done()
	close(p.rawChannel)
}

func (s *SenderWorker) Run(ctx context.Context) {
	go func() {
		speedometer := speedo.NewSpeedometer(speedo.Config{Name: s.GetName() + " Sender", Log: s.cfg.Debug})
		// 消费 URL 通道
		for task := range s.TaskChannel {
			_, err := s.Client.AddTasks(context.Background(), &api.AddTasksReq{Tasks: []*api.Task{task}})
			if err != nil {
				log.Println(err)
				time.Sleep(time.Second)
				// retry, use goroutines to prevent from deadlock
				go func() {
					s.TaskChannel <- task
				}()
				continue
			} else {
				speedometer.AddCount(1)
			}
		}
	}()
	go func() {
		// 填 URL 通道
		s.Sender(s.TaskChannel)
	}()
	<-ctx.Done()
	close(s.TaskChannel)
}
func (w *BothWorker) Run(ctx context.Context) {
	w.ParserWorker.cfg = w.cfg
	w.SenderWorker.cfg = w.cfg
	switch w.cfg.Mode {
	case ModeOnlyParser:
		go w.ParserWorker.Run(ctx)
	case ModeOnlySender:
		go w.SenderWorker.Run(ctx)
	default:
		go w.SenderWorker.Run(ctx)
		go w.ParserWorker.Run(ctx)
	}
	<-ctx.Done()
}

func getConfig(options []WorkerOptions) *WorkerConfig {
	cfg := WorkerConfig{}
	for _, opt := range options {
		opt(&cfg)
	}
	return &cfg
}
