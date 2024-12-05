import Logging
import NIOCore
import NIOPosix
import Vapor

// Set default configuration values
let defaultConfig = [
    "port": 80,
    "sslport": 443,
]

// Status API handler
@Sendable func statusAPI(_ req: Request) throws -> EventLoopFuture<Response> {
    do {
        let statusResult = try getStatus()
        let status: String
        switch statusResult {
        case .success(let noMachineStatus):
            status = "\(noMachineStatus)"  // Assuming NoMachineStatus conforms to CustomStringConvertible or has a suitable description
        case .failure(let error):
            throw error
        }
        let response = Response(status: .ok, body: .init(string: status))
        response.headers.replaceOrAdd(name: .contentType, value: "application/json")
        return req.eventLoop.makeSucceededFuture(response)
    } catch {
        // TODO: log this error
        throw error
    }
}

@main
enum Entrypoint {
    static func main() async throws {
        var env = try Environment.detect()
        try LoggingSystem.bootstrap(from: &env)

        let app = try await Application.make(env)

        // This attempts to install NIO as the Swift Concurrency global executor.
        // You can enable it if you'd like to reduce the amount of context switching between NIO and Swift Concurrency.
        // Note: this has caused issues with some libraries that use `.wait()` and cleanly shutting down.
        // If enabled, you should be careful about calling async functions before this point as it can cause assertion failures.
        // let executorTakeoverSuccess = NIOSingletons.unsafeTryInstallSingletonPosixEventLoopGroupAsConcurrencyGlobalExecutor()
        // app.logger.debug("Tried to install SwiftNIO's EventLoopGroup as Swift's global concurrency executor", metadata: ["success": .stringConvertible(executorTakeoverSuccess)])

        do {
            try await configure(app)
        } catch {
            app.logger.report(error: error)
            try? await app.asyncShutdown()
            throw error
        }
        try await app.execute()
        try await app.asyncShutdown()
    }
}

// configures your application
public func configure(_ app: Application) async throws {
    let corsConfiguration = CORSMiddleware.Configuration(
        allowedOrigin: .all,
        allowedMethods: [.GET, .POST, .PUT, .OPTIONS, .DELETE, .PATCH],
        allowedHeaders: [.accept, .authorization, .contentType, .origin, .xRequestedWith]
    )

    // Configure middleware
    let cors = CORSMiddleware(configuration: corsConfiguration)
    app.middleware.use(cors)
    app.middleware.use(FileMiddleware(publicDirectory: "dist"))

    // Configure routes
    app.get("api", use: statusAPI)

    // Configure logging
    app.logger.logLevel = .info
    app.logger.info("*** STARTING PID \(ProcessInfo.processInfo.processIdentifier)")
}
