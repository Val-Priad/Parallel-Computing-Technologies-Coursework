#import "@preview/algorithmic:1.0.7"
#import algorithmic: algorithm-figure, style-algorithm
#show: style-algorithm

#pagebreak()

#counter(figure.where(kind: "algorithm")).update(0)


#algorithm-figure(
  supplement: "PACO",
  "Parallel Ant Colony Optimization",
  {
    import algorithmic: *

    Procedure(
      "Solve-PACO",
      ("instance", "cfg"),
      {
        If(`length(instance.Customers) = 0 OR instance.Vehicles = 0`, {
          Return[EmptySolution()]
        })

        Call[Apply-Config-Defaults][cfg.BaseConfig]

        If(`NOT Validate-Instance(instance)`, {
          Return[InfeasibleSolution()]
        })

        Comment[Normalize worker count]
        Assign[`cfg.NumWorkers`][Resolve-PACO-Workers(cfg.NumWorkers, cfg.BaseConfig.NumAnts)]

        Comment[Split ants between workers]
        Assign[`ants_per_worker`][Split-Ants(cfg.BaseConfig.NumAnts, cfg.NumWorkers)]
        Assign[`results`][Array(cfg.NumWorkers)]

        Comment[Parallel execution (parallel for loop) in implementation: loop body runs concurrently for workers]
        For([`w = 0` to `cfg.NumWorkers - 1`], {
          Assign[`local_cfg`][cfg.BaseConfig]
          Assign[`local_cfg.Seed`][`cfg.BaseConfig.Seed + w * 1000`]
          Assign[`local_cfg.NumAnts`][`ants_per_worker[w]`]

          Assign[`rng`][Random(local_cfg.Seed)]
          Assign[`n`][length(instance.Dist)]
          Assign[`τ`][Matrix(n, n, local_cfg.InitialPheromone)]

          Assign[`best`][InfeasibleSolution()]

          For([`iter = 1` to `local_cfg.Iterations`], {
            Assign[`ants`][[]]

            Comment[Construct solutions]
            For([`k = 1` to `local_cfg.NumAnts`], {
              Assign[`sol, feasible`][Build-Solution(instance, τ, local_cfg, rng)]

              Line[append `(sol, feasible)` to `ants`]

              If(`feasible AND sol.cost < best.cost`, {
                Assign[`best`][Clone(sol)]
              })
            })

            Comment[Evaporation]
            Call[Evaporate][τ, local_cfg.Evaporation]

            Comment[Update from ants]
            For([`each ant in ants`], {
              If(`ant.feasible AND ant.solution.cost > 0`, {
                Call[Deposit-Solution][
                  τ,
                  ant.solution,
                  `local_cfg.Q / ant.solution.cost`
                ]
              })
            })

            Comment[Elite update]
            If(`best.cost < +∞ AND best.cost > 0`, {
              Call[Deposit-Solution][
                τ,
                best,
                `local_cfg.EliteWeight * local_cfg.Q / best.cost`
              ]
            })
          })

          Assign[`results[w]`][best]
        })

        Comment[Global reduction]
        Assign[`global_best`][InfeasibleSolution()]

        For([`each r in results`], {
          If(`r.cost < global_best.cost`, {
            Assign[`global_best`][Clone(r)]
          })
        })

        If(`global_best.cost = +∞`, {
          Return[InfeasibleSolution()]
        })

        Return[global_best]
      },
    )
  },
)

#pagebreak()

#algorithm-figure(
  supplement: "PACO",
  "Ant Distribution",
  {
    import algorithmic: *

    Function(
      "Split-Ants",
      ("total", "workers"),
      {
        Assign[`result`][Array(workers)]

        Assign[`base`][`total / workers`]
        Assign[`rem`][`total mod workers`]

        For([`i = 0` to `workers - 1`], {
          Assign[`result[i]`][base]

          If(`i < rem`, {
            Assign[`result[i]`][`result[i] + 1`]
          })
        })

        Return[result]
      },
    )
  },
)
