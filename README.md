# Poker Odds Engine

A high-performance poker odds calculator built with **Go** and **Gin framework**. Calculates Texas Hold'em winning probabilities using Monte Carlo simulation with concurrent goroutines.

## Tech Stack

- **Go 1.21**
- **Gin Framework**
- **Docker**

## Performance

- **10-50x faster** than Python implementations
- **Hand Evaluation**: <1ms per hand
- **Odds Calculation**: 50-100ms for 10,000 simulations
- **Throughput**: 100+ concurrent requests/second
- **True Parallelism**: No GIL limitations

## Quick Start

### Run Locally

```bash
# Clone repository
git clone https://github.com/KyleKDang/poker-odds-engine.git
cd poker-odds-engine

# Download dependencies
go mod tidy

# Run server
go run cmd/server/main.go
```

Server starts on `http://localhost:8001`

### Run with Docker

```bash
# Build and run
docker-compose up -d

# View logs
docker-compose logs -f

# Stop
docker-compose down
```

### Using Makefile

```bash
make run          # Run locally
make docker-run   # Run with Docker
make docker-logs  # View logs
make docker-stop  # Stop containers
```

## API Documentation

### Health Check

```http
GET /health
```

**Response:**
```json
{
  "status": "ok",
  "service": "poker-odds-engine"
}
```

### Evaluate Hand

Evaluates the best 5-card poker hand.

```http
POST /evaluate
Content-Type: application/json
```

**Request:**
```json
{
  "hole_cards": ["AS", "KH"],
  "board_cards": ["QS", "JS", "TS", "2D", "3C"]
}
```

**Response:**
```json
{
  "hand": "Straight",
  "rank": 5
}
```

### Calculate Odds

Calculates winning probability via Monte Carlo simulation.

```http
POST /odds
Content-Type: application/json
```

**Request:**
```json
{
  "hole_cards": ["AS", "AH"],
  "board_cards": [],
  "num_opponents": 1,
  "simulations": 10000,
  "workers": 4
}
```

**Parameters:**
- `hole_cards` (required): Array of 2 cards
- `board_cards` (required): Array of 0-5 cards
- `num_opponents` (required): Number of opponents (1-9)
- `simulations` (optional): Number of simulations (default: 10000)
- `workers` (optional): Number of parallel workers (default: 4)

**Response:**
```json
{
  "win": 0.8523,
  "tie": 0.0077,
  "loss": 0.1400
}
```

## Usage Examples

### cURL

```bash
# Evaluate hand
curl -X POST http://localhost:8001/evaluate \
  -H "Content-Type: application/json" \
  -d '{"hole_cards":["AS","KH"],"board_cards":["QS","JS","TS"]}'

# Calculate odds
curl -X POST http://localhost:8001/odds \
  -H "Content-Type: application/json" \
  -d '{"hole_cards":["AS","AH"],"board_cards":[],"num_opponents":1}'
```

### JavaScript

```javascript
const response = await fetch('http://localhost:8001/odds', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    hole_cards: ['AS', 'AH'],
    board_cards: [],
    num_opponents: 1
  })
});
const data = await response.json();
console.log(`Win: ${(data.win * 100).toFixed(2)}%`);
```

## Project Structure

```
poker-odds-engine/
├── cmd/
│   └── server/          # Main application entry point
├── internal/            # Private application code
│   ├── api/             # Gin handlers and routing
│   ├── card/            # Card model and deck operations
│   ├── evaluator/       # Hand evaluation logic
│   └── simulator/       # Monte Carlo simulation
├── pkg/
│   └── models/          # API request/response models
├── Dockerfile
├── docker-compose.yml
├── Makefile
└── go.mod
```

## Card Format

Cards are 2-character strings: `[Rank][Suit]`

**Ranks:** `2-9`, `T` (Ten), `J` (Jack), `Q` (Queen), `K` (King), `A` (Ace)  
**Suits:** `S` (Spades), `H` (Hearts), `D` (Diamonds), `C` (Clubs)

**Examples:** `AS` (Ace of Spades), `KH` (King of Hearts), `TC` (Ten of Clubs)

## Configuration

Environment variables:

- `PORT` - Server port (default: 8001)
- `GIN_MODE` - Gin mode: `debug` or `release` (default: debug)

## Development

```bash
# Run locally
make run

# Format code
make fmt

# Build binary
make build

# Run tests (when added)
make test
```

## Performance Tuning

### Simulation Count

- **1,000 sims**: Fast (~10ms), less accurate
- **10,000 sims**: Balanced (default, ~100ms)
- **100,000 sims**: Slow (~1s), very accurate

### Worker Count

Optimal: Number of CPU cores (typically 4-8). More workers = faster parallel processing.

## Integration with FastAPI

This service is designed as a microservice to complement a FastAPI backend:

```python
# In your FastAPI routes
import httpx

GO_SERVICE = "http://localhost:8001"

@router.post("/odds")
async def calculate_odds(request: OddsRequest):
    async with httpx.AsyncClient() as client:
        response = await client.post(
            f"{GO_SERVICE}/odds",
            json=request.dict()
        )
        return response.json()
```
