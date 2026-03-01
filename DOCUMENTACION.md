# Documentación del Proyecto Chatvist-Chat

## Descripción General

Chatvist-Chat es un sistema de chat en tiempo real desarrollado en Go que integra modelos de Inteligencia Artificial (IA) como participantes activos en las conversaciones. El sistema permite a usuarios humanos y agentes de IA interactuar simultáneamente en grupos de chat.

---

## Estructura del Proyecto

```
chatvist-chat/
├── main.go                    # Punto de entrada de la aplicación
├── go.mod                     # Dependencias del proyecto
├── go.sum                     # Checksums de dependencias
├── .env                       # Variables de entorno
├── tmp/                       # Archivos temporales
├── config/                    # Configuración externa
│   └── db/                    # Gestión de base de datos
│       ├── connection.go      # Conexión a PostgreSQL
│       └── inicialized.go     # Inicialización de tablas
└── internal/
    ├── auth/              # Autenticación y autorización (Clean Architecture)
    │   ├── delivery/      # Controladores HTTP
    │   └── usecase/       # Casos de uso de autenticación
    ├── domain/            # Modelos de dominio e interfaces
    ├── grupo/             # Gestión de grupos de chat
    │   ├── delivery/      # Controladores HTTP
    │   ├── repository/    # Acceso a datos
    │   └── usecase/       # Casos de uso
    ├── grupousario/       # Relación usuarios-grupos
    ├── usuario/           # Gestión de usuarios
    ├── mensaje/           # Gestión de mensajes
    ├── ia/                # Servicios de Inteligencia Artificial
    │   ├── iaConfig.go    # Configuración de modelos IA
    │   └── service.go     # Servicio de procesamiento IA
    ├── llm/               # Cliente de modelos de lenguaje
    │   └── chatllm.go     # Comunicación con LLMs
    ├── websocket/         # Comunicación en tiempo real
    │   ├── handler.go     # Manejadores de WebSocket
    │   └── hub.go         # Hub central de mensajería
    └── pkg/               # Utilidades compartidas
        ├── brcrypt.go         # Encriptación de contraseñas
        ├── generateJwt.go     # Generación de tokens JWT
        ├── generateid.go      # Generación de identificadores
        ├── middleware/        # Middleware de autenticación
        ├── params.go          # Validación de parámetros
        ├── responde.go        # Respuestas HTTP estandarizadas
        └── validate.go        # Validación de datos
```

---

## Componentes Principales

### 1. **WebSocket Hub (`internal/websocket/`)**

- **Función**: Centro de distribución de mensajes en tiempo real
- **Características**:
  - Gestiona conexiones de usuarios mediante WebSocket
  - Distribuye mensajes a usuarios según pertenencia a grupos
  - Canal especial `aiChannel` para enrutar mensajes a IAs
  - Mantiene mapeo de usuarios a grupos (`userGroups`)
  - Thread-safe mediante mutex

### 2. **Sistema de IA (`internal/ia/`)**

#### Configuración de IAs (`iaConfig.go`)

Define múltiples instancias de IA con diferentes modelos:

```go
type IAConfig struct {
    UserID     string    // ID del usuario IA en el sistema
    LLMBaseURL string    // URL del servidor de modelo de lenguaje
    LLMName    string    // Nombre del modelo (ej: "gpt-oss", "deepseek-r1")
    LLMAPIKey  string    // API Key para autenticación
}
```

**Configuraciones predefinidas:**

- **IA Usuario 1**: Puerto 11434, modelo "gpt-oss"
- **IA Usuario 2**: Puerto 11435, modelo "deepseek-r1"

#### Servicio de IA (`service.go`)

