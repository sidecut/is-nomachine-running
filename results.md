# Results

## Go version

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

- 0:05
- 13M
- 12M