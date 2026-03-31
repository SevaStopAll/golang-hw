#!/bin/bash

# Calendar API Testing Script
# Usage: ./test-api.sh

BASE_URL="http://localhost:8080"
API_URL="$BASE_URL/api"

echo "Calendar API Test"
echo ""

# 1. Health Check
echo "1. Health Check"
curl -s "$BASE_URL/health"
echo ""

# 2. Get all events
echo "2. GET /api/events"
curl -s "$API_URL/events"
echo ""

# 3. Create event
echo "3. POST /api/events - Create event"
RESPONSE=$(curl -s -X POST \
    -H "Content-Type: application/json" \
    -d '{
      "title": "Team Meeting",
      "description": "Weekly sync",
      "start_time": "2026-01-15T10:00:00Z",
      "end_time": "2026-01-15T11:00:00Z",
      "user_id": "user_1",
      "notify_before": 900
    }' \
    "$API_URL/events")
echo "$RESPONSE"
EVENT_ID=$(echo "$RESPONSE" | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
echo ""

# 4. Get event by ID
if [ ! -z "$EVENT_ID" ]; then
    echo "4. GET /api/events/$EVENT_ID"
    curl -s "$API_URL/events/$EVENT_ID"
    echo ""
fi

# 5. Update event
if [ ! -z "$EVENT_ID" ]; then
    echo "5. PUT /api/events/$EVENT_ID - Update event"
    curl -s -X PUT \
        -H "Content-Type: application/json" \
        -d '{
          "title": "Updated Meeting",
          "description": "Updated description",
          "start_time": "2026-01-15T14:00:00Z",
          "end_time": "2026-01-15T15:00:00Z",
          "user_id": "user_1",
          "notify_before": 1800
        }' \
        "$API_URL/events/$EVENT_ID"
    echo ""
fi

# 6. Get all events
echo "6. GET /api/events - All events"
curl -s "$API_URL/events"
echo ""

# 7. Delete event
#if [ ! -z "$EVENT_ID" ]; then
#    echo "7. DELETE /api/events/$EVENT_ID"
#    curl -s -X DELETE "$API_URL/events/$EVENT_ID"
#    echo ""
#fi

# 8. Final events list
echo "8. GET /api/events - Final list"
curl -s "$API_URL/events"
echo ""

echo "Testing completed"
echo ""
echo "Press any key to exit..."
read -n 1 -s -r
