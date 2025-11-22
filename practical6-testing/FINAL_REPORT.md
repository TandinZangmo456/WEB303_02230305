# Practical 6 Final Report: Comprehensive Testing for Microservices

## Executive Summary

Successfully implemented a comprehensive testing strategy for a microservices-based cafe management system. The project demonstrates proficiency in unit testing, integration testing, and test automation using Go, gRPC, and modern testing frameworks.

**Final Result**:  **14/14 tests passing (100% success rate)**


## 1. Project Overview

### System Architecture
Built a microservices system consisting of:
- **User Service** - Manages user accounts and cafe owners
- **Menu Service** - Handles menu items and pricing
- **Order Service** - Processes orders with validation
- **API Gateway** - HTTP interface (not tested in this phase)

### Technology Stack
- **Language**: Go 1.24
- **Communication**: gRPC with Protocol Buffers
- **Database**: PostgreSQL (production), SQLite (testing)
- **ORM**: GORM
- **Testing**: Go testing framework + Testify
- **Automation**: Makefile

## 2. Implementation Steps Completed

### Phase 1: Project Setup 
1. Created folder structure for all services
2. Initialized Go modules for each service
3. Created Protocol Buffer definitions for 3 services
4. Generated Go code from proto files

### Phase 2: Service Implementation 
1. Implemented database models (User, MenuItem, Order, OrderItem)
2. Created database connection and migration logic
3. Built gRPC server implementations for all services
4. Implemented main.go files with environment configuration

### Phase 3: Unit Testing 
**User Service Tests** (3 tests, 5 scenarios):
- TestCreateUser - Creates users and cafe owners
- TestGetUser - Retrieves existing user, handles not found error
- TestGetUsers - Lists all users

**Menu Service Tests** (4 tests, 9 scenarios):
- TestCreateMenuItem - Creates items, handles zero prices
- TestGetMenuItem - Retrieves items, handles not found error
- TestGetMenuItems - Lists all menu items
- TestPriceHandling - Tests float precision handling

**Order Service Tests** (3 tests with mocks):
- TestCreateOrder_Success - Creates order with valid data
- TestCreateOrder_InvalidUser - Validates user existence
- TestCreateOrder_InvalidMenuItem - Validates menu item existence

**Key Achievements**:
- Used in-memory SQLite for fast, isolated tests
- Implemented mock objects for external service dependencies
- Applied table-driven test pattern for maintainability
- Achieved proper gRPC error code validation

### Phase 4: Integration Testing 
**Integration Tests** (4 tests):
- TestIntegration_CreateUser - User service standalone
- TestIntegration_CreateMenuItem - Menu service standalone  
- TestIntegration_UserAndMenuServices - Multi-service workflow
- TestIntegration_MultipleUsers - Concurrent operations

**Key Achievements**:
- Used bufconn for in-memory gRPC connections (no network overhead)
- Isolated database per test to prevent data contamination
- Validated cross-service communication
- Tested complete user workflows

### Phase 5: Test Automation 
Created comprehensive Makefile with commands:
- `make test-unit` - Run all unit tests
- `make test-integration` - Run integration tests
- `make test-all` - Run complete test suite
- `make test-coverage` - Generate HTML coverage reports

Also created:
- `docker-compose.yml` - For E2E testing environment
- `generate_proto.sh` - Automated proto code generation

## 3. Test Results

### Unit Test Results
```
User Service:     PASS (3 tests, 0.00s)
Menu Service:     PASS (4 tests, 0.00s)  
Order Service:    PASS (3 tests, 0.00s)
```

### Integration Test Results
```
Integration:      PASS (4 tests, 0.275s)
```

### Overall Statistics
- **Total Tests**: 14
- **Passing**: 14 (100%)
- **Failing**: 0 (0%)
- **Execution Time**: ~0.3 seconds
- **Test Types**: Unit (10) + Integration (4)


## 4. Testing Best Practices Demonstrated

### Test Isolation
Each test uses its own database instance, preventing test interference.
```go
dbName := fmt.Sprintf("file:test_%d.db?mode=memory", time.Now().UnixNano())
```

