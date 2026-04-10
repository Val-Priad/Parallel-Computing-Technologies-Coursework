
= Структури даних

*Point*:
- id: Integer
- x: Float
- y: Float
- demand: Integer

*VRPInstance*:
- depot: Point
- customers: Array\<Point>
- vehicles: Integer
- vehicleCapacity: Integer
- capacityMode: CapacityMode
- dist: Matrix\<Float>

*Route*:
- vehicleId: Integer
- nodes: Array\<Integer>

*Solution*:
- routes: Array\<Route>
- cost: Float
- durationMS: Float
