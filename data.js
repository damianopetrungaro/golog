window.BENCHMARK_DATA = {
  "lastUpdate": 1651371820709,
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
          "id": "044e160a6f0512de475b441f87d321c0b883fe68",
          "message": "ci: test github pages deployment",
          "timestamp": "2022-04-30T23:07:06Z",
          "url": "https://github.com/damianopetrungaro/golog/pull/31/commits/044e160a6f0512de475b441f87d321c0b883fe68"
        },
        "date": 1651371820226,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkLogger/logrus",
            "value": 9177,
            "unit": "ns/op\t    6125 B/op\t      69 allocs/op",
            "extra": "121102 times\n2 procs"
          },
          {
            "name": "BenchmarkLogger/golog",
            "value": 2382,
            "unit": "ns/op\t    2840 B/op\t      27 allocs/op",
            "extra": "471121 times\n2 procs"
          },
          {
            "name": "BenchmarkLogger/zap",
            "value": 2491,
            "unit": "ns/op\t    2825 B/op\t      20 allocs/op",
            "extra": "465554 times\n2 procs"
          }
        ]
      }
    ]
  }
}