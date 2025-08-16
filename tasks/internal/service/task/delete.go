package task

import "context"

func (s *Serv) Delete(ctx context.Context, id int64) error {
	const op = "service.task.Delete"

	err := s.taskRepo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
