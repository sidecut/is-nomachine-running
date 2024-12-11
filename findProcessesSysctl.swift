import Foundation

func findProcesses(named processName: String) -> [Int32] {
    var pids: [Int32] = []
    var taskInfo = kinfo_proc()
    var taskInfoCount = Int(MemoryLayout<kinfo_proc>.size)
    var mib: [Int32] = [CTL_KERN, KERN_PROC, KERN_PROC_ALL]

    let result = sysctl(&mib, UInt32(mib.count), &taskInfo, &taskInfoCount, nil, 0)

    if result == 0 {
        let numberOfProcesses = taskInfoCount / Int(MemoryLayout<kinfo_proc>.size)
        let processList = UnsafeBufferPointer(start: &taskInfo, count: Int(numberOfProcesses))

        for var proc in processList {
            let processID = proc.kp_proc.p_pid
            let processName = withUnsafePointer(to: &proc.kp_proc.p_comm) {
                $0.withMemoryRebound(to: CChar.self, capacity: Int(MAXCOMLEN)) {
                    String(cString: $0)
                }
            }

            if processName == processName {
                pids.append(processID)
            }
        }
    } else {
        print("Error: sysctl failed")
    }

    return pids
}

let pids = findProcesses(named: "fred")
print("PIDs of processes named 'fred': \(pids)")
