import Foundation

struct NoMachineStatus: Encodable {
    var hostName: String?
    var noMachineRunning: Bool = false
    var clientAttached: Bool = false
}

func getFirstProcessByName(_ name: String) -> Int32? {
    // var processes = ProcessInfo.processInfo.processIdentifier
    let task = Process()
    task.launchPath = "/bin/ps"
    task.arguments = ["-A"]

    let pipe = Pipe()
    task.standardOutput = pipe
    task.launch()

    let data = pipe.fileHandleForReading.readDataToEndOfFile()
    if let output = String(data: data, encoding: .utf8) {
        for line in output.components(separatedBy: "\n") {
            if line.contains(name) {
                let components = line.trimmingCharacters(in: .whitespaces)
                    .components(separatedBy: " ")
                    .filter { !$0.isEmpty }
                if let pid = Int32(components[0]) {
                    return pid
                }
            }
        }
    }
    return nil
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

    if let noMachinePid = getFirstProcessByName("nxserver.bin"), noMachinePid > 0 {
        status.noMachineRunning = true
    }

    if let noMachineClientPid = getFirstProcessByName("nxexec"), noMachineClientPid > 0 {
        status.clientAttached = true
    }

    return .success(status)
}
