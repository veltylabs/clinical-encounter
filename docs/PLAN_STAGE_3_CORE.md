# Stage 3 — Module Core (`mcp.go`)

← [Stage 2](PLAN_STAGE_2_ORM.md) | Next → [Stage 4](PLAN_STAGE_4_VISIT.md)

## Code

```go
//go:build !wasm

package patientvisit

import (
    "github.com/tinywasm/orm"
    "github.com/tinywasm/unixid"
)

// EventPublisher is implemented by the host application (SSE, websocket, etc.).
// Pass nil to disable event publishing (no-op).
type EventPublisher interface {
    Publish(event string, payload any) error
}

type Module struct {
    db  *orm.DB
    uid *unixid.UnixID
    pub EventPublisher
}

func New(db *orm.DB, pub EventPublisher) (*Module, error) {
    u, err := unixid.NewUnixID()
    if err != nil {
        return nil, err
    }
    return &Module{db: db, uid: u, pub: pub}, nil
}

// publish fires an event if a publisher is configured.
func (m *Module) publish(event string, payload any) {
    if m.pub != nil {
        m.pub.Publish(event, payload)
    }
}
```
