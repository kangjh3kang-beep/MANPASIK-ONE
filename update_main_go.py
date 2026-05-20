#!/usr/bin/env python3
"""13개 서비스의 cmd/main.go에 DB 스위치 로직을 추가합니다."""
import re, os

BASE = "/home/kangjh3kang/Manpasik/backend/services"

# DB switch block template for services with single repo
DB_SWITCH_SINGLE = '''	var repo service.{iface}

	if _, dbHostSet := os.LookupEnv("DB_HOST"); dbHostSet && cfg.DB.Host != "" && cfg.DB.DBName != "" {{
		connCtx, connCancel := context.WithTimeout(context.Background(), 5*time.Second)
		pool, poolErr := pgxpool.New(connCtx, cfg.DB.DSN())
		connCancel()
		if poolErr != nil {{
			log.Printf("[%s] DB connection failed, using memory: %v", serviceName, poolErr)
			repo = memory.New{memCtor}()
		}} else {{
			pingCtx, pingCancel := context.WithTimeout(context.Background(), 3*time.Second)
			if pingErr := pool.Ping(pingCtx); pingErr != nil {{
				pingCancel()
				pool.Close()
				log.Printf("[%s] DB ping failed, using memory: %v", serviceName, pingErr)
				repo = memory.New{memCtor}()
			}} else {{
				pingCancel()
				defer pool.Close()
				log.Printf("[%s] Connected to PostgreSQL", serviceName)
				repo = postgres.New{pgCtor}(pool)
			}}
		}}
	}} else {{
		repo = memory.New{memCtor}()
	}}'''

def update_grpc_service_single(svc_dir, svc_name, iface, memCtor, pgCtor, old_repo_lines, import_path_prefix):
    """Update a gRPC service with shared/config and single repo."""
    path = os.path.join(BASE, svc_dir, "cmd", "main.go")
    with open(path, 'r') as f:
        content = f.read()

    # 1. Add pgxpool + postgres imports
    old_import = f'"{import_path_prefix}/internal/repository/memory"'
    new_import = f'"github.com/jackc/pgx/v5/pgxpool"\n\t"{import_path_prefix}/internal/repository/memory"\n\t"{import_path_prefix}/internal/repository/postgres"'
    content = content.replace(old_import, new_import)

    # 2. Replace repo initialization
    db_switch = DB_SWITCH_SINGLE.format(iface=iface, memCtor=memCtor, pgCtor=pgCtor)
    content = content.replace(old_repo_lines, db_switch)

    with open(path, 'w') as f:
        f.write(content)
    print(f"  [OK] {svc_dir}/cmd/main.go updated")

# ============================================================================
# 1. audit-service (gRPC, has config, has Seed())
# ============================================================================
path = os.path.join(BASE, "audit-service", "cmd", "main.go")
with open(path, 'r') as f:
    content = f.read()

old_import = '"github.com/manpasik/backend/services/audit-service/internal/repository/memory"'
new_import = '"github.com/jackc/pgx/v5/pgxpool"\n\t"github.com/manpasik/backend/services/audit-service/internal/repository/memory"\n\t"github.com/manpasik/backend/services/audit-service/internal/repository/postgres"'
content = content.replace(old_import, new_import)

