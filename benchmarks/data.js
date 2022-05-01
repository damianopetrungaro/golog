window.BENCHMARK_DATA = {
  "lastUpdate": 1651372355689,
  "repoUrl": "https://github.com/damianopetrungaro/golog",
  "entries": {
    "Log benchmarks": [
      {
        "commit": {
          "author": {
            "name": "damianopetrungaro",
            "username": "damianopetrungaro"
          },
          "committer": {
            "name": "damianopetrungaro",
            "username": "damianopetrungaro"
          },
          "id": "5d9e7264b5d22d917302c00d2bed9218c740af44",
          "message": "ci: test github pages deployment",
          "timestamp": "2022-04-30T23:07:06Z",
          "url": "https://github.com/damianopetrungaro/golog/pull/31/commits/5d9e7264b5d22d917302c00d2bed9218c740af44"
        },
        "date": 1651372355200,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkLogger/logrus",
            "value": 8964,
            "unit": "ns/op\t    6124 B/op\t      69 allocs/op",
            "extra": "131916 times\n2 procs"
          },
          {
            "name": "BenchmarkLogger/golog",
            "value": 2244,
            "unit": "ns/op\t    2840 B/op\t      27 allocs/op",
            "extra": "507265 times\n2 procs"
          },
          {
            "name": "BenchmarkLogger/zap",
            "value": 2428,
            "unit": "ns/op\t    2825 B/op\t      20 allocs/op",
            "extra": "487958 times\n2 procs"
          }
        ]
      }
    ]
  }
}