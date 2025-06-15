package task

import (
	"context"
	"github.com/Gustcat/task-server/internal/model"
	"github.com/Gustcat/task-server/internal/repository/converter"

	//"sync"
	"time"
	//"google.golang.org/grpc"
	//"google.golang.org/grpc/credentials/insecure"
)

const (
	adressAuthService = "auth-service:50051"
)

func (s *Serv) Create(ctx context.Context, task *model.Task, author int64) (int64, error) {
	users := make([]int64, 1, 2)
	users[0] = author
	if task.Operator != nil {
		users = append(users, *task.Operator)
	}
	//err := validateUsers(users)
	//if err != nil {
	//	return 0, err //TODO: обработать как ошибку валидации
	//}

	task.Author = author

	if task.Status == model.StatusDone {
		*task.CompletedAt = time.Now()
	}

	insertTask := converter.TaskToRepo(task)

	id, err := s.taskRepo.Create(ctx, insertTask)
	if err != nil {
		return 0, err

	}

	return id, nil

}

//func validateUsers(users []int64) error {
//	wg := &sync.WaitGroup{}
//	mu := &sync.Mutex{}
//	chanErr := make(chan error, len(users))
//
//	conn, err := grpc.NewClient(adressAuthService, grpc.WithTransportCredentials(insecure.NewCredentials()))
//	if err != nil {
//		log.Printf("failed to connect to server: %v", err)
//	}
//	defer func(conn *grpc.ClientConn) {
//		err := conn.Close()
//		if err != nil {
//			log.Printf("failed to close to server: %v", err)
//		}
//	}(conn)
//
//	wg.Add(len(users))
//
//	for _, user := range users {
//		go func() {
//			defer wg.Done()
//			//TODO: Get-обращение к auth-серверу
//		}
//	}
//
//	wg.Wait()
//
//}
