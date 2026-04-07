package solver

import (
	"math"
	"math/rand"
	"parallel-aco/internal/logging"
	"parallel-aco/internal/vrp"
	"time"
)

type PACOConfig struct {
	Colonies      int
	ExchangeEvery int
	BaseACOConfig ACOConfig
}

type exchangeMsg struct {
	From int
	Best vrp.Solution
}

type colonyResult struct {
	sol vrp.Solution
}

func SolvePACO(
	instance vrp.VRPInstance,
	logger *logging.Logger,
	pacoCfg PACOConfig,
) vrp.Solution {
	_ = logger

	startTime := time.Now()

	if len(instance.Customers) == 0 || instance.Vehicles == 0 {
		return vrp.Solution{
			Routes: []vrp.Route{},
			Cost:   0,
			Metrics: vrp.SearchMetrics{
				DurationMS: float64(time.Since(startTime).Nanoseconds()) / 1e6,
			},
		}
	}

	applyConfigDefaults(&pacoCfg.BaseACOConfig)
	if !validateInstance(instance) {
		return solutionWithMetrics([]vrp.Route{}, math.Inf(1), startTime)
	}

	if pacoCfg.Colonies <= 0 {
		pacoCfg.Colonies = 4
	}
	if pacoCfg.ExchangeEvery <= 0 {
		pacoCfg.ExchangeEvery = 25
	}

	updatesCh := make(chan exchangeMsg, pacoCfg.Colonies*2)
	inboxes := make([]chan exchangeMsg, pacoCfg.Colonies)
	for i := range inboxes {
		inboxes[i] = make(chan exchangeMsg, 1)
	}

	doneBroker := make(chan struct{})
	go brokerExchanges(updatesCh, inboxes, doneBroker)

	resultCh := make(chan colonyResult, pacoCfg.Colonies)

	for colonyID := 0; colonyID < pacoCfg.Colonies; colonyID++ {
		go runColony(
			colonyID,
			instance,
			pacoCfg,
			updatesCh,
			inboxes[colonyID],
			resultCh,
		)
	}

	best := vrp.Solution{Cost: math.Inf(1)}
	for i := 0; i < pacoCfg.Colonies; i++ {
		res := <-resultCh
		if res.sol.Cost < best.Cost {
			best = res.sol
		}
	}

	close(doneBroker)

	if math.IsInf(best.Cost, 1) {
		return solutionWithMetrics([]vrp.Route{}, math.Inf(1), startTime)
	}

	return solutionWithMetrics(best.Routes, best.Cost, startTime)
}

func brokerExchanges(updatesCh <-chan exchangeMsg, inboxes []chan exchangeMsg, done <-chan struct{}) {
	best := vrp.Solution{Cost: math.Inf(1)}

	for {
		select {
		case <-done:
			return
		case msg := <-updatesCh:
			if math.IsInf(msg.Best.Cost, 1) || msg.Best.Cost >= best.Cost {
				continue
			}

			best = cloneSolution(msg.Best)
			bestMsg := exchangeMsg{From: msg.From, Best: best}
			for colonyID, inbox := range inboxes {
				if colonyID == msg.From {
					continue
				}
				select {
				case inbox <- bestMsg:
				default:
				}
			}
		}
	}
}

func runColony(
	id int,
	instance vrp.VRPInstance,
	pacoCfg PACOConfig,
	updatesCh chan<- exchangeMsg,
	inbox <-chan exchangeMsg,
	resultCh chan<- colonyResult,
) {
	cfg := pacoCfg.BaseACOConfig
	cfg.Seed += int64(id * 1000)

	n := len(instance.Dist)
	pheromone := makeMatrix(n, n, cfg.InitialPheromone)
	rng := rand.New(rand.NewSource(cfg.Seed))

	best := vrp.Solution{Cost: math.Inf(1)}

	for iter := 0; iter < cfg.Iterations; iter++ {
		ants := make([]antSolution, 0, cfg.NumAnts)

		for ant := 0; ant < cfg.NumAnts; ant++ {
			sol, feasible := buildSolution(instance, pheromone, cfg, rng)
			if feasible && sol.Cost < best.Cost {
				best = cloneSolution(sol)
			}
			ants = append(ants, antSolution{Solution: sol, Feasible: feasible})
		}

		evaporate(pheromone, cfg.Evaporation)

		for _, ant := range ants {
			if !ant.Feasible || ant.Solution.Cost <= 0 {
				continue
			}
			depositSolution(pheromone, ant.Solution, cfg.Q/ant.Solution.Cost)
		}

		if !math.IsInf(best.Cost, 1) && best.Cost > 0 {
			depositSolution(pheromone, best, cfg.EliteWeight*cfg.Q/best.Cost)
		}

		if iter%pacoCfg.ExchangeEvery == 0 && !math.IsInf(best.Cost, 1) {
			select {
			case updatesCh <- exchangeMsg{From: id, Best: cloneSolution(best)}:
			default:
			}

			select {
			case msg := <-inbox:
				if msg.Best.Cost < best.Cost {
					depositSolution(pheromone, msg.Best, cfg.Q/msg.Best.Cost)
					best = cloneSolution(msg.Best)
				}
			default:
			}
		}
	}

	resultCh <- colonyResult{sol: best}
}
