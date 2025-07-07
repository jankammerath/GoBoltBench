# GoBoltBench

This is a benchmark using Shakespeare and bbolt written in Go to benchmark AWS Graviton2 (`arm64`) with IBM LinuxONE (`s390x`).

## IBM LinuxONE

```bash
[linux1@linuxonedev goboltbench]$ ./gobbench-s390x
GoBoltBench — IBM/S390 (3.65 GB)
2025/07/07 15:29:48 Processing 899588 lines with 16 workers
2025/07/07 15:30:09 All workers completed successfully
2025/07/07 15:30:09 Total processing time: 21.515265237s
```