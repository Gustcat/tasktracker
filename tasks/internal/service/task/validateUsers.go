package task

import (
	"context"
	"errors"
	"github.com/Gustcat/task-server/internal/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type userResult struct {
	userType string // "author" или "operator"
	user     *model.User
	err      error
}

func (s *Serv) validateUsers(ctx context.Context, authorId int64, operatorId *int64) (*model.User, *model.User, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	resultCh := make(chan userResult, 2)

	findUser := func(userId int64, userType string) {
		user, err := s.authClient.GetUser(ctx, userId)
		if err != nil {
			st, ok := status.FromError(err)
			if ok && st.Code() == codes.NotFound {
				err = ErrUserNotFound
			}
		}
		resultCh <- userResult{userType: userType, user: user, err: err}
	}

	go findUser(authorId, "author")
	expected := 1

	if operatorId != nil {
		go findUser(*operatorId, "operator")
		expected += 1
	}

	var author *model.User
	var operator *model.User

	for i := 0; i < expected; i++ {
		result := <-resultCh
		if result.err != nil {
			cancel()
			return nil, nil, result.err
		}
		switch result.userType {
		case "author":
			author = result.user
		case "operator":
			operator = result.user
		}
	}

	return author, operator, nil
}