- **Pool de Workers**: Cada IA tiene 5 workers concurrentes para procesar mensajes
- **Flujo de procesamiento**:
  1. Escucha mensajes del canal `inputChannel`
  2. Verifica pertenencia al grupo
  3. Obtiene historial completo del grupo desde BD
  4. Construye prompt con contexto (system + historial)
  5. Envía al LLM correspondiente
  6. Parsea respuesta JSON del modelo
  7. Guarda mensaje de IA en base de datos
  8. Difunde respuesta a través del Hub

### 3. **Cliente LLM (`internal/llm/chatllm.go`)**

- **Función**: Comunicación HTTP con servidores de modelos de lenguaje
- **Protocolo**: Envía conversaciones completas en formato JSON
- **Respuesta esperada**: `{"content": "...", "answer_id": "..."}`

### 4. **Modelos de Datos (`internal/domain/`)**

#### Usuarios

```go
type Usuarios struct {
    Id       uint64
    Nombre   string
    Apodo    string
    Email    string
    Password string    // Hash bcrypt
    Fecha    time.Time
    Token    string
}
```

#### Grupos

```go
type Grupos struct {
    Id          uint64
    Clave       string    // Identificador único del grupo
    Nombre      string
    Fecha       time.Time
    CreatedById uint64    // Usuario creador
}
```

#### Mensajes

```go
type Mensajes struct {
    Id         uint64
    Contenido  string
    Fecha      time.Time
    GrupoId    uint64
    UsuarioId  uint64
    ResponseId *uint64   // ID del mensaje al que responde (nullable)
}
```

#### Relación Grupos-Usuarios

```go
type GruposUsuarios struct {
    IdGrupo   uint64
    IdUsuario uint64
}
```

---

## Librerías y Dependencias

### Dependencias Principales (`go.mod`)

#### Framework Web

- **`github.com/gofiber/fiber/v2` (v2.52.9)**
  - Framework HTTP de alto rendimiento
  - Manejo de rutas, middleware y respuestas

#### WebSocket

- **`github.com/gofiber/websocket/v2` (v2.2.1)**
  - Implementación de WebSocket para Fiber
  - Comunicación bidireccional en tiempo real

#### Base de Datos

- **`gorm.io/gorm` (v1.30.1)**
  - ORM para Go
  - Migraciones y gestión de modelos

- **`gorm.io/driver/postgres` (v1.6.0)**
  - Driver de PostgreSQL para GORM

#### Seguridad

- **`golang.org/x/crypto` (v0.31.0)**
  - Encriptación bcrypt para contraseñas

- **`github.com/golang-jwt/jwt/v5` (v5.3.0)**
  - Generación y validación de tokens JWT

#### Utilidades

- **`github.com/google/uuid` (v1.6.0)**
  - Generación de identificadores únicos

- **`github.com/joho/godotenv` (v1.5.1)**
  - Carga de variables de entorno desde `.env`

### Instalación de Dependencias

```bash
go mod download
```

---

## Puertos y Configuración

### Puertos del Sistema

#### Backend (Chatvist-Chat)

- **Puerto 3100**: Servidor HTTP/WebSocket principal
  - Rutas API REST: `http://localhost:3100/api/`
  - WebSocket: `ws://localhost:3100/api/public/ws/chat`

#### Base de Datos

- **Puerto 5432**: PostgreSQL
  - Host: `ip_bd`
  - Base de datos: `name_bd`
  - Usuario: `name_bd`
  - Contraseña: configurada en `.env`

#### Modelos de IA (LLMs)

##### IA Usuario 1 (gpt-oss)

- **Puerto 11434**: `http://localhost:11434/api/chat`
- Modelo: `gpt-oss`
- Requiere servidor Ollama u compatible

##### IA Usuario 2 (deepseek-r1)

- **Puerto 11435**: `http://localhost:11435/api/chat`
- Modelo: `deepseek-r1`
- Requiere servidor Ollama u compatible

#### Frontend (si aplica)

- **Puerto 5173**: Cliente web (configurado en CORS)

### Configuración CORS

```go
AllowOrigins: "http://localhost:5173"
AllowMethods: "GET,POST,PUT,PATCH,DELETE,OPTIONS"
```

