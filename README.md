# AlgoPay - Algorand Payment SDK & Gateway

AlgoPay is a complete payment SDK and gateway that allows developers and merchants to accept payments in ALGO and ASA tokens on the Algorand TestNet. Built with Go and designed for production use.

## 🚀 Features

- **REST API** for payment initialization and status checking
- **Real-time blockchain monitoring** using Algorand Indexer
- **Webhook notifications** for payment confirmations
- **QR code generation** for easy mobile payments
- **ASA token support** for custom tokens
- **Payment expiration** and cleanup
- **SQLite database** for persistent storage
- **Production-ready** with proper error handling and logging

## 📋 Table of Contents

- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Configuration](#configuration)
- [Running the Server](#running-the-server)
- [API Endpoints](#api-endpoints)
- [Testing Payment Flow](#testing-payment-flow)
- [Webhook Integration](#webhook-integration)
- [Deployment](#deployment)
- [Development](#development)
- [Troubleshooting](#troubleshooting)

## Prerequisites

Before you begin, ensure you have the following installed:

- **Go 1.20+** - [Download and install Go](https://golang.org/doc/install)
- **Git** - [Download and install Git](https://git-scm.com/downloads)
- **SQLite3** (usually pre-installed on most systems)

## Installation

1. **Clone the repository:**
   ```bash
   git clone https://github.com/ItsOrganic/algorand-payment-sdk &&  cd algopay
   ```

2. **Install dependencies:**
   ```bash
   make install-deps
   ```

3. **Create environment file:**
   ```bash
   make init-env
   ```

4. **Build the application:**
   ```bash
   make build
   ```

## Configuration

The application uses environment variables for configuration. Copy `.env.example` to `.env` and modify as needed:

```bash
cp .env.example .env
```

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | `8080` |
| `DATABASE_PATH` | SQLite database file path | `./algopay.db` |
| `ALGO_NODE_URL` | Algorand node URL | `https://testnet-api.algonode.cloud` |
| `ALGO_INDEXER_URL` | Algorand indexer URL | `https://testnet-idx.algonode.cloud` |
| `ALGO_TOKEN` | Algorand API token (optional for public nodes) | `` |
| `PAYMENT_TIMEOUT` | Payment timeout in minutes | `30` |

## Running the Server

### Development Mode
```bash
make run
```

### Production Build
```bash
make build
./build/algopay
```

### Custom Port
```bash
make run-port  # Runs on port 9000
# OR
./build/algopay -port=9000
```

The server will start and display:
```
🚀 AlgoPay Server Starting
📡 Port: 8080
🗄️ Database: ./algopay.db
🌐 Algorand Node: https://testnet-api.algonode.cloud
🔍 Algorand Indexer: https://testnet-idx.algonode.cloud
⏰ Payment Timeout: 30 minutes

📋 API Endpoints:
   POST /api/v1/init-payment     - Initialize new payment
   GET  /api/v1/check-payment/:id - Check payment status
   GET  /api/v1/payment/:id       - Get payment details
   GET  /health                   - Health check

🌟 Server running on http://localhost:8080
```

## API Endpoints

### 1. Initialize Payment
**POST** `/api/v1/init-payment`

Create a new payment request.

**Request Body:**
```json
{
  "merchant_address": "MERCHANT_ALGORAND_ADDRESS",
  "amount": 1000000,
  "asset_id": 0,
  "callback_url": "https://your-domain.com/webhook"
}
```

**Response:**
```json
{
  "payment_id": "uuid-string",
  "merchant_address": "MERCHANT_ALGORAND_ADDRESS",
  "amount": 1000000,
  "asset_id": 0,
  "qr_code": "base64-encoded-qr-image",
  "expires_at": "2024-01-15T10:30:00Z",
  "status": "pending"
}
```

### 2. Check Payment Status
**GET** `/api/v1/check-payment/:id`

Check the status of a payment.

**Response:**
```json
{
  "payment_id": "uuid-string",
  "status": "completed",
  "txn_id": "transaction-id",
  "created_at": "2024-01-15T10:00:00Z",
  "updated_at": "2024-01-15T10:05:00Z"
}
```

### 3. Get Payment Details
**GET** `/api/v1/payment/:id`

Get full payment details.

**Response:**
```json
{
  "id": "uuid-string",
  "merchant_address": "MERCHANT_ALGORAND_ADDRESS",
  "amount": 1000000,
  "asset_id": 0,
  "callback_url": "https://your-domain.com/webhook",
  "status": "completed",
  "txn_id": "transaction-id",
  "created_at": "2024-01-15T10:00:00Z",
  "updated_at": "2024-01-15T10:05:00Z",
  "expires_at": "2024-01-15T10:30:00Z"
}
```

### 4. Health Check
**GET** `/health`

Check if the server is running.

**Response:**
```json
{
  "status": "ok"
}
```

## Testing Payment Flow

### Step 1: Get Test ALGO

1. **Visit the TestNet Faucet:**
   - Go to [TestNet Faucet](https://testnet.algoexplorer.io/dispenser)
   - Enter your Algorand address
   - Request test ALGO

2. **Using Pera Wallet:**
   - Download [Pera Wallet](https://perawallet.app/)
   - Create/import your wallet
   - Switch to TestNet in settings
   - Use the faucet to get test ALGO

### Step 2: Initialize a Payment

```bash
curl -X POST http://localhost:8080/api/v1/init-payment \
  -H "Content-Type: application/json" \
  -d '{
    "merchant_address": "YOUR_MERCHANT_ADDRESS",
    "amount": 1000000,
    "asset_id": 0,
    "callback_url": "https://your-domain.com/webhook"
  }'
```

**Example Response:**
```json
{
  "payment_id": "123e4567-e89b-12d3-a456-426614174000",
  "merchant_address": "ABCD...XYZ",
  "amount": 1000000,
  "asset_id": 0,
  "qr_code": "iVBORw0KGgoAAAANSUhEUgAAAQAAAAEA...",
  "expires_at": "2024-01-15T10:30:00Z",
  "status": "pending"
}
```

### Step 3: Make the Payment

**Option A: Using the QR Code**
- Decode the base64 QR code image
- Scan with your mobile wallet
- Complete the transaction

**Option B: Manual Transaction**
- Send exactly the specified amount
- To the merchant address
- From any TestNet wallet

### Step 4: Check Payment Status

```bash
curl http://localhost:8080/api/v1/check-payment/123e4567-e89b-12d3-a456-426614174000
```

**Expected Response (when completed):**
```json
{
  "payment_id": "123e4567-e89b-12d3-a456-426614174000",
  "status": "completed",
  "txn_id": "ABC123...XYZ789",
  "created_at": "2024-01-15T10:00:00Z",
  "updated_at": "2024-01-15T10:05:00Z"
}
```

## Webhook Integration

When a payment is completed, AlgoPay will send a POST request to your callback URL.

### Webhook Payload

```json
{
  "payment_id": "uuid-string",
  "status": "completed",
  "merchant_address": "MERCHANT_ADDRESS",
  "amount": 1000000,
  "asset_id": 0,
  "txn_id": "transaction-id",
  "timestamp": "2024-01-15T10:05:00Z"
}
```

### Example Webhook Handler (Node.js/Express)

```javascript
app.post('/webhook', (req, res) => {
  const payment = req.body;
  
  if (payment.status === 'completed') {
    console.log(`Payment completed: ${payment.payment_id}`);
    console.log(`Transaction ID: ${payment.txn_id}`);
    console.log(`Amount: ${payment.amount} microALGO`);
    
    // Process the successful payment
    // Update your database, send confirmations, etc.
  }
  
  res.status(200).send('OK');
});
```

## Deployment

### Local Deployment with ngrok

For local testing with webhooks:

1. **Install ngrok:**
   ```bash
   # macOS
   brew install ngrok
   
   # Or download from https://ngrok.com/
   ```

2. **Run AlgoPay:**
   ```bash
   make run
   ```

3. **Expose with ngrok:**
   ```bash
   ngrok http 8080
   ```

4. **Use the ngrok URL for webhooks:**
   ```
   https://abc123.ngrok.io/your-webhook-endpoint
   ```

### Production Deployment

1. **Build for production:**
   ```bash
   make build-all
   ```

2. **Set environment variables:**
   ```bash
   export PORT=8080
   export DATABASE_PATH=/path/to/production.db
   export ALGO_NODE_URL=https://testnet-api.algonode.cloud
   export ALGO_INDEXER_URL=https://testnet-idx.algonode.cloud
   ```

3. **Run the binary:**
   ```bash
   ./build/algopay-linux-amd64
   ```

### Docker Deployment

Create a `Dockerfile`:

```dockerfile
FROM golang:1.20-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o algopay cmd/algopay/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/algopay .
COPY --from=builder /app/.env.example .env

EXPOSE 8080
CMD ["./algopay"]
```

Build and run:
```bash
docker build -t algopay .
docker run -p 8080:8080 algopay
```

### Network Configuration

- **Firewall:** Open port 8080 (or your chosen port)
- **HTTPS:** Use a reverse proxy like Nginx for SSL termination
- **Load Balancing:** Use multiple instances behind a load balancer for high availability

### Environment Variables for Production

```bash
# Production environment
PORT=8080
DATABASE_PATH=/var/lib/algopay/production.db
ALGO_NODE_URL=https://testnet-api.algonode.cloud
ALGO_INDEXER_URL=https://testnet-idx.algonode.cloud
PAYMENT_TIMEOUT=30
```

## Development

### Hot Reload Development

Install Air for hot reloading:
```bash
go install github.com/cosmtrek/air@latest
```

Create `.air.toml`:
```toml
root = "."
testdata_dir = "testdata"
tmp_dir = "tmp"

[build]
  bin = "./tmp/main"
  cmd = "go build -o ./tmp/main cmd/algopay/main.go"
  delay = 1000
  exclude_dir = ["assets", "tmp", "vendor", "testdata"]
  exclude_file = []
  exclude_regex = ["_test.go"]
  exclude_unchanged = false
  follow_symlink = false
  full_bin = ""
  include_dir = []
  include_ext = ["go", "tpl", "tmpl", "html"]
  kill_delay = "0s"
  log = "build-errors.log"
  send_interrupt = false
  stop_on_root = false

[color]
  app = ""
  build = "yellow"
  main = "magenta"
  runner = "green"
  watcher = "cyan"

[log]
  time = false

[misc]
  clean_on_exit = false
```

Run with hot reload:
```bash
make dev
```

### Running Tests

```bash
make test
```

### Code Organization

```
algopay/
├── cmd/algopay/        # Main application
├── api/                # HTTP handlers and routing
├── algorand/           # Algorand client and blockchain logic
├── models/             # Data models and structures
├── db/                 # Database operations
├── config/             # Configuration management
├── build/              # Compiled binaries
└── README.md
```

### Adding New Features

1. **Models:** Add new structs in `models/`
2. **Database:** Add new operations in `db/`
3. **API:** Add new endpoints in `api/handlers.go`
4. **Blockchain:** Add new Algorand operations in `algorand/client.go`

## Troubleshooting

### Common Issues

1. **"Failed to connect to Algorand node"**
   - Check if the node URL is correct
   - Verify internet connectivity
   - Try using public nodes: `https://testnet-api.algonode.cloud`

2. **"Database locked" error**
   - Ensure only one instance is running
   - Check file permissions on the database file
   - Use absolute path for `DATABASE_PATH`

3. **"Payment not detected"**
   - Verify the transaction was sent to the correct address
   - Check if the amount matches exactly (including fees)
   - Ensure you're using TestNet addresses and ALGO

4. **Webhook not working**
   - Check if the callback URL is accessible
   - Verify the endpoint accepts POST requests
   - Use ngrok for local testing

### Debug Mode

Enable debug logging by setting:
```bash
export GIN_MODE=debug
```

### Logs

The application logs important events:
- Payment creations
- Payment confirmations
- Webhook deliveries
- Errors and warnings

### Support

For additional support:
1. Check the logs for detailed error messages
2. Verify your Algorand addresses are valid TestNet addresses
3. Ensure proper network connectivity
4. Test with small amounts first

---

## User Manual

### 🧪 Complete Testing Guide

This section provides a step-by-step guide to test the AlgoPay system from start to finish.

#### Prerequisites Setup

1. **Install Go 1.20+**
   ```bash
   # macOS
   brew install go
   
   # Ubuntu/Debian
   sudo apt install golang-go
   
   # Or download from: https://golang.org/doc/install
   ```

2. **Verify Go Installation**
   ```bash
   go version
   # Should output: go version go1.20+ ...
   ```

#### Installation and Setup

1. **Clone and Build**
   ```bash
   git clone <repository-url>
   cd algopay
   make install-deps
   make init-env
   make build
   ```

2. **Start the Server**
   ```bash
   make run
   ```
   
   You should see:
   ```
   🚀 AlgoPay Server Starting
   📡 Port: 8080
   🌟 Server running on http://localhost:8080
   ```

#### Getting Test ALGO

1. **Get a TestNet Wallet**
   - Download [Pera Wallet](https://perawallet.app/) on your mobile device
   - Create a new wallet or import existing
   - Go to Settings → Developer Settings → Node Settings
   - Switch to TestNet

2. **Fund Your Wallet**
   - Visit [TestNet Faucet](https://testnet.algoexplorer.io/dispenser)
   - Enter your wallet address
   - Request 10 ALGO (takes 1-2 minutes)
   - Verify balance in your wallet

#### Complete Payment Flow Test

**Step 1: Initialize Payment**
```bash
curl -X POST http://localhost:8080/api/v1/init-payment \
  -H "Content-Type: application/json" \
  -d '{
    "merchant_address": "7ZUECA7HFLZTXENRV24SHLU4AVPUTMTTDUFUBNBD64C73F3UHRTHAIOF6Q",
    "amount": 1000000,
    "asset_id": 0,
    "callback_url": "https://httpbin.org/post"
  }'
```

**Expected Response:**
```json
{
  "payment_id": "550e8400-e29b-41d4-a716-446655440000",
  "merchant_address": "7ZUECA7HFLZTXENRV24SHLU4AVPUTMTTDUFUBNBD64C73F3UHRTHAIOF6Q",
  "amount": 1000000,
  "asset_id": 0,
  "qr_code": "iVBORw0KGgoAAAANSUhEUgAAAQAAAAEA...",
  "expires_at": "2024-01-15T11:00:00Z",
  "status": "pending"
}
```

**Step 2: Check Initial Status**
```bash
curl http://localhost:8080/api/v1/check-payment/550e8400-e29b-41d4-a716-446655440000
```

**Expected Response:**
```json
{
  "payment_id": "550e8400-e29b-41d4-a716-446655440000",
  "status": "pending",
  "txn_id": "",
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z"
}
```

**Step 3: Make the Payment**
- Open Pera Wallet
- Tap "Send"
- Enter merchant address: `7ZUECA7HFLZTXENRV24SHLU4AVPUTMTTDUFUBNBD64C73F3UHRTHAIOF6Q`
- Enter amount: `1` ALGO (1000000 microALGO)
- Confirm and send transaction

**Step 4: Wait for Confirmation**
The system checks every 10 seconds. Wait 1-2 minutes, then check status again:

```bash
curl http://localhost:8080/api/v1/check-payment/550e8400-e29b-41d4-a716-446655440000
```

**Expected Response (after payment):**
```json
{
  "payment_id": "550e8400-e29b-41d4-a716-446655440000",
  "status": "completed",
  "txn_id": "ABC123DEF456...",
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:32:15Z"
}
```

#### Webhook Testing

**Step 1: Set up webhook.site**
1. Go to [webhook.site](https://webhook.site)
2. Copy your unique URL (e.g., `https://webhook.site/12345678-1234-1234-1234-123456789012`)

**Step 2: Initialize Payment with Webhook**
```bash
curl -X POST http://localhost:8080/api/v1/init-payment \
  -H "Content-Type: application/json" \
  -d '{
    "merchant_address": "YOUR_MERCHANT_ADDRESS",
    "amount": 1000000,
    "asset_id": 0,
    "callback_url": "https://webhook.site/12345678-1234-1234-1234-123456789012"
  }'
```

**Step 3: Complete Payment and Check Webhook**
- Send the payment as before
- Go back to webhook.site
- You should see a POST request with payload:

```json
{
  "payment_id": "uuid-string",
  "status": "completed",
  "merchant_address": "YOUR_MERCHANT_ADDRESS",
  "amount": 1000000,
  "asset_id": 0,
  "txn_id": "transaction-hash",
  "timestamp": "2024-01-15T10:32:15Z"
}
```

#### ASA Token Testing

**Step 1: Find a TestNet ASA**
- Visit [TestNet Explorer](https://testnet.algoexplorer.io/assets)
- Find an ASA with available balance (e.g., Asset ID: 10458941)

**Step 2: Opt-in to ASA (if needed)**
- In Pera Wallet, go to "Manage Assets"
- Add the ASA by ID
- Confirm opt-in transaction

**Step 3: Test ASA Payment**
```bash
curl -X POST http://localhost:8080/api/v1/init-payment \
  -H "Content-Type: application/json" \
  -d '{
    "merchant_address": "YOUR_MERCHANT_ADDRESS",
    "amount": 100,
    "asset_id": 10458941,
    "callback_url": "https://webhook.site/your-url"
  }'
```

#### Error Testing

**Test 1: Invalid Address**
```bash
curl -X POST http://localhost:8080/api/v1/init-payment \
  -H "Content-Type: application/json" \
  -d '{
    "merchant_address": "invalid-address",
    "amount": 1000000,
    "asset_id": 0
  }'
```

**Expected Response:**
```json
{
  "error": "Invalid merchant address"
}
```

**Test 2: Invalid Payment ID**
```bash
curl http://localhost:8080/api/v1/check-payment/invalid-id
```

**Expected Response:**
```json
{
  "error": "Payment not found"
}
```

**Test 3: Zero Amount**
```bash
curl -X POST http://localhost:8080/api/v1/init-payment \
  -H "Content-Type: application/json" \
  -d '{
    "merchant_address": "VALID_ADDRESS",
    "amount": 0,
    "asset_id": 0
  }'
```

**Expected Response:**
```json
{
  "error": "Amount must be greater than 0"
}
```

#### Performance Testing

**Test Multiple Payments**
```bash
for i in {1..5}; do
  curl -X POST http://localhost:8080/api/v1/init-payment \
    -H "Content-Type: application/json" \
    -d "{
      \"merchant_address\": \"YOUR_MERCHANT_ADDRESS\",
      \"amount\": $((1000000 * i)),
      \"asset_id\": 0,
      \"callback_url\": \"https://httpbin.org/post\"
    }"
  echo
done
```

#### Local Development with ngrok

**Step 1: Install ngrok**
```bash
# macOS
brew install ngrok

# Or download from https://ngrok.com/download
```

**Step 2: Expose Local Server**
```bash
# In one terminal, run AlgoPay
make run

# In another terminal, expose it
ngrok http 8080
```

**Step 3: Use ngrok URL**
```
https://abc123.ngrok.io/api/v1/init-payment
```

#### Health Monitoring

**Check Server Health**
```bash
curl http://localhost:8080/health
```

**Monitor Logs**
The server provides detailed logs for:
- Payment creation
- Blockchain monitoring
- Payment confirmation
- Webhook delivery
- Error handling

#### Production Deployment Checklist

- [ ] Environment variables configured
- [ ] Database path set correctly
- [ ] Firewall rules configured
- [ ] SSL certificate installed
- [ ] Webhook URLs use HTTPS
- [ ] Monitoring and alerting setup
- [ ] Backup strategy implemented
- [ ] Load balancing configured (if needed)

#### Troubleshooting Common Issues

1. **Payment Not Detected**
   - Verify exact amount was sent
   - Check transaction on [TestNet Explorer](https://testnet.algoexplorer.io)
   - Ensure payment was sent to correct address
   - Wait up to 2 minutes for blockchain confirmation

2. **Webhook Not Received**
   - Test webhook URL manually
   - Check server logs for delivery attempts
   - Verify webhook endpoint accepts POST requests
   - Use webhook.site for testing

3. **Database Issues**
   - Check file permissions
   - Ensure only one instance running
   - Verify disk space available
   - Use absolute path for database file

#### Support and Resources

- **Algorand Developer Portal**: [developer.algorand.org](https://developer.algorand.org)
- **TestNet Explorer**: [testnet.algoexplorer.io](https://testnet.algoexplorer.io)
- **Pera Wallet**: [perawallet.app](https://perawallet.app)
- **Go Documentation**: [golang.org/doc](https://golang.org/doc)

This completes the comprehensive user manual for AlgoPay. The system is now ready for production use with proper testing and deployment procedures.
