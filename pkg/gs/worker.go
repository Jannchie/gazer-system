package gs

import (
	"context"
	"github.com/jannchie/gazer-system/api"
	"log"
	"net/url"
	"time"

	"github.com/jannchie/speedo"
)

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
}

type WorkerGroup struct {
	WorkerList []WorkUnit
	Server     string
	Client     *Client
}

func NewWorkerGroup(servers ...string) *WorkerGroup {
	return &WorkerGroup{
		Client: NewClientWithLB(servers...),
	}
}

func (g *WorkerGroup) AddWorker(targetURL string, sender func(chan<- *api.Task), parser func(*api.Raw, *Client) error) {
	g.WorkerList = append(g.WorkerList, NewBothWorker(g.Client, targetURL, sender, parser))
}

func (g *WorkerGroup) AddByWorkUnit(w WorkUnit) {
	g.WorkerList = append(g.WorkerList, w)
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
}

func (w *BothWorker) GetName() string {
	return w.name
}

type SenderWorker struct {
	name        string
	Client      *Client
	TaskChannel chan *api.Task
	Sender      func(chan<- *api.Task)
}

type ParserWorker struct {
	name       string
	Client     *Client
	Tag        string
	rawChannel chan *api.Raw
	Parser     func(*api.Raw, *Client) error
}

func NewSenderWorker(client *Client, name string, sender func(chan<- *api.Task)) *SenderWorker {
	return &SenderWorker{
		Client:      client,
		Sender:      sender,
		TaskChannel: make(chan *api.Task),
		name:        name,
	}
}

func NewParserWorker(gazerSystemClient *Client, tag string, parser func(*api.Raw, *Client) error) *ParserWorker {
	return &ParserWorker{
		Client:     gazerSystemClient,
		Parser:     parser,
		rawChannel: make(chan *api.Raw),
		Tag:        tag,
	}
}
func (p *ParserWorker) GetName() string {
	return p.Tag
}

func (s *SenderWorker) GetName() string {
	return s.name
}

func NewBothWorker(client *Client, tag string, sender func(chan<- *api.Task), parser func(*api.Raw, *Client) error) *BothWorker {
	return &BothWorker{
		Client:       client,
		SenderWorker: NewSenderWorker(client, tag, sender),
		ParserWorker: NewParserWorker(client, tag, parser),
	}
}

func (p *ParserWorker) Run(ctx context.Context) {
	if p.Tag != "" {
		go func() {
			// 填 Raw 通道
			for {
				dataList, err := p.Client.ListRaws(context.Background(), &api.ListRawsReq{Tag: p.Tag})
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
			s := speedo.NewSpeedometer(speedo.Config{Name: p.GetName() + " Parser", Log: false})
			for {
				select {
				case data, ok := <-p.rawChannel:
					if !ok {
						return
					}
					err := p.Parser(data, p.Client)
					if err != nil {
						log.Println(err)
						continue
					} else {
						if err != nil {
							log.Println(err)
							continue
						}
						_, err = p.Client.ConsumeRaws(ctx, &api.ConsumeRawsReq{IdList: []uint64{data.GetId()}})
						if err != nil {
							log.Println(err)
							continue
						}
						s.AddCount(1)
					}
				}
			}
		}()
	}
	<-ctx.Done()
	close(p.rawChannel)
}

func (s *SenderWorker) Run(ctx context.Context) {
	go func() {
		speedometer := speedo.NewSpeedometer(speedo.Config{Name: s.GetName() + " Sender", Log: true})
		// 消费 URL 通道
		for {
			select {
			case task, ok := <-s.TaskChannel:
				if !ok {
					return
				}
				_, err := s.Client.AddTasks(context.Background(), &api.AddTasksReq{Tasks: []*api.Task{task}})
				if err != nil {
					log.Println(err)
					time.Sleep(time.Second)
					continue
				} else {
					speedometer.AddCount(1)
				}
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
	go w.ParserWorker.Run(ctx)
	go w.SenderWorker.Run(ctx)
	<-ctx.Done()
}
