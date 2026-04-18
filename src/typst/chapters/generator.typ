#import "@preview/algorithmic:1.0.7"
#import algorithmic: algorithm-figure, style-algorithm
#show: style-algorithm

#algorithm-figure(
  supplement: "Generator",
  "VRP Instance Generation",
  {
    import algorithmic: *

    Procedure(
      "Generate-Instance",
      "cfg",
      {
        Call[Apply-Generator-Defaults][cfg]

        Comment[Initialization]
        Assign[`rng`][Random(cfg.seed)]

        Comment[Locate depot at map center]
        Assign[`depot.id`][0]
        Assign[`depot.x`][`cfg.Width / 2`]
        Assign[`depot.y`][`cfg.Height / 2`]
        Assign[`depot.demand`][0]

        Assign[`points`][[depot]]

        Assign[`totalDemand`][0]
        Assign[`maxDemand`][0]

        Comment[Generate customers and update demand statistics]
        For([`i = 1` to `cfg.NumCustomers`], {
          Assign[`d`][RandomInt(cfg.MinDemand, cfg.MaxDemand)]
          Assign[`x`][RandomFloat(0, cfg.Width)]
          Assign[`y`][RandomFloat(0, cfg.Height)]

          Assign[`totalDemand`][`totalDemand + d`]

          If(`d > maxDemand`, {
            Assign[`maxDemand`][d]
          })

          Line[append customer $(i, x, y, d)$ to `points`]
        })

        Comment[Compute vehicle capacity under selected policy]
        Assign[`capacity`][Compute-Capacity(cfg, totalDemand, maxDemand)]

        Comment[Construct pairwise distance matrix]
        Assign[`dist`][Build-Distance-Matrix(points)]

        Comment[Assemble VRP instance]
        Assign[`instance.depot`][`points[0]`]
        Assign[`instance.customers`][`points[1..]`]
        Assign[`instance.vehicles`][`cfg.Vehicles`]
        Assign[`instance.capacity`][`capacity`]
        Assign[`instance.mode`][`cfg.CapacityMode`]
        Assign[`instance.dist`][`dist`]

        Return[`points, instance`]
      },
    )
  },
)

#pagebreak()

#algorithm-figure(
  supplement: "Generator",
  "Vehicle Capacity Computation",
  {
    import algorithmic: *

    Procedure(
      "Compute-Capacity",
      ("cfg", "totalDemand", "maxDemand"),
      {
        Comment[Degenerate cases]
        If(`totalDemand <= 0`, {
          Return[1]
        })

        If(`cfg.Vehicles <= 0`, {
          Return[max(maxDemand, 1)]
        })

        Assign[`required`][ceil(totalDemand / cfg.Vehicles)]

        Comment[Policy-dependent capacity proposal]

        IfElseChain(
          `cfg.CapacityMode = TIGHT`,
          {
            Assign[`capacity`][required]
          },
          [`cfg.CapacityMode = LOOSE`],
          {
            Assign[`capacity`][max(totalDemand, maxDemand)]
          },
          [`cfg.CapacityMode = FIXED`],
          {
            If(`cfg.FixedCapacity <= 0`, {
              Return[max(maxDemand, 1)]
            })
            Return[max(cfg.FixedCapacity, maxDemand)]
          },
          {
            Assign[`capacity`][ceil(totalDemand \* cfg.CapacitySlack / cfg.Vehicles)]
          },
        )

        Comment[Feasibility guards]

        If(`capacity < required`, {
          Assign[`capacity`][required]
        })

        If(`capacity < maxDemand`, {
          Assign[`capacity`][maxDemand]
        })

        Return[max(capacity, 1)]
      },
    )
  },
)

#pagebreak()

#algorithm-figure(
  supplement: "Generator",
  "Distance Matrix Construction",
  {
    import algorithmic: *

    Procedure(
      "Build-Distance-Matrix",
      "points",
      {
        Assign[`n`][length(points)]
        Assign[`dist`][Matrix(n, n, 0)]

        For([`i = 0` to `n-1`], {
          For([`j = 0` to `n-1`], {
            Assign[`dist[i][j]`][Distance(points[i], points[j])]
          })
        })

        Return[`dist`]
      },
    )
  },
)

#pagebreak()

#algorithm-figure(
  supplement: "Generator",
  "Euclidean Distance",
  {
    import algorithmic: *

    Function(
      "Distance",
      ("a", "b"),
      {
        Comment[Euclidean metric in the plane]
        Return[`sqrt((a.x - b.x)^2 + (a.y - b.y)^2)`]
      },
    )
  },
)

#pagebreak()
