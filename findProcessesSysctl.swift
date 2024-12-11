import Foundation

func findProcesses(named processName: String) -> [Int32] {
    var pids: [Int32] = []
    // var taskInfo = kinfo_proc()
    // var taskInfoCount = Int(MemoryLayout<kinfo_proc>.size)

    // largely from https://gaitatzis.medium.com/listing-running-system-processes-using-swift-43e24c20789c
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
        print("bufferSizeResult < 0")
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

    for index in 0..<entryCount {
        let process = processList![index]
        let processId = process.kp_proc.p_pid
        if processId == 0 {
            continue
        }
        let comm = process.kp_proc.p_comm
        let name = String(cString: Mirror(reflecting: comm).children.map { $0.value as! CChar })
        print("PID: \(processId), App: \(name)")
        pids.append(processId)
    }

    return pids
}

let pids = findProcesses(named: "zsh")
print("PIDs of processes named 'zsh': \(pids)")
