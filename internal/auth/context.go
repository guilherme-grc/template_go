package auth

import "context"

type contextKey string

// UserContextKey - Key used to store/retrieve the user from the context
const UserContextKey contextKey = "authenticated_user"

// InjectUserIntoContext — Saves the authenticated user into the request context
// Equivalent to Laravel's auth()->user()
func InjectUserIntoContext(ctx context.Context, claims *Claims) context.Context {
	return context.WithValue(ctx, UserContextKey, claims)
}

// UserFromContext — Retrieves the authenticated user from the context
// Equivalent to Laravel's auth()->user() or $request->user()
func UserFromContext(ctx context.Context) (*Claims, bool) {
	claims, ok := ctx.Value(UserContextKey).(*Claims)
	return claims, ok
}
