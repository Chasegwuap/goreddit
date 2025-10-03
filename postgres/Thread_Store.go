package postgres

import (
	"fmt"

	"github.com/Chasegwuap/goreddit"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ThreadStore struct {
	*sqlx.DB
}

func (s *ThreadStore) Thread(id uuid.UUID) (goreddit.Thread, error) {
	var t goreddit.Thread
	if err := s.Get(&t, "SELECT * FROM threads WHERE id=$1", id); err != nil {
		return goreddit.Thread{}, fmt.Errorf("getting thread by id: %w", err)
	}
	return t, nil
}

func (s *ThreadStore) Threads() ([]*goreddit.Thread, error) {
	var tt []goreddit.Thread
	if err := s.Select(&tt, "SELECT * FROM threads"); err != nil {
		return []*goreddit.Thread{}, fmt.Errorf("getting all threads: %w", err)
	}

	// Convert []goreddit.Thread -> []*goreddit.Thread
	threadPtrs := make([]*goreddit.Thread, len(tt))
	for i := range tt {
		threadPtrs[i] = &tt[i]
	}
	return threadPtrs, nil
}

func (s *ThreadStore) CreateThread(t *goreddit.Thread) error {
	if err := s.Get(t, "INSERT INTO threads VALUES ($1, $2, $3) RETURNING *",
		t.ID,
		t.Title,
		t.Description); err != nil {
		return fmt.Errorf("inserting thread: %w", err)
	}
	return nil
}

func (s *ThreadStore) UpdateThread(t *goreddit.Thread) error {
	if err := s.Get(t, "UPDATE threads SET title=$2, description=$3 WHERE id=$1 RETURNING *",
		t.ID,
		t.Title,
		t.Description); err != nil {
		return fmt.Errorf("error updating thread: %w", err)
	}
	return nil
}

func (s *ThreadStore) DeleteThread(id uuid.UUID) error {
	if _, err := s.Exec("DELETE FROM threads WHERE id=$1", id); err != nil {
		return fmt.Errorf("deleting thread: %w", err)
	}
	return nil
}
