namespace GoLean.Compress.Flate

/-- The full backing suffix visible by reslicing to capacity. -/
structure HistogramResult where
  panicked : Bool
  backing : Array UInt16

private def bump (h : Array UInt16) (b : UInt8) : Array UInt16 :=
  h.modify b.toNat (· + 1)

/-- Go 1.27.0 src/compress/flate/huffman_code.go lines 395–417. -/
def histogramSplit (b : Array UInt8) (h : Array UInt16) : HistogramResult := Id.run do
  if h.size < 256 then
    return ⟨true, h⟩
  let remainder := b.size % 4
  let n := b.size / 4
  let mut bins := h
  for i in [:remainder] do
    bins := bump bins b[i]!
  for i in [:n] do
    bins := bump bins b[remainder + i]!
    bins := bump bins b[remainder + i + n]!
    bins := bump bins b[remainder + i + n + n]!
    bins := bump bins b[remainder + i + n + n + n]!
  return ⟨false, bins⟩

/-- Go 1.27.0 src/compress/flate/huffman_code.go lines 384–393. -/
def histogram (b : Array UInt8) (h : Array UInt16) : HistogramResult :=
  if b.size ≥ 8192 then histogramSplit b h
  else if h.size < 256 then ⟨true, h⟩
  else ⟨false, b.foldl bump h⟩

theorem histogramSplit_short_capacity (b : Array UInt8) (h : Array UInt16)
    (short : h.size < 256) :
    histogramSplit b h = ⟨true, h⟩ := by
  simp [histogramSplit, short]

theorem histogram_short_capacity (b : Array UInt8) (h : Array UInt16)
    (short : h.size < 256) :
    histogram b h = ⟨true, h⟩ := by
  simp [histogram, short, histogramSplit_short_capacity b h short]

end GoLean.Compress.Flate
