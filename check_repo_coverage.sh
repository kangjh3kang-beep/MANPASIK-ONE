#!/bin/bash
cd /home/kangjh3kang/Manpasik/backend/services
for svc in */; do
  svc=${svc%/}
  mem_dir="$svc/internal/repository/memory"
  pg_dir="$svc/internal/repository/postgres"
  if [ -d "$mem_dir" ] && [ -d "$pg_dir" ]; then
    mem_funcs=$(grep -c '^func ' $mem_dir/*.go 2>/dev/null || echo 0)
    pg_funcs=$(grep -c '^func ' $pg_dir/*.go 2>/dev/null || echo 0)
    echo "$svc: mem=$mem_funcs pg=$pg_funcs"
  elif [ -d "$mem_dir" ]; then
    mem_funcs=$(grep -c '^func ' $mem_dir/*.go 2>/dev/null || echo 0)
    echo "$svc: mem=$mem_funcs pg=MISSING"
  elif [ -d "$pg_dir" ]; then
    pg_funcs=$(grep -c '^func ' $pg_dir/*.go 2>/dev/null || echo 0)
    echo "$svc: mem=MISSING pg=$pg_funcs"
  fi
done
