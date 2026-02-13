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
