package office

import (
	"github.com/icrowley/fake"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

const (
	// 해고 작업 시작 후 마치기까지 걸리는 시간
	// to_be_fired => fired로 처리되는데에 걸리는 시간.
	FireDuration = time.Second * 1
)
var (
	Freelancers = make([]*FreelancerGopher, 0)
	Tasks = make(chan Task)
	FreelancerStateReports = make(chan FreelancerGopher) // 개별 freelancer의 상태에 대한 보고서
	FreelancerFireReports = make(chan FreelancerGopher) // 인사팀의 freelancer 해고에 대한 보고서
	HRMutex = new(sync.Mutex) // 인사팀의 Freelancer 해고 작업이 한 명씩 이루어지도록 하기위한 Mutex Lock
	MiniFreelancer = 3 // 어떻게 구현하지..?
)
// 초기 오피스의 프리랜서들을 고용
func NewOffice(){
	HireFreelancers(5)
}

func HireFreelancers(num int){
	for i := 0; i < num; i++{
		var id int
		length := len(Freelancers)
		if length == 0{
			id = 1
		} else{
			id = Freelancers[length-1].ID + 1
		}
		freelancer := &FreelancerGopher{
			ID: id,
			Name: createUniqueName(),
			State: "idle",
			TasksOut: Tasks,
		}
		// 새로 고용한 프리랜서를 보고함.
		// Reports의 수신자가 ready 될 때 까지 block 되므로 go routine으로 실
		Freelancers = append(Freelancers, freelancer)
		go func(){
			freelancerID := freelancer.ID
			freelancerIndex := -1
			freelancer.Start()

			HRMutex.Lock()
			logrus.WithField("부서", "인사과").WithField("name", freelancer.Name).Warn(freelancer.Name, "를 자르기 시작합니다. 한 명씩만 순서대로 잘려야합니다.")
			time.Sleep(FireDuration)
			for i, freelancer := range Freelancers{
				if freelancerID == freelancer.ID{
					freelancerIndex = i
					break
				}
			}

			Freelancers = append(Freelancers[:freelancerIndex], Freelancers[freelancerIndex+1:]...)
			logrus.WithField("부서", "인사과").WithField("name", freelancer.Name).Warn(freelancer.Name, "를 잘랐습니다.")
			freelancer.State = FreelancerGopherFiredState
			FreelancerFireReports <- *freelancer
			HRMutex.Unlock()
		}()
	}
}

func AddTasks(num int){
	for i := 0; i < num; i++{
		Tasks <- Task(i)
		logrus.Println("작업을 추가했습니다.")
	}
}

func createUniqueName() string{
	var name string
	for exists:= true; exists;{
		exists = false
		name = fake.FirstName()
		for _, freelancer := range Freelancers{
			if freelancer.Name == name{
				exists = true
			}
		}
	}

	return name
}