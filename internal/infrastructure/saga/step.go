package saga

import "context"

type SagaStep struct {
	Name       string
	Action     func(ctx context.Context,sagaState interface{}) error
	Compensate func(ctx context.Context,sagaState interface{}) error
}


