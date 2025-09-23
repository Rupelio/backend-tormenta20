package middleware

import (
	"net"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	UserSessionCookie = "user_session_id"
	UserSessionHeader = "X-User-Session-ID"
	UserIPHeader      = "X-User-IP"
)

// UserSessionMiddleware gerencia identificação de usuários via cookie/sessão
func UserSessionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Tenta obter sessão existente do cookie
		sessionID, err := c.Cookie(UserSessionCookie)
		if err != nil || sessionID == "" {
			// 2. Se não existe, cria nova sessão
			sessionID = uuid.New().String()

			// Define cookie com sessão (httpOnly para segurança)
			c.SetCookie(
				UserSessionCookie,
				sessionID,
				60*60*24*30, // 30 dias
				"/",
				"",
				false, // secure (deve ser true em HTTPS)
				true,  // httpOnly
			)
		}

		// 3. Obtém IP do usuário (considerando proxies)
		userIP := getUserIP(c)

		// 4. Adiciona informações ao contexto
		c.Set("user_session_id", sessionID)
		c.Set("user_ip", userIP)

		// 5. Adiciona headers para facilitar debug
		c.Header(UserSessionHeader, sessionID)
		c.Header(UserIPHeader, userIP)

		c.Next()
	}
}

// getUserIP obtém o IP real do usuário considerando proxies e load balancers
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
