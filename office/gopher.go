package office

import (
	"github.com/sirupsen/logrus"
	"math/rand"
	"strconv"
	"time"
)

type FreelancerGopher struct{
	ID int
	Name string
	State string
	TasksDone int
	WorkingHour time.Duration `json:"-"`
	WorkingHourString string `json:"WorkingHour"`
	TasksOut <-chan Task `json:"-"` // 채널은 json으로 직렬화될 수 없음.
}

type Task int


// a freelancer gopher is hired by you
func (freelancer *FreelancerGopher) Start(){
	logger := logrus.WithFields(logrus.Fields{"name": freelancer.ID})
	logger.Println("Start")
	defer logger.Println("Finish")

	for task := range freelancer.TasksOut{
		// 복사를 이용하므로 값 안전

		freelancer.Work(task)
		Reports <- *freelancer
	}
}

// a freelancer keeps working
func (freelancer *FreelancerGopher) Work(task Task){
	logger := logrus.WithFields(logrus.Fields{
		"id": freelancer.ID,
		"name": freelancer.Name, "task": task})
	defer func(){
		//일을 마쳤다는 보고
		freelancer.WorkingHourString = strconv.Itoa(int(freelancer.WorkingHour.Seconds())) + " s"
		Reports <- *freelancer
		logger.Println("Finish work")
	}()
	freelancer.State = "working"
	// 일을 시작한다는 보고
	Reports <- *freelancer
	logger.Println("Start work")
	startTime := time.Now()

	freelancer.HandleTask()

	freelancer.State = "idle"

	endTime := time.Now()
	elapsed := endTime.Sub(startTime)
	freelancer.WorkingHour += elapsed
	freelancer.TasksDone += 1

}

func (freelancer *FreelancerGopher) HandleTask(){
	estimatedTime := time.Millisecond * time.Duration(rand.Int() % 6000)
	//estimatedTime := time.Millisecond * time.Duration(3 + rand.Int() % 6)
	time.Sleep(estimatedTime)
}