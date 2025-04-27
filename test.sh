#!/bin/bash

set -e

GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[0;33m'
NC='\033[0m'

BASE_URL="http://localhost:8080"

# Test GET /health
echo -e "${YELLOW}Checking health endpoint...${NC}"
GET_HEALTH_RESPONSE=`curl -s -X GET "${BASE_URL}/health"`
echo ${GET_HEALTH_RESPONSE} | jq -c .

STATUS=`echo ${GET_HEALTH_RESPONSE} | jq -r '.status'`

if [[ ${STATUS} == "ok" ]]; then
  echo -e "${GREEN}Health check passed!${NC}"
else
  echo -e "${RED}Health check failed!${NC}"
  exit 1
fi

# Test POST /todos
echo -e "${YELLOW}Creating a TODO item...${NC}"
POST_TODO=`curl -i -s -X POST "${BASE_URL}/todos" \
  -H "Content-Type: application/json" \
  -d '{"title": "Sample Todo", "description": "This is a test todo"}'`

LOCATION_HEADER_VALUE=$(echo "${POST_TODO}" | grep -i location | cut -d ' ' -f2 | tr -d '\r')

echo -e "${GREEN}TODO item created successfully! The created item location: ${LOCATION_HEADER_VALUE}${NC}"

# Test GET /todos
echo -e "${YELLOW}Fetch all TODO items...${NC}"
curl -s -X GET "${BASE_URL}/todos" | jq -c .

# Test PATCH /todos/:id
echo -e "${YELLOW}Updating a TODO item...${NC}"
PATCH_HTTP_CODE=`curl -s -w "%{http_code}" -X PATCH "${BASE_URL}${LOCATION_HEADER_VALUE}" \
  -H "Content-Type: application/json" \
  -d '{"title": "Updated Todo", "description": "This is an updated test todo","completed":true}' | jq -c .`
if [[ ${PATCH_HTTP_CODE} -eq 204 ]]; then
  echo -e "${GREEN}TODO item updated successfully!${NC}"
else
  echo -e "${RED}Failed to update TODO item!${NC}"
  exit 1
fi

# Test GET /todos/:id
echo -e "${YELLOW}Fetching an updated TODO item...${NC}"
curl -s -X GET "${BASE_URL}${LOCATION_HEADER_VALUE}" | jq -c .

# Test DELETE /todos/:id
echo -e "${YELLOW}Deleting a TODO item...${NC}"
DELETE_HTTP_CODE=`curl -s -w "%{http_code}" -X DELETE "${BASE_URL}${LOCATION_HEADER_VALUE}"`
if [[ ${DELETE_HTTP_CODE} -eq 204 ]]; then
  echo -e "${GREEN}TODO item deleted successfully!${NC}"
else
  echo -e "${RED}Failed to delete TODO item!${NC}"
  exit 1
fi
echo -e "${YELLOW}All tests completed successfully!${NC}"
