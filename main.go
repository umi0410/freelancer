package main

import (
	"fmt"
	"github.com/umi0410/freelancer/http"
	"github.com/umi0410/freelancer/office"
	"math/rand"
	"time"
)


func init(){
	//logrus.SetLevel(logrus.WarnLevel)
	rand.Seed(time.Now().UnixNano())
}

func main(){
	e := http.NewEcho()

	go func(){
		office.NewOffice()
	}()
	go func(){
		office.AddTasks(3)
	}()

	e.Logger.Fatal(e.Start(":1323"))
	fmt.Println("hello")
}
