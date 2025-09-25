// internal/middleware/user_session.go

package middleware

import (
	"net"
	"strings"
	"tormenta20-builder/internal/database"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	UserSessionCookie = "user_session_id"
	UserSessionHeader = "X-User-Session-ID"
	UserIPHeader      = "X-User-IP"
)

func UserSessionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var sessionID string
		userIP := getUserIP(c)

		// 1. Tenta obter sessão do cabeçalho (MAIOR PRIORIDADE)
		sessionID = c.GetHeader(UserSessionHeader)

		// 2. Se não veio no cabeçalho, tenta obter do cookie
		if sessionID == "" {
			cookieSessionID, err := c.Cookie(UserSessionCookie)
			if err == nil && cookieSessionID != "" {
				sessionID = cookieSessionID
			}
		}

		// 3. Se não existe em nenhum lugar, cria uma nova sessão
		if sessionID == "" {
			sessionID = uuid.New().String()

			if userIP != "unknown" {
				// A lógica de re-associação por IP pode ser mantida como um fallback
				go func(sid, ip string) {
					database.DB.Exec(
						"UPDATE personagens SET user_session_id = ? WHERE user_ip = ? AND (user_session_id IS NULL OR user_session_id != ?)",
						sid, ip, sid,
					)
				}(sessionID, userIP)
			}
		}

		// 4. Define o cookie (ou atualiza o tempo de expiração do existente)
		// Isso garante que o cookie seja criado/atualizado mesmo que a sessão venha pelo header
		isSecure := c.GetHeader("X-Forwarded-Proto") == "https" ||
			c.Request.TLS != nil ||
			strings.HasPrefix(c.Request.Host, "backend-tormenta20.fly.dev")
		domain := "" // Domínio vazio para permitir cross-origin

		c.SetCookie(
			UserSessionCookie,
			sessionID,
			60*60*24*30, // 30 dias
			"/",
			domain,
			isSecure,
			false, // httpOnly: false para permitir acesso via JS se necessário
		)

		// 5. Adiciona informações ao contexto
		c.Set("user_session_id", sessionID)
		c.Set("user_ip", userIP)

		// 6. Adiciona headers na RESPOSTA para o frontend sempre ter o ID mais atual
		c.Header(UserSessionHeader, sessionID)
		c.Header(UserIPHeader, userIP)

		c.Next()
	}
}

// ... (o resto do arquivo user_session.go permanece o mesmo)
// getUserIP, GetUserSessionID, GetUserIP, GetUserIdentification
func getUserIP(c *gin.Context) string {
	// Verifica headers de proxy mais comuns
	headers := []string{
		"X-Forwarded-For",
		"X-Real-IP",
		"X-Client-IP",
		"CF-Connecting-IP", // Cloudflare
		"True-Client-IP",   // Akamai
	}

	for _, header := range headers {
		ip := c.GetHeader(header)
		if ip != "" {
			// X-Forwarded-For pode ter múltiplos IPs separados por vírgula
			if strings.Contains(ip, ",") {
				ip = strings.TrimSpace(strings.Split(ip, ",")[0])
			}

			// Valida se é um IP válido
			if net.ParseIP(ip) != nil {
				return ip
			}
		}
	}

	// Se não encontrou nos headers, usa o IP da conexão direta
	clientIP := c.ClientIP()
	if clientIP != "" {
		return clientIP
	}

	// Fallback para IP da conexão remota
	if ip, _, err := net.SplitHostPort(c.Request.RemoteAddr); err == nil {
		return ip
	}

	return "unknown"
}

// GetUserSessionID obtém o session ID do contexto
func GetUserSessionID(c *gin.Context) string {
	if sessionID, exists := c.Get("user_session_id"); exists {
		if id, ok := sessionID.(string); ok {
			return id
		}
	}
	return ""
}

// GetUserIP obtém o IP do usuário do contexto
func GetUserIP(c *gin.Context) string {
	if userIP, exists := c.Get("user_ip"); exists {
		if ip, ok := userIP.(string); ok {
			return ip
		}
	}
	return ""
}

// GetUserIdentification retorna both session ID and IP for flexibility
func GetUserIdentification(c *gin.Context) (sessionID, userIP string) {
	return GetUserSessionID(c), GetUserIP(c)
}
