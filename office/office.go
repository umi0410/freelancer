package office

import (
	"context"
	"github.com/icrowley/fake"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

const (
	// 해고 작업 시작 후 마치기까지 걸리는 시간
	// to_be_fired => fired로 처리되는데에 걸리는 시간.
	FireDuration = time.Nanosecond * 1
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
		idleSign := make(chan bool)
		//finishedChan := make(chan bool)
		ctx, notifyFired := context.WithCancel(context.Background())
		freelancer := &FreelancerGopher{
			Context: ctx,
			ID: id,
			Name: createUniqueName(),
			State: "idle",
			TasksOut: Tasks,
			IdleSign: idleSign,
			Mutex: new(sync.Mutex),
			//Finished: finishedChan,
		}
		// 새로 고용한 프리랜서를 보고함.
		HRMutex.Lock()
		Freelancers = append(Freelancers, freelancer)
		HRMutex.Unlock()

		go func(){
			freelancer.Start()
		}()
		go func(){
			for _ = range idleSign{
				HRMutex.Lock()

				logrus.WithField("부서", "인사과").WithField("name", freelancer.Name).Println(freelancer.Name, "의 고용을 검토합니다. 한 명씩만 순서대로 잘려야합니다.")
				if len(Freelancers) <= MiniFreelancer{
					logrus.WithField("부서", "인사과").WithField("name", freelancer.Name).Println("현재 최소 인력을 유지 중이므로", freelancer.Name, "는 잘리지 않습니다.")
				} else{
					logrus.WithField("부서", "인사과").WithField("name", freelancer.Name).Println(freelancer.Name, "를 자르는 방향으로 작업 중입니다.")
					time.Sleep(FireDuration)
					logrus.WithField("부서", "인사과").WithField("name", freelancer.Name).Println(freelancer.Name, "를 자르기로 결정했습니다.")

					notifyFired() // context.Done()

					freelancerIndex := -1

					for i, f := range Freelancers{
						if freelancer.ID == f.ID{
							freelancerIndex = i
							break
						}
					}
					if freelancerIndex == -1{
						logrus.WithField("부서", "인사과").WithField("name", freelancer.Name).Println(freelancer.Name, "은 이미 잘렸군요. idle 신호 후 자르는 동안 다시 idle 신호를 보내고 작업을 했나봅니다.")
					} else{
						Freelancers = append(Freelancers[:freelancerIndex], Freelancers[freelancerIndex+1:]...)
						logrus.WithField("부서", "인사과").WithField("name", freelancer.Name).Println(freelancer.Name, "를 잘랐습니다.")
					}
				}

				HRMutex.Unlock()
			}

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