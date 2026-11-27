# Independent review of `math.archAsin`

**Verdict: FAITHFUL**

Evidence SHA-256: `50c4c44a57a9174841462de43da4c8de372194e220c5d6339549ef03dc3dbf51`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `da9ad0d515ceb3e86a8fdae4cf2da9ccd08035a2e2f80c0d448f812fd7715751`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `07c57d44f28347f0da4badb76c28eaf9a6eb57279cd16dad7b73a7505c2b4f8f`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `f71777730e05faae7c1c5c50ddc458ec243e33372b1acaea537655a6dabad9b5`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `e6fcf573e1a27bf8cecb273743931e52dafa206fa29860d7f52cbabcb66d693a`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `c3d9d2386932817b2171849cdaf4e23de462cd29e4e7b3e1dfcade7d408d6464`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `e4471190a11532c3b70bdaf5c55489221df808b937ba933c6d45e585d802c620`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `6215cafc821e2206ddd3482a179795b0ed5d0adb53fc0949fb68346ebf152fe8`.

I checked the pinned Go 1.27.0 declaration, its availability constant and
the non-s390x build constraint. The function body unconditionally panics
with the Go string `not implemented`. The shared independent control
parses the actual pinned source and checks five operand rows plus the
excluded s390x profile for each of the 25 declarations.

The author invoked the exact private symbol with a same-package Go
bridge and checked the panic value and availability flag on darwin/arm64.
Each declaration received operand boundary cases and 60,000 seeded
cases using 17, 2718 and 65537. All 1,503,048 comparisons passed across
the batch. Native non-arm64 ports and the s390x implementations remain
outside the reviewed observations.
