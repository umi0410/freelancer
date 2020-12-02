package main

import (
	"github.com/umi0410/freelancer/http"
	"github.com/umi0410/freelancer/office"
	"math/rand"
	"time"
)


func init(){
	rand.Seed(time.Now().UnixNano())
}

func main(){
	e := http.NewEcho()
	office.AddTasks(3)
	e.Logger.Fatal(e.Start(":1323"))
}
