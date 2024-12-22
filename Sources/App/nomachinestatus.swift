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

struct processResult {
    var pid: Int
    var name: String
}

func getRunningProcesses() -> [processResult] {
    var memoryInformationBase = [
        KERN_PROC,  // get the process id
        KERN_PROC_ALL,  // get everything including the name
        // KERN_USER,  // get the user who started the process
    ]

    var bufferSize = 0
    let bufferSizeResult = sysctl(
        &memoryInformationBase,
        UInt32(memoryInformationBase.count),
        nil,
        &bufferSize,
        nil,
        0
    )
    if bufferSizeResult < 0 {
        perror(&errno)
        return []
    }

    let entryCount = bufferSize / MemoryLayout<kinfo_proc>.stride
    var processList: UnsafeMutablePointer<kinfo_proc>?
    processList = UnsafeMutablePointer.allocate(capacity: entryCount)
    defer { processList?.deallocate() }

    let populateProcessListResult = sysctl(
        &memoryInformationBase,
        UInt32(memoryInformationBase.count),
        processList,
        &bufferSize,
        nil,
        0
    )
    if populateProcessListResult < 0 {
        perror(&errno)
        return []
    }

    var processes = [processResult]()

    for index in 0..<entryCount {
        let process = processList![index]
        let processId = process.kp_proc.p_pid
        if processId == 0 {
            continue
        }
        let comm = process.kp_proc.p_comm
        let name = String(cString: Mirror(reflecting: comm).children.map { $0.value as! CChar })
        // print("PID: \(processId), App: \(name)")
        processes.append(processResult(pid: Int(processId), name: name))
    }

    return processes
}
