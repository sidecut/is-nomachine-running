import Foundation

func getProcessList() -> [(pid: pid_t, name: String, username: String)] {
    var name: [Int32] = [CTL_KERN, KERN_PROC, KERN_PROC_ALL, 0]
    var size: size_t = 0

    // Get size needed for buffer
    guard sysctl(&name, 3, nil, &size, nil, 0) == 0 else {
        print("Error getting process list size: \(String(cString: strerror(errno)))")
        return []
    }

    // Allocate memory for process list
    let count = size / MemoryLayout<kinfo_proc>.stride
    var procList = Array(repeating: kinfo_proc(), count: count)

    // Get actual process list
    guard sysctl(&name, 3, &procList, &size, nil, 0) == 0 else {
        print("Error getting process list: \(String(cString: strerror(errno)))")
        return []
    }

    var processes: [(pid: pid_t, name: String, username: String)] = []

    for proc in procList {
        let pid = proc.kp_proc.p_pid

        // Get process name
        if pid == 0 { continue }  // Skip kernel process

        // Get username from uid
        let uid = proc.kp_eproc.e_ucred.cr_uid
        var username = "unknown"
        if let pwd = getpwuid(uid) {
            username = String(cString: pwd.pointee.pw_name)
        }

        // Get process name using proc_name
        var name = [CChar](repeating: 0, count: Int(256))
        proc_name(pid, &name, UInt32(256))
        let processName = String(cString: name)

        if !processName.isEmpty {
            processes.append((pid: pid, name: processName, username: username))
        }
    }

    return processes
}

// // Usage example:
// let processes = getProcessList()
// for process in processes {
//     print("PID: \(process.pid), Name: \(process.name), User: \(process.username)")
// }

func main() {
    // Get processes whose name matches any of the command line arguments
    let processes = getProcessList()
    let args = CommandLine.arguments.dropFirst()
    let matchingProcesses = processes.filter { args.contains($0.name) }
    for process in matchingProcesses {
        print("PID: \(process.pid), Name: \(process.name), User: \(process.username)")
    }
    if matchingProcesses.isEmpty {
        print("No matching processes found.")
    }
}