---

## Variables de Entorno (`.env`)

```env
# Base de datos PostgreSQL
HOST=ip_bd
PORT=5432
DBUSER=name_bd
PASSWORD=password_bd
DBNAME=name_bd

# Autenticación JWT
SECRET_KEY_JWT=jwt_secret_key

# LLM Base (puede ser usado como fallback)
LLM_BASE_URL=http://localhost:1234

# API Keys para modelos IA
LLM_API_KEY_1=<tu_api_key_modelo_1>
LLM_API_KEY_2=<tu_api_key_modelo_2>

# Feature flag para apagar/encender la Inteligencia Artificial
ENABLE_AI_MODELS=false
```

---

## Funcionamiento del Sistema de IA

### 1. **Inicialización de IAs**

Al iniciar la aplicación (`main.go`):

```go
// Se crean servicios de IA para cada configuración
aiServices := make(map[string]*ia.AIService)
for _, config := range ia.AiConfigurations {
    aiService := ia.NewAIService(wsHub, mensajeService, grupoService, config)
    aiServices[config.UserID] = aiService
    go aiService.Start(ctx, 5) // 5 workers concurrentes
}
```

### 2. **Suscripción a Grupos**

Cada IA se suscribe automáticamente a los grupos donde está agregada:

```go
func (s *AIService) SuscribeToGroup() error {
    // Obtiene todos los grupos del usuario IA
    aiUserID, err := s.GrupoService.GetAllGruposByUsuarioIdToIds(idUsuario)
    // Suscribe al Hub
    s.Hub.SubscribeUserToGroups(s.UserID, aiUserID)
}
```

### 3. **Flujo de Mensajes con IA**

#### Paso 1: Usuario envía mensaje

```
Usuario → WebSocket → Hub.Broadcast()
```

#### Paso 2: Hub distribuye el mensaje

```go
// Envía a todos los usuarios del grupo
for userID, conn := range h.clients {
    if h.userGroups[userID][msg.GroupID] {
        conn.WriteMessage(websocket.TextMessage, jsonMsg)
    }
}
// Envía también al canal de IA
h.aiChannel <- msg
```

#### Paso 3: Router de IA procesa el mensaje

```go
for msg := range wsHub.AIChannel() {
    // Ignora mensajes de las propias IAs
    if _, ok := aiServices[msg.SenderID]; ok {
        continue
    }

    // Enruta al servicio de IA correspondiente del grupo
    for _, service := range aiServices {
        if wsHub.CheckUserInGroup(service.UserID, msg.GroupID) {
            service.InputChannel() <- msg
        }
    }
}
```

#### Paso 4: Worker de IA procesa

```go
func (s *AIService) worker(ctx context.Context, id int) {
    for msg := range s.jobs {
        // 1. Obtiene historial del grupo
        allGroupMessages := s.MensajeRepo.GetAllByGrupoClave(msg.GroupID)

        // 2. Construye prompt con contexto
        llmMessages := s.buildPromptFromHistory(allGroupMessages)

        // 3. Llama al LLM
        aiResponse := llm.PostCompletion(llmMessages, s.LLMBaseURL, ...)

        // 4. Parsea respuesta
        aiMsg := s.ParseAIResponse(aiResponse, msg.GroupID, s.UserID)

        // 5. Guarda en BD
        gormMsg := s.saveAIToDB(aiMsg)

        // 6. Difunde respuesta
        s.Hub.Broadcast(*aiMsg)
    }
}
```

#### Paso 5: Hub distribuye respuesta de IA

```
Hub.Broadcast() → WebSocket → Todos los usuarios del grupo
```

### 4. **Formato de Comunicación con LLM**

#### Request al LLM

