package security

import (
	"context"
	"github.com/ngdangkietswe/swe-go-common-shared/domain"
	"github.com/ngdangkietswe/swe-go-common-shared/grpc/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"
)

// SecuredAuth is a function that checks if the user has the required permission.
func SecuredAuth[Req any, Resp any](
	ctx context.Context,
	req Req,
	permission domain.Permission,
	serviceFunc func(context.Context, Req) (Resp, error)) (Resp, error) {
	principal := util.GetGrpcPrincipal(ctx)

	if hasPermission(permission, principal.UserPermission) {
		return serviceFunc(ctx, req)
	}

	return *new(Resp), accessDenied()
}

// hasPermission checks if the user has the required permission.
func hasPermission(permissionRequired domain.Permission, userPermission *domain.UserPermission) bool {
	if userPermission == nil {
		return false
	}

	permissions := userPermission.Permissions
	if permissions == nil || len(permissions) == 0 {
		return false
	}

	for _, permission := range permissions {
		if strings.EqualFold(permission.Action, permissionRequired.Action) && strings.EqualFold(permission.Resource, permissionRequired.Resource) {
			return true
		}
	}

	return false
}

// accessDenied returns an error with code PermissionDenied.
func accessDenied() error {
	return status.Errorf(codes.PermissionDenied, "Access denied")
}
