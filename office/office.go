package office

import (
	"github.com/icrowley/fake"
	"github.com/sirupsen/logrus"
)

var Freelancers = make([]*FreelancerGopher, 0)
var Tasks = make(chan Task)
var Reports = make(chan FreelancerGopher)
// 초기 오피스의 프리랜서들을 고용
func init(){
	fake.SetLang("en")
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
		// Reports의 수신자가 ready 될 때 까지 block 되므로 go routine으로 실행
		go func(){Reports <- *freelancer}()
		Freelancers = append(Freelancers, freelancer)
		go freelancer.Start()
	}
}

func AddTasks(num int){
	for i := 0; i < num; i++{
		logrus.Println("add a task")
		Tasks <- Task(i)
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