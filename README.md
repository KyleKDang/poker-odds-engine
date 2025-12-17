# Poker Odds Engine

A high-performance poker odds calculator built with **Go** and **Gin framework**. Calculates Texas Hold'em winning probabilities using Monte Carlo simulation with concurrent goroutines.

## Tech Stack

- **Go 1.21**
- **Gin Framework**
- **Docker**

## Performance

Benchmarked on MacBook Air M-series (8 cores), all times averaged over 5 trials:

### Summary
- **3-4x faster** than Python FastAPI implementation
- **Hand Evaluation**: <1ms per hand
- **Default (10k sims)**: ~172ms average
- **High accuracy (100k sims)**: ~1.63s average
- **Goroutine speedup**: 3.1x with 8 workers vs sequential
- **Scales linearly** with number of opponents

### Performance by Simulation Count
*1 opponent, 4 workers*

| Simulations | Avg Time | Min | Max | Std Dev |
|-------------|----------|-----|-----|---------|
| 1,000 | 25ms | 25ms | 25ms | 0ms |
| 10,000 | 172ms | 171ms | 174ms | 1ms |
| 50,000 | 817ms | 811ms | 830ms | 7ms |
| 100,000 | 1.63s | 1.62s | 1.66s | 0.02s |

**Linear scaling**: ~17μs per simulation

### Performance by Number of Opponents
*10,000 simulations, 4 workers*

| Opponents | Avg Time | Win Rate | Notes |
|-----------|----------|----------|-------|
| 1 | 171ms | 85% | Fast heads-up |
| 3 | 333ms | 64% | Small table |
| 5 | 481ms | 49% | Medium table |
| 9 | 806ms | 31% | Full table |

**Linear scaling**: ~80ms per additional opponent

### Performance by Worker Count (Parallelism)
*100,000 simulations, 1 opponent*

| Workers | Avg Time | Speedup | Efficiency |
|---------|----------|---------|------------|
| 1 | 4.30s | 1.0x | 100% |
| 2 | 2.44s | 1.76x | 88% |
| 4 | 1.63s | 2.64x | 66% |
| 8 | 1.38s | 3.12x | 39% |

**Optimal**: 4 workers for most scenarios (best efficiency)

### Performance by Board Stage
*1 opponent, 10,000 simulations, 4 workers*

| Stage | Cards | Avg Time | Notes |
|-------|-------|----------|-------|
| Pre-flop | 0 | 189ms | Most uncertainty |
| Flop | 3 | 180ms | Minimal difference |
| Turn | 4 | 182ms | Consistent |
| River | 5 | 181ms | Fast (no randomness) |

**Insight**: Board cards have negligible performance impact

### Stress Tests
*9 opponents, 8 workers*

| Simulations | Avg Time | Use Case |
|-------------|----------|----------|
| 50,000 | 3.47s | High accuracy, full table |
| 100,000 | 6.79s | Maximum accuracy |

## Performance vs Python

Direct comparison with Python FastAPI implementation (both running 10,000 simulations):

| Scenario | Python (FastAPI) | Go (Gin) | Speedup |
|----------|------------------|----------|---------|
| 1 opponent | 764ms | 172ms | **4.4x faster** ⚡ |
| 3 opponents | 1,206ms | 333ms | **3.6x faster** ⚡ |
| 9 opponents | 2,700ms | 806ms | **3.3x faster** ⚡ |

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

Evaluates the best 5-card poker hand from 1-7 cards.

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
- `hole_cards` (required): Array of exactly 2 cards
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
  -d '{"hole_cards":["AS","KS"],"board_cards":["QS","JS","TS"]}'

# Calculate odds (default settings)
curl -X POST http://localhost:8001/odds \
  -H "Content-Type: application/json" \
  -d '{"hole_cards":["AS","AH"],"board_cards":[],"num_opponents":1}'

# Calculate odds (high accuracy)
curl -X POST http://localhost:8001/odds \
  -H "Content-Type: application/json" \
  -d '{
    "hole_cards":["AS","AH"],
    "board_cards":[],
    "num_opponents":1,
    "simulations":100000,
    "workers":8
  }'
```

### Python (FastAPI Integration)

```python
import httpx

GO_SERVICE_URL = "http://localhost:8001"

async def calculate_odds(hole_cards, board_cards, num_opponents):
    async with httpx.AsyncClient(timeout=30.0) as client:
        response = await client.post(
            f"{GO_SERVICE_URL}/odds",
            json={
                "hole_cards": hole_cards,
                "board_cards": board_cards,
                "num_opponents": num_opponents,
            }
        )
        return response.json()

# Usage
result = await calculate_odds(["AS", "AH"], [], 1)
print(f"Win probability: {result['win']:.1%}")
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

Cards are represented as 2-character strings: `[Rank][Suit]`

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
```

## Performance Tuning

### Simulation Count Recommendations

| Simulations | Time (1 opp) | Accuracy | Use Case |
|-------------|--------------|----------|----------|
| 1,000 | 25ms | ±3% | Quick estimates, UI responsiveness |
| 10,000 | 172ms | ±1% | **Default** - Good balance |
| 50,000 | 817ms | ±0.5% | High accuracy analysis |
| 100,000 | 1.63s | ±0.3% | Maximum precision |

### Worker Count Recommendations

| Workers | Best For | Notes |
|---------|----------|-------|
| 1 | Testing | Sequential, slowest |
| 2 | Low-end CPUs | Good efficiency (88%) |
| 4 | **Most systems** | Optimal balance (66% efficiency) |
| 8 | High-end CPUs | Diminishing returns (39% efficiency) |

**Rule of thumb**: Use number of CPU cores, max 8 workers

### Optimization Tips

1. **For UI/Real-time**: Use 1,000-10,000 simulations
2. **For Analysis**: Use 50,000-100,000 simulations
3. **For Full Table (9 players)**: Expect ~5x slower than heads-up
4. **Workers**: 4 workers optimal for most cases

## License

MIT
