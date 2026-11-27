# Independent review of `math.archLog1p`

**Verdict: FAITHFUL**

Evidence SHA-256: `8e4813cc56421a63bc5b9d96b416025704639919ed83ed4021f1fe8633d1e467`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `077dfbf5d34e0f4f26db89faafa5c3397ff9c071fd3c66f742a5a6e9bd39869c`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `3c1e4e0e9e8c8081c15d3d858896df924dbb2411403aaccd9e1d66897007b1de`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `611a4456c9255bf2a12fe2e35b7152738e30cefd815ec0edee3e4bf2850fef29`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `2742002651a52b6d15f2e0b7b939dafa2c09b5a596e1237f0f9d5d61203fd865`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `a419acd25d0792a09578054e898cad4717641ad9989bf1a962506b041d3c26ed`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `809fbfb4948046a154182c7b9ca8a79fe81672ec927d0f2caf4419fca8153add`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `5339a1b7d146ab13df286b373ceecfe84f307cc701f4a197bf8a6bbcc647b952`.

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
