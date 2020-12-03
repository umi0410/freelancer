package office

import (
	"github.com/sirupsen/logrus"
	"math/rand"
	"strconv"
	"time"
)

const(
	FreelancerGopherIdleState string = "idle"
	FreelancerGopherWorkingState string = "working"
	FreelancerGopherToBeFiredState string = "to_be_fired" // CSS에서 바로 이용하기위해 _ case
	FreelancerGopherFiredState string = "fired"
	IdleTimeout time.Duration = time.Second * 5 // idle timeout 을 넘어서도 idle하면 짤림.
)

type FreelancerGopher struct{
	ID int
	Name string
	State string
	TasksDone int
	CurrentTask *Task
	WorkingHour time.Duration `json:"-"`
	WorkingHourString string `json:"WorkingHour"`
	TasksOut <-chan Task `json:"-"` // 채널은 json으로 직렬화될 수 없음.
}

type Task int

//type FreelancerStateReport struct{
//	Freelancer FreelancerGopher
//}

// a freelancer gopher is hired by you
func (freelancer *FreelancerGopher) Start(){
	logger := logrus.WithFields(logrus.Fields{"id": freelancer.ID, "name": freelancer.Name})
	logger.Println("직장을 다니기 시작합니다.")
	defer func() {
		freelancer.State = FreelancerGopherToBeFiredState
		FreelancerStateReports <- *freelancer
		logger.Warn("직장에서 잘릴 예정입니다.")
	}()

	FreelancerStateReports <- *freelancer // 입사 신고
	// 복사를 이용하므로 값 안전
	for	isFired := false; !isFired;{
		select {
		case task := <-freelancer.TasksOut:
			freelancer.Work(task)
		case <-time.After(IdleTimeout):
			isFired = true
		}
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
		FreelancerStateReports <- *freelancer
		logger.Println("작업을 마쳤습니다.")
	}()
	freelancer.State = FreelancerGopherWorkingState
	// 일을 시작한다는 보고
	FreelancerStateReports <- *freelancer
	logger.Println("작업을 시작합니다.")
	startTime := time.Now()

	freelancer.handleTask()

	freelancer.State = FreelancerGopherIdleState

	endTime := time.Now()
	elapsed := endTime.Sub(startTime)
	freelancer.WorkingHour += elapsed
	freelancer.TasksDone += 1

}

// 실제 작업
func (freelancer *FreelancerGopher) handleTask(){
	estimatedTime := time.Millisecond * time.Duration(rand.Int() % 6000)
	//estimatedTime := time.Duration(0)
	time.Sleep(estimatedTime)
}