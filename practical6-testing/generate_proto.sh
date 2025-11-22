#!/bin/bash

# Generate user service proto
protoc --go_out=user-service --go_opt=paths=source_relative \
  --go-grpc_out=user-service --go-grpc_opt=paths=source_relative \
  --go_opt=Mproto/user.proto=user-service/proto/userv1 \
  -I=. \
  proto/user.proto

# Generate menu service proto
protoc --go_out=menu-service --go_opt=paths=source_relative \
  --go-grpc_out=menu-service --go-grpc_opt=paths=source_relative \
  --go_opt=Mproto/menu.proto=menu-service/proto/menuv1 \
  -I=. \
  proto/menu.proto

# Generate order service proto
protoc --go_out=order-service --go_opt=paths=source_relative \
  --go-grpc_out=order-service --go-grpc_opt=paths=source_relative \
  --go_opt=Mproto/order.proto=order-service/proto/orderv1 \
  -I=. \
  proto/order.proto

# Move generated files to correct locations
mkdir -p user-service/proto/userv1
mkdir -p menu-service/proto/menuv1
mkdir -p order-service/proto/orderv1

mv user-service/proto/user.pb.go user-service/proto/userv1/ 2>/dev/null
mv user-service/proto/user_grpc.pb.go user-service/proto/userv1/ 2>/dev/null

mv menu-service/proto/menu.pb.go menu-service/proto/menuv1/ 2>/dev/null
mv menu-service/proto/menu_grpc.pb.go menu-service/proto/menuv1/ 2>/dev/null

mv order-service/proto/order.pb.go order-service/proto/orderv1/ 2>/dev/null
mv order-service/proto/order_grpc.pb.go order-service/proto/orderv1/ 2>/dev/null

echo "Proto code generation complete!"
