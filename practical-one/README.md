# Practical 1: Microservices with gRPC and Docker

## Overview

Built two microservices that communicate using gRPC:
- **Time Service** - Returns current timestamp
- **Greeter Service** - Gets time from Time Service and returns personalized greeting

## What I Learned

### Part 1: Development Environment Setup

**Go Programming Language**
- Installed Go 1.23.2 on Linux
- Learned Go is ideal for microservices (fast, concurrent, small binaries)
- Configured GOPATH and workspace structure

**Protocol Buffers & gRPC**
- Installed protoc compiler and Go plugins
- Learned Protocol Buffers define service contracts in language-agnostic way
- Understood how protoc generates type-safe code automatically
- Challenge: Had to create bin directory first before installing plugins

**Docker**
- Installed Docker for containerization
- Learned containers package apps with all dependencies
- Understood Docker provides consistency across environments

### Part 2: Building Microservices

**Service Contracts**
- Defined two .proto files for Time and Greeter services
- Learned how Protocol Buffers act as contracts between services
- Generated 4 Go files from proto definitions

**Implementing Services**
- Built Time Service: Simple gRPC server on port 50052
- Built Greeter Service: Acts as both server and client
- Learned services can call other services using gRPC
- Understood inter-service communication patterns

**Go Modules**
- Created separate go.mod for each service
- Learned about module dependencies and replace directives
- Challenge: Had to fix Go version from 1.24 to 1.23
- Challenge: Generated go.sum files using `go mod tidy`

**Containerization**
- Created multi-stage Dockerfiles for both services
- Learned multi-stage builds = smaller images (10MB vs 300MB)
- First stage: Build with full Go image
- Second stage: Run with minimal Alpine image

**Docker Compose**
- Orchestrated both services in docker-compose.yml
- Learned Docker creates internal network for service communication
- Services find each other by hostname (time-service:50052)
- Used depends_on to control startup order

**Testing**
- Used grpcurl to test the services
- Successfully called Greeter Service which internally called Time Service
- Verified inter-service communication in logs

## Key Concepts I Understood

**Microservices Architecture**
- Small, independent services doing one thing well
- Can develop, deploy, and scale each service independently
- Failures in one service don't crash everything

**gRPC vs REST**
- gRPC uses binary format (faster than JSON)
- Strongly typed (errors caught at compile time)
- Auto-generates client code

**Service Discovery**
- Docker provides internal DNS
- Services use hostnames instead of IP addresses
- greeter-service finds time-service automatically

**Multi-Stage Builds**
- Separates build environment from runtime
- Final image only has compiled binary, not source code
- More secure and efficient

## Problems I Solved

1. **Go plugins not found** → Created bin directory first
2. **Docker network conflict** → Ran `docker-compose down` to cleanup
3. **Go version mismatch** → Changed go.mod from 1.24 to 1.23
4. **Missing go.sum** → Ran `go mod tidy` in each directory
5. **Module path issues** → Created go.mod in proto/gen package

## Project Structure

```
practical-one/
├── proto/              # Service contracts
├── time-service/       # Time microservice
├── greeter-service/    # Greeter microservice
├── docker-compose.yml  # Orchestration
└── screenshots/        # Evidence (12 screenshots)
```

## How to Run

```bash
# Build and start
docker-compose up --build

# Test (in new terminal)
grpcurl -plaintext \
    -import-path ./proto -proto greeter.proto \
    -d '{"name": "Your Name"}' \
    0.0.0.0:50051 greeter.GreeterService/SayHello

# Expected response
{
  "message": "Hello Your Name! The current time is 2025-11-23T..."
}
```

## Screenshots Evidence

1. Go installation verification
2. Protoc installation
3. Go plugins verification
4. Docker hello-world test
5. Project directory structure
6. Generated proto files
7. Docker build output
8. Services running logs
9. grpcurl test success
10. Inter-service communication logs

## Conclusion

This practical taught me:
- How to set up microservices development environment
- How services communicate using gRPC
- How to containerize and orchestrate multiple services
- How to troubleshoot dependency and version issues

Most important: Microservices are independent but work together through well-defined contracts.

