package gs

import (
	"context"
	"github.com/jannchie/gazer-system/api"
	"log"
	"net/url"

	"github.com/jannchie/speedo"
)

type SenderUnit interface {
	Sender(chan<- string)
}

type ParserUnit interface {
	Parser(raw *api.Raw) error
	GetTargetURL() *url.URL
}

type WorkUnit interface {
	Sender(chan<- *api.Task)
	Parser(*api.Raw) error
	GetTag() string
}

type WorkerGroup struct {
	WorkerList []*Worker
	Server     string
}

func NewWorkerGroup(server string) *WorkerGroup {
	return &WorkerGroup{
		Server: server,
	}
}

func (g *WorkerGroup) AddWorker(targetURL string, sender func(chan<- *api.Task), parser func(*api.Raw) error) {
	g.WorkerList = append(g.WorkerList, NewWorker(g.Server, targetURL, sender, parser))
}

func (g *WorkerGroup) AddByWorkUnit(w WorkUnit) {
	g.WorkerList = append(g.WorkerList, NewWorker(g.Server, w.GetTag(), w.Sender, w.Parser))
}
func (g *WorkerGroup) Run(ctx context.Context) {
	for i := range g.WorkerList {
		go g.WorkerList[i].Run(ctx)
	}
	<-ctx.Done()
}

type Worker struct {
	Worker       api.GazerSystemClient
	SenderWorker *SenderWorker
	ParserWorker *ParserWorker
}

type SenderWorker struct {
	GazerSystemClient api.GazerSystemClient
	taskChannel       chan *api.Task
	Sender            func(chan<- *api.Task)
}

type ParserWorker struct {
	GazerSystemClient api.GazerSystemClient
	Tag               string
	rawChannel        chan *api.Raw
	Parser            func(*api.Raw) error
}

func NewSenderWorker(gazerSystemClient api.GazerSystemClient, sender func(chan<- *api.Task)) *SenderWorker {
	return &SenderWorker{
		GazerSystemClient: gazerSystemClient,
		Sender:            sender,
		taskChannel:       make(chan *api.Task),
	}
}
func NewParserWorker(gazerSystemClient api.GazerSystemClient, tag string, parser func(*api.Raw) error) *ParserWorker {
	return &ParserWorker{
		GazerSystemClient: gazerSystemClient,
		Parser:            parser,
		rawChannel:        make(chan *api.Raw),
		Tag:               tag,
	}
}
func NewWorker(server string, tag string, sender func(chan<- *api.Task), parser func(*api.Raw) error) *Worker {
	core := NewClient(server)
	return &Worker{
		Worker:       core,
		SenderWorker: NewSenderWorker(core, sender),
		ParserWorker: NewParserWorker(core, tag, parser),
	}
}

func (w *ParserWorker) Run(ctx context.Context) {
	if w.Tag != "" {
		go func() {
			for {
				dataList, err := w.GazerSystemClient.ListRaws(context.Background(), &api.ListRawsReq{Tag: w.Tag})
				if err != nil {
					log.Println(err)
					continue
				}
				for _, data := range dataList.GetRaws() {
					w.rawChannel <- data
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
						_, err = w.GazerSystemClient.ConsumeRaws(ctx, &api.ConsumeRawsReq{IdList: []uint64{data.GetId()}})
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
		// 消费 URL 通道
		for {
			select {
			case task, ok := <-w.taskChannel:
				if !ok {
					return
				}
				_, err := w.GazerSystemClient.AddTasks(context.Background(), &api.AddTasksReq{Tasks: []*api.Task{task}})
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
func (w *Worker) Run(ctx context.Context) {
	go w.ParserWorker.Run(ctx)
	go w.SenderWorker.Run(ctx)
	<-ctx.Done()
}
