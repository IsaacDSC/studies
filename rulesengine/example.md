# Rule Engine - API Examples

## Start the Server

```bash
go run cmd/server/main.go
```

Server will start on `http://localhost:8080`

---

## 1. Blocklist Operations

### Block a User

```bash
curl -X POST http://localhost:8080/backoffice/block-list \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "550e8400-e29b-41d4-a716-446655440000"
  }'
```

**Response (201 Created):**
```json
{
  "message": "user blocked successfully",
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "blocked_at": "2025-12-29T19:00:00Z"
}
```

### Unblock a User

```bash
curl -X DELETE http://localhost:8080/backoffice/block-list/550e8400-e29b-41d4-a716-446655440000
```

**Response (200 OK):**
```json
{
  "message": "user unblocked successfully",
  "user_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

---

## 2. Campaign Management

### Create a Credit Campaign (Basic)

```bash
curl -X PATCH "http://localhost:8080/backoffice/campaigns/basic_credit?operations=credit" \
  -H "Content-Type: application/json" \
  -d '{
    "contentType": "application/vnd.gorules.decision",
    "nodes": [
      {"id": "input", "type": "inputNode", "name": "Request", "position": {"x": 0, "y": 0}},
      {"id": "output", "type": "outputNode", "name": "Response", "position": {"x": 0, "y": 0}},
      {
        "id": "credit_rules",
        "type": "decisionTableNode",
        "name": "Credit Rules",
        "position": {"x": 0, "y": 0},
        "content": {
          "hitPolicy": "first",
          "inputs": [
            {"id": "amount", "name": "Amount", "field": "requested_amount"}
          ],
          "outputs": [
            {"id": "approved", "name": "Approved", "field": "approved"},
            {"id": "rate", "name": "Rate", "field": "rate"},
            {"id": "installments", "name": "Installments", "field": "installments"},
            {"id": "value", "name": "Value", "field": "value"}
          ],
          "rules": [
            {"id": "r1", "amount": "<= 100000", "approved": "true", "rate": "1.99", "installments": "[6, 12, 24]", "value": "requested_amount"},
            {"id": "r2", "amount": "<= 300000", "approved": "true", "rate": "2.49", "installments": "[6, 12]", "value": "requested_amount"},
            {"id": "r3", "amount": "> 300000", "approved": "false", "rate": "0", "installments": "[]", "value": "0"}
          ]
        }
      }
    ],
    "edges": [
      {"id": "e1", "sourceId": "input", "targetId": "credit_rules"},
      {"id": "e2", "sourceId": "credit_rules", "targetId": "output"}
    ]
  }'
```

**Response (200 OK):**
```json
{
  "message": "rules updated successfully",
  "campaign": "basic_credit",
  "operations": ["credit"],
  "updated_at": "2025-12-29T19:00:00Z"
}
```

### Create a Black Friday Campaign (Credit + Anticipation)

```bash
curl -X PATCH "http://localhost:8080/backoffice/campaigns/black_friday_2025?operations=credit,anticipation" \
  -H "Content-Type: application/json" \
  -d '{
    "contentType": "application/vnd.gorules.decision",
    "nodes": [
      {"id": "input", "type": "inputNode", "name": "Request", "position": {"x": 0, "y": 0}},
      {"id": "output", "type": "outputNode", "name": "Response", "position": {"x": 0, "y": 0}},
      {
        "id": "promo_rules",
        "type": "decisionTableNode",
        "name": "Black Friday Promo",
        "position": {"x": 0, "y": 0},
        "content": {
          "hitPolicy": "first",
          "inputs": [
            {"id": "amount", "name": "Amount", "field": "requested_amount"}
          ],
          "outputs": [
            {"id": "approved", "name": "Approved", "field": "approved"},
            {"id": "rate", "name": "Rate", "field": "rate"},
            {"id": "installments", "name": "Installments", "field": "installments"},
            {"id": "value", "name": "Value", "field": "value"}
          ],
          "rules": [
            {"id": "r1", "amount": "<= 200000", "approved": "true", "rate": "0.99", "installments": "[6, 12, 24, 36]", "value": "requested_amount"},
            {"id": "r2", "amount": "<= 500000", "approved": "true", "rate": "1.49", "installments": "[6, 12, 24]", "value": "requested_amount"},
            {"id": "r3", "amount": "> 500000", "approved": "true", "rate": "1.99", "installments": "[6, 12]", "value": "500000"}
          ]
        }
      }
    ],
    "edges": [
      {"id": "e1", "sourceId": "input", "targetId": "promo_rules"},
      {"id": "e2", "sourceId": "promo_rules", "targetId": "output"}
    ]
  }'