old_repo = '''\trepo := memory.NewAuditRepository()
\trepo.Seed()
\tsvc := service.NewAuditService(logger, repo)'''
new_repo = '''\tvar repo service.AuditRepository

\tif _, dbHostSet := os.LookupEnv("DB_HOST"); dbHostSet && cfg.DB.Host != "" && cfg.DB.DBName != "" {
\t\tconnCtx, connCancel := context.WithTimeout(context.Background(), 5*time.Second)
\t\tpool, poolErr := pgxpool.New(connCtx, cfg.DB.DSN())
\t\tconnCancel()
\t\tif poolErr != nil {
\t\t\tlog.Printf("[%s] DB connection failed, using memory: %v", serviceName, poolErr)
\t\t\tmemRepo := memory.NewAuditRepository()
\t\t\tmemRepo.Seed()
\t\t\trepo = memRepo
\t\t} else {
\t\t\tpingCtx, pingCancel := context.WithTimeout(context.Background(), 3*time.Second)
\t\t\tif pingErr := pool.Ping(pingCtx); pingErr != nil {
\t\t\t\tpingCancel()
\t\t\t\tpool.Close()
\t\t\t\tlog.Printf("[%s] DB ping failed, using memory: %v", serviceName, pingErr)
\t\t\t\tmemRepo := memory.NewAuditRepository()
\t\t\t\tmemRepo.Seed()
\t\t\t\trepo = memRepo
\t\t\t} else {
\t\t\t\tpingCancel()
\t\t\t\tdefer pool.Close()
\t\t\t\tlog.Printf("[%s] Connected to PostgreSQL", serviceName)
\t\t\t\trepo = postgres.NewAuditRepository(pool)
\t\t\t}
\t\t}
\t} else {
\t\tmemRepo := memory.NewAuditRepository()
\t\tmemRepo.Seed()
\t\trepo = memRepo
\t}

\tsvc := service.NewAuditService(logger, repo)'''
content = content.replace(old_repo, new_repo)

with open(path, 'w') as f:
    f.write(content)
print("  [OK] audit-service/cmd/main.go updated")

# ============================================================================
# 2. digital-twin-service (gRPC, has config)
# ============================================================================
update_grpc_service_single(
    "digital-twin-service", "digital-twin-service",
    "TwinRepository", "TwinRepository", "TwinRepository",
    "\trepo := memory.NewTwinRepository()",
    "github.com/manpasik/backend/services/digital-twin-service"
)

# ============================================================================
# 3. vision-service (gRPC, has config) - uses analysisRepo var name
# ============================================================================
path = os.path.join(BASE, "vision-service", "cmd", "main.go")
with open(path, 'r') as f:
    content = f.read()

old_import = '"github.com/manpasik/backend/services/vision-service/internal/repository/memory"'
new_import = '"github.com/jackc/pgx/v5/pgxpool"\n\t"github.com/manpasik/backend/services/vision-service/internal/repository/memory"\n\t"github.com/manpasik/backend/services/vision-service/internal/repository/postgres"'
content = content.replace(old_import, new_import)

old_repo = '''\t// FoodAnalysisRepository: 인메모리 (향후 PostgreSQL 추가)
\tanalysisRepo := memory.NewFoodAnalysisRepository()'''
new_repo = '''\t// FoodAnalysisRepository: DB 스위치 (PostgreSQL / 인메모리)
\tvar analysisRepo service.FoodAnalysisRepository

\tif _, dbHostSet := os.LookupEnv("DB_HOST"); dbHostSet && cfg.DB.Host != "" && cfg.DB.DBName != "" {
\t\tconnCtx, connCancel := context.WithTimeout(context.Background(), 5*time.Second)
\t\tpool, poolErr := pgxpool.New(connCtx, cfg.DB.DSN())
\t\tconnCancel()
\t\tif poolErr != nil {
\t\t\tlog.Printf("[%s] DB connection failed, using memory: %v", serviceName, poolErr)
\t\t\tanalysisRepo = memory.NewFoodAnalysisRepository()
\t\t} else {
\t\t\tpingCtx, pingCancel := context.WithTimeout(context.Background(), 3*time.Second)
\t\t\tif pingErr := pool.Ping(pingCtx); pingErr != nil {
\t\t\t\tpingCancel()
\t\t\t\tpool.Close()
\t\t\t\tlog.Printf("[%s] DB ping failed, using memory: %v", serviceName, pingErr)
\t\t\t\tanalysisRepo = memory.NewFoodAnalysisRepository()
\t\t\t} else {
\t\t\t\tpingCancel()
\t\t\t\tdefer pool.Close()
\t\t\t\tlog.Printf("[%s] Connected to PostgreSQL", serviceName)
\t\t\t\tanalysisRepo = postgres.NewFoodAnalysisRepository(pool)
\t\t\t}
\t\t}
\t} else {
\t\tanalysisRepo = memory.NewFoodAnalysisRepository()
\t}'''
content = content.replace(old_repo, new_repo)

with open(path, 'w') as f:
    f.write(content)
print("  [OK] vision-service/cmd/main.go updated")

# ============================================================================
# 4. assistant-service (gRPC, has config)
# ============================================================================
update_grpc_service_single(
    "assistant-service", "assistant-service",
    "AssistantRepository", "AssistantRepository", "AssistantRepository",
    "\trepo := memory.NewAssistantRepository()",
    "github.com/manpasik/backend/services/assistant-service"
)

# ============================================================================
# 5. concept-service (gRPC, has config, TWO repos)
# ============================================================================
path = os.path.join(BASE, "concept-service", "cmd", "main.go")
with open(path, 'r') as f:
    content = f.read()

old_import = '"github.com/manpasik/backend/services/concept-service/internal/repository/memory"'
new_import = '"github.com/jackc/pgx/v5/pgxpool"\n\t"github.com/manpasik/backend/services/concept-service/internal/repository/memory"\n\t"github.com/manpasik/backend/services/concept-service/internal/repository/postgres"'
content = content.replace(old_import, new_import)

old_repo = '''\tconceptRepo := memory.NewConceptRepository()
\torgRepo := memory.NewOrganizationRepository()
\tsvc := service.NewConceptService(logger, conceptRepo, orgRepo)'''
new_repo = '''\tvar conceptRepo service.ConceptRepository
\tvar orgRepo service.OrganizationRepository

\tif _, dbHostSet := os.LookupEnv("DB_HOST"); dbHostSet && cfg.DB.Host != "" && cfg.DB.DBName != "" {
\t\tconnCtx, connCancel := context.WithTimeout(context.Background(), 5*time.Second)
\t\tpool, poolErr := pgxpool.New(connCtx, cfg.DB.DSN())
\t\tconnCancel()
\t\tif poolErr != nil {
\t\t\tlog.Printf("[%s] DB connection failed, using memory: %v", serviceName, poolErr)
\t\t\tconceptRepo = memory.NewConceptRepository()
\t\t\torgRepo = memory.NewOrganizationRepository()
\t\t} else {
\t\t\tpingCtx, pingCancel := context.WithTimeout(context.Background(), 3*time.Second)
\t\t\tif pingErr := pool.Ping(pingCtx); pingErr != nil {
\t\t\t\tpingCancel()
\t\t\t\tpool.Close()
\t\t\t\tlog.Printf("[%s] DB ping failed, using memory: %v", serviceName, pingErr)
\t\t\t\tconceptRepo = memory.NewConceptRepository()
\t\t\t\torgRepo = memory.NewOrganizationRepository()
\t\t\t} else {
\t\t\t\tpingCancel()
\t\t\t\tdefer pool.Close()
\t\t\t\tlog.Printf("[%s] Connected to PostgreSQL", serviceName)
\t\t\t\tconceptRepo = postgres.NewConceptRepository(pool)
\t\t\t\torgRepo = postgres.NewOrganizationRepository(pool)
\t\t\t}
\t\t}
\t} else {
\t\tconceptRepo = memory.NewConceptRepository()
\t\torgRepo = memory.NewOrganizationRepository()
\t}

\tsvc := service.NewConceptService(logger, conceptRepo, orgRepo)'''
content = content.replace(old_repo, new_repo)

with open(path, 'w') as f:
    f.write(content)
print("  [OK] concept-service/cmd/main.go updated")

# ============================================================================
# 6. cartridge-store-service (gRPC, has config)
# ============================================================================
update_grpc_service_single(
    "cartridge-store-service", "cartridge-store-service",
    "StoreRepository", "StoreRepository", "StoreRepository",
    "\trepo := memory.NewStoreRepository()",
    "github.com/manpasik/backend/services/cartridge-store-service"
)

# ============================================================================
# 7. data-platform-service (gRPC, has config)
# ============================================================================
update_grpc_service_single(
    "data-platform-service", "data-platform-service",
    "PlatformRepository", "PlatformRepository", "PlatformRepository",
    "\trepo := memory.NewPlatformRepository()",
    "github.com/manpasik/backend/services/data-platform-service"
)