```json
{
  "messages": [
    { "role": "system", "content": "Eres asistente en español." },
    { "role": "user", "content": "Hola, ¿cómo estás?" },
    { "role": "assistant", "content": "¡Hola! Estoy bien, gracias." },
    { "role": "user", "content": "¿Puedes ayudarme?" }
  ],
  "model": "gpt-oss"
}
```

#### Response del LLM

```json
{
  "content": "Por supuesto, ¿en qué necesitas ayuda?",
  "answer_id": "123" // Opcional: ID del mensaje respondido
}
```

### 5. **Diferenciación de IAs por Colores/Identificación**

Aunque el código no implementa explícitamente colores, el sistema permite diferenciar IAs mediante:

#### Por ID de Usuario (`UserID`)

- IA 1: `UserID = "1"` → Modelo `gpt-oss`
- IA 2: `UserID = "2"` → Modelo `deepseek-r1`

#### En el Frontend (implementación sugerida)

```javascript
// Asignar colores según UserID
const iaColors = {
  1: "#FF6B6B", // Rojo para IA 1 (gpt-oss)
  2: "#4ECDC4", // Verde azulado para IA 2 (deepseek-r1)
  // Usuarios normales: color por defecto
};

function renderMessage(message) {
  const color = iaColors[message.SenderID] || "#333333";
  // Aplicar color al mensaje en la UI
}
```

#### Identificación en Base de Datos

```sql
-- Las IAs son usuarios normales con IDs específicos
SELECT * FROM usuarios WHERE id IN (1, 2);
-- Se les puede agregar campos adicionales:
-- is_ia BOOLEAN, ia_model VARCHAR, color_hex VARCHAR
```

### 6. **Concurrencia y Escalabilidad**

- **Pool de Workers**: Cada IA procesa 5 mensajes simultáneamente
- **Canales con Buffer**: `inputChannel` y `jobs` tienen buffer de 100 mensajes
- **Separación de Contextos**: Cada IA mantiene su propia instancia de `AIService`
- **Thread-Safety**: Uso de `sync.Map` para conversaciones y mutex en Hub

---

## Rutas API

### Públicas (`/api/public`)

- `POST /usuario` - Registro de usuario
- `POST /auth/login` - Inicio de sesión
- `GET /ws/chat` - Conexión WebSocket

### Protegidas (`/api`) - Requieren JWT

#### Autenticación

- `GET /auth/verify` - Verificar token

#### Usuarios

- `GET /usuario/:id` - Obtener usuario
- `GET /usuario` - Listar usuarios
- `PUT /usuario/:id` - Actualizar usuario
- `DELETE /usuario/:id` - Eliminar usuario

#### Grupos

- `POST /grupo` - Crear grupo
- `GET /grupo/:id` - Obtener grupo
- `GET /grupo` - Listar grupos
- `PUT /grupo/:id` - Actualizar grupo
- `DELETE /grupo/:id` - Eliminar grupo

#### Mensajes

- `POST /mensaje` - Enviar mensaje
- `GET /mensaje/:id` - Obtener mensaje
- `GET /mensaje/grupo/:groupId` - Mensajes del grupo

#### Grupo-Usuario

- `POST /grupo-usuario` - Agregar usuario a grupo
- `DELETE /grupo-usuario` - Remover usuario de grupo

---

## Ejecución del Proyecto

### Requisitos Previos

1. **Go 1.24.4+**
2. **PostgreSQL 12+** corriendo en puerto 5432
3. **Ollama u otro servidor LLM** en puertos 11434 y 11435
4. **Modelos descargados**: `gpt-oss` y `deepseek-r1`

### Pasos de Instalación

```bash
# 1. Clonar repositorio
cd /home/moy45/Proyectos_go/chatvist-chat

# 2. Configurar variables de entorno
cp .env.example .env
# Editar .env con tus credenciales

# 3. Instalar dependencias
go mod download

# 4. Iniciar servidores LLM (ejemplo con Ollama)
ollama serve --host 0.0.0.0:11434 --model gpt-oss
ollama serve --host 0.0.0.0:11435 --model deepseek-r1

# 5. Crear base de datos PostgreSQL
createdb chatvist_chat

# 6. Ejecutar aplicación
go run main.go
```

