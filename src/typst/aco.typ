#import "@preview/algorithmic:1.0.7"
#import algorithmic: algorithm-figure, style-algorithm
#show: style-algorithm

#counter(figure.where(kind: "algorithm")).update(0)

#pagebreak()

#algorithm-figure(
  supplement: "ACO",
  "Ant Colony Optimization",
  {
    import algorithmic: *

    Procedure(
      "Solve-ACO",
      ("instance", "cfg"),
      {
        If(`instance.Customers = [] OR instance.Vehicles = 0`, {
          Return[EmptySolution()]
        })

        Call[Apply-Config-Defaults][cfg]

        If(`NOT Validate-Instance(instance)`, {
          Return[InfeasibleSolution()]
        })

        Comment[Initialize search state]
        Assign[`n`][length(instance.Dist)]
        Assign[`pheromone`][Matrix(n, n, cfg.InitialPheromone)]
        Assign[`rng`][Random(cfg.Seed)]

        Assign[`best.cost`][+∞]
        Assign[`best.routes`][[]]

        For([`iter = 1` to `cfg.Iterations`], {
          Comment[Phase 1: construct one candidate per ant]
          Assign[`ants`][[]]

          For([`k = 1` to `cfg.NumAnts`], {
            Assign[`sol, feasible`][Build-Solution(instance, pheromone, cfg, rng)]

            Line[Append `sol, feasible` to `ants`]

            If(`feasible AND sol.cost < best.cost`, {
              Assign[`best`][sol]
            })
          })

          Comment[Phase 2: evaporate stale pheromone]
          Call[Evaporate][pheromone, cfg.Evaporation]

          Comment[Phase 3: reinforce feasible ant solutions]
          For([`each ant in ants`], {
            If(`ant.feasible AND ant.solution.cost > 0`, {
              Call[Deposit-Solution][
                pheromone,
                ant.solution,
                `cfg.Q / ant.solution.cost`
              ]
            })
          })

          Comment[Phase 4: apply elite reinforcement from global best]
          If(`best.cost < +∞ AND best.cost > 0`, {
            Call[Deposit-Solution][
              pheromone,
              best,
              `cfg.EliteWeight * cfg.Q / best.cost`
            ]
          })
        })

        Return[best]
      },
    )
  },
)

#pagebreak()

#algorithm-figure(
  supplement: "ACO",
  "Solution Construction",
  {
    import algorithmic: *

    Procedure(
      "Build-Solution",
      ("instance", "pheromone", "cfg", "rng"),
      {
        Assign[`remaining`][length(instance.Customers)]
        Assign[`visited`][all false]
        Assign[`routes`][[]]

        Comment[Construct routes until all customers are served or vehicles are exhausted]
        While(`remaining > 0 AND length(routes) < instance.Vehicles`, {
          Assign[`route`][[]]
          Assign[`load`][0]
          Assign[`current`][0]

          Comment[Greedily extend current route using probabilistic ACO choice]
          While(`true`, {
            Assign[`candidates`][[]]

            For([`each c in instance.Customers`], {
              If(`NOT visited[c.id] AND load + c.demand <= instance.capacity`, {
                Line[Append `c.id` to `candidates`]
              })
            })

            If(`candidates = []`, {
              Break
            })

            Assign[`next`][
              Select-Next(
              current,
              candidates,
              instance.Dist,
              pheromone,
              cfg,
              rng
              )
            ]

            If(`next = -1`, {
              Break
            })

            Line[Append `next` to `route`]

            Assign[`load`][`load + Demand(next)`]
            Assign[`current`][next]
            Assign[`visited[next]`][true]
            Assign[`remaining`][`remaining - 1`]
          })

          If(`route != []`, {
            Line[Append `route` to `routes`]
          })
        })

        If(`remaining > 0`, {
          Return[InfeasibleSolution(), false]
        })

        Assign[`cost`][Compute-Total-Cost(routes)]

        Assign[`solution.routes`][routes]
        Assign[`solution.cost`][cost]

        Return[`solution`, true]
      },
    )
  },
)

#pagebreak()

#algorithm-figure(
  supplement: "ACO",
  "Next Customer Selection",
  {
    import algorithmic: *

    Function(
      "Select-Next",
      ("current", "candidates", "dist", "pheromone", "cfg", "rng"),
      {
        If(`length(candidates) = 0`, {
          Return[-1]
        })

        If(`length(candidates) = 1`, {
          Return[candidates[0]]
        })

        Comment[Attraction model: $w_j = tau_(i,j)^alpha dot eta_(i,j)^beta$]
        Assign[`weights`][[]]
        Assign[`total`][0]

        For([`each j in candidates`], {
          Assign[`d`][dist[current][j]]

          If(`d <= 0`, {
            Assign[`d`][ε]
          })

          Assign[`w`][`(pheromone[current][j])^cfg.Alpha * (1 / d)^cfg.Beta`]

          If(`w <= 0`, {
            Assign[`w`][ε]
          })

          Line[Append `w` to `weights`]
          Assign[`total`][`total + w`]
        })

        If(`total <= 0`, {
          Return[Nearest-Neighbor(current, candidates)]
        })

        Comment[Sample next node with roulette-wheel selection]

        Assign[`r`][RandomFloat(0, total)]
        Assign[`acc`][0]

        For([`i = 0` to `length(candidates)-1`], {
          Assign[`acc`][`acc + weights[i]`]

          If(`r <= acc`, {
            Return[candidates[i]]
          })
        })

        Return[candidates[last]]
      },
    )
  },
)

#pagebreak()

#algorithm-figure(
  supplement: "ACO",
  "Pheromone Evaporation",
  {
    import algorithmic: *

    Procedure(
      "Evaporate",
      ("pheromone", "rate"),
      {
        Comment[Uniform decay to reduce influence of old paths]
        Assign[`n`][length(pheromone)]
        Assign[`factor`][`1 - rate`]

        For([`i = 0` to `n-1`], {
          For([`j = 0` to `n-1`], {
            If(`i = j`, {
              Assign[`pheromone[i][j]`][0]
            })
            Else({
              Assign[`pheromone[i][j]`][`pheromone[i][j] * factor`]
            })
          })
        })
      },
    )
  },
)

#pagebreak()

#algorithm-figure(
  supplement: "ACO",
  "Pheromone Update",
  {
    import algorithmic: *

    Procedure(
      "Deposit-Solution",
      ("pheromone", "solution", "amount"),
      {
        Comment[Reinforce each traversed edge in both directions]
        If(`amount <= 0`, {
          Return[]
        })

        For([`each route in solution.routes`], {
          If(`route != []`, {
            Assign[`prev`][0]

            For([`each node in route`], {
              Assign[`pheromone[prev][node]`][`pheromone[prev][node] + amount`]
              Assign[`pheromone[node][prev]`][`pheromone[node][prev] + amount`]
              Assign[`prev`][node]
            })

            Assign[`pheromone[prev][0]`][`pheromone[prev][0] + amount`]
            Assign[`pheromone[0][prev]`][`pheromone[0][prev] + amount`]
          })
        })
      },
    )
  },
)