# ============================================================================
# 8. voice-profile-service (gRPC, has config)
# ============================================================================
update_grpc_service_single(
    "voice-profile-service", "voice-profile-service",
    "VoiceProfileRepository", "VoiceProfileRepository", "VoiceProfileRepository",
    "\trepo := memory.NewVoiceProfileRepository()",
    "github.com/manpasik/backend/services/voice-profile-service"
)

# ============================================================================
# HTTP-based services (no shared/config) - need more extensive changes
# ============================================================================

# 9. analytics-service
path = os.path.join(BASE, "analytics-service", "cmd", "main.go")
with open(path, 'r') as f:
    content = f.read()

content = content.replace(
    '''import (
\t"log"
\t"net/http"
\t"os"
\t"os/signal"
\t"syscall"

\t"github.com/manpasik/backend/services/analytics-service/internal/repository/memory"
\t"github.com/manpasik/backend/services/analytics-service/internal/service"
)''',
    '''import (
\t"context"
\t"log"
\t"net/http"
\t"os"
\t"os/signal"
\t"syscall"
\t"time"

\t"github.com/jackc/pgx/v5/pgxpool"
\t"github.com/manpasik/backend/services/analytics-service/internal/repository/memory"
\t"github.com/manpasik/backend/services/analytics-service/internal/repository/postgres"
\t"github.com/manpasik/backend/services/analytics-service/internal/service"
\t"github.com/manpasik/backend/shared/config"
)'''
)

content = content.replace(
    '''\tlog.Printf("[%s] Starting...", serviceName)

\trepo := memory.NewAnalyticsRepository()
\tsvc := service.NewAnalyticsService(repo)
\t_ = svc''',
    '''\tcfg := config.LoadFromEnv(serviceName)
\tlog.Printf("[%s] Starting v%s...", serviceName, cfg.Version)

\tvar repo service.AnalyticsRepository

\tif _, dbHostSet := os.LookupEnv("DB_HOST"); dbHostSet && cfg.DB.Host != "" && cfg.DB.DBName != "" {
\t\tconnCtx, connCancel := context.WithTimeout(context.Background(), 5*time.Second)
\t\tpool, poolErr := pgxpool.New(connCtx, cfg.DB.DSN())
\t\tconnCancel()
\t\tif poolErr != nil {
\t\t\tlog.Printf("[%s] DB connection failed, using memory: %v", serviceName, poolErr)
\t\t\trepo = memory.NewAnalyticsRepository()
\t\t} else {
\t\t\tpingCtx, pingCancel := context.WithTimeout(context.Background(), 3*time.Second)
\t\t\tif pingErr := pool.Ping(pingCtx); pingErr != nil {
\t\t\t\tpingCancel()
\t\t\t\tpool.Close()
\t\t\t\tlog.Printf("[%s] DB ping failed, using memory: %v", serviceName, pingErr)
\t\t\t\trepo = memory.NewAnalyticsRepository()
\t\t\t} else {
\t\t\t\tpingCancel()
\t\t\t\tdefer pool.Close()
\t\t\t\tlog.Printf("[%s] Connected to PostgreSQL", serviceName)
\t\t\t\trepo = postgres.NewAnalyticsRepository(pool)
\t\t\t}
\t\t}
\t} else {
\t\trepo = memory.NewAnalyticsRepository()
\t}

\tsvc := service.NewAnalyticsService(repo)
\t_ = svc'''
)

with open(path, 'w') as f:
    f.write(content)
print("  [OK] analytics-service/cmd/main.go updated")

# 10. emergency-service
path = os.path.join(BASE, "emergency-service", "cmd", "main.go")
with open(path, 'r') as f:
    content = f.read()

content = content.replace(
    '''import (
\t"log"
\t"net/http"
\t"os"
\t"os/signal"
\t"syscall"

\t"github.com/manpasik/backend/services/emergency-service/internal/handler"
\t"github.com/manpasik/backend/services/emergency-service/internal/repository/memory"
\t"github.com/manpasik/backend/services/emergency-service/internal/service"
)''',
    '''import (
\t"context"
\t"log"
\t"net/http"
\t"os"
\t"os/signal"
\t"syscall"
\t"time"

\t"github.com/jackc/pgx/v5/pgxpool"
\t"github.com/manpasik/backend/services/emergency-service/internal/handler"
\t"github.com/manpasik/backend/services/emergency-service/internal/repository/memory"
\t"github.com/manpasik/backend/services/emergency-service/internal/repository/postgres"
\t"github.com/manpasik/backend/services/emergency-service/internal/service"
\t"github.com/manpasik/backend/shared/config"
)'''
)

