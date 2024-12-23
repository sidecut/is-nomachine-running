# Results

30 seconds, 10 users

See the last few results. Calling ps or pgrep SUCKS for throughput.

### Summary

Swift version had same CPU time as Go version, but 5x the throughput and 47% of the memory usage of Go.

Maybe Rust would be better, but Swift is so easy to work with.

## Go version

- 46784 requests sent
- 3:21 CPU time
- 66M peak during
- 28M after a GC after stress test end

## Swift version

_NOTE_ Without using sysctl

- 0:06 CPU time
- 30M peak during
- 25M a few seconds later -- it slowly ramps down

## Swift version with sysctl

- 4:33 CPU time
- 33M peak during
- 24M a few seconds later -- it slowly ramps down

### Second test

- 4:31
- 41M

## Swift version after code improvments from Copilot

- 3:40
- 34M
- 21M

## After fixing bug -- not much optimization

`22341cf (Fix bug where getStatus always returned noMachineRunning and clientAttached true, 2024-12-22)`

- 3:33
- 29M
- 27M

## Using pgrep

`1edaa9e (Use pgrep, 2024-12-22)`

- 3074 requests sent
- 0:05
- 13M
- 12M

## Revert pgrep

- 237979 requests sent
- 3:34
- 31M
