# Ed25591 Signatures Under One Minute

The MPCL can now compute Ed25519 signatures in about 90 seconds. The
next goal is to optimize the signature computation (and key
generation) under 60 seconds. The two main optimizations are:

- Parallel garbling
- Circuit optimization with AND Invert Graph (AIG)
- [Three Halves Make a Whole?](https://eprint.iacr.org/2021/749)
  Beating the Half-Gates Lower Bound for Garbled Circuits

## Parallel Garbling

```
┌─────────┬─────┬─────────┬────────┬─────────┬──────┬─────┬─────────┬─────┬────┐
│Instr    │Count│      XOR│    XNOR│      AND│    OR│  INV│     !XOR│    L│   W│
├─────────┼─────┼─────────┼────────┼─────────┼──────┼─────┼─────────┼─────┼────┤
│imult/64 │63390│494970948│26561052│252766410│     0│    0│252766410│  218│ 917│
│iadd/64  │73254│ 18098114│       0│  4615002│     0│    0│  4615002│  189│   2│
│umult/32 │ 4270│  8437520│  273280│  4496310│     0│    0│  4496310│  121│ 397│
│umult/64 │ 1350│  7695000│  481950│  3901500│     0│    0│  3901500│  218│ 522│
│bor/1024 │  467│        0│       0│        0│478208│    0│   478208│    1│1024│
│isub/64  │ 9173│  1146625│ 1174144│   577899│     0│    0│   577899│  191│   2│
│band/32  │17475│        0│       0│   557592│     0│    0│   557592│    1│  32│
│circ/1024│    4│  1146896│       0│   231788│     0│19784│   251572│10548│1796│
│isub/32  │ 4556│   277916│  291584│   141236│     0│    0│   141236│   95│   2│
│band/512 │  192│      192│       0│    98304│     0│    0│    98304│    2│ 480│
│iadd/32  │ 2745│   331587│       0│    83583│     0│    0│    83583│   93│   2│
│bor/64   │  242│        0│       0│        0│ 15488│    0│    15488│    1│  64│
│usub/32  │  512│    31232│   32768│    15872│     0│    0│    15872│   95│   2│
│band/64  │   79│       79│       0│     5056│     0│    0│     5056│    2│  33│
│bor/32   │   17│        0│       0│        0│   520│    0│      520│    1│  32│
│iadd/8   │   64│     1728│       0│      448│     0│    0│      448│   21│   2│
│isub/8   │   63│      819│    1008│      441│     0│    0│      441│   23│   2│
│bxor/32  │35072│  1122304│       0│        0│     0│    0│        0│    1│  32│
│bxor/8   │   65│      520│       0│        0│     0│    0│        0│    1│   8│
└─────────┴─────┴─────────┴────────┴─────────┴──────┴─────┴─────────┴─────┴────┘
```

## AND Invert Graph (AIG)

[Implementing other gates with NAND gates](https://en.wikipedia.org/wiki/NAND_logic)

Q = A NAND B

| A | B | Q |
|---|---|---|
| 0 | 0 | 1 |
| 0 | 1 | 1 |
| 1 | 0 | 1 |
| 1 | 1 | 0 |


### NOT

Q = A NAND A

| A | A | Q |
|---|---|---|
| 0 | 0 | 1 |
| 1 | 1 | 0 |

### AND

Q = A AND B	= ( A NAND B ) NAND ( A NAND B )

## Three Halves Make a Whole?
