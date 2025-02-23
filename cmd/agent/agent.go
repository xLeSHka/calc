package main

import (
	"github.com/xLeSHka/calc/app/agent"
	agent2 "github.com/xLeSHka/calc/internal/agent"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		agent.Agent,
		fx.Invoke(func(*agent2.Agent) {}),
	).Run()
}
