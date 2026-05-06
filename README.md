# Relief Atlas

A disaster management system for tracking active disasters, affected areas, shelters, aid distribution, and donations.

## Features

- Live disaster tracking — monitor active disasters and their status
- Shelter management — view capacity and occupancy across shelters
- Aid distribution — log and track material distribution by area
- Donations — record and display donor contributions
- Admin portal — secure login for administrators
- Analytics — charts and reports on relief operations

## Tech Stack

| Layer    | Technology                      |
|----------|---------------------------------|
| Frontend | HTML, CSS, JavaScript           |
| Charts   | Chart.js                        |
| Icons    | Phosphor Icons                  |
| Backend  | Go 1.22 (net/http)              |
| Database | Oracle XE 21c (PL/SQL)          |
| Runtime  | Docker / Docker Compose         |

## Project Structure

```
relief-atlas/
├── backend/
│   ├── main.go                  # entry point + routing
│   ├── go.mod
│   ├── Dockerfile
│   ├── db/
│   │   └── connection.go        # Oracle connection with retry
│   ├── models/
│   │   └── models.go            # all data structs
│   └── handlers/
│       ├── handler.go           # shared helpers
│       ├── disaster.go
│       ├── area.go
│       ├── shelter.go
│       ├── admin.go
│       ├── distribution.go
│       ├── user.go
│       └── donation.go
├── db/
│   └── init/
│       ├── 01_schema.sql        # Oracle DDL (tables)
│       ├── 02_plsql.sql         # procedures, functions, triggers, views
│       └── 03_seed.sql          # sample data
├── docker-compose.yml
├── .env.example
├── schema.sql                   # original MySQL schema (reference)
├── erd.jpeg                     # entity-relationship diagram
├── index.html                   # dashboard
├── login.html                   # login page
└── styles.css / script.js / login.*
```

## Database Schema

Seven tables modelled from the ER diagram:

| Table           | Description                               |
|-----------------|-------------------------------------------|
| `Disaster`      | Disaster events and current status        |
| `Affected_Areas`| Areas impacted by each disaster           |
| `Shelter`       | Shelters linked to affected areas         |
| `Admins`        | Administrator accounts                    |
| `Distribution`  | Aid distribution records per area         |
| `Users`         | Registered donor/public accounts          |
| `Donations`     | Donation transactions by users            |

### PL/SQL Objects

| Object                    | Type      | Purpose                                       |
|---------------------------|-----------|-----------------------------------------------|
| `trg_shelter_capacity`    | Trigger   | Prevents occupied > capacity                  |
| `trg_distribution_date`   | Trigger   | Auto-sets distribution_date if NULL           |
| `sp_register_disaster`    | Procedure | Creates a disaster + first affected area atomically |
| `sp_distribute_aid`       | Procedure | Creates a distribution record with validation |
| `sp_update_occupancy`     | Procedure | Safely updates shelter occupancy with locking |
| `fn_total_donations`      | Function  | Returns total donation amount                 |
| `fn_user_donations`       | Function  | Returns total donations for a specific user   |
| `fn_shelter_available`    | Function  | Returns available spots in a shelter          |
| `v_disaster_summary`      | View      | Disaster + affected area count + population   |
| `v_shelter_status`        | View      | Shelter occupancy % with area context         |
| `v_distribution_log`      | View      | Distribution records with area and admin info |

## API Endpoints

Base URL: `http://localhost:8080`

### Disasters
| Method | Path                    | Description          |
|--------|-------------------------|----------------------|
| GET    | /api/disasters          | List all disasters   |
| POST   | /api/disasters          | Create disaster      |
| GET    | /api/disasters/{id}     | Get by ID            |
| PUT    | /api/disasters/{id}     | Update               |
| DELETE | /api/disasters/{id}     | Delete               |

### Affected Areas
| Method | Path              | Description                          |
|--------|-------------------|--------------------------------------|
| GET    | /api/areas        | List all (optional `?disaster_id=N`) |
| POST   | /api/areas        | Create                               |
| GET    | /api/areas/{id}   | Get by ID                            |
| PUT    | /api/areas/{id}   | Update                               |
| DELETE | /api/areas/{id}   | Delete                               |

### Shelters
| Method | Path                | Description                       |
|--------|---------------------|-----------------------------------|
| GET    | /api/shelters       | List all (optional `?area_id=N`)  |
| POST   | /api/shelters       | Create                            |
| GET    | /api/shelters/{id}  | Get by ID                         |
| PUT    | /api/shelters/{id}  | Update                            |
| DELETE | /api/shelters/{id}  | Delete                            |

### Admins
| Method | Path                   | Description            |
|--------|------------------------|------------------------|
| GET    | /api/admins            | List all admins        |
| POST   | /api/admins            | Register admin         |
| POST   | /api/admins/login      | Login (returns admin)  |
| GET    | /api/admins/{id}       | Get by ID              |
| PUT    | /api/admins/{id}       | Update                 |
| DELETE | /api/admins/{id}       | Delete                 |

### Distributions
| Method | Path                       | Description                        |
|--------|----------------------------|------------------------------------|
| GET    | /api/distributions         | List all (optional `?area_id=N`)   |
| POST   | /api/distributions         | Create                             |
| GET    | /api/distributions/{id}    | Get by ID                          |
| PUT    | /api/distributions/{id}    | Update                             |
| DELETE | /api/distributions/{id}    | Delete                             |

### Users
| Method | Path                | Description           |
|--------|---------------------|-----------------------|
| GET    | /api/users          | List all users        |
| POST   | /api/users          | Register user         |
| POST   | /api/users/login    | Login (returns user)  |
| GET    | /api/users/{id}     | Get by ID             |
| PUT    | /api/users/{id}     | Update                |
| DELETE | /api/users/{id}     | Delete                |

### Donations
| Method | Path                  | Description                       |
|--------|-----------------------|-----------------------------------|
| GET    | /api/donations        | List all (optional `?user_id=N`)  |
| POST   | /api/donations        | Create                            |
| GET    | /api/donations/{id}   | Get by ID                         |
| DELETE | /api/donations/{id}   | Delete                            |

## Quick Start

### Prerequisites

- [Docker Desktop](https://www.docker.com/products/docker-desktop/) installed and running

### 1. Configure environment

```bash
cp .env.example .env
# Edit .env if you want to change passwords
```

### 2. Start everything

```bash
docker compose up -d
```

Oracle XE takes **3-5 minutes** to initialise on first run. The backend waits automatically with retry logic. Watch progress with:

```bash
docker compose logs -f
```

### 3. Verify

```bash
# Health check
curl http://localhost:8080/api/disasters

# Create an admin (required before assigning distributions)
curl -X POST http://localhost:8080/api/admins \
  -H "Content-Type: application/json" \
  -d '{"admin_name":"John Smith","email":"john@relief.org","password":"Admin123!"}'

# Admin login
curl -X POST http://localhost:8080/api/admins/login \
  -H "Content-Type: application/json" \
  -d '{"email":"john@relief.org","password":"Admin123!"}'
```

### 4. Frontend

Open `login.html` in a browser, or serve statically:

```bash
python -m http.server 3000
# then visit http://localhost:3000/login.html
```

### Stop

```bash
docker compose down           # keep database volume
docker compose down -v        # also delete database data
```

## Development — run backend locally

```bash
# Start only the database
docker compose up oracle -d

cd backend
go mod tidy
go run .
```

Environment variables default to `localhost:1521/XEPDB1` — matches the Docker Oracle port.
