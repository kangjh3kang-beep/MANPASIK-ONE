#!/bin/bash
cd /home/kangjh3kang/Manpasik/backend/services
echo "=== Memory vs Postgres Method Coverage ==="
printf "%-30s %6s %6s %8s\n" "SERVICE" "MEM" "PG" "STATUS"
printf "%-30s %6s %6s %8s\n" "-------" "---" "---" "------"
for svc in */; do
  svc=${svc%/}
  mem_dir="$svc/internal/repository/memory"
  pg_dir="$svc/internal/repository/postgres"
  mem_n=0; pg_n=0
  if [ -d "$mem_dir" ]; then
    mem_n=$(grep -rch '^func ' "$mem_dir"/*.go 2>/dev/null | awk '{s+=$1}END{print s+0}')
  fi
  if [ -d "$pg_dir" ]; then
    pg_n=$(grep -rch '^func ' "$pg_dir"/*.go 2>/dev/null | awk '{s+=$1}END{print s+0}')
  fi
  if [ "$pg_n" -eq 0 ]; then
    status="NO_PG"
  elif [ "$pg_n" -ge "$mem_n" ]; then
    status="OK"
  else
    status="GAP"
  fi
  printf "%-30s %6d %6d %8s\n" "$svc" "$mem_n" "$pg_n" "$status"
done
