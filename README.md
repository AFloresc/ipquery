# IPQuery Service 🚀

Un microservicio de alto rendimiento desarrollado en **Go** para consultar información detallada de direcciones IP en tiempo real. Este servicio actúa como una capa de abstracción entre una API externa y tu aplicación, proporcionando una estructura de datos unificada, caché en memoria y resiliencia ante fallos.

## 🛠 Características
* **Caché Eficiente:** Utiliza `sync.Map` para almacenar resultados en memoria por 1 hora, reduciendo latencia y consumo de API.
* **Resiliencia:** Incluye `timeouts`, manejo de errores y un endpoint de `health check` para despliegues en la nube.
* **Mapeo de Datos:** Implementa el patrón adaptador para asegurar que la respuesta cumpla con un contrato JSON estricto.
* **Listo para Cloud:** Configurado para desplegarse fácilmente en plataformas como **Render** o Docker.

## 🚀 Instalación y Ejecución

1. **Clonar el repositorio:**
   ```bash
   git clone [https://github.com/tu-usuario/ipquery.git](https://github.com/tu-usuario/ipquery.git)
   cd ipquery
   
2. **Instalar dependencias:**
   ```bash
   go mod tidy

3. **Instalar dependencias:**
   ```bash
   go run cmd/server/main.go

## 🌐 Endpoints

### GET /v1/health
Verifica si el servicio está operativo.
#### Respuesta: 200 OK
### GET /v1/ip/{ip}
Obtiene información detallada de la IP proporcionada.
#### Respuesta (JSON):
 ```json
    {
        "ip": "185.107.108.208",
        "isp": { "asn": "AS200434", "org": "Estabanell Impulsa S.A", "isp": "Estabanell Impulsa S.A" },
        "location": {
            "country": "Spain",
            "country_code": "ES",
            "city": "Òrrius",
            "state": "Catalonia",
            "zipcode": "08317",
            "latitude": 41.5341,
            "longitude": 2.3513,
            "timezone": "Europe/Madrid",
            "localtime": "2026-07-05T08:03:22"
        },
        "risk": { "is_mobile": false, "is_vpn": false, "is_tor": false, "is_proxy": false, "is_datacenter": false, "risk_score": 0 }
    }
 ```

 
## 🐳 Despliegue en Render

1. Sube tu código a un repositorio en GitHub.
2. En Render, crea un nuevo Web Service conectando tu repositorio.
3. Render detectará automáticamente el puerto mediante la variable de entorno PORT.
4. Añade un Health Check Path configurado en /health.

## 🏗 Arquitectura
El proyecto sigue el estándar de diseño de aplicaciones en Go, separando responsabilidades:

* **cmd/** : Punto de entrada de la aplicación.
* **internal/ipinfo/** : Lógica de negocio, caché y proveedor de datos.
* **internal/handlers/** : Controladores HTTP.