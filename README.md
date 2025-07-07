# GoBoltBench

This is a benchmark is written in Go and uses [shakespeare.txt](https://gist.github.com/blakesanie/dde3a2b7e698f52f389532b4b52bc254) with [bbolt](https://github.com/etcd-io/bbolt) to benchmark AWS Graviton2 on `arm64` systems with IBM LinuxONE `s390x`. The application performs insertion of messages extracted from Shakespeare with random users assigned to them and the `sha256` hash calculated for each message. By doing so, the benchmark intents to mimic usual social media operation using the [etcd-io/bbolt](https://github.com/etcd-io/bbolt) key-value store.

## IBM LinuxONE

```bash
[linux1@linuxonedev goboltbench]$ ./gobbench-s390x
GoBoltBench — IBM/S390 (3.65 GB)
Red Hat Enterprise Linux 9.6 (Plow)
2025/07/07 15:47:41 Processing 114634 lines with 16 workers
2025/07/07 15:47:44 Total processing time: 2.826166479s
```

## Apple Macbook Pro M1 Pro

```bash
jan@MacBook-Pro-von-Jan GoBoltBench % ./bin/gobbench
GoBoltBench — Apple M1 Pro (17.18 GB)
2025/07/07 22:48:01 Processing 114634 lines with 16 workers
2025/07/07 22:48:13 Total processing time: 12.516969333s
```