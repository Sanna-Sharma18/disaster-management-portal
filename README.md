# Relief Atlas

A disaster management dashboard for tracking active disasters, affected areas, shelters, aid distribution, and donations.

## Features

- **Live disaster tracking** — monitor active disasters and their status
- **Shelter management** — view capacity and occupancy across shelters
- **Aid distribution** — log and track material distribution by area
- **Donations** — record and display donor contributions
- **Admin portal** — secure login for administrators
- **Analytics** — charts and reports on relief operations

## Tech Stack

| Layer    | Technology          |
|----------|---------------------|
| Frontend | HTML, CSS, JavaScript |
| Charts   | Chart.js            |
| Icons    | Phosphor Icons      |
| Database | MySQL               |

## Database Schema

The schema (`schema.sql`) defines six tables:

- `Disaster` — disaster events and their status
- `Affected_Areas` — areas impacted by each disaster
- `Shelter` — shelters linked to affected areas
- `Admins` — administrator accounts
- `Distribution` — aid distribution records per area
- `Users` & `Donations` — donor accounts and contribution history

See `erd.jpeg` for the full entity-relationship diagram.

## Setup

### Database

1. Create the database and tables:
   ```bash
   mysql -u root -p < schema.sql
   ```

### Frontend

No build step required. Open `login.html` in a browser to start, or serve with any static file server:

```bash
# Python
python -m http.server 8080

# Node.js (npx)
npx serve .
```

Then navigate to `http://localhost:8080/login.html`.

## Project Structure

```
relief-atlas/
├── index.html       # Main dashboard
├── styles.css       # Dashboard styles
├── script.js        # Dashboard logic
├── login.html       # Login page
├── login.css        # Login styles
├── login.js         # Auth logic
├── hero-bg.png      # Hero background image
├── schema.sql       # MySQL schema
└── erd.jpeg         # Entity-relationship diagram
```
