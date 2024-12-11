import Foundation

func findProcesses(named processName: String) -> [Int32] {
    let task = Process()
    let pipe = Pipe()

    task.launchPath = "/bin/ps"
    task.arguments = ["-eo", "pid,comm"]
    task.standardOutput = pipe

    task.launch()

    let data = pipe.fileHandleForReading.readDataToEndOfFile()
    let output = String(data: data, encoding: .utf8)!
    let lines = output.split(separator: "\n")

    var pids: [Int32] = []
    for line in lines {
        let parts = line.split(separator: " ")
        if parts.count >= 2 && parts[1] == processName {
            if let pid = Int32(parts[0]) {
                pids.append(pid)
            }
        }
    }

    return pids
}

// Get process name from command line arguments
guard CommandLine.arguments.count > 1 else {
    print("Please provide a process name as argument")
    exit(1)
}

let processName = CommandLine.arguments[1]
let pids = findProcesses(named: processName)
print("PIDs of processes named '\(processName)': \(pids)")
