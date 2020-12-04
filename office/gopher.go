package office

import (
	"github.com/sirupsen/logrus"
	"context"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

const(
	FreelancerGopherIdleState string = "idle"
	FreelancerGopherWorkingState string = "working"
	FreelancerGopherToBeFiredState string = "to_be_fired" // CSS에서 바로 이용하기위해 _ case
	FreelancerGopherFiredState string = "fired"
	IdleTimeout time.Duration = time.Second * 3 // idle timeout 을 넘어서도 idle하면 짤림.
	BaseTaskDuration time.Duration = time.Millisecond * 3000
	RandomTaskBlock = 3
)

type FreelancerGopher struct{
	context.Context
	ID int
	Name string
	State string
	TasksDone int
	CurrentTask *Task
	WorkingHour time.Duration `json:"-"`
	WorkingHourString string `json:"WorkingHour"`
	TasksOut <-chan Task `json:"-"` // 채널은 json으로 직렬화될 수 없음.
	Mutex *sync.Mutex
	IdleSign chan<- bool `json:"-"` // Office에게 내가 idle하다고 신호.
}

type Task int

//type FreelancerStateReport struct{
//	Freelancer FreelancerGopher
//}

// a freelancer gopher is hired by you
func (freelancer *FreelancerGopher) Start() {
	logger := logrus.WithFields(logrus.Fields{"id": freelancer.ID, "name": freelancer.Name})
	defer func() {
		freelancer.State = FreelancerGopherFiredState
		FreelancerFireReports <- *freelancer
		close(freelancer.IdleSign)
	}()

	logger.Println("직장을 다니기 시작합니다.")

	FreelancerStateReports <- *freelancer // 입사 신고
	// 복사를 이용하므로 값 안전
	for toBeFired := false; !toBeFired;{
		freelancer.Mutex.Lock()
		select{
		case <- freelancer.Done():
			logger.Println("잘렸다는 통보를 받았습니다...")
			toBeFired = true
		default:
			select {
			case task := <-freelancer.TasksOut:
				freelancer.HandleTask(task)
			case <-time.After(IdleTimeout):
				logger.Debug("office에게 idle함을 신호...")
				freelancer.IdleSign <- true
				logger.Debug("office에게 idle함을 신호끝")
			}
		}
		freelancer.Mutex.Unlock()
	}
}

// a freelancer keeps working
func (freelancer *FreelancerGopher) HandleTask(task Task){
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

	freelancer.task()

	freelancer.State = FreelancerGopherIdleState

	endTime := time.Now()
	elapsed := endTime.Sub(startTime)
	freelancer.WorkingHour += elapsed
	freelancer.TasksDone += 1

}

// 실제 작업 내용이라고 가정.
func (freelancer *FreelancerGopher) task(){
	estimatedTime := BaseTaskDuration * time.Duration(rand.Int() % RandomTaskBlock)
	//estimatedTime := time.Duration(0)
	time.Sleep(estimatedTime)
}