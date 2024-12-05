import Foundation
import Vapor

// Set default configuration values
let defaultConfig = [
    "port": 80,
    "sslport": 443,
]

// Status API handler
func statusAPI(_ req: Request) throws -> EventLoopFuture<Response> {
    do {
        let status = try getStatus()
        return try JSONResponse(status).encodeResponse(for: req)
    } catch {
        // TODO: log this error
        throw error
    }
}

func main() throws {
    // Main application
    let app = try Application()
    let corsConfiguration = CORSMiddleware.Configuration(
        allowedOrigin: .all,
        allowedMethods: [.GET, .POST, .PUT, .OPTIONS, .DELETE, .PATCH],
        allowedHeaders: [.accept, .authorization, .contentType, .origin, .xRequestedWith]
    )

    // Configure middleware
    let cors = CORSMiddleware(configuration: corsConfiguration)
    app.middleware.use(cors)
    app.middleware.use(FileMiddleware(publicDirectory: "dist"))

    // Configure environment variables
    let env = Environment.get("ISNO_PORT") ?? String(defaultConfig["port"]!)
    let sslPort = Environment.get("ISNO_SSLPORT") ?? String(defaultConfig["sslport"]!)
    let port = Int(env) ?? defaultConfig["port"]!
    let sslPortNum = Int(sslPort) ?? defaultConfig["sslport"]!

    // Configure routes
    app.get("api", use: statusAPI)

    // Configure logging
    app.logger.logLevel = .info
    app.logger.info("*** STARTING PID \(ProcessInfo.processInfo.processIdentifier)")

    // Start SSL server
    DispatchQueue.global().async {
        do {
            try app.start(tls: .automatic, port: sslPortNum)
        } catch {
            app.logger.error("SSL Server failed to start: \(error)")
            app.shutdown()
        }
    }

    // Start regular server
    try app.start(port: port)

    // Handle shutdown signals
    signal(SIGTERM) { signal in
        app.logger.warning(
            "*** STOPPING PID \(ProcessInfo.processInfo.processIdentifier) with signal \(signal)")
        app.shutdown()
    }

    signal(SIGINT) { signal in
        app.logger.warning(
            "*** STOPPING PID \(ProcessInfo.processInfo.processIdentifier) with signal \(signal)")
        app.shutdown()
    }

    // Run event loop
    try app.run()
}
