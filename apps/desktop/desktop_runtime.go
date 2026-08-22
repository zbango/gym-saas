package main

import "context"

// DesktopRuntime owns process-level state shared by Wails delivery adapters.
// It deliberately contains no feature methods or business logic.
type DesktopRuntime struct {
	ctx context.Context
}

func (r *DesktopRuntime) startup(ctx context.Context) {
	r.ctx = ctx
}

func (r *DesktopRuntime) requestContext() context.Context {
	if r.ctx != nil {
		return r.ctx
	}
	return context.Background()
}
