package auth

import "context"

type contextKey string

const usuarioContextKey contextKey = "usuario_autenticado"

// InjetarUsuarioNoContexto — salva o usuário autenticado no contexto da requisição
// Equivalente ao auth()->user() do Laravel
func InjetarUsuarioNoContexto(ctx context.Context, claims *Claims) context.Context {
	return context.WithValue(ctx, usuarioContextKey, claims)
}

// UsuarioDoContexto — recupera o usuário autenticado do contexto
// Equivalente ao auth()->user() ou $request->user() do Laravel
func UsuarioDoContexto(ctx context.Context) (*Claims, bool) {
	claims, ok := ctx.Value(usuarioContextKey).(*Claims)
	return claims, ok
}