content = content.replace(
    '''\tlog.Printf("[%s] Starting...", serviceName)

\t// 인메모리 저장소 초기화
\trepo := memory.NewEmergencyRepository()

\t// 서비스 레이어 초기화
\tsvc := service.NewEmergencyService(repo)''',
    '''\tcfg := config.LoadFromEnv(serviceName)
\tlog.Printf("[%s] Starting v%s...", serviceName, cfg.Version)

\tvar repo service.EmergencyRepository

\tif _, dbHostSet := os.LookupEnv("DB_HOST"); dbHostSet && cfg.DB.Host != "" && cfg.DB.DBName != "" {
\t\tconnCtx, connCancel := context.WithTimeout(context.Background(), 5*time.Second)
\t\tpool, poolErr := pgxpool.New(connCtx, cfg.DB.DSN())
\t\tconnCancel()
\t\tif poolErr != nil {
\t\t\tlog.Printf("[%s] DB connection failed, using memory: %v", serviceName, poolErr)
\t\t\trepo = memory.NewEmergencyRepository()
\t\t} else {
\t\t\tpingCtx, pingCancel := context.WithTimeout(context.Background(), 3*time.Second)
\t\t\tif pingErr := pool.Ping(pingCtx); pingErr != nil {
\t\t\t\tpingCancel()
\t\t\t\tpool.Close()
\t\t\t\tlog.Printf("[%s] DB ping failed, using memory: %v", serviceName, pingErr)
\t\t\t\trepo = memory.NewEmergencyRepository()
\t\t\t} else {
\t\t\t\tpingCancel()
\t\t\t\tdefer pool.Close()
\t\t\t\tlog.Printf("[%s] Connected to PostgreSQL", serviceName)
\t\t\t\trepo = postgres.NewEmergencyRepository(pool)
\t\t\t}
\t\t}
\t} else {
\t\trepo = memory.NewEmergencyRepository()
\t}

\tsvc := service.NewEmergencyService(repo)'''
)

with open(path, 'w') as f:
    f.write(content)
print("  [OK] emergency-service/cmd/main.go updated")

# 11. iot-gateway-service - need to read it first
path = os.path.join(BASE, "iot-gateway-service", "cmd", "main.go")
with open(path, 'r') as f:
    content = f.read()

content = content.replace(
    '"github.com/manpasik/backend/services/iot-gateway-service/internal/repository/memory"',
    '"github.com/jackc/pgx/v5/pgxpool"\n\t"github.com/manpasik/backend/services/iot-gateway-service/internal/repository/memory"\n\t"github.com/manpasik/backend/services/iot-gateway-service/internal/repository/postgres"\n\t"github.com/manpasik/backend/shared/config"'
)

# Add context/time imports if missing
if '"context"' not in content:
    content = content.replace(
        'import (\n\t"log"',
        'import (\n\t"context"\n\t"log"'
    )
if '"time"' not in content:
    content = content.replace(
        '"syscall"',
        '"syscall"\n\t"time"'
    )

