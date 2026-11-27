import GoLean.Archive.Zip.WriterIO
import GoDiff.Math.Bits.Tuple

open GoLean.Archive.Zip

def main : IO Unit := runTuple fun row => do
  let count :: written :: failed :: existing :: cells := row |
    throw (IO.userError "expected writer state and write fixture")
  let count := BitVec.ofInt 64 count
  let reply : WriteReply := ⟨written.toNat, failed == 1⟩
  let counted := countWriterWrite count reply
  let first := zipWriterFlush count cells.length false reply
  let second := zipWriterFlush count first.buffered first.failed reply
  let created := zipNewWriter (if existing > 0 then some (existing.toNat, 1) else none)
  let close := nopCloserClose 0
  return s!"{counted.count.toInt}:{counted.returned.count}:{counted.returned.failed}:{counted.calls}|{first.buffered}:{first.failed}:{first.calls}:{first.archiveCount.toInt}|{second.buffered}:{second.failed}:{second.calls}:{second.archiveCount.toInt}|{created.capacity}:{created.buffered}:{created.reused}:{created.archiveCount}:{created.calls}|{close.1}:{close.2}"
