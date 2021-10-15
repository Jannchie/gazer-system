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
	Sender(chan<- string)
	Run(ctx context.Context)
}

type ParserUnit interface {
	Parser(raw *api.Raw) error
	GetTargetURL() *url.URL
	Run(ctx context.Context)
}

type WorkUnit interface {
	Run(ctx context.Context)
}

type WorkerGroup struct {
	WorkerList []WorkUnit
	Server     string
	Client     api.GazerSystemClient
}

func NewWorkerGroup(server string) *WorkerGroup {
	return &WorkerGroup{
		Server: server,
		Client: NewClient(server),
	}
}

func (g *WorkerGroup) AddWorker(targetURL string, sender func(chan<- *api.Task), parser func(*api.Raw) error) {
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
	Client       api.GazerSystemClient
	SenderWorker *SenderWorker
	ParserWorker *ParserWorker
}

type SenderWorker struct {
	Client      api.GazerSystemClient
	taskChannel chan *api.Task
	Sender      func(chan<- *api.Task)
}

type ParserWorker struct {
	Client     api.GazerSystemClient
	Tag        string
	rawChannel chan *api.Raw
	Parser     func(*api.Raw) error
}

func NewSenderWorker(gazerSystemClient api.GazerSystemClient, sender func(chan<- *api.Task)) *SenderWorker {
	return &SenderWorker{
		Client:      gazerSystemClient,
		Sender:      sender,
		taskChannel: make(chan *api.Task),
	}
}

func NewParserWorker(gazerSystemClient api.GazerSystemClient, tag string, parser func(*api.Raw) error) *ParserWorker {
	return &ParserWorker{
		Client:     gazerSystemClient,
		Parser:     parser,
		rawChannel: make(chan *api.Raw),
		Tag:        tag,
	}
}
func NewBothWorker(client api.GazerSystemClient, tag string, sender func(chan<- *api.Task), parser func(*api.Raw) error) *BothWorker {
	return &BothWorker{
		Client:       client,
		SenderWorker: NewSenderWorker(client, sender),
		ParserWorker: NewParserWorker(client, tag, parser),
	}
}

func (w *ParserWorker) Run(ctx context.Context) {
	if w.Tag != "" {
		go func() {
			for {
				dataList, err := w.Client.ListRaws(context.Background(), &api.ListRawsReq{Tag: w.Tag})
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
						w.rawChannel <- data
					}
				}
			}
			// 填 Raw 通道
		}()
		go func() {
			// 消费 RAW 通道
			s := speedo.NewSpeedometer(speedo.Config{Name: "Parser", Log: true})
			for {
				select {
				case data, ok := <-w.rawChannel:
					if !ok {
						return
					}
					err := w.Parser(data)
					if err != nil {
						log.Println(err)
						continue
					} else {
						if err != nil {
							log.Println(err)
							continue
						}
						_, err = w.Client.ConsumeRaws(ctx, &api.ConsumeRawsReq{IdList: []uint64{data.GetId()}})
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
	close(w.rawChannel)
}

func (w *SenderWorker) Run(ctx context.Context) {
	go func() {
		s := speedo.NewSpeedometer(speedo.Config{Name: "Sender", Log: true})
		// 消费 URL 通道
		for {
			select {
			case task, ok := <-w.taskChannel:
				if !ok {
					return
				}
				s.AddCount(1)
				_, err := w.Client.AddTasks(context.Background(), &api.AddTasksReq{Tasks: []*api.Task{task}})
				if err != nil {
					log.Println(err)
					continue
				}
			}
		}
	}()
	go func() {
		// 填 URL 通道
		w.Sender(w.taskChannel)
	}()
	<-ctx.Done()
	close(w.taskChannel)
}
func (w *BothWorker) Run(ctx context.Context) {
	go w.ParserWorker.Run(ctx)
	go w.SenderWorker.Run(ctx)
	<-ctx.Done()
}
