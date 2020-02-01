package controller

import (
	"fmt"
	"log"
	"net/http"

	"github.com/JabinGP/demo-liveserver/tool"
	"github.com/fsnotify/fsnotify"
	"github.com/kataras/neffos"
	"github.com/kataras/neffos/gorilla"

	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/websocket"

	// 这里由于命名冲突，进行重命名导入
	gorillaWs "github.com/gorilla/websocket"
)

// WSWatcher WebSocket Watcher
var WSWatcher WebSocketWatcher

func init() {
	WSWatcher = NewWebSocketWatcher()
}

// WebSocketWatcher ...
type WebSocketWatcher struct {
	wsServer *neffos.Server
	watcher  Watcher
}

// NewWebSocketWatcher 初始化一个WebSocketWatcher
func NewWebSocketWatcher() WebSocketWatcher {
	w := WebSocketWatcher{}

	myUpgrader := gorilla.Upgrader(gorillaWs.Upgrader{CheckOrigin: func(*http.Request) bool {
		return true
	}})

	w.wsServer = websocket.New(myUpgrader, websocket.Events{
		// 这个事件其实并没有什么用
		// 这个服务器更多的功能是将监听到的事件主动推送到客户端
		// 用不到接收消息的功能
		websocket.OnNativeMessage: func(nsConn *websocket.NSConn, msg websocket.Message) error {
			log.Printf("Watcher server got: %s from [%s]", msg.Body, nsConn.Conn.ID())
			return nil
		},
	})
	// Disconnect事件触发，关闭监听器
	w.wsServer.OnDisconnect = func(c *websocket.Conn) {
		log.Printf("[%s] Disconnected from watcher server", c.ID())
		// 关闭监听器并通知事件接收线程关闭
		w.watcher.Close()
	}
	// Connect事件触发的方法，在这里获取请求中的路径
	// 启动监听器监听该路径
	w.wsServer.OnConnect = func(c *websocket.Conn) error {
		log.Printf("[%s] Connected to watcher server!", c.ID())
		ctx := websocket.GetContext(c)
		path := fmt.Sprintf(".%s", ctx.URLParam("pathname"))
		w.watcher = NewWatcher()
		w.watcher.Send = func(msg string) {
			c.Write(neffos.Message{IsNative: true, Body: []byte(msg)})
		}
		w.watcher.Add(path)
		w.watcher.Start()
		return nil
	}
	return w
}

// Serve 返回Handler用于Route
func (w *WebSocketWatcher) Serve() context.Handler {
	return websocket.Handler(w.wsServer)
}

// Watcher 文件变化监听器
type Watcher struct {
	watcher     *fsnotify.Watcher
	CloseSignal chan bool
	Send        func(msg string)
}

// NewWatcher 初始化一个监听器
func NewWatcher() Watcher {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Println(err)
	}
	closeSignal := make(chan bool)

	return Watcher{
		watcher:     watcher,
		CloseSignal: closeSignal,
	}
}

// Add 添加监听路径
func (w *Watcher) Add(path string) {
	err := w.watcher.Add(path)
	if err != nil {
		log.Println(err)
	}
}

// Start 启动监听器与接收器
func (w *Watcher) Start() {
	log.Println("[Watcher] started.")
	w.ReceiveEventAndSend()
	w.WaitCloseSignal()
	log.Println("[Watcher] stoped.")
}

// Close 关闭监听器
func (w *Watcher) Close() {
	w.SendCloseSignal()
	w.watcher.Close()
}

// SendCloseSignal 发送关闭信号，关闭信息收集进程
func (w *Watcher) SendCloseSignal() {
	w.CloseSignal <- true
}

// WaitCloseSignal 阻塞等待关闭信号，保持线程不结束
func (w *Watcher) WaitCloseSignal() {
	<-w.CloseSignal
}

// ReceiveEventAndSend 新开线程接收事件，并通过SendEvent方法将实现发送出去
func (w *Watcher) ReceiveEventAndSend() {
	go func() {
		log.Println("[Watcher.EventReceiver] started.")
		log.Println("[Watcher.EventReceiver] 过滤小于200ms间隔的事件.")
		timer := tool.NewIntervalTimer()
		for {
			select {
			case event, ok := <-w.watcher.Events:
				if !ok {
					log.Println("[Watcher.EventReceiver] stoped.")
					return
				}
				// 跳过间隔过短的事件
				if timer.Millisecond() < 200 {
					continue
				}
				log.Println("[Watcher.EventReceiver]", event.Op.String(), event.Name)
				w.Send(event.Op.String())
			case err, ok := <-w.watcher.Errors:
				if !ok {
					log.Println("[Watcher.EventReceiver] stoped.")
					return
				}
				log.Println("[Watcher.EventReceiver]", "err: "+err.Error())
				w.Send("err: " + err.Error())
			}
		}
	}()
}
