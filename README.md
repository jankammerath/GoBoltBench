# GoBoltBench

This is a benchmark is written in Go and uses [shakespeare.txt](https://gist.github.com/blakesanie/dde3a2b7e698f52f389532b4b52bc254) with [bbolt](https://github.com/etcd-io/bbolt) to benchmark AWS Graviton3 on `arm64` systems with IBM LinuxONE `s390x`. The application performs insertion of messages extracted from Shakespeare with random users assigned to them and the `sha256` hash calculated for each message. By doing so, the benchmark intents to mimic usual social media operation using the [etcd-io/bbolt](https://github.com/etcd-io/bbolt) key-value store.

## IBM LinuxONE

This benchmark used a 4GB instance from [Marist University](https://www.marist.edu).

Cost for this instance on IBM Cloud in `Frankfurt (eu-de)` would be **$78.36/month**.

```

 ___ ____  __  __      _     _                   ___  _   _ _____
|_ _| __ )|  \/  |    | |   (_)_ __  _   ___  __/ _ \| \ | | ____|
 | ||  _ \| |\/| |    | |   | | '_ \| | | \ \/ / | | |  \| |  _|
 | || |_) | |  | |    | |___| | | | | |_| |>  <| |_| | |\  | |___
|___|____/|_|  |_|    |_____|_|_| |_|\__,_/_/\_\\___/|_| \_|_____|

=================================================================================
Welcome to the IBM LinuxONE Community Cloud!

This server is for authorized users only. All activity is logged and monitored.
Individuals using this server must abide to the Terms and Conditions listed here:
https://www.ibm.com/community/z/ibm-linuxone-community-cloud-terms-and-conditions/
Your access will be revoked for any non-compliance.
==================================================================================

[linux1@linuxonedev goboltbench]$ ./gobbench-s390x
GoBoltBench — IBM/S390 (3.65 GB)
Red Hat Enterprise Linux 9.6 (Plow)
2025/07/07 15:47:41 Processing 114634 lines with 16 workers
2025/07/07 15:47:44 Total processing time: 2.826166479s
```

## AWS Graviton3

This benchmark used a `c7g.large` using and `io2` EBS (8,000 iops) instance in `us-east-1`.

Cost for this instance on AWS in `Frankfurt (eu-central-1)` would be **$67.67/month**.

```

   ,     #_
   ~\_  ####_        Amazon Linux 2023
  ~~  \_#####\
  ~~     \###|
  ~~       \#/ ___   https://aws.amazon.com/linux/amazon-linux-2023
   ~~       V~' '->
    ~~~         /
      ~~._.   _/
         _/ _/
       _/m/'

[ec2-user@ip-172-31-81-25 goboltbench]$ ./gobbench-graviton3
GoBoltBench — Unknown CPU (4.00 GB)
Amazon Linux 2023.7.20250623
2025/07/07 21:07:45 Processing 114634 lines with 16 workers
2025/07/07 21:07:47 Total processing time: 1.850152923s
```

## Apple Macbook Pro M1 Pro

This is the local development reference system.

```bash
jan@MacBook-Pro-von-Jan GoBoltBench % ./bin/gobbench
GoBoltBench — Apple M1 Pro (17.18 GB)
2025/07/07 22:48:01 Processing 114634 lines with 16 workers
2025/07/07 22:48:13 Total processing time: 12.516969333s
```