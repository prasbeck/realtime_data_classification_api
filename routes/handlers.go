package routes

import (
    "database/sql"
    "net/http"
    "rtmc/websocket"
)

func RegisterRoutes(db *sql.DB, wsManager *websocket.Manager) {
    http.HandleFunc("/signup", SignupHandler(db))
    http.HandleFunc("/login", LoginHandler(db))

    // Example of a protected route
    http.Handle("/protected", JWTMiddleware(http.HandlerFunc(ProtectedHandler)))

    // Rule endpoints
    http.HandleFunc("/rules/create", JWTMiddleware(http.HandlerFunc(CreateRuleHandler)))
    http.HandleFunc("/rules", JWTMiddleware(http.HandlerFunc(GetRulesHandler)))
    http.HandleFunc("/rules/get", JWTMiddleware(http.HandlerFunc(GetRuleHandler)))
    http.HandleFunc("/rules/update", JWTMiddleware(http.HandlerFunc(UpdateRuleHandler)))
    http.HandleFunc("/rules/delete", JWTMiddleware(http.HandlerFunc(DeleteRuleHandler)))

    // WebSocket endpoint
    http.HandleFunc("/ws", wsManager.HandleConnections)
}

// ProtectedHandler is an example of a protected route
func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("This is a protected route"))
}


