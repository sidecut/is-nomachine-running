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

    let nxServerProcess = try getRunningProcesses(searchForNameExact: "nxserver.bin")
    if nxServerProcess.count > 0 {
        status.noMachineRunning = true
    }

    let nxExecProcess = try getRunningProcesses(searchForNameExact: "nxexec")
    if nxExecProcess.count > 0 {
        status.clientAttached = true
    }

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

func getRunningProcesses(searchForNameExact: String? = nil) throws -> [ProcessResult] {
    var mib = [CTL_KERN, KERN_PROC, KERN_PROC_ALL]
    var size = 0

    guard sysctl(&mib, UInt32(mib.count), nil, &size, nil, 0) == 0 else {
        throw SysCtlError.FailedToGetProcessList1()
    }

    let entryCount = size / MemoryLayout<kinfo_proc>.stride
    let processList = UnsafeMutablePointer<kinfo_proc>.allocate(capacity: entryCount)
    defer { processList.deallocate() }

    guard sysctl(&mib, UInt32(mib.count), processList, &size, nil, 0) == 0 else {
        throw SysCtlError.FailedToGetProcessList2()
    }

    var processes = [ProcessResult]()

    for index in 0..<entryCount {
        var process = processList[index]
        let processId = process.kp_proc.p_pid
        if processId == 0 {
            continue
        }
        let processPcom = process.kp_proc.p_comm
        let name = withUnsafePointer(to: &process.kp_proc.p_comm) {
            $0.withMemoryRebound(
                to: CChar.self, capacity: MemoryLayout.size(ofValue: processPcom)
            ) {
                String(cString: $0)
            }
        }
        let processResult = ProcessResult(pid: Int(processId), name: name)

        if let searchForName = searchForNameExact, name.lowercased() == searchForName.lowercased() {
            return [processResult]
        }

        processes.append(processResult)
    }

    return processes
}
