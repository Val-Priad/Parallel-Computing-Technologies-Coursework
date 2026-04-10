
#import "@preview/algorithmic:1.0.7"
#import algorithmic: algorithm-figure, style-algorithm
#show: style-algorithm

#counter(figure.where(kind: "algorithm")).update(0)

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

        Comment[Initialize pheromone trails and the best solution tracker.]
        Assign[`n`][length(instance.Dist)]
        Assign[`pheromone`][Matrix(n, n, cfg.InitialPheromone)]
        Assign[`rng`][Random(cfg.Seed)]

        Assign[`best.cost`][+∞]
        Assign[`best.routes`][[]]

        For([`iter = 1` to `cfg.Iterations`], {
          Comment[Build one candidate solution for each ant in the colony.]
          Assign[`ants`][[]]

          For([`k = 1` to `cfg.NumAnts`], {
            Assign[`sol, feasible`][Build-Solution(instance, pheromone, cfg, rng)]

            Line[Append `sol, feasible` to `ants`]

            If(`feasible AND sol.cost < best.cost`, {
              Assign[`best`][sol]
            })
          })

          Call[Evaporate][pheromone, cfg.Evaporation]

          For([`each ant in ants`], {
            If(`ant.feasible AND ant.solution.cost > 0`, {
              Call[Deposit-Solution][
                pheromone,
                ant.solution,
                `cfg.Q / ant.solution.cost`
              ]
            })
          })

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

        Comment[Start a new vehicle route and keep adding feasible customers.]
        While(`remaining > 0 AND length(routes) < instance.Vehicles`, {
          Assign[`route`][[]]
          Assign[`load`][0]
          Assign[`current`][0]

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

            Assign[`load`][`load + demand(next)`]
            Assign[`current`][next]
            Assign[`visited[next]`][true]
            Assign[`remaining`][`remaining - 1`]
          })

          Line[Append `route` to `routes`]
        })

        If(`remaining > 0`, {
          Return[InfeasibleSolution(), false]
        })

        Assign[`cost`][Compute-Total-Cost(routes)]

        Return[`routes, cost`, true]
      },
    )
  },
)
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

        Comment[Compute attraction weights from pheromone and distance.]
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
        Comment[Use roulette-wheel sampling to pick the next customer.]

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

#algorithm-figure(
  supplement: "ACO",
  "Pheromone Evaporation",
  {
    import algorithmic: *

    Procedure(
      "Evaporate",
      ("pheromone", "rate"),
      {
        Comment[Decay all pheromone trails to reduce the influence of older solutions.]
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

#algorithm-figure(
  supplement: "ACO",
  "Pheromone Update",
  {
    import algorithmic: *

    Procedure(
      "Deposit-Solution",
      ("pheromone", "solution", "amount"),
      {
        Comment[Reinforce the edges used by the selected solution.]
        If(`amount <= 0`, {
          Return[]
        })

        For([`each route in solution.routes`], {
          If(`route != []`, {
            Assign[`prev`][0]

            For([`each node in route`], {
              Assign[`pheromone[prev][node]`][`+ amount`]
              Assign[`pheromone[node][prev]`][`+ amount`]
              Assign[`prev`][node]
            })

            Assign[`pheromone[prev][0]`][`+ amount`]
            Assign[`pheromone[0][prev]`][`+ amount`]
          })
        })
      },
    )
  },
)
