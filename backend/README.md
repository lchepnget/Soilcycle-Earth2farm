
# Backend structure 

```
                 SOILCYCLE DATABASE
                         |
          +--------------+--------------+
          |              |              |
        USERS           FARMS          CROPS
          |              |              |
          |              |              |
          |          FARM_CROPS         |
          |              |              |
          |         +----+-----+
          |         |          |
          |      SOIL_DATA  CLIMATE_DATA
          |         |          |
          +---------+----------+
                    |
               ROTATION_PLANS
                    |
                    |
              ROTATION_RESULTS

```
# Role collaboration`

```
                 SOILCYCLE BY EARTH2FARM
                           |
                 Backend Lead (Member 1)
                           |
          +----------------+----------------+
          |                |                |
      Frontend         NASA Data       Agricultural
      Developer        Engineer        Researcher
      Member 2         Member 3        Member 4
          |                |                |
          +----------------+----------------+
                           |
                   Integration & QA
                      Member 5
                           |
                 Complete Application
                           |
                    Final Testing
                           |
                      Deployment

```

# Planning and 
```
Farmer Input
     |
     v
Frontend
     |
     v
Go Backend
     |
     +------------------+
     |                  |
     v                  v
NASA Data API       Farm Database
     |                  |
     +------------------+
              |
              v
     Crop Rotation Engine
              |
              v
       Recommendation
              |
              v
         Dashboard

```
# Structure

```
backend/
├── cmd/
│   └── api/
│       └── main.go
│
├── internal/
│   ├── config/
│   │   └── config.go
│   │
│   ├── database/
│   │   └── postgres.go
│   │
│   ├── server/
│   │   └── routes.go
│   │
│   └── modules/
│       ├── farmers/
│       ├── farms/
│       ├── earthdata/
│       ├── crops/
│       ├── soil/
│       └── rotation/
│
├── .env
├── .env.example
├── go.mod
└── cmd/
```
# Crops done

```
internal/modules/
│
├── farmers/
│   ├── handler.go
│   ├── service.go
│   ├── repository.go
│   ├── model.go
│   ├── dto.go
│   └── routes.go
│
├── farms/
│   ├── handler.go
│   ├── service.go
│   ├── repository.go
│   ├── model.go
│   ├── dto.go
│   └── routes.go
│
└── crops/
    ├── handler.go
    ├── service.go
    ├── repository.go
    ├── model.go
    ├── dto.go
    └── routes.go

```