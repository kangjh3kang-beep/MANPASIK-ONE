#!/bin/bash
export PATH=/usr/local/go/bin:/usr/bin:/bin
cd /home/kangjh3kang/Manpasik/backend/services
total=0
for svc in admin-service ai-inference-service analytics-service auth-service calibration-service cartridge-service coaching-service community-service device-service emergency-service family-service gateway health-record-service iot-gateway-service marketplace-service measurement-service nlp-service notification-service payment-service prescription-service reservation-service shop-service subscription-service telemedicine-service translation-service user-service video-service vision-service; do
  cnt=$(cd "/home/kangjh3kang/Manpasik/backend/services/$svc" && GOWORK=off go test ./... -v 2>&1 | grep -c 'PASS:')
  total=$((total + cnt))
  echo "$svc: $cnt"
done
echo "=== TOTAL: $total ==="
