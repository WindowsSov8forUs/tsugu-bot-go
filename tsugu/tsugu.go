package tsugu

import (
	"fmt"
	"strings"
	"sync"

	"github.com/WindowsSov8forUs/tsugu-bot-go/adapter"
	"github.com/WindowsSov8forUs/tsugu-bot-go/config"
	log "github.com/WindowsSov8forUs/tsugu-bot-go/mylog"
)

type tsuguHandler func(session adapter.Session, bot adapter.Bot, conf *config.Config) (bool, error)

type tsuguTempHandler func(session adapter.Session, bot adapter.Bot, conf *config.Config) (bool, bool, error)

type TempHandlerQueue struct {
	mutex    *sync.Mutex
	handlers []tsuguTempHandler
}

func NewTempHandlerQueue() *TempHandlerQueue {
	return &TempHandlerQueue{
		mutex:    &sync.Mutex{},
		handlers: make([]tsuguTempHandler, 0),
	}
}

func (q *TempHandlerQueue) Push(handler tsuguTempHandler) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	q.handlers = append(q.handlers, handler)
}

func (q *TempHandlerQueue) isEmpty() bool {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	return len(q.handlers) == 0
}

func (q *TempHandlerQueue) Pop() (tsuguTempHandler, bool) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	if len(q.handlers) == 0 {
		return nil, true
	}
	handler := q.handlers[0]
	q.handlers = q.handlers[1:]
	return handler, false
}

func (q *TempHandlerQueue) Clear() {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	q.handlers = make([]tsuguTempHandler, 0)
}

func (q *TempHandlerQueue) Copy() *TempHandlerQueue {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	handlersCopy := make([]tsuguTempHandler, len(q.handlers))
	copy(handlersCopy, q.handlers)

	return &TempHandlerQueue{
		mutex:    &sync.Mutex{},
		handlers: handlersCopy,
	}

}

var handlers = make([]tsuguHandler, 0)

var tempQueue = NewTempHandlerQueue()

func registerHandler(handler tsuguHandler) {
	handlers = append(handlers, handler)
}

func registerTempHandler(handler tsuguTempHandler) {
	fmt.Printf("handlers: %v\n", len(tempQueue.handlers))
	tempQueue.Push(handler)
	fmt.Printf("handlers: %v\n", len(tempQueue.handlers))
}

func Handler(session adapter.Session, bot adapter.Bot, conf *config.Config) error {
	message := session.Message()

	// 帮助信息
	if !conf.Tsugu.Functions.Help {
		if strings.HasPrefix(message, "帮助") {
			return helpCommand(session, strings.TrimPrefix(message, "帮助"), bot)
		} else if strings.HasPrefix(message, "help") {
			return helpCommand(session, strings.TrimPrefix(message, "help"), bot)
		} else if strings.HasSuffix(message, "-h") {
			return helpCommand(session, strings.TrimSuffix(message, "-h"), bot)
		}
	}

	// 临时处理函数
	if !tempQueue.isEmpty() {
		// 复制一份 tempQueue
		handlers := tempQueue.Copy()
		fmt.Printf("handlers: %v\n", len(handlers.handlers))
		tempQueue.Clear()

		for !handlers.isEmpty() {
			handler, empty := handlers.Pop()
			fmt.Printf("handlers: %v\n", len(handlers.handlers))
			if empty {
				break
			}
			_, retryable, err := handler(session, bot, conf)
			if retryable {
				registerTempHandler(handler)
			}
			if err != nil {
				log.Warnf("<Tsugu> 临时处理函数出错: %v", err)
			}
		}
	}

	// 车牌转发
	if conf.Tsugu.Functions.CarForward {
		forwarded, err := submitCarMessage(session, bot, conf)
		if err != nil {
			log.Errorf("<Tsugu> 车牌转发失败: %v", err)
		}
		if forwarded {
			return nil
		}
	}

	// 进行命令匹配
	for _, handler := range handlers {
		if ok, err := handler(session, bot, conf); err != nil {
			return err
		} else if ok {
			return nil
		}
	}
	return nil
}

func sendMessage(session adapter.Session, bot adapter.Bot, response []*ResponseData) error {
	message := &adapter.Message{}
	for _, data := range response {
		switch data.Type {
		case "string":
			message.Text(data.String)
		case "base64":
			message.Image(data.String)
		}
	}
	return bot.Send(session, message)
}