# Replace repo initialization
content = content.replace(
    '''\tlog.Printf("[%s] Starting...", serviceName)

\trepo := memory.NewIoTRepository()
\tsvc := service.NewIoTGatewayService(repo)
\t_ = svc''',
    '''\tcfg := config.LoadFromEnv(serviceName)
\tlog.Printf("[%s] Starting v%s...", serviceName, cfg.Version)

\tvar repo service.IoTRepository

\tif _, dbHostSet := os.LookupEnv("DB_HOST"); dbHostSet && cfg.DB.Host != "" && cfg.DB.DBName != "" {
\t\tconnCtx, connCancel := context.WithTimeout(context.Background(), 5*time.Second)
\t\tpool, poolErr := pgxpool.New(connCtx, cfg.DB.DSN())
\t\tconnCancel()
\t\tif poolErr != nil {
\t\t\tlog.Printf("[%s] DB connection failed, using memory: %v", serviceName, poolErr)
\t\t\trepo = memory.NewIoTRepository()
\t\t} else {
\t\t\tpingCtx, pingCancel := context.WithTimeout(context.Background(), 3*time.Second)
\t\t\tif pingErr := pool.Ping(pingCtx); pingErr != nil {
\t\t\t\tpingCancel()
\t\t\t\tpool.Close()
\t\t\t\tlog.Printf("[%s] DB ping failed, using memory: %v", serviceName, pingErr)
\t\t\t\trepo = memory.NewIoTRepository()
\t\t\t} else {
\t\t\t\tpingCancel()
\t\t\t\tdefer pool.Close()
\t\t\t\tlog.Printf("[%s] Connected to PostgreSQL", serviceName)
\t\t\t\trepo = postgres.NewIoTRepository(pool)
\t\t\t}
\t\t}
\t} else {
\t\trepo = memory.NewIoTRepository()
\t}

\tsvc := service.NewIoTGatewayService(repo)
\t_ = svc'''
)

with open(path, 'w') as f:
    f.write(content)
print("  [OK] iot-gateway-service/cmd/main.go updated")

# 12. marketplace-service (HTTP, 3 repos)
path = os.path.join(BASE, "marketplace-service", "cmd", "main.go")
with open(path, 'r') as f:
    content = f.read()

content = content.replace(
    '''import (
\t"encoding/json"
\t"log"
\t"net/http"
\t"os"
\t"os/signal"
\t"syscall"

\t"github.com/manpasik/backend/services/marketplace-service/internal/repository/memory"
\t"github.com/manpasik/backend/services/marketplace-service/internal/service"
)''',
    '''import (
\t"context"
\t"encoding/json"
\t"log"
\t"net/http"
\t"os"
\t"os/signal"
\t"syscall"
\t"time"

\t"github.com/jackc/pgx/v5/pgxpool"
\t"github.com/manpasik/backend/services/marketplace-service/internal/repository/memory"
\t"github.com/manpasik/backend/services/marketplace-service/internal/repository/postgres"
\t"github.com/manpasik/backend/services/marketplace-service/internal/service"
\t"github.com/manpasik/backend/shared/config"
)'''
)

content = content.replace(
    '''\tlog.Printf("[%s] Starting...", serviceName)

\t// 인메모리 저장소 초기화
\tproductRepo := memory.NewProductRepository()
\tpartnerRepo := memory.NewPartnerRepository()
\tstatsRepo := memory.NewStatsRepository()

\t// 서비스 레이어 초기화
\tsvc := service.NewMarketplaceService(productRepo, partnerRepo, statsRepo)''',
    '''\tcfg := config.LoadFromEnv(serviceName)
\tlog.Printf("[%s] Starting v%s...", serviceName, cfg.Version)

\tvar productRepo service.ProductRepository
\tvar partnerRepo service.PartnerRepository
\tvar statsRepo service.StatsRepository

\tif _, dbHostSet := os.LookupEnv("DB_HOST"); dbHostSet && cfg.DB.Host != "" && cfg.DB.DBName != "" {
\t\tconnCtx, connCancel := context.WithTimeout(context.Background(), 5*time.Second)
\t\tpool, poolErr := pgxpool.New(connCtx, cfg.DB.DSN())
\t\tconnCancel()
\t\tif poolErr != nil {
\t\t\tlog.Printf("[%s] DB connection failed, using memory: %v", serviceName, poolErr)
\t\t\tproductRepo = memory.NewProductRepository()
\t\t\tpartnerRepo = memory.NewPartnerRepository()
\t\t\tstatsRepo = memory.NewStatsRepository()
\t\t} else {
\t\t\tpingCtx, pingCancel := context.WithTimeout(context.Background(), 3*time.Second)
\t\t\tif pingErr := pool.Ping(pingCtx); pingErr != nil {
\t\t\t\tpingCancel()
\t\t\t\tpool.Close()
\t\t\t\tlog.Printf("[%s] DB ping failed, using memory: %v", serviceName, pingErr)
\t\t\t\tproductRepo = memory.NewProductRepository()
\t\t\t\tpartnerRepo = memory.NewPartnerRepository()
\t\t\t\tstatsRepo = memory.NewStatsRepository()
\t\t\t} else {
\t\t\t\tpingCancel()
\t\t\t\tdefer pool.Close()
\t\t\t\tlog.Printf("[%s] Connected to PostgreSQL", serviceName)
\t\t\t\tproductRepo = postgres.NewProductRepository(pool)
\t\t\t\tpartnerRepo = postgres.NewPartnerRepository(pool)
\t\t\t\tstatsRepo = postgres.NewStatsRepository(pool)
\t\t\t}
\t\t}
\t} else {
\t\tproductRepo = memory.NewProductRepository()
\t\tpartnerRepo = memory.NewPartnerRepository()
\t\tstatsRepo = memory.NewStatsRepository()
\t}

\tsvc := service.NewMarketplaceService(productRepo, partnerRepo, statsRepo)'''
)

