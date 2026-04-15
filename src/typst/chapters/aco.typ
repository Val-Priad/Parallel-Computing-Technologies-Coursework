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
        Call[Apply-Config-Defaults][cfg]

        If(`NOT Validate-Instance(instance)`, {
          Return[InfeasibleSolution()]
        })

        Comment[Initialization]
        Assign[`n`][length(instance.Dist)]
        Assign[`τ`][Matrix(n, n, cfg.InitialPheromone)]
        Assign[`rng`][Random(cfg.Seed)]

        Assign[`best_solution`][InfeasibleSolution()]

        For([`iter = 1` to `cfg.Iterations`], {
          Comment[Phase 1: solution construction]
          Assign[`ant_solutions`][[]]

          For([`k = 1` to `cfg.NumAnts`], {
            Assign[`candidate, feasible`][Build-Solution(instance, τ, cfg, rng)]

            Line[append `candidate, feasible` to `ant_solutions`]

            If(`feasible AND candidate.cost < best_solution.cost`, {
              Assign[`best_solution`][candidate]
            })
          })

          Comment[Phase 2: pheromone evaporation]
          Call[Evaporate][τ, cfg.Evaporation]

          Comment[Phase 3: pheromone update from iteration solutions]
          For([`each ant in ant_solutions`], {
            If(`ant.feasible AND ant.solution.cost > 0`, {
              Call[Deposit-Solution][
                τ,
                ant.solution,
                `cfg.Q / ant.solution.cost`
              ]
            })
          })

          Comment[Phase 4: elite reinforcement by global best]
          If(`best_solution.cost < +∞ AND best_solution.cost > 0`, {
            Call[Deposit-Solution][
              τ,
              best_solution,
              `cfg.EliteWeight * cfg.Q / best_solution.cost`
            ]
          })
        })

        Return[best_solution]
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
      ("instance", "τ", "cfg", "rng"),
      {
        Assign[`remaining`][length(instance.Customers)]
        Assign[`visited`][all false]
        Assign[`routes`][[]]

        Comment[Construct routes until all customers are served or no vehicle remains]
        While(`remaining > 0 AND length(routes) < instance.Vehicles`, {
          Assign[`route`][[]]
          Assign[`load`][0]
          Assign[`current`][0]

          Comment[Extend current route while feasible customers exist]
          While(`exists c in instance.Customers such that NOT visited[c.id] AND load + c.demand <= instance.capacity`, {
            Assign[`candidates`][[]]

            For([`each c in instance.Customers`], {
              If(`NOT visited[c.id] AND load + c.demand <= instance.capacity`, {
                Line[append `c.id` to `candidates`]
              })
            })

            Assign[`next`][
              Select-Next(
              current,
              candidates,
              instance.Dist,
              τ,
              cfg,
              rng
              )
            ]

            If(`next = -1`, {
              Break
            })

            Line[append `next` to `route`]

            Assign[`load`][`load + Demand(next)`]
            Assign[`current`][next]
            Assign[`visited[next]`][true]
            Assign[`remaining`][`remaining - 1`]
          })

          If(`route != []`, {
            Line[append `route` to `routes`]
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
      ("current", "candidates", "dist", "τ", "cfg", "rng"),
      {
        If(`length(candidates) = 0`, {
          Return[-1]
        })

        If(`length(candidates) = 1`, {
          Return[candidates[0]]
        })

        Comment[Attraction model: $w_j = τ_(i,j)^α · η_(i,j)^β$, with $η_(i,j) = 1 / d_(i,j)$]
        Assign[`weights`][[]]
        Assign[`total`][0]

        For([`each j in candidates`], {
          Assign[`d_(i,j)`][dist[current][j]]

          If(`d_(i,j) <= 0`, {
            Assign[`d_(i,j)`][ε]
          })

          Assign[`τ_(i,j)`][`τ[current][j]`]
          Assign[`η_(i,j)`][`1 / d_(i,j)`]
          Assign[`w`][`(τ_(i,j))^cfg.Alpha * (η_(i,j))^cfg.Beta`]

          If(`w <= 0`, {
            Assign[`w`][ε]
          })

          Line[append `w` to `weights`]
          Assign[`total`][`total + w`]
        })

        If(`total <= 0`, {
          Return[Nearest-Neighbor(current, candidates)]
        })

        Comment[Roulette-wheel sampling]

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
        Comment[Uniform evaporation of all trail values]
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
        Comment[Reinforce each traversed edge symmetrically]

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
