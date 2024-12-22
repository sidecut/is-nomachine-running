import Testing
import XCTVapor

@testable import IsNoMachineRunning

@Suite("App Tests")
struct AppTests {
    private func withApp(_ test: (Application) async throws -> Void) async throws {
        let app = try await Application.make(.testing)
        do {
            try await configure(app)
            try await test(app)
        } catch {
            try await app.asyncShutdown()
            throw error
        }
        try await app.asyncShutdown()
    }

    @Test("Test api")
    func api() async throws {
        try await withApp { app in
            try await app.test(
                .GET, "api",
                afterResponse: { res async in
                    #expect(res.status == .ok)
                    #expect(res.body.string.contains("hostName"))
                })
        }
    }
}
