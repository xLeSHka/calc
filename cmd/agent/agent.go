package main

import (
	agent2 "github.com/xLeSHka/calc/internal/agent"
	"github.com/xLeSHka/calc/internal/app/agent"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		agent.Agent,
		fx.Invoke(func(*agent2.Agent) {}),
	).Run()
}