### Table-Driven Tests
Easy to add new test cases without code duplication.
```go
tests := []struct {
    name    string
    request *userv1.CreateUserRequest
    wantErr bool
}{
    {"successful creation", req1, false},
    {"duplicate email", req2, true},
}
```

### Mock Objects
Simulated external dependencies for unit testing.
```go
mockUserClient := new(MockUserServiceClient)
mockUserClient.On("GetUser", mock.Anything, req).Return(resp, nil)
```

### Proper Error Handling
Validated gRPC error codes and messages.
```go
st, ok := status.FromError(err)
assert.Equal(t, codes.NotFound, st.Code())
```

### Fast Execution
All tests run in under 1 second using in-memory databases.

## 5. Challenges Faced and Solutions

### Challenge 1: Proto File Compatibility
**Problem**: Order service imports user and menu protos, causing type conflicts in integration tests.

**Solution**: Used separate proto copies for order-service and imported from each service's own proto package in integration tests.

### Challenge 2: Shared Database State
**Problem**: Integration tests were sharing database, causing test failures.

**Solution**: Created unique database instance per test using timestamps.
```go
dbName := fmt.Sprintf("file:test_%d.db?mode=memory", time.Now().UnixNano())
```

### Challenge 3: gRPC Connection Management
**Problem**: Network overhead in integration tests.

**Solution**: Used bufconn for in-memory gRPC connections, eliminating network latency.

## 6. Key Learnings

### Technical Skills Acquired
1. **Testing Pyramid**: Implemented 70% unit, 30% integration tests
2. **Mocking**: Created mock gRPC clients for isolated testing
3. **gRPC Testing**: Used bufconn for efficient integration tests
4. **Test Automation**: Built Makefile for one-command test execution
5. **CI/CD Ready**: Tests can easily integrate into GitHub Actions

### Best Practices Learned
1. Always isolate test data
2. Use table-driven tests for maintainability
3. Test error cases, not just happy paths
4. Keep tests fast with in-memory databases
5. Automate everything with make commands

## 7. Project Structure
```
practical6-testing/
├── user-service/
│   ├── grpc/
│   │   ├── server.go
│   │   └── server_test.go     
│   ├── database/
│   ├── models/
│   ├── proto/userv1/
│   └── main.go
├── menu-service/
│   ├── grpc/
│   │   ├── server.go
│   │   └── server_test.go      
│   ├── database/
│   ├── models/
│   ├── proto/menuv1/
│   └── main.go
├── order-service/
│   ├── grpc/
│   │   ├── server.go
│   │   └── server_test.go     
│   ├── database/
│   ├── models/
│   ├── proto/orderv1/
│   └── main.go
├── tests/
│   └── integration/
│       ├── integration_test.go 
│       └── go.mod
├── proto/                       
│   ├── user.proto
│   ├── menu.proto
│   └── order.proto
├── Makefile                     
├── docker-compose.yml           
├── generate_proto.sh           
├── README.md
├── TEST_RESULTS.md
└── FINAL_REPORT.md             
```


## 8. How to Run Tests

### Prerequisites
- Go 1.24+
- Protocol Buffers compiler

### Quick Start
```bash
# Run all tests
make test-all

# Run specific test suites
make test-unit           # Unit tests only
make test-integration    # Integration tests only
make test-unit-user      # User service only

# Generate coverage
make test-coverage
```

### Expected Output
```
Running Unit Tests
PASS user-service/grpc
PASS menu-service/grpc
PASS order-service/grpc

Running Integration Tests
PASS integration-tests

All Tests Completed 
```

## 9. Future Enhancements

While the core testing is complete, potential improvements include:

### Not Implemented (But Prepared)
1. **E2E Tests**: Infrastructure ready (docker-compose.yml created)
2. **API Gateway**: HTTP handlers need implementation
3. **Coverage Reporting**: Command exists (`make test-coverage`)
4. **CI/CD Pipeline**: Tests ready for GitHub Actions integration

### Possible Extensions
1. Performance/load testing with k6
2. Contract testing with Pact
3. Chaos engineering tests
4. Security testing (SQL injection, XSS)
5. Database migration testing

