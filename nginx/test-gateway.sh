#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'

# Base URL
BASE_URL="http://localhost"
API_KEY="key123"

# Function to check HTTP status
check_status() {
    if [[ $1 -eq 200 || $1 -eq 201 ]]; then
        echo -e "${GREEN}✔ PASS${NC}"
    else
        echo -e "${RED}✘ FAIL (Status: $1)${NC}"
    fi
}

echo "Testing API Gateway..."

# Health check
echo -e "\nTesting health endpoint:"
STATUS=$(curl -s -o /dev/null -w "%{http_code}" ${BASE_URL}/health)
check_status $STATUS

# User Service
echo -e "\nTesting user service (create user):"
STATUS=$(curl -s -o /dev/null -w "%{http_code}" -X POST \
    -H "Content-Type: application/json" \
    -H "X-Api-Key: ${API_KEY}" \
    -d '{"username":"test_user","email":"test3@example.com","password":"password123"}' \
    ${BASE_URL}/api/users)
check_status $STATUS

# Order Service
echo -e "\nTesting order service (create order):"
STATUS=$(curl -s -o /dev/null -w "%{http_code}" -X POST \
    -H "Content-Type: application/json" \
    -H "X-Api-Key: ${API_KEY}" \
    -d '{"user_id":1,"status":"pending","total":99.99,"items":[{"name":"Item 1","price":49.99,"quantity":2}]}' \
    ${BASE_URL}/api/orders/)
check_status $STATUS

# Rate limiting test
echo -e "\nTesting rate limiting (11 requests in 1 second):"
for i in $(seq 1 11); do
    STATUS=$(curl -s -o /dev/null -w "%{http_code}" \
        -H "X-Api-Key: ${API_KEY}" \
        ${BASE_URL}/api/users/1 &)
    sleep 0.1
done
wait
echo "Rate limit test completed."

# Invalid API key
echo -e "\nTesting invalid API key:"
STATUS=$(curl -s -o /dev/null -w "%{http_code}" \
    -H "X-Api-Key: invalid_key" \
    ${BASE_URL}/api/users/1)
check_status $STATUS

# CORS Test
echo -e "\nTesting CORS headers:"
curl -s -I \
    -H "Origin: http://example.com" \
    -H "X-Api-Key: ${API_KEY}" \
    ${BASE_URL}/api/users/1 | grep -i "Access-Control-Allow-Origin"

echo -e "\nAll tests completed."