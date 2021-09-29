# IPC Benchmark (Named Pipes, TCP, UDP)

### Requirements

- Go (1.17+)
- Bash shell (optional)

### Usage

```shell
cd <project-directory>
go run ./cmd -s <message-size> -c <roundtrip-count>
# optionally you can run predefined benchmarks
./reproduce.sh
# example:
go run ./cmd -c 50000 -s 10
# result:
+------------+------------+----------------------+--------------+--------------+
|    IPC     |    TYPE    | TRANSFERRED MESSAGES |   DURATION   | SUCCESS RATE |
+------------+------------+----------------------+--------------+--------------+
| Named Pipe | Throughput |               506081 | 1s           |       100.00 |
|            | Latency    |                50000 | 1.793193991s |       100.00 |
| TCP        | Throughput |               113159 | 1s           |       100.00 |
|            | Latency    |                50000 | 1.98712208s  |       100.00 |
| UDP        | Throughput |               166389 | 1s           |       100.00 |
|            | Latency    |                50000 | 787.750032ms |       100.00 |
+------------+------------+----------------------+--------------+--------------+
```

-s <message-size> is in bytes

-c <roundtrip-count> is message count

P.S. Tested on Arch Linux.

## Benchmarks

In my personal laptop the results are as the following:

CPU: Intel i7-4700HQ (8) @ 3.400GHz

Memory: 15 GB

```shell
128 bytes, count 1
+------------+------------+----------------------+-----------+--------------+
|    IPC     |    TYPE    | TRANSFERRED MESSAGES | DURATION  | SUCCESS RATE |
+------------+------------+----------------------+-----------+--------------+
| Named Pipe | Throughput |                 7693 | 1s        |       100.00 |
|            | Latency    |                    1 | 158.496µs |       100.00 |
| TCP        | Throughput |                 2803 | 1s        |       100.00 |
|            | Latency    |                    1 | 286.77µs  |       100.00 |
| UDP        | Throughput |                12102 | 1s        |       100.00 |
|            | Latency    |                    1 | 106.432µs |       100.00 |
+------------+------------+----------------------+-----------+--------------+

512 bytes, count 10
+------------+------------+----------------------+-----------+--------------+
|    IPC     |    TYPE    | TRANSFERRED MESSAGES | DURATION  | SUCCESS RATE |
+------------+------------+----------------------+-----------+--------------+
| Named Pipe | Throughput |                59743 | 1s        |       100.00 |
|            | Latency    |                   10 | 400.245µs |       100.00 |
| TCP        | Throughput |                21586 | 1s        |       100.00 |
|            | Latency    |                   10 | 609.316µs |       100.00 |
| UDP        | Throughput |                71897 | 1s        |       100.00 |
|            | Latency    |                   10 | 390.56µs  |       100.00 |
+------------+------------+----------------------+-----------+--------------+

1024 bytes, count 100
+------------+------------+----------------------+------------+--------------+
|    IPC     |    TYPE    | TRANSFERRED MESSAGES |  DURATION  | SUCCESS RATE |
+------------+------------+----------------------+------------+--------------+
| Named Pipe | Throughput |               296190 | 1s         |       100.00 |
|            | Latency    |                  100 | 4.322611ms |       100.00 |
| TCP        | Throughput |                56434 | 1s         |       100.00 |
|            | Latency    |                  100 | 5.809666ms |       100.00 |
| UDP        | Throughput |               107315 | 1s         |       100.00 |
|            | Latency    |                  100 | 2.805755ms |       100.00 |
+------------+------------+----------------------+------------+--------------+

4096 bytes, count 1000
+------------+------------+----------------------+--------------+--------------+
|    IPC     |    TYPE    | TRANSFERRED MESSAGES |   DURATION   | SUCCESS RATE |
+------------+------------+----------------------+--------------+--------------+
| Named Pipe | Throughput |               200771 | 1s           |       100.00 |
|            | Latency    |                 1000 | 99.432648ms  |       100.00 |
| TCP        | Throughput |                91678 | 1s           |       100.00 |
|            | Latency    |                 1000 | 113.745119ms |       100.00 |
| UDP        | Throughput |               121068 | 1s           |       100.00 |
|            | Latency    |                 1000 | 23.711583ms  |       100.00 |
+------------+------------+----------------------+--------------+--------------+
```
