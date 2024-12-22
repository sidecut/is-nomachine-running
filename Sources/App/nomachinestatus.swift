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

    let nxServerProcess = getRunningProcesses(searchForNameExact: "nxserver.bin")
    if nxServerProcess.count > 0 {
        status.noMachineRunning = true
    }

    let nxExecProcess = getRunningProcesses(searchForNameExact: "nxexec")
    if nxExecProcess.count > 0 {
        status.clientAttached = true
    }

    return .success(status)
}

struct processResult: Content {
    var pid: Int
    var name: String
}

func getRunningProcesses(searchForNameExact: String? = nil) -> [processResult] {
    var memoryInformationBase = [
        CTL_KERN,  // kernel-related parameters and information
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
    let processList: UnsafeMutablePointer<kinfo_proc>?
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
        let processResult = processResult(pid: Int(processId), name: name)

        if let searchForName = searchForNameExact {
            // Do case insensitive search
            if name.lowercased() == searchForName.lowercased() {
                return [processResult]
            }
        }

        processes.append(processResult)
    }

    return processes
}
