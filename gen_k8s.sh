#!/bin/bash
# Phase G: 누락된 12개 K8s 서비스 매니페스트 생성
BASE=/home/kangjh3kang/Manpasik/infrastructure/kubernetes/base/services

# gRPC 서비스 (7개)
declare -A GRPC_SERVICES=(
  ["audit-service"]="50072"
  ["digital-twin-service"]="50073"
  ["assistant-service"]="50077"
  ["concept-service"]="50078"
  ["cartridge-store-service"]="50079"
  ["data-platform-service"]="50080"
  ["voice-profile-service"]="50081"
)

for svc in "${!GRPC_SERVICES[@]}"; do
  port="${GRPC_SERVICES[$svc]}"
  cat > "$BASE/$svc.yaml" << YAML
apiVersion: apps/v1
kind: Deployment
metadata:
  name: $svc
  labels:
    app: $svc
    tier: backend
spec:
  replicas: 2
  selector:
    matchLabels:
      app: $svc
  template:
    metadata:
      labels:
        app: $svc
        tier: backend
    spec:
      containers:
        - name: $svc
          image: manpasik/$svc:latest
          ports:
            - containerPort: $port
              name: grpc
            - containerPort: 9100
              name: metrics
          envFrom:
            - configMapRef:
                name: manpasik-config
            - secretRef:
                name: manpasik-secrets
          resources:
            requests:
              cpu: 100m
              memory: 128Mi
            limits:
              cpu: 500m
              memory: 512Mi
          readinessProbe:
            exec:
              command: ["/bin/grpc_health_probe", "-addr=:$port"]
            initialDelaySeconds: 5
            periodSeconds: 10
          livenessProbe:
            exec:
              command: ["/bin/grpc_health_probe", "-addr=:$port"]
            initialDelaySeconds: 15
            periodSeconds: 20
---
apiVersion: v1
kind: Service
metadata:
  name: $svc
spec:
  selector:
    app: $svc
  ports:
    - port: $port
      targetPort: $port
      name: grpc
    - port: 9100
      targetPort: 9100
      name: metrics
YAML
  echo "Created: $svc.yaml (gRPC :$port)"
done

# HTTP-only 서비스 (5개)
HTTP_SERVICES=("analytics-service" "emergency-service" "iot-gateway-service" "marketplace-service" "nlp-service")

for svc in "${HTTP_SERVICES[@]}"; do
  cat > "$BASE/$svc.yaml" << YAML
apiVersion: apps/v1
kind: Deployment
metadata:
  name: $svc
  labels:
    app: $svc
    tier: backend
spec:
  replicas: 2
  selector:
    matchLabels:
      app: $svc
  template:
    metadata:
      labels:
        app: $svc
        tier: backend
    spec:
      containers:
        - name: $svc
          image: manpasik/$svc:latest
          ports:
            - containerPort: 8080
              name: http
            - containerPort: 9100
              name: metrics
          envFrom:
            - configMapRef:
                name: manpasik-config
            - secretRef:
                name: manpasik-secrets
          resources:
            requests:
              cpu: 100m
              memory: 128Mi
            limits:
              cpu: 500m
              memory: 512Mi
          readinessProbe:
            httpGet:
              path: /health
              port: 9100
            initialDelaySeconds: 5
            periodSeconds: 10
          livenessProbe:
            httpGet:
              path: /health
              port: 9100
            initialDelaySeconds: 15
            periodSeconds: 20
---
apiVersion: v1
kind: Service
metadata:
  name: $svc
spec:
  selector:
    app: $svc
  ports:
    - port: 8080
      targetPort: 8080
      name: http
    - port: 9100
      targetPort: 9100
      name: metrics
YAML
  echo "Created: $svc.yaml (HTTP :8080)"
done

echo "Done: 12 K8s service manifests created"
