package http

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	"github.com/umi0410/freelancer/office"
	netHttp "net/http"
	"strconv"
)

type AddRequest struct{
	Number int
}

var (
	upgrader = websocket.Upgrader{}
	WebsocketConnections = []*websocket.Conn{}
)

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
	e.GET("/ws/", wsHandler)

	//go broadcast()

	e.Static("/", "./public")
	return e
}

func listFreelancersHandler(context echo.Context)  error {
	return context.JSON(200, office.Freelancers)
}


func createFreelancersHandler(context echo.Context)  error {
	body := &AddRequest{}
	err := context.Bind(body)
	if err != nil{
		logrus.Panic(err)
	}
	logrus.Info(body)
	office.HireFreelancers(body.Number)

	return context.String(201, strconv.Itoa(body.Number) + " freelancers has been hired.")
}

func createTasksHandler(context echo.Context) error {
	body := &AddRequest{}
	err := context.Bind(body)
	if err != nil {
		logrus.Panic(err)
	}
	logrus.Info(body)
	// 받아주는 채널이 있어야만 집어넣을 수 있음 => 집어넣는 작업도 goroutine으로
	go office.AddTasks(body.Number)

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
	defer ws.Close()
	for report := range office.Reports {
		// Write
		//data, err := json.Marshal(map[string]interface{}{"name": "report", "data": report})
		//for _, ws := range WebsocketConnections{
			data, err := json.Marshal(report)
			if err != nil {
				logrus.Error(err)
			}
			err = ws.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				logrus.Println("A user has been disconnected")
				return err
			}
		//}
	}

	return nil
}

