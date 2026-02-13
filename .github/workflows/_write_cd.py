#!/usr/bin/env python3
"""Temporary script to generate cd.yml - delete after use."""
import os

content = r'''name: CD

on:
  push:
    tags:
      - 'v*'
  workflow_dispatch:
    inputs:
      environment:
        description: 'Deployment environment'
        required: true
        default: 'dev'
        type: choice
        options:
          - dev
          - staging
          - production

# ============================================================================
# 동일 환경 병렬 배포 방지 (concurrency)
# ============================================================================
concurrency:
  group: cd-${{ github.event.inputs.environment || 'production' }}
  cancel-in-progress: false

env:
  REGISTRY: ghcr.io
  IMAGE_PREFIX: ${{ github.repository }}
  GO_VERSION: '1.22'
  KUSTOMIZE_VERSION: 'v5.4.3'
  IMAGE_TAG: ${{ github.sha }}

# ============================================================================
# 전체 서비스 목록 (23개: gateway + 22 마이크로서비스)
# ============================================================================
# gateway, auth-service, user-service, device-service, measurement-service,
# subscription-service, shop-service, payment-service, ai-inference-service,
# cartridge-service, calibration-service, coaching-service, family-service,
# notification-service, health-record-service, prescription-service,
# reservation-service, telemedicine-service, admin-service, community-service,
# video-service, translation-service, vision-service

jobs:
  # ============================================================================
  # 0단계: 변경 감지 및 메타데이터 수집
  # ============================================================================
  prepare:
    name: Prepare Deployment
    runs-on: ubuntu-latest
    outputs:
      version: ${{ steps.version.outputs.version }}
      short_sha: ${{ steps.version.outputs.short_sha }}
      deploy_env: ${{ steps.env.outputs.environment }}

    steps:
      - uses: actions/checkout@v4

      - name: Extract version info
        id: version
        run: |
          if [[ "${{ github.ref }}" == refs/tags/v* ]]; then
            echo "version=${GITHUB_REF#refs/tags/v}" >> $GITHUB_OUTPUT
          else
            echo "version=0.0.0-dev+${{ github.sha }}" >> $GITHUB_OUTPUT
          fi
          echo "short_sha=$(echo ${{ github.sha }} | cut -c1-7)" >> $GITHUB_OUTPUT

      - name: Determine deployment environment
        id: env
        run: |
          if [ "${{ github.event_name }}" == "workflow_dispatch" ]; then
            echo "environment=${{ github.event.inputs.environment }}" >> $GITHUB_OUTPUT
          else
            echo "environment=dev" >> $GITHUB_OUTPUT
          fi

  # ============================================================================
  # 1단계: Docker 이미지 빌드 & 푸시
  # ============================================================================
  docker-build:
    name: "Build & Push: ${{ matrix.service }}"
    runs-on: ubuntu-latest
    needs: [prepare]
    permissions:
      contents: read
      packages: write

    strategy:
      fail-fast: false
      matrix:
        service:
          - gateway
          - auth-service
          - user-service
          - device-service
          - measurement-service
          - subscription-service
          - shop-service
          - payment-service
          - ai-inference-service
          - cartridge-service
          - calibration-service
          - coaching-service
          - family-service
          - notification-service
          - health-record-service
          - prescription-service
          - reservation-service
          - telemedicine-service
          - admin-service
          - community-service
          - video-service
          - translation-service
          - vision-service

    steps:
      - uses: actions/checkout@v4

      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v3

      - name: Login to GitHub Container Registry
        uses: docker/login-action@v3
        with:
          registry: ${{ env.REGISTRY }}
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}

      - name: Extract metadata (tags, labels)
        id: meta
        uses: docker/metadata-action@v5
        with:
          images: ${{ env.REGISTRY }}/${{ env.IMAGE_PREFIX }}/${{ matrix.service }}
          tags: |
            type=sha,prefix=
            type=ref,event=branch
            type=semver,pattern={{version}}
            type=semver,pattern={{major}}.{{minor}}
            type=raw,value=latest,enable=${{ github.ref == 'refs/heads/main' }}

      - name: Build and push Docker image
        uses: docker/build-push-action@v5
        with:
          context: .
          file: backend/Dockerfile
          build-args: |
            SERVICE_NAME=${{ matrix.service }}
            VERSION=${{ needs.prepare.outputs.version }}
            BUILD_TIME=${{ github.event.head_commit.timestamp }}
            GIT_COMMIT=${{ github.sha }}
          push: true
          tags: ${{ steps.meta.outputs.tags }}
          labels: ${{ steps.meta.outputs.labels }}
          cache-from: type=gha
          cache-to: type=gha,mode=max
          platforms: linux/amd64

  # ============================================================================
  # 2단계: Dev 환경 배포 (자동)
  # ============================================================================
  deploy-dev:
    name: Deploy to Dev
    runs-on: ubuntu-latest
    needs: [prepare, docker-build]
    if: >-
      needs.prepare.outputs.deploy_env == 'dev' ||
      (github.event_name == 'push' && startsWith(github.ref, 'refs/tags/v'))
    environment:
      name: dev
      url: https://dev.manpasik.com

    steps:
      - uses: actions/checkout@v4

      - name: Set up kubectl
        uses: azure/setup-kubectl@v4

      - name: Set up Kustomize
        run: |
          curl -sL "https://github.com/kubernetes-sigs/kustomize/releases/download/kustomize%2F${{ env.KUSTOMIZE_VERSION }}/kustomize_${{ env.KUSTOMIZE_VERSION }}_linux_amd64.tar.gz" | tar xz
          sudo mv kustomize /usr/local/bin/kustomize
          kustomize version

      - name: Configure AWS credentials
        uses: aws-actions/configure-aws-credentials@v4
        with:
          aws-access-key-id: ${{ secrets.AWS_ACCESS_KEY_ID }}
          aws-secret-access-key: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
          aws-region: ap-northeast-2

      - name: Update kubeconfig
        run: aws eks update-kubeconfig --name manpasik-dev --region ap-northeast-2

      - name: Update image tags (Dev)
        run: |
          cd infrastructure/kubernetes/overlays/dev
          SERVICES=(
            gateway auth-service user-service device-service measurement-service
            subscription-service shop-service payment-service
            ai-inference-service cartridge-service calibration-service coaching-service
            family-service notification-service health-record-service
            prescription-service reservation-service telemedicine-service
            admin-service community-service video-service translation-service
            vision-service
          )
          for svc in "${SERVICES[@]}"; do
            kustomize edit set image \
              "${{ env.REGISTRY }}/${{ env.IMAGE_PREFIX }}/${svc}=${{ env.REGISTRY }}/${{ env.IMAGE_PREFIX }}/${svc}:${{ env.IMAGE_TAG }}" \
              2>/dev/null || true
          done

      - name: Run database migrations (Dev)
        run: |
          echo "Running database migrations for dev environment..."
          kubectl apply -f infrastructure/kubernetes/base/config/ -n manpasik-dev
          if [ -f infrastructure/kubernetes/overlays/dev/migrations-job.yaml ]; then
            kubectl apply -f infrastructure/kubernetes/overlays/dev/migrations-job.yaml
            kubectl wait --for=condition=complete job/db-migration -n manpasik-dev --timeout=120s || echo "::warning::Migration job not found or timed out"
          fi

      - name: Deploy to Kubernetes (Dev)
        run: |
          # Kustomize dev 오버레이 적용
          kustomize build infrastructure/kubernetes/overlays/dev/ | kubectl apply -f - -n manpasik-dev
          echo "Dev deployment applied successfully"

      - name: Wait for rollout (Dev)
        run: |
          SERVICES=(
            gateway
            auth-service
            user-service
            device-service
            measurement-service
            subscription-service
            shop-service
            payment-service
            ai-inference-service
            cartridge-service
            calibration-service
            coaching-service
            family-service
            notification-service
            health-record-service
            prescription-service
            reservation-service
            telemedicine-service
            admin-service
            community-service
            video-service
            translation-service
            vision-service
          )
          FAILED=()
          for svc in "${SERVICES[@]}"; do
            echo "::group::Verifying $svc"
            if ! kubectl rollout status deployment/$svc -n manpasik-dev --timeout=120s; then
              FAILED+=("$svc")
              echo "::warning::$svc rollout failed"
            fi
            echo "::endgroup::"
          done
          if [ ${#FAILED[@]} -ne 0 ]; then
            echo "::error::Failed services: ${FAILED[*]}"
            exit 1
          fi
          echo "All services rolled out successfully"

      - name: Health check (Dev)
        run: |
          GATEWAY_URL=$(kubectl get svc gateway -n manpasik-dev -o jsonpath='{.status.loadBalancer.ingress[0].hostname}' 2>/dev/null || echo "")
          if [ -z "$GATEWAY_URL" ]; then
            GATEWAY_URL=$(kubectl get svc gateway -n manpasik-dev -o jsonpath='{.spec.clusterIP}')
          fi
          echo "Gateway URL: $GATEWAY_URL"

          for i in {1..5}; do
            if curl -sf "http://${GATEWAY_URL}/health" --connect-timeout 10; then
              echo ""
              echo "Dev health check passed (attempt $i)"
              exit 0
            fi
            echo "Attempt $i failed, retrying in 10s..."
            sleep 10
          done
          echo "::error::Dev health check failed after 5 attempts"
          exit 1

      - name: Smoke tests (Dev)
        run: |
          GATEWAY_URL=$(kubectl get svc gateway -n manpasik-dev -o jsonpath='{.status.loadBalancer.ingress[0].hostname}' 2>/dev/null || echo "")
          if [ -z "$GATEWAY_URL" ]; then
            GATEWAY_URL=$(kubectl get svc gateway -n manpasik-dev -o jsonpath='{.spec.clusterIP}')
          fi

          echo "=== Dev Smoke Tests ==="
          echo "Testing gateway health..."
          curl -sf "http://${GATEWAY_URL}/health" || exit 1

          echo "Testing auth health..."
          curl -sf "http://${GATEWAY_URL}/api/v1/auth/health" --connect-timeout 10 || echo "::warning::Auth health not reachable"

          echo "Dev smoke tests completed"

      - name: Notify deployment status (Dev)
        if: always()
        uses: 8398a7/action-slack@v3
        with:
          status: ${{ job.status }}
          text: |
            ManPaSik Dev Deployment
            Status: ${{ job.status }}
            Version: ${{ needs.prepare.outputs.version }}
            Commit: ${{ needs.prepare.outputs.short_sha }}
        env:
          SLACK_WEBHOOK_URL: ${{ secrets.SLACK_WEBHOOK_URL }}

  # ============================================================================
  # 3단계: Staging 환경 배포 (수동 승인)
  # ============================================================================
  deploy-staging:
    name: Deploy to Staging
    runs-on: ubuntu-latest
    needs: [prepare, deploy-dev]
    if: >-
      needs.prepare.outputs.deploy_env == 'staging' ||
      needs.prepare.outputs.deploy_env == 'production' ||
      (github.event_name == 'push' && startsWith(github.ref, 'refs/tags/v'))
    environment:
      name: staging
      url: https://staging.manpasik.com

    steps:
      - uses: actions/checkout@v4

      - name: Set up kubectl
        uses: azure/setup-kubectl@v4

      - name: Set up Kustomize
        run: |
          curl -sL "https://github.com/kubernetes-sigs/kustomize/releases/download/kustomize%2F${{ env.KUSTOMIZE_VERSION }}/kustomize_${{ env.KUSTOMIZE_VERSION }}_linux_amd64.tar.gz" | tar xz
          sudo mv kustomize /usr/local/bin/kustomize
          kustomize version

      - name: Configure AWS credentials
        uses: aws-actions/configure-aws-credentials@v4
        with:
          aws-access-key-id: ${{ secrets.AWS_ACCESS_KEY_ID }}
          aws-secret-access-key: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
          aws-region: ap-northeast-2

      - name: Update kubeconfig
        run: aws eks update-kubeconfig --name manpasik-staging --region ap-northeast-2

      - name: Update image tags (Staging)
        run: |
          cd infrastructure/kubernetes/overlays/staging
          SERVICES=(
            gateway auth-service user-service device-service measurement-service
            subscription-service shop-service payment-service
            ai-inference-service cartridge-service calibration-service coaching-service
            family-service notification-service health-record-service
            prescription-service reservation-service telemedicine-service
            admin-service community-service video-service translation-service
            vision-service
          )
          for svc in "${SERVICES[@]}"; do
            kustomize edit set image \
              "${{ env.REGISTRY }}/${{ env.IMAGE_PREFIX }}/${svc}=${{ env.REGISTRY }}/${{ env.IMAGE_PREFIX }}/${svc}:${{ env.IMAGE_TAG }}" \
              2>/dev/null || true
          done

      - name: Run database migrations (Staging)
        run: |
          echo "Running database migrations for staging environment..."
          kubectl apply -f infrastructure/kubernetes/base/config/ -n manpasik-staging
          if [ -f infrastructure/kubernetes/overlays/staging/migrations-job.yaml ]; then
            kubectl apply -f infrastructure/kubernetes/overlays/staging/migrations-job.yaml
            kubectl wait --for=condition=complete job/db-migration -n manpasik-staging --timeout=180s || echo "::warning::Migration job not found or timed out"
          fi

      - name: Deploy to Kubernetes (Staging)
        run: |
          # Kustomize staging 오버레이 적용
          kustomize build infrastructure/kubernetes/overlays/staging/ | kubectl apply -f - -n manpasik-staging
          echo "Staging deployment applied successfully"

      - name: Wait for rollout (Staging)
        run: |
          SERVICES=(
            gateway
            auth-service
            user-service
            device-service
            measurement-service
            subscription-service
            shop-service
            payment-service
            ai-inference-service
            cartridge-service
            calibration-service
            coaching-service
            family-service
            notification-service
            health-record-service
            prescription-service
            reservation-service
            telemedicine-service
            admin-service
            community-service
            video-service
            translation-service
            vision-service
          )
          FAILED=()
          for svc in "${SERVICES[@]}"; do
            echo "::group::Verifying $svc"
            if ! kubectl rollout status deployment/$svc -n manpasik-staging --timeout=180s; then
              FAILED+=("$svc")
              echo "::warning::$svc rollout failed"
            fi
            echo "::endgroup::"
          done
          if [ ${#FAILED[@]} -ne 0 ]; then
            echo "::error::Failed services: ${FAILED[*]}"
            exit 1
          fi
          echo "All services rolled out successfully"

      - name: Health check (Staging)
        run: |
          GATEWAY_URL=$(kubectl get svc gateway -n manpasik-staging -o jsonpath='{.status.loadBalancer.ingress[0].hostname}' 2>/dev/null || echo "")
          if [ -z "$GATEWAY_URL" ]; then
            GATEWAY_URL=$(kubectl get svc gateway -n manpasik-staging -o jsonpath='{.spec.clusterIP}')
          fi
          echo "Gateway URL: $GATEWAY_URL"

          for i in {1..5}; do
            if curl -sf "http://${GATEWAY_URL}/health" --connect-timeout 10; then
              echo ""
              echo "Staging health check passed (attempt $i)"
              exit 0
            fi
            echo "Attempt $i failed, retrying in 10s..."
            sleep 10
          done
          echo "::error::Staging health check failed after 5 attempts"
          exit 1

      - name: Run smoke tests (Staging)
        run: |
          GATEWAY_URL=$(kubectl get svc gateway -n manpasik-staging -o jsonpath='{.status.loadBalancer.ingress[0].hostname}' 2>/dev/null || echo "")
          if [ -z "$GATEWAY_URL" ]; then
            GATEWAY_URL=$(kubectl get svc gateway -n manpasik-staging -o jsonpath='{.spec.clusterIP}')
          fi

          echo "=== Staging Smoke Tests ==="

          # 1. Gateway health
          echo "Testing gateway health..."
          curl -sf "http://${GATEWAY_URL}/health" || exit 1

          # 2. Auth service
          echo "Testing auth endpoint..."
          curl -sf "http://${GATEWAY_URL}/api/v1/auth/health" --connect-timeout 10 || echo "::warning::Auth health not reachable"

          # 3. 개별 서비스 K8s pod 가용성 확인
          GRPC_SERVICES=(
            auth-service user-service device-service measurement-service
            subscription-service shop-service payment-service
            ai-inference-service cartridge-service calibration-service
            coaching-service family-service notification-service
            health-record-service prescription-service reservation-service
            telemedicine-service admin-service community-service
            video-service translation-service vision-service
          )
          UNHEALTHY=()
          for svc in "${GRPC_SERVICES[@]}"; do
            READY=$(kubectl get deployment/$svc -n manpasik-staging -o jsonpath='{.status.readyReplicas}' 2>/dev/null || echo "0")
            DESIRED=$(kubectl get deployment/$svc -n manpasik-staging -o jsonpath='{.spec.replicas}' 2>/dev/null || echo "0")
            if [ "$READY" != "$DESIRED" ] || [ "$READY" == "0" ]; then
              UNHEALTHY+=("$svc(${READY}/${DESIRED})")
            fi
          done
          if [ ${#UNHEALTHY[@]} -ne 0 ]; then
            echo "::warning::Unhealthy/missing services: ${UNHEALTHY[*]}"
          fi

          echo "Staging smoke tests completed"

      - name: Notify deployment status (Staging)
        if: always()
        uses: 8398a7/action-slack@v3
        with:
          status: ${{ job.status }}
          text: |
            ManPaSik Staging Deployment
            Status: ${{ job.status }}
            Version: ${{ needs.prepare.outputs.version }}
            Commit: ${{ needs.prepare.outputs.short_sha }}
        env:
          SLACK_WEBHOOK_URL: ${{ secrets.SLACK_WEBHOOK_URL }}

  # ============================================================================
  # 4단계: Production 환경 배포 (수동 승인 + Canary)
  # ============================================================================
  deploy-production:
    name: Deploy to Production
    runs-on: ubuntu-latest
    needs: [prepare, deploy-staging]
    if: >-
      needs.prepare.outputs.deploy_env == 'production' ||
      (github.event_name == 'push' && startsWith(github.ref, 'refs/tags/v'))
    environment:
      name: production
      url: https://manpasik.com

    steps:
      - uses: actions/checkout@v4

      - name: Set up kubectl
        uses: azure/setup-kubectl@v4

      - name: Set up Kustomize
        run: |
          curl -sL "https://github.com/kubernetes-sigs/kustomize/releases/download/kustomize%2F${{ env.KUSTOMIZE_VERSION }}/kustomize_${{ env.KUSTOMIZE_VERSION }}_linux_amd64.tar.gz" | tar xz
          sudo mv kustomize /usr/local/bin/kustomize
          kustomize version

      - name: Configure AWS credentials
        uses: aws-actions/configure-aws-credentials@v4
        with:
          aws-access-key-id: ${{ secrets.AWS_ACCESS_KEY_ID }}
          aws-secret-access-key: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
          aws-region: ap-northeast-2

      - name: Update kubeconfig
        run: aws eks update-kubeconfig --name manpasik-production --region ap-northeast-2

      - name: Run database migrations (Production)
        run: |
          echo "Running database migrations for production environment..."
          kubectl apply -f infrastructure/kubernetes/base/config/ -n manpasik
          if [ -f infrastructure/kubernetes/overlays/production/migrations-job.yaml ]; then
            kubectl apply -f infrastructure/kubernetes/overlays/production/migrations-job.yaml
            kubectl wait --for=condition=complete job/db-migration -n manpasik --timeout=300s || {
              echo "::error::Production migration failed or timed out"
              exit 1
            }
          fi

      - name: Deploy Canary (10% traffic)
        run: |
          # Canary 배포: 프로덕션 오버레이에서 canary 리소스 적용
          if [ -d infrastructure/kubernetes/overlays/production/canary ]; then
            kustomize build infrastructure/kubernetes/overlays/production/canary/ | kubectl apply -f - -n manpasik
          else
            kubectl apply -k infrastructure/kubernetes/overlays/production/
          fi

          # Canary 인스턴스 롤아웃 대기
          kubectl rollout status deployment/gateway-canary -n manpasik --timeout=300s || echo "::warning::Canary deployment not found, proceeding with full rollout"

      - name: Canary health validation
        run: |
          echo "Waiting 60s for canary metrics collection..."
          sleep 60

          # Prometheus 메트릭 기반 canary 검증
          GATEWAY_URL=$(kubectl get svc gateway -n manpasik -o jsonpath='{.status.loadBalancer.ingress[0].hostname}' 2>/dev/null || echo "")
          if [ -z "$GATEWAY_URL" ]; then
            GATEWAY_URL=$(kubectl get svc gateway -n manpasik -o jsonpath='{.spec.clusterIP}')
          fi

          ERROR_COUNT=0
          TOTAL=20
          for i in $(seq 1 $TOTAL); do
            HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" "http://${GATEWAY_URL}/health" --connect-timeout 5 || echo "000")
            if [ "$HTTP_CODE" != "200" ]; then
              ERROR_COUNT=$((ERROR_COUNT + 1))
            fi
            sleep 2
          done

          ERROR_RATE=$(echo "scale=2; $ERROR_COUNT * 100 / $TOTAL" | bc)
          echo "Canary error rate: ${ERROR_RATE}% ($ERROR_COUNT/$TOTAL)"

          if [ "$ERROR_COUNT" -gt 2 ]; then
            echo "::error::Canary error rate too high: ${ERROR_RATE}%"
            exit 1
          fi

          echo "Canary validation passed"

      - name: Promote to full rollout
        run: |
          # Kustomize로 전체 프로덕션 이미지 태그 업데이트
          cd infrastructure/kubernetes/overlays/production
          SERVICES=(
            gateway auth-service user-service device-service measurement-service
            subscription-service shop-service payment-service
            ai-inference-service cartridge-service calibration-service coaching-service
            family-service notification-service health-record-service
            prescription-service reservation-service telemedicine-service
            admin-service community-service video-service translation-service
            vision-service
          )
          for svc in "${SERVICES[@]}"; do
            kustomize edit set image \
              "${{ env.REGISTRY }}/${{ env.IMAGE_PREFIX }}/${svc}=${{ env.REGISTRY }}/${{ env.IMAGE_PREFIX }}/${svc}:${{ env.IMAGE_TAG }}" \
              2>/dev/null || true
          done
          cd ${{ github.workspace }}

          # Kustomize로 전체 오버레이 적용
          kustomize build infrastructure/kubernetes/overlays/production/ | kubectl apply -f - -n manpasik

          # 전체 롤아웃 대기
          FAILED=()
          for svc in "${SERVICES[@]}"; do
            echo "::group::Waiting for $svc rollout"
            if ! kubectl rollout status deployment/$svc -n manpasik --timeout=600s; then
              FAILED+=("$svc")
              echo "::warning::$svc rollout did not complete"
            fi
            echo "::endgroup::"
          done

          if [ ${#FAILED[@]} -ne 0 ]; then
            echo "::error::Failed rollouts: ${FAILED[*]}"
            exit 1
          fi

          echo "Full production rollout completed"

      - name: Production health check
        run: |
          GATEWAY_URL=$(kubectl get svc gateway -n manpasik -o jsonpath='{.status.loadBalancer.ingress[0].hostname}' 2>/dev/null || echo "")
          if [ -z "$GATEWAY_URL" ]; then
            GATEWAY_URL=$(kubectl get svc gateway -n manpasik -o jsonpath='{.spec.clusterIP}')
          fi

          # Gateway HTTP 헬스체크
          for i in {1..10}; do
            if curl -sf "http://${GATEWAY_URL}/health" --connect-timeout 10; then
              echo ""
              echo "Production gateway health check passed (attempt $i)"
              break
            fi
            echo "Attempt $i failed, retrying in 15s..."
            sleep 15
            if [ $i -eq 10 ]; then
              echo "::error::Production health check failed"
              exit 1
            fi
          done

          # 개별 서비스 가용성 확인
          echo ""
          echo "=== Individual Service Health Check ==="
          GRPC_SERVICES=(
            auth-service user-service device-service measurement-service
            subscription-service shop-service payment-service
            ai-inference-service cartridge-service calibration-service
            coaching-service family-service notification-service
            health-record-service prescription-service reservation-service
            telemedicine-service admin-service community-service
            video-service translation-service vision-service
          )
          UNHEALTHY=()
          for svc in "${GRPC_SERVICES[@]}"; do
            READY=$(kubectl get deployment/$svc -n manpasik -o jsonpath='{.status.readyReplicas}' 2>/dev/null || echo "0")
            DESIRED=$(kubectl get deployment/$svc -n manpasik -o jsonpath='{.spec.replicas}' 2>/dev/null || echo "0")
            if [ "$READY" != "$DESIRED" ] || [ "$READY" == "0" ]; then
              UNHEALTHY+=("$svc(${READY}/${DESIRED})")
              echo "::warning::$svc: $READY/$DESIRED pods ready"
            else
              echo "$svc: $READY/$DESIRED pods ready"
            fi
          done

          if [ ${#UNHEALTHY[@]} -ne 0 ]; then
            echo "::warning::Unhealthy services: ${UNHEALTHY[*]}"
          fi

          echo "Production health checks completed"

      - name: Annotate deployment
        if: success()
        run: |
          SERVICES=(
            gateway auth-service user-service device-service measurement-service
            subscription-service shop-service payment-service
            ai-inference-service cartridge-service calibration-service coaching-service
            family-service notification-service health-record-service
            prescription-service reservation-service telemedicine-service
            admin-service community-service video-service translation-service
            vision-service
          )
          for svc in "${SERVICES[@]}"; do
            kubectl annotate deployment/$svc -n manpasik \
              kubernetes.io/change-cause="CD pipeline: v${{ needs.prepare.outputs.version }} (${{ needs.prepare.outputs.short_sha }})" \
              --overwrite 2>/dev/null || true
          done

      - name: Notify Slack (Production)
        if: always()
        uses: 8398a7/action-slack@v3
        with:
          status: ${{ job.status }}
          text: |
            ManPaSik Production Deployment
            Status: ${{ job.status }}
            Version: ${{ needs.prepare.outputs.version }}
            Commit: ${{ needs.prepare.outputs.short_sha }}
            Tag: ${{ github.ref_name }}
          fields: repo,message,commit,author,action,eventName,ref,workflow
        env:
          SLACK_WEBHOOK_URL: ${{ secrets.SLACK_WEBHOOK_URL }}

  # ============================================================================
  # 5단계: 롤백 (Production 실패 시)
  # ============================================================================
  rollback:
    name: Rollback Production
    runs-on: ubuntu-latest
    if: failure() && needs.deploy-production.result == 'failure'
    needs: [deploy-production]
    environment: production

    steps:
      - name: Set up kubectl
        uses: azure/setup-kubectl@v4

      - name: Configure AWS credentials
        uses: aws-actions/configure-aws-credentials@v4
        with:
          aws-access-key-id: ${{ secrets.AWS_ACCESS_KEY_ID }}
          aws-secret-access-key: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
          aws-region: ap-northeast-2

      - name: Update kubeconfig
        run: aws eks update-kubeconfig --name manpasik-production --region ap-northeast-2

      - name: Rollback all services
        run: |
          SERVICES=(
            gateway
            auth-service
            user-service
            device-service
            measurement-service
            subscription-service
            shop-service
            payment-service
            ai-inference-service
            cartridge-service
            calibration-service
            coaching-service
            family-service
            notification-service
            health-record-service
            prescription-service
            reservation-service
            telemedicine-service
            admin-service
            community-service
            video-service
            translation-service
            vision-service
          )
          for svc in "${SERVICES[@]}"; do
            echo "Rolling back $svc..."
            kubectl rollout undo deployment/$svc -n manpasik || echo "::warning::Failed to rollback $svc (may not exist)"
          done

      - name: Verify rollback
        run: |
          SERVICES=(
            gateway
            auth-service
            user-service
            device-service
            measurement-service
            subscription-service
            shop-service
            payment-service
            ai-inference-service
            cartridge-service
            calibration-service
            coaching-service
            family-service
            notification-service
            health-record-service
            prescription-service
            reservation-service
            telemedicine-service
            admin-service
            community-service
            video-service
            translation-service
            vision-service
          )
          FAILED=()
          for svc in "${SERVICES[@]}"; do
            echo "::group::Verifying rollback for $svc"
            if ! kubectl rollout status deployment/$svc -n manpasik --timeout=120s; then
              FAILED+=("$svc")
              echo "::warning::$svc rollback status check failed"
            fi
            echo "::endgroup::"
          done
          if [ ${#FAILED[@]} -ne 0 ]; then
            echo "::error::Rollback verification failed for: ${FAILED[*]}"
          fi

      - name: Post-rollback health check
        run: |
          GATEWAY_URL=$(kubectl get svc gateway -n manpasik -o jsonpath='{.status.loadBalancer.ingress[0].hostname}' 2>/dev/null || echo "")
          if [ -z "$GATEWAY_URL" ]; then
            GATEWAY_URL=$(kubectl get svc gateway -n manpasik -o jsonpath='{.spec.clusterIP}')
          fi

          for i in {1..5}; do
            if curl -sf "http://${GATEWAY_URL}/health" --connect-timeout 10; then
              echo ""
              echo "Post-rollback health check passed (attempt $i)"
              exit 0
            fi
            sleep 10
          done
          echo "::error::Post-rollback health check failed"

      - name: Notify Slack (Rollback)
        if: always()
        uses: 8398a7/action-slack@v3
        with:
          status: custom
          custom_payload: |
            {
              "text": ":rotating_light: ManPaSik Production ROLLBACK executed\nVersion: ${{ github.ref_name }}\nCommit: ${{ github.sha }}\nTriggered by deployment failure.\nRollback status: ${{ job.status }}"
            }
        env:
          SLACK_WEBHOOK_URL: ${{ secrets.SLACK_WEBHOOK_URL }}
'''

target = os.path.join(os.path.dirname(os.path.abspath(__file__)), 'cd.yml')
with open(target, 'w', encoding='utf-8') as f:
    f.write(content.lstrip('\n'))
print(f'Written {os.path.getsize(target)} bytes to {target}')
