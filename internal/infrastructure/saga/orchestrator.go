package saga

import (
	"context"
	"fmt"

	"github.com/Junaidmdv/goalcircle-team_service/pkg/logger"
)

type Orchestrator struct {
	steps  []SagaStep
	logger logger.Logger
}

func NewOrchestrator(steps []SagaStep, logger logger.Logger) *Orchestrator {
	return &Orchestrator{
		steps:  steps,
		logger: logger,
	}
}

func (o *Orchestrator) Execute(ctx context.Context, SagaState interface{}) error {
	completed := []SagaStep{}

	for _, step := range o.steps {
		o.logger.Info(fmt.Sprintf("[SAGA] executing step: %s", step.Name))

		if err := step.Action(ctx, SagaState); err != nil {
			o.logger.Error(fmt.Sprintf("[SAGA] step %s failed: %v — starting rollback", step.Name, err))
			o.compensate(ctx, SagaState, completed)
			return err
		}

		completed = append(completed, step)
		o.logger.Info(fmt.Sprintf("[SAGA] step %s succeeded", step.Name))
	}

	o.logger.Info("[SAGA] all steps completed successfully")
	return nil
}

func (o *Orchestrator) compensate(ctx context.Context, SagaState interface{}, completed []SagaStep) {
	for i := len(completed) - 1; i >= 0; i-- {
		step := completed[i]
		o.logger.Info(fmt.Sprintf("[SAGA] compensating: %s", step.Name))

		if err := step.Compensate(ctx, SagaState); err != nil {
			o.logger.Error(fmt.Sprintf("[SAGA] compensation failed for %s: %v", step.Name, err))
			// continue — best-effort, attempt all remaining compensations
		}
	}
}