### Logs de Inicio Exitoso

```
¡Hola, mundo desde Go!
Hello, World!
IA suscrita a sus grupos correctamente
Worker 0 iniciado
Worker 1 iniciado
Worker 2 iniciado
Worker 3 iniciado
Worker 4 iniciado
Server is running on port 3100
```

---

## Integración de Nuevas IAs

Para agregar un nuevo modelo de IA:

### 1. Actualizar configuración (`internal/ia/iaConfig.go`)

```go
var AiConfigurations = []IAConfig{
    // IAs existentes...
    {
        UserID:     "3",
        LLMBaseURL: "http://localhost:11436/api/chat",
        LLMName:    "llama3-70b",
        LLMAPIKey:  os.Getenv("LLM_API_KEY_3"),
    },
}
```

### 2. Crear usuario en base de datos

```sql
INSERT INTO usuarios (id, nombre, apodo, email, password, fecha)
VALUES (3, 'IA Llama3', 'Llama', 'ia3@chatvist.com', 'hash_bcrypt', NOW());
```

### 3. Agregar IA a grupos deseados

```sql
INSERT INTO grupos_usuarios (id_grupo, id_usuario)
VALUES (1, 3), (2, 3);
```

### 4. Reiniciar servidor

```bash
go run main.go
```

---

## Troubleshooting

### IA no responde

- Verificar que el servidor LLM esté corriendo en el puerto correcto
- Revisar logs: `Worker X procesando mensaje`
- Confirmar que la IA está agregada al grupo

### Error de conexión a base de datos

- Verificar credenciales en `.env`
- Confirmar que PostgreSQL está corriendo: `pg_isready`

### WebSocket se desconecta

- Verificar configuración CORS
- Revisar que el token JWT sea válido

### LLM responde con error

- Confirmar formato de respuesta JSON: `{"content": "..."}`
- Verificar logs de Ollama/servidor LLM

---

## Arquitectura de Seguridad

### Autenticación

- **Passwords**: Hash bcrypt con salt
- **Tokens**: JWT con expiración configurable
- **Middleware**: Validación de tokens en rutas protegidas

### WebSocket

- Autenticación por token en query params
- Validación de pertenencia a grupo antes de difundir mensajes

### Base de Datos

- Prepared statements (GORM protege contra SQL injection)
- Validación de entradas en capa de servicio

---

## Monitoreo y Logs

El sistema registra:

- Conexiones/desconexiones de usuarios
- Mensajes procesados por IAs
- Errores de LLM
- Estados de workers

Ejemplo de logs:

```
Hub: Usuario 123 registrado.
Worker 2 procesando mensaje {SenderID:123 GroupID:abc Content:Hola}
AIService: recibido mensaje en AIChannel: {Id: SenderID:123 GroupID:abc}
Hub: Usuario 123 desconectado.
```

---

## Mejoras Futuras Sugeridas

1. **Sistema de colores persistente**: Agregar campo `color_hex` en tabla `usuarios`
2. **Rate limiting**: Limitar mensajes por usuario/IA
3. **Historial parcial**: Enviar solo últimos N mensajes al LLM para optimizar
4. **Typing indicators**: Mostrar cuando IA está "escribiendo"
5. **Respuestas en streaming**: Mostrar respuesta de IA en tiempo real
6. **Métricas**: Prometheus/Grafana para monitorear latencia de IAs
7. **Fallback de IAs**: Si una IA falla, usar modelo alternativo
8. **Cache de respuestas**: Redis para respuestas frecuentes

---

## Contacto y Soporte

Para dudas o contribuciones al proyecto, consultar con el equipo de desarrollo.

**Versión de documentación**: 2.0  
**Fecha**: Marzo 2026
