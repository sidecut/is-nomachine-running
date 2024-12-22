import Darwin  // Import necessary C functions for system calls
import Foundation

struct NoMachineStatus: Encodable {
    var hostName: String?
    var noMachineRunning: Bool = false
    var clientAttached: Bool = false
}

func getRunningProcesses() -> [String] {
    var processList: [String] = []

    // Use sysctl to get process information
    var mib = [CTL_KERN, KERN_PROC, KERN_PROC_ALL, 0]
    var size = 0

    // Get the size of the process list
    sysctl(&mib, UInt32(mib.count), nil, &size, nil, 0)

    // Allocate memory for the process list
    var processData = UnsafeMutablePointer<kinfo_proc>.allocate(capacity: Int(size))

    // Get the process list
    sysctl(&mib, UInt32(mib.count), processData, &size, nil, 0)

    // Iterate through each process
    var currentProcess = processData
    while currentProcess.pointee.kp_proc.p_flag != 0 {
        let processName = String(cString: currentProcess.pointee.kp_proc.p_comm)
        processList.append(processName)
        currentProcess = currentProcess.advanced(by: 1)
    }

    // Deallocate memory
    processData.deallocate()

    return processList
}

// // Usage example
// let runningProcesses = getRunningProcesses()
// print(runningProcesses)

func getFirstProcessByName(_ name: String) -> Int32? {
    let processes = getRunningProcesses()
    for process in processes {
        if process.contains(name) {
            // Assuming the process name is followed by its PID
            let components = process.split(separator: " ")
            if let pidString = components.first, let pid = Int32(pidString) {
                return pid
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
