# Student Cafe Microservices - README

## Project Overview
This project demonstrates the refactoring of a monolithic Student Cafe application into independent microservices, showcasing modern distributed system architecture patterns.

## Architecture

### Services
- **API Gateway** (Port 8080) - Single entry point for all client requests
- **User Service** (Port 8081) - Manages user accounts and profiles
- **Menu Service** (Port 8082) - Handles menu items and pricing
- **Order Service** (Port 8083) - Processes orders with inter-service communication
- **Monolith** (Port 8090) - Original application for comparison
- **Consul** (Port 8500) - Service discovery and health monitoring

### Databases
Each service has its own dedicated database (database-per-service pattern):
- `user_db` - User data
- `menu_db` - Menu items
- `order_db` - Orders and order items
- `student_cafe` - Monolith database

## Architecture Diagram
```
┌─────────────┐
│   Clients   │
└──────┬──────┘
       │
       ▼
┌──────────────────┐
│  API Gateway     │ :8080
└────────┬─────────┘
         │
    ┌────┴────┬─────────┐
    ▼         ▼         ▼
┌─────────┐ ┌─────────┐ ┌─────────┐
│  User   │ │  Menu   │ │  Order  │
│ Service │ │ Service │ │ Service │
│  :8081  │ │  :8082  │ │  :8083  │
└────┬────┘ └────┬────┘ └────┬────┘
     │           │           │
     ▼           ▼           ▼
┌─────────┐ ┌─────────┐ ┌─────────┐
│ user_db │ │ menu_db │ │order_db │
└─────────┘ └─────────┘ └─────────┘

        ┌────────────┐
        │  Consul    │ :8500
        │  (Service  │
        │ Discovery) │
        └────────────┘
```

## Service Boundaries Justification

### Why We Split This Way

**User Service**
- **Responsibility**: User registration, authentication, profile management
- **Independence**: User data doesn't depend on menu or orders
- **Scaling**: Can scale independently during registration spikes
- **Changes**: User features evolve separately from business logic

**Menu Service**
- **Responsibility**: Menu catalog, pricing, item descriptions
- **Independence**: Menu is read-heavy and doesn't need order/user data
- **Scaling**: High read traffic during browsing - scale horizontally
- **Changes**: Menu updates don't affect other services

**Order Service**
- **Responsibility**: Order creation, status tracking, order history
- **Dependencies**: Calls user-service and menu-service via HTTP
- **Data Ownership**: Snapshots prices at order time (historical accuracy)
- **Complexity**: Orchestrates business logic across services

## Key Design Patterns Implemented

### 1. Database-Per-Service
Each microservice owns its data. No shared databases.

**Benefits:**
- Services deploy independently
- Database schema changes don't affect others
- Technology diversity (could use different DB types)

**Trade-offs:**
- More complex queries across services
- Data consistency challenges
- Increased infrastructure cost

### 2. API Gateway Pattern
Single entry point routing requests to backend services.

**Benefits:**
- Simplified client interface
- Centralized authentication/logging
- Load balancing
- Backend services hidden from clients

### 3. Inter-Service Communication
Order service calls user-service and menu-service via HTTP REST.

**How it works:**
```
Client → API Gateway → Order Service
                         ├→ User Service (validate user)
                         ├→ Menu Service (get price)
                         └→ Creates order with data
```

### 4. Service Discovery with Consul
Services register with Consul for dynamic discovery.

**Benefits:**
- No hardcoded service addresses
- Health monitoring
- Automatic failover
- Service registry

## Technologies Used
- **Language**: Go 1.23
- **Web Framework**: Chi (lightweight HTTP router)
- **ORM**: GORM
- **Database**: PostgreSQL 13
- **Service Discovery**: HashiCorp Consul
- **Containerization**: Docker & Docker Compose

## Running the Project

### Prerequisites
- Docker Desktop
- Docker Compose
- Go 1.23+ (for local development)

### Start All Services
```bash
cd practicals/practical5
docker-compose up --build
```

### Access Points
- API Gateway: http://localhost:8080
- Consul UI: http://localhost:8500
- User Service (direct): http://localhost:8081
- Menu Service (direct): http://localhost:8082
- Order Service (direct): http://localhost:8083
- Monolith: http://localhost:8090

### Test the System
```bash
# Create a user
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"name": "Alice", "email": "alice@test.com"}'

# Create a menu item
curl -X POST http://localhost:8080/api/menu \
  -H "Content-Type: application/json" \
  -d '{"name": "Coffee", "description": "Hot coffee", "price": 2.50}'

# Create an order (inter-service communication)
curl -X POST http://localhost:8080/api/orders \
  -H "Content-Type: application/json" \
  -d '{"user_id": 1, "items": [{"menu_item_id": 1, "quantity": 2}]}'

# Get all orders
curl http://localhost:8080/api/orders
```

## Challenges Encountered

### 1. Docker Network Configuration
**Problem**: Services couldn't find each other initially
**Solution**: Used Docker Compose networking with service names as DNS

### 2. Database-Per-Service Data Access
**Problem**: Order service needed user and menu data
**Solution**: Implemented HTTP-based inter-service communication

### 3. Port Conflicts
**Problem**: Multiple services competing for same ports
**Solution**: Assigned unique ports to each service (8080-8083, 8090, 8500)

### 4. Service Startup Order
**Problem**: Services started before databases were ready
**Solution**: Added health checks and `depends_on` with conditions in docker-compose

### 5. Go Module Dependencies
**Problem**: Docker builds failed with module checksum errors
**Solution**: Ran `go mod tidy` to regenerate correct checksums

## What I Learned

### Technical Skills
- **Microservices Architecture**: Understanding service boundaries and decomposition
- **Docker & Containerization**: Multi-container orchestration with Docker Compose
- **Service Communication**: REST APIs, HTTP clients, inter-service calls
- **Database Design**: Database-per-service pattern and data isolation
- **Service Discovery**: Dynamic service registration with Consul
- **API Gateway Pattern**: Request routing and single entry point

### Design Principles
- **Separation of Concerns**: Each service has single responsibility
- **Loose Coupling**: Services communicate via APIs, not shared databases
- **High Cohesion**: Related functionality grouped together
- **Independent Deployment**: Services can be updated without affecting others

### Trade-offs Understanding
- **Complexity vs Scalability**: More services = more complexity but better scaling
- **Data Consistency vs Independence**: Distributed data is harder to keep consistent
- **Network Overhead**: Inter-service calls add latency
- **Development Speed**: More initial setup but faster parallel development

## Screenshots for Submission

### Required Screenshots

1. **consul-ui.png**: Consul dashboard at http://localhost:8500
2. **order-creation.png**: Terminal showing successful order creation
3. **inter-service-logs.png**: Logs showing order service calling user/menu services
4. **api-gateway-logs.png**: API Gateway routing logs
5. **all-services-running.png**: `docker-compose ps` output
6. **architecture-overview.png**: All services running together
