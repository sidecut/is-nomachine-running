import Foundation
import Vapor

struct NoMachineStatus: Content {
    var hostName: String?
    var noMachineRunning: Bool = false
    var clientAttached: Bool = false
}

func getStatus() throws -> Result<NoMachineStatus, Error> {
    var status = NoMachineStatus()

    guard let hostName = Host.current().localizedName else {
        return .failure(
            NSError(
                domain: "", code: -1,
                userInfo: [NSLocalizedDescriptionKey: "Could not get hostname"]))
    }

    status.hostName = hostName
    status.noMachineRunning = try isProcessRunning(name: "nxserver.bin")
    status.clientAttached = try isProcessRunning(name: "nxexec")

    return .success(status)
}

struct ProcessResult: Content {
    var pid: Int
    var name: String
}

enum SysCtlError: Error {
    case FailedToGetProcessList1(message: String = "Failed to get process list step 1")
    case FailedToGetProcessList2(message: String = "Failed to get process list step 2")
}

func isProcessRunning(name: String) throws -> Bool {
    let process = Process()
    process.executableURL = URL(fileURLWithPath: "/usr/bin/pgrep")
    process.arguments = [name]

    let pipe = Pipe()
    process.standardOutput = pipe

    try process.run()
    process.waitUntilExit()
    let data = pipe.fileHandleForReading.readDataToEndOfFile()
    return !data.isEmpty
}
