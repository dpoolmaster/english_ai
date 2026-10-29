package worker

import (
	"context"
	"log/slog"
)

// Worker is the background process lifecycle. River jobs are registered in a
// later slice; for now the process only owns start/stop semantics.
type Worker struct {
	log *slog.Logger
}

func New(log *slog.Logger) *Worker {
	return &Worker{log: log}
}

func (w *Worker) Run(ctx context.Context) error {
	w.log.Info("worker started")
	<-ctx.Done()
	w.log.Info("worker stopped")
	return nil
}
