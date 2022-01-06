# Multi-Party Computation Language

I heard about multi-party computation (MPC) sometime around 2018. The
idea sounded exciting, and I wanted to learn how it works. I like to
learn by doing so the natural next step was to write a simple [Garbled
Circuit](https://github.com/markkurossi/mpc) evaluator. The Garbled
Circuit [Wikipedia](https://en.wikipedia.org/wiki/Garbled_circuit)
article provides an excellent explanation about garbled circuits. I
followed it in my first PoC implementation in March 2019. The
essential parts were:

 - Oblivious transfer with RSA encryption.
 - AES encryption for circuit garbling.
 - Bristol-style binary circuits.

The PoC worked fine, but it used very trivial circuits (2-3 gates) or
some pre-compiled ones I found from the Internet. Since I wanted to
write different programs, I needed a higher-level language and
compiler for creating the binary circuits. The end goal was to compute
the RSA signature using MPC.

The compiler implementation was an educating journey. Garbled circuits
can't implement looping or branches based on any computed value. The
circuit evaluation goes through all gates, so the program must be a
directed acyclic graph (DAG). I have written several compilers over
the years, but this acyclic memoryless model was still different. One
of the learnings was Static Single Assignment (SSA) form assembly.

After one year of hobby coding, I implemented modular exponentiation
to any bit size integer numbers. And as soon as that, I also learned
that the algorithm would not work on numbers larger than 128 bits. So
what next? The Ed25519 algorithm is faster; would that work?

## Ed25519

I borrowed the
[Ed25519](https://github.com/markkurossi/mpc/tree/master/pkg/crypto/ed25519)
algorithm from Go's Ed25591 package. The MPCL language implements most
of the necessary Go language constructs, but the compiler still needs
more work. Compared to the original Go code, the current Ed25519
implementation has about ten lines or constructs that do not compile
without modifications. Nevertheless, it computes correct Ed25519
signatures. The first proper algorithm had the following statistics:

| Gate | Count     |
|------|----------:|
| XOR  | 616368261 |
| XNOR | 29253505  |
| AND  | 292577583 |
| OR   | 494216    |
| INV  | 19784     |

The signature computation transferred 26 GB of data and took 177
seconds. The garbler and evaluator evaluated 5,303,465 gates per
second.

So the Ed25519 was indeed much faster compared to the RSA
algorithm. And Ed25519 signatures are as
[strong](https://ed25519.cr.yp.to/index.html) as signatures with about
3000 bit RSA keys. I have benchmarked the RSA signature computations
up to 512 bits. Doubling the private key size makes the signature
computation time approximately seven times longer:

| Input | MODP |     Gates | Non-XOR  | Stream Gates | Stream !XOR | Stream   |
|------:|-----:|----------:|---------:|-------------:|------------:|---------:|
|     2 |    4 |       708 |      201 |          740 |         271 | 367.66ms |
|     4 |    8 |      5596 |     1571 |         5548 |        1719 | 115.51ms |
|     8 |   16 |     44796 |    12423 |        46252 |       13199 | 218.26ms |
|    16 |   32 |    374844 |   102255 |       364892 |      101535 | 245.80ms |
|    32 |   64 |   2986556 |   801887 |      2895932 |      788799 | 563.39ms |
|    64 |  128 |  23171068 |  6137023 |     22494524 |     6029311 |  2.4991s |
|   128 |  256 | 177580028 | 46495359 |    172945532 |    45732095 | 14.2368s |
|   256 |  512 |           |          |   1326461180 |   346797567 |  1m40.9s |
|   512 | 1024 |           |          |  10197960188 |  2641252351 | 13m3.86s |

2048 bit RSA signature computation would take 10 hours, compared to
equally strong Ed25519 signature, computed in 3 minutes.

## Optimizations Since the First Release

I made a few optimizations since the first public
[announcement](https://twitter.com/markkurossi/status/1436755119857623051)
last September:

 - Optimized `smov` and `srshift` instructions. When moving signed
   integer values to bigger variables, or bitwise-shifting values
   right, the instruction must expand the sign-bit to the most
   significant bits. The previous version had one and zero expansion
   versions of the `mov` and `rshift` instructions. The code
   generation used comparison and the `phi` instruction to use the
   proper versions. The new optimized implementation uses the sign bit
   as the expansion bit for signed integer operations, eliminating the
   comparison.
 - I removed memory move operations from the p2p communication
   protocol. The circuit evaluation does not allocate memory, and the
   sender and receiver do not copy memory in the transfer buffers when
   transferring data.
 - The [Half
   AND](https://en.wikipedia.org/wiki/Garbled_circuit#Half_And)
   optimization reduced the data needed to implement the AND
   gates. The reduction is 50%, which contributed significantly to the
   overall transfer data reduction. The signature computation now uses
   16GB compared to the original 25GB of data.
 - I also implemented the [Row
   Reduction](https://en.wikipedia.org/wiki/Garbled_circuit#Row_reduction)
   optimization for the OR and INV gates. However, since we prioritize
   XOR gates over OR gates, this optimization didn't significantly
   impact the performance.
 - Many programs utilize constant values as part of their
   evaluation. The Constant Propagation pushes these values forward in
   the circuit and removes unnecessary gates. This optimization can
   reduce some programs significantly. The current constant
   propagation algorithm works well; however, it only works with the
   known input wire values. A more generic optimization is to use [AND
   Inverter Graph](https://en.wikipedia.org/wiki/NAND_gate) (AIG),
   which would apply optimizations globally for the circuit. I'll
   implement the AIG optimization as a more comprehensive circuit
   optimization strategy soon.

These optimizations had a very positive impact on the system
performance. The Ed25519 signature computation time went down from 177
seconds to 90 seconds (49%). The network transfer decreased from 26 GB
to 15 GB (42%).

## Links and References

 - MPCL at [Github](https://github.com/markkurossi/mpc).
 - MPCL [documentation](https://www.markkurossi.com/mpcl/index.html)
   and API.
 - General info about garbled circuits at
   [Wikipedia](https://en.wikipedia.org/wiki/Garbled_circuit).

And finally, please join the discussion on
[Twitter](https://twitter.com/markkurossi/status/1479000199062294534).
