package middleware

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Role constants
const (
	RoleAdmin        = "admin"
	RoleMedicalStaff = "medical_staff"
	RoleUser         = "user"
	RoleFamilyMember = "family_member"
	RoleResearcher   = "researcher"
)

// RBACConfig defines which roles can access which methods
type RBACConfig struct {
	// MethodRoles maps gRPC full method names to allowed roles
	// If a method is not in the map, any authenticated user can access it
	MethodRoles map[string][]string
}

// RBACInterceptor creates a role-based access control interceptor
func RBACInterceptor(config *RBACConfig) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// If no config or method not restricted, allow access
		if config == nil {
			return handler(ctx, req)
		}

		allowedRoles, restricted := config.MethodRoles[info.FullMethod]
		if !restricted {
			return handler(ctx, req)
		}

		// Get user role from context (set by AuthInterceptor)
		role, ok := UserRoleFromContext(ctx)
		if !ok {
			return nil, status.Error(codes.PermissionDenied, "역할 정보 없음")
		}

		// Admin always has access
		if role == RoleAdmin {
			return handler(ctx, req)
		}

		// Check if role is in allowed list
		for _, allowed := range allowedRoles {
			if role == allowed {
				return handler(ctx, req)
			}
		}

		return nil, status.Errorf(codes.PermissionDenied, "역할 '%s'은(는) 이 작업에 대한 권한이 없습니다", role)
	}
}

// RBACStreamInterceptor creates a role-based access control stream interceptor
func RBACStreamInterceptor(config *RBACConfig) grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		if config == nil {
			return handler(srv, ss)
		}

		allowedRoles, restricted := config.MethodRoles[info.FullMethod]
		if !restricted {
			return handler(srv, ss)
		}

		role, ok := UserRoleFromContext(ss.Context())
		if !ok {
			return status.Error(codes.PermissionDenied, "역할 정보 없음")
		}

		if role == RoleAdmin {
			return handler(srv, ss)
		}

		for _, allowed := range allowedRoles {
			if role == allowed {
				return handler(srv, ss)
			}
		}

		return status.Errorf(codes.PermissionDenied, "역할 '%s'은(는) 이 작업에 대한 권한이 없습니다", role)
	}
}

// DefaultRBACConfig returns the default RBAC configuration for ManPaSik services.
// Admin-only methods are explicitly listed; all others are available to any authenticated user.
func DefaultRBACConfig() *RBACConfig {
	adminOnly := []string{RoleAdmin}
	medicalOrAdmin := []string{RoleAdmin, RoleMedicalStaff}

	return &RBACConfig{
		MethodRoles: map[string][]string{
			// Admin Service — admin only
			"/manpasik.v1.AdminService/ListUsers":           adminOnly,
			"/manpasik.v1.AdminService/ChangeUserRole":      adminOnly,
			"/manpasik.v1.AdminService/BulkUserAction":      adminOnly,
			"/manpasik.v1.AdminService/GetSystemStats":      adminOnly,
			"/manpasik.v1.AdminService/GetAuditLog":         adminOnly,
			"/manpasik.v1.AdminService/UpdateSystemConfig":  adminOnly,
			"/manpasik.v1.AdminService/GetSystemConfig":     adminOnly,

			// Telemedicine — medical staff or admin
			"/manpasik.v1.TelemedicineService/CreatePrescription": medicalOrAdmin,
			"/manpasik.v1.TelemedicineService/UpdatePrescription": medicalOrAdmin,
			"/manpasik.v1.TelemedicineService/EndConsultation":    medicalOrAdmin,

			// Translation — admin only for management
			"/manpasik.v1.TranslationService/UpdateTranslation":  adminOnly,
			"/manpasik.v1.TranslationService/BulkImport":         adminOnly,
		},
	}
}
