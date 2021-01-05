package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/umi0410/freelancer/http"
	"github.com/umi0410/freelancer/office"
	"math/rand"
	"time"
)


func init(){
	//logrus.SetLevel(logrus.ErrorLevel)
	logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true, DisableColors: false})
	logrus.SetLevel(logrus.DebugLevel)
	rand.Seed(time.Now().UnixNano())
}

func main(){
	e := http.NewEcho()

	office.MainOffice = office.NewOffice()
	go func(){
		office.MainOffice.HireFreelancers(4)
	}()
	go func(){
		office.MainOffice.AddTasks(0)
	}()

	e.Logger.Fatal(e.Start(":1323"))
	fmt.Println("hello")
}
