// freelancer의 http 통신과 관련된 내용을 담습니다.
// 클라이언트들은 이 패키지의 문서를 통해 어떠한 형식의 데이터를
// 전달해야하는지, 자신이 받는 데이터가 어떠한 형식인지를 파악할 수 있습니다.
package http

import (
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	"github.com/umi0410/freelancer/office"
	netHttp "net/http"
	"strconv"
)

// websocket을 통해 클라이언트에게 메시지를 보낼 때 사용하는 타입.
//
// Type 은 주로 타입의 이벤트가 발생하였는지를 의미하고,
// Data 는 현재 발생한 이벤트가 전달하고싶은 정보를 의미한다.
// Data 에는 어떠한 자료형도 올 수 있지만, 전송될 때에는 JSON 형태로
// 전송된다.
//
// 클라이언트는 어떤 이벤트 핸들러가 어떤 Type에 의해 어떤 동작을 할 지
// 정의하고, 그 동작 시에 필요한 data는 Data에서 얻을 수 있다.
type SocketMessage struct{
	Type string
	Data interface{} // 어떤 타입이든 가능
}

// Task나 Freelancer를 추가할 때 클라이언트가 보내주어야하는 
// body 형식
type AddRequest struct{
	Number int
}

var (
	// 현재 연결되어있는 websocket connection들
	WebsocketConnections = []*websocket.Conn{}
	upgrader = websocket.Upgrader{}
)

// Echo 서버를 run할 수 있는 *echo.Echo type을 생성합니다.
func NewEcho() *echo.Echo{
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
	  AllowOrigins: []string{"*"},
	  AllowHeaders: []string{"*"},
	}))

	group := e.Group("/api")
	group.GET("/freelancers", listFreelancersHandler)
	group.POST("/freelancers", createFreelancersHandler)
	group.POST("/tasks", createTasksHandler)

	e.GET("/ws", wsHandler)

	e.Static("/", "./public")
	return e
}

func listFreelancersHandler(context echo.Context)  error {
	return context.JSON(200, office.MainOffice.Freelancers)
}


func createFreelancersHandler(context echo.Context)  error {
	body := &AddRequest{}
	err := context.Bind(body)
	if err != nil{
		logrus.Panic(err)
	}
	logrus.Infof("%#v", body)
	office.MainOffice.HireFreelancers(body.Number)

	return context.String(201, strconv.Itoa(body.Number) + " freelancers has been hired.")
}

func createTasksHandler(context echo.Context) error {
	body := &AddRequest{}
	err := context.Bind(body)
	if err != nil {
		logrus.Panic(err)
	}
	logrus.Infof("%#v", body)
	// 받아주는 채널이 있어야만 집어넣을 수 있음 => 집어넣는 작업도 goroutine으로
	go office.MainOffice.AddTasks(body.Number)

	return context.String(201, strconv.Itoa(body.Number)+" tasks has been added.")
}

//수정 중..
//func broadcast(){
//	for report := range office.Reports {
//		// Write
//		//data, err := json.Marshal(map[string]interface{}{"name": "report", "data": report})
//		for _, ws := range WebsocketConnections{
//			data, err := json.Marshal(report)
//			if err != nil {
//				logrus.Error(err)
//			}
//			err = ws.WriteMessage(websocket.TextMessage, data)
//			if err != nil {
//				logrus.Println("A user has been disconnected")
//				return err
//			}
//		}
//	}
//}

func wsHandler(c echo.Context) error {
	upgrader.CheckOrigin= func(r *netHttp.Request) bool {return true} // origin이 달라도 허용
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		logrus.Panic(err)
	}

	//ws.send
	defer ws.Close()
	for {
		var message interface{}
		var messageType string
		select{
		case message = <- office.MainOffice.FreelancerStateReports:
			messageType = "freelancer_state_report"
		case message = <- office.MainOffice.FreelancerFireReports:
			messageType = "freelancer_fire_report"
		}

		err = ws.WriteJSON(SocketMessage{Type: messageType, Data: message})
		if err != nil {
			logrus.Println("A user has been disconnected")
			logrus.Error(err)
			return err
		}
	}
}

