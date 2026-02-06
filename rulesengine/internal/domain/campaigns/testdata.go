package campaigns

import "encoding/json"

// ValidGraph is a minimal valid gorules decision graph for testing
var ValidGraph = json.RawMessage(`{
  "contentType": "application/vnd.gorules.decision",
  "nodes": [
    {"id": "input", "type": "inputNode", "name": "Request", "position": {"x": 0, "y": 0}},
    {"id": "output", "type": "outputNode", "name": "Response", "position": {"x": 0, "y": 0}},
    {
      "id": "credit_check",
      "type": "decisionTableNode",
      "name": "Credit Check",
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
          {"id": "r1", "amount": "<= 100000", "approved": "true", "rate": "1.99", "installments": "[6, 12]", "value": "requested_amount"},
          {"id": "r2", "amount": "> 100000", "approved": "false", "rate": "0", "installments": "[]", "value": "0"}
        ]
      }
    }
  ],
  "edges": [
    {"id": "e1", "sourceId": "input", "targetId": "credit_check"},
    {"id": "e2", "sourceId": "credit_check", "targetId": "output"}
  ]
}`)

// InvalidGraph is an invalid graph for testing
var InvalidGraph = json.RawMessage(`{"invalid": "graph"}`)