with open(path, 'w') as f:
    f.write(content)
print("  [OK] marketplace-service/cmd/main.go updated")

# 13. nlp-service
path = os.path.join(BASE, "nlp-service", "cmd", "main.go")
with open(path, 'r') as f:
    content = f.read()

content = content.replace(
    '''import (
\t"log"
\t"net/http"
\t"os"
\t"os/signal"
\t"syscall"

\t"github.com/manpasik/backend/services/nlp-service/internal/repository/memory"
\t"github.com/manpasik/backend/services/nlp-service/internal/service"
)''',
    '''import (
\t"context"
\t"log"
\t"net/http"
\t"os"
\t"os/signal"
\t"syscall"
\t"time"

\t"github.com/jackc/pgx/v5/pgxpool"
\t"github.com/manpasik/backend/services/nlp-service/internal/repository/memory"
\t"github.com/manpasik/backend/services/nlp-service/internal/repository/postgres"
\t"github.com/manpasik/backend/services/nlp-service/internal/service"
\t"github.com/manpasik/backend/shared/config"
)'''
)

content = content.replace(
    '''\tlog.Printf("[%s] Starting...", serviceName)

\trepo := memory.NewNLPRepository()
\tsvc := service.NewNLPService(repo)
\t_ = svc''',
    '''\tcfg := config.LoadFromEnv(serviceName)
\tlog.Printf("[%s] Starting v%s...", serviceName, cfg.Version)

\tvar repo service.NLPRepository

\tif _, dbHostSet := os.LookupEnv("DB_HOST"); dbHostSet && cfg.DB.Host != "" && cfg.DB.DBName != "" {
\t\tconnCtx, connCancel := context.WithTimeout(context.Background(), 5*time.Second)
\t\tpool, poolErr := pgxpool.New(connCtx, cfg.DB.DSN())
\t\tconnCancel()
\t\tif poolErr != nil {
\t\t\tlog.Printf("[%s] DB connection failed, using memory: %v", serviceName, poolErr)
\t\t\trepo = memory.NewNLPRepository()
\t\t} else {
\t\t\tpingCtx, pingCancel := context.WithTimeout(context.Background(), 3*time.Second)
\t\t\tif pingErr := pool.Ping(pingCtx); pingErr != nil {
\t\t\t\tpingCancel()
\t\t\t\tpool.Close()
\t\t\t\tlog.Printf("[%s] DB ping failed, using memory: %v", serviceName, pingErr)
\t\t\t\trepo = memory.NewNLPRepository()
\t\t\t} else {
\t\t\t\tpingCancel()
\t\t\t\tdefer pool.Close()
\t\t\t\tlog.Printf("[%s] Connected to PostgreSQL", serviceName)
\t\t\t\trepo = postgres.NewNLPRepository(pool)
\t\t\t}
\t\t}
\t} else {
\t\trepo = memory.NewNLPRepository()
\t}

\tsvc := service.NewNLPService(repo)
\t_ = svc'''
)

with open(path, 'w') as f:
    f.write(content)
print("  [OK] nlp-service/cmd/main.go updated")

print("\n=== 13개 main.go DB 스위치 로직 추가 완료 ===")