```

### Create an Anticipation-Only Campaign

```bash
curl -X PATCH "http://localhost:8080/backoffice/campaigns/anticipation_standard?operations=anticipation" \
  -H "Content-Type: application/json" \
  -d '{
    "contentType": "application/vnd.gorules.decision",
    "nodes": [
      {"id": "input", "type": "inputNode", "name": "Request", "position": {"x": 0, "y": 0}},
      {"id": "output", "type": "outputNode", "name": "Response", "position": {"x": 0, "y": 0}},
      {
        "id": "anticipation_rules",
        "type": "decisionTableNode",
        "name": "Anticipation Rules",
        "position": {"x": 0, "y": 0},
        "content": {
          "hitPolicy": "first",
          "inputs": [
            {"id": "amount", "name": "Amount", "field": "requested_amount"}
          ],
          "outputs": [
            {"id": "approved", "name": "Approved", "field": "approved"},
            {"id": "rate", "name": "Rate", "field": "rate"},
            {"id": "installments", "name": "Installments", "field": "installments"},
            {"id": "value", "name": "Value", "field": "value"}
          ],
          "rules": [
            {"id": "r1", "amount": "<= 50000", "approved": "true", "rate": "2.99", "installments": "[1]", "value": "requested_amount"},
            {"id": "r2", "amount": "<= 150000", "approved": "true", "rate": "3.49", "installments": "[1]", "value": "requested_amount"},
            {"id": "r3", "amount": "> 150000", "approved": "false", "rate": "0", "installments": "[]", "value": "0"}
          ]
        }
      }
    ],
    "edges": [
      {"id": "e1", "sourceId": "input", "targetId": "anticipation_rules"},
      {"id": "e2", "sourceId": "anticipation_rules", "targetId": "output"}
    ]
  }'
```

### List All Campaigns

```bash
curl http://localhost:8080/backoffice/campaigns/list
```

**Response (200 OK):**
```json
{
  "campaigns": [
    {
      "name": "basic_credit",
      "operations": ["credit"],
      "created_at": "2025-12-29T19:00:00Z",
      "updated_at": "2025-12-29T19:00:00Z"
    },
    {
      "name": "black_friday_2025",
      "operations": ["credit", "anticipation"],
      "created_at": "2025-12-29T19:00:00Z",
      "updated_at": "2025-12-29T19:00:00Z"
    },
    {
      "name": "anticipation_standard",
      "operations": ["anticipation"],
      "created_at": "2025-12-29T19:00:00Z",
      "updated_at": "2025-12-29T19:00:00Z"
    }
  ],
  "total": 3
}
```

### Delete a Campaign

```bash
curl -X DELETE http://localhost:8080/backoffice/campaigns/basic_credit/credit
```

**Response (200 OK):**
```json
{
  "message": "campaign deleted successfully",
  "campaign": "basic_credit",
  "operation": "credit"
}
```

---

## 3. Rule Evaluation (Credit/Anticipation)

### Request Credit - Low Amount (Approved)

```bash
curl -X POST http://localhost:8080/rule-engine/rule/credit \
  -H "Content-Type: application/json" \
  -d '{
    "requested_amount": 50000,
    "user_id": "123e4567-e89b-12d3-a456-426614174000"
  }'
```

**Response (200 OK - Approved):**
```json
{
  "status": "approved",
  "user_id": "123e4567-e89b-12d3-a456-426614174000",
  "requested_amount": 50000,
  "offers": [
    {
      "campaign": "basic_credit",
      "approved": true,
      "rate": 1.99,
      "installments": [6, 12, 24],
      "value": 50000
    },
    {
      "campaign": "black_friday_2025",
      "approved": true,
      "rate": 0.99,
      "installments": [6, 12, 24, 36],
      "value": 50000
    }
  ],
  "evaluated_at": "2025-12-29T19:00:00Z"
}
```

### Request Credit - Medium Amount

```bash
curl -X POST http://localhost:8080/rule-engine/rule/credit \
  -H "Content-Type: application/json" \
  -d '{
    "requested_amount": 250000,
    "user_id": "123e4567-e89b-12d3-a456-426614174000"
  }'
```

**Response (200 OK):**
```json
{
  "status": "approved",
  "user_id": "123e4567-e89b-12d3-a456-426614174000",
  "requested_amount": 250000,
  "offers": [
    {
      "campaign": "basic_credit",
      "approved": true,
      "rate": 2.49,
      "installments": [6, 12],
      "value": 250000
    },
    {
      "campaign": "black_friday_2025",
      "approved": true,
      "rate": 1.49,
      "installments": [6, 12, 24],
      "value": 250000
    }
  ],
  "evaluated_at": "2025-12-29T19:00:00Z"
}
```

### Request Credit - High Amount (Partial Approval)

```bash
curl -X POST http://localhost:8080/rule-engine/rule/credit \
  -H "Content-Type: application/json" \
  -d '{
    "requested_amount": 600000,
    "user_id": "123e4567-e89b-12d3-a456-426614174000"
  }'
```

**Response (200 OK - Some denied, some approved with cap):**
```json
{
  "status": "approved",
  "user_id": "123e4567-e89b-12d3-a456-426614174000",
  "requested_amount": 600000,
  "offers": [
    {
      "campaign": "basic_credit",
      "approved": false,
      "rate": 0,
      "installments": [],
      "value": 0
    },
    {
      "campaign": "black_friday_2025",
      "approved": true,
      "rate": 1.99,
      "installments": [6, 12],
      "value": 500000
    }
  ],
  "evaluated_at": "2025-12-29T19:00:00Z"
}
```

### Request Anticipation

```bash
curl -X POST http://localhost:8080/rule-engine/rule/anticipation \
  -H "Content-Type: application/json" \
  -d '{
    "requested_amount": 30000,
    "user_id": "123e4567-e89b-12d3-a456-426614174000"
  }'
```

**Response (200 OK):**
```json
{
  "status": "approved",
  "user_id": "123e4567-e89b-12d3-a456-426614174000",
  "requested_amount": 30000,
  "offers": [
    {
      "campaign": "black_friday_2025",
      "approved": true,
      "rate": 0.99,
      "installments": [6, 12, 24, 36],
      "value": 30000
    },
    {
      "campaign": "anticipation_standard",
      "approved": true,
      "rate": 2.99,
      "installments": [1],
      "value": 30000
    }
  ],
  "evaluated_at": "2025-12-29T19:00:00Z"
}
```

### Request Credit - Blocked User

```bash
# First, block the user
curl -X POST http://localhost:8080/backoffice/block-list \
  -H "Content-Type: application/json" \
  -d '{"user_id": "blocked-user-e89b-12d3-a456-426614174000"}'

# Then try to request credit
curl -X POST http://localhost:8080/rule-engine/rule/credit \
  -H "Content-Type: application/json" \
  -d '{
    "requested_amount": 50000,
    "user_id": "blocked-user-e89b-12d3-a456-426614174000"
  }'
```

**Response (200 OK - Blocked):**
```json
{
  "status": "blocked",
  "user_id": "blocked-user-e89b-12d3-a456-426614174000",
  "requested_amount": 50000,
  "reason": "user_blocked",
  "message": "User is blocked for credit operations",
  "offers": [],
  "evaluated_at": "2025-12-29T19:00:00Z"
}
```

---

## 4. Error Examples

### Invalid Operation

```bash
curl -X POST http://localhost:8080/rule-engine/rule/loan \
  -H "Content-Type: application/json" \
  -d '{
    "requested_amount": 50000,
    "user_id": "123e4567-e89b-12d3-a456-426614174000"
  }'
```

**Response (400 Bad Request):**
```json
{
  "error": "invalid_operation",
  "message": "operation must be 'credit' or 'anticipation'"
}
```

### Invalid UUID

```bash
curl -X POST http://localhost:8080/rule-engine/rule/credit \
  -H "Content-Type: application/json" \
  -d '{
    "requested_amount": 50000,
    "user_id": "invalid-uuid"
  }'
```

**Response (400 Bad Request):**
```json
{
  "error": "invalid_user_id",
  "message": "user_id must be a valid UUID"
}
```

### Invalid Amount

```bash
curl -X POST http://localhost:8080/rule-engine/rule/credit \
  -H "Content-Type: application/json" \
  -d '{
    "requested_amount": -100,
    "user_id": "123e4567-e89b-12d3-a456-426614174000"
  }'
```

**Response (400 Bad Request):**
```json
{
  "error": "invalid_amount",
  "message": "requested_amount must be a positive integer (cents)"
}
```

### No Campaigns Configured

```bash
# Delete all campaigns first, then request credit
curl -X POST http://localhost:8080/rule-engine/rule/credit \
  -H "Content-Type: application/json" \
  -d '{
    "requested_amount": 50000,
    "user_id": "123e4567-e89b-12d3-a456-426614174000"
  }'
```

**Response (422 Unprocessable Entity):**
```json
{
  "error": "no_campaigns",
  "message": "no campaigns configured for this operation"
}
```

---

## 5. Health Check

```bash
curl http://localhost:8080/health
```

**Response (200 OK):**
```json
{"status":"ok"}
```

---

## Quick Test Script

```bash
#!/bin/bash

BASE_URL="http://localhost:8080"
USER_ID="123e4567-e89b-12d3-a456-426614174000"

echo "=== Creating Campaign ==="
curl -s -X PATCH "$BASE_URL/backoffice/campaigns/test_campaign?operations=credit" \
  -H "Content-Type: application/json" \
  -d '{
    "contentType": "application/vnd.gorules.decision",
    "nodes": [
      {"id": "input", "type": "inputNode", "name": "Request", "position": {"x": 0, "y": 0}},
      {"id": "output", "type": "outputNode", "name": "Response", "position": {"x": 0, "y": 0}},
      {
        "id": "rules",
        "type": "decisionTableNode",
        "name": "Rules",
        "position": {"x": 0, "y": 0},
        "content": {
          "hitPolicy": "first",
          "inputs": [{"id": "amount", "name": "Amount", "field": "requested_amount"}],
          "outputs": [
            {"id": "approved", "name": "Approved", "field": "approved"},
            {"id": "rate", "name": "Rate", "field": "rate"},
            {"id": "installments", "name": "Installments", "field": "installments"},
            {"id": "value", "name": "Value", "field": "value"}
          ],
          "rules": [
            {"id": "r1", "amount": "<= 100000", "approved": "true", "rate": "1.99", "installments": "[6, 12]", "value": "requested_amount"},
            {"id": "r2", "amount": "> 100000", "approved": "false", "rate": "0", "installments": "[]", "value": "0"}
          ]
        }
      }
    ],
    "edges": [
      {"id": "e1", "sourceId": "input", "targetId": "rules"},
      {"id": "e2", "sourceId": "rules", "targetId": "output"}
    ]
  }' | jq .

echo -e "\n=== Requesting Credit (should be approved) ==="
curl -s -X POST "$BASE_URL/rule-engine/rule/credit" \
  -H "Content-Type: application/json" \
  -d "{\"requested_amount\": 50000, \"user_id\": \"$USER_ID\"}" | jq .

echo -e "\n=== Requesting Credit (should be denied) ==="
curl -s -X POST "$BASE_URL/rule-engine/rule/credit" \
  -H "Content-Type: application/json" \
  -d "{\"requested_amount\": 200000, \"user_id\": \"$USER_ID\"}" | jq .

echo -e "\n=== Listing Campaigns ==="
curl -s "$BASE_URL/backoffice/campaigns/list" | jq .
```


