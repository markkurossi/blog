# Measure First

## Processing Pipeline

```
Compile ──> Inputs ──> Garble ──> Write ───> Results
                         ^          │
                         │          │
                         └──────────┘
```

Signature computation times:

```
┌─────────────┬─────────────────┬────────┬───────┐
│ Op          │            Time │      % │  Xfer │
├─────────────┼─────────────────┼────────┼───────┤
│ Compile     │    2.257813981s │  2.57% │       │
│ Init        │      1.681455ms │  0.00% │    0B │
│ OT Init     │    281.670059ms │  0.32% │  16kB │
│ Peer Inputs │    4.026191293s │  4.58% │ 667kB │
│ Garble      │ 1m21.278656046s │ 92.52% │  15GB │
│ Result      │      1.402231ms │  0.00% │   8kB │
│ Total       │ 1m27.847415065s │        │  15GB │
└─────────────┴─────────────────┴────────┴───────┘
```

Raw gate garbling speed for Ed25519:

| Gate  | Count     |  ns/Gate | Time |
|:------|----------:|---------:|-----:|
| XOR   | 533261481 |   20.64  | 11.0 |
| XNOR  | 28815787  |   20.40  |  0.6 |
| AND   | 267491441 |   176.5  | 47.2 |
| OR    | 494216    |   152.0  |  0.1 |
| INV   | 19784     |   83.47  |  0.0 |
| Total |           |          | 58.9 |

Time to transmit 15GB of data:

```
$ time ./io
2428928	762	0	2
49152	12	524288	85
12288	3	524288	81
24576	6	524288	80
24576	6	524288	81
16384	4	524288	83
20480	5	524288	86
4096	1	524288	80

real	0m14.025s
user	0m4.578s
sys	0m16.632s
```

The theoretical minimum time is 2.3s + 4.0s + 58.9s + 14.0s = 1m19.2s.

The theoretical minimum garble time is 58.9s+14s = 72.9s. We have
measured 69.8s (stream time, 1m19s out of which circuit garbling plus
transmit takes 69.8s). So the overhead is 3s.


## Optimizing the Pipeline

Optimized AND garbling:

| Gate  | Count     |  ns/Gate | Time |
|:------|----------:|---------:|-----:|
| XOR   | 533261481 |   20.64  | 11.0 |
| XNOR  | 28815787  |   20.40  |  0.6 |
| AND   | 267491441 |   158.0  | 42.3 |
| OR    | 494216    |   152.0  |  0.1 |
| INV   | 19784     |   83.47  |  0.0 |
| Total |           |          | 54.0 |

Optimizing wire encoding:

| Encoding      |      Bytes |      Save | Save % |   -Tx | ns/op | +Garble |
|:--------------|-----------:|----------:|-------:|------:|------:|--------:|
| 16/32bit      | 5619637622 |         0 |   0.0% | 0.00s | 2.035 |   0.00s |
| 8/16/24/32bit | 5523203621 |  96434001 |   1.7% | 0.24s | 2.362 |   0.27s |
| 7bit          | 5320775630 | 298861992 |   5.6% | 0.78s | 9.587 |   6.27s |
| 7bit  inline  | 5320775630 | 298861992 |   5.6% | 0.78s | 4.153 |   1.76s |
