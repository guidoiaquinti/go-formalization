namespace GoLean.Encoding.Binary

def bigEndianString : String := "BigEndian"

theorem bigEndianString_value : bigEndianString = "BigEndian" := rfl

def bigEndianGoString : String := "binary.BigEndian"

theorem bigEndianGoString_value : bigEndianGoString = "binary.BigEndian" := rfl

def littleEndianString : String := "LittleEndian"

theorem littleEndianString_value : littleEndianString = "LittleEndian" := rfl

def littleEndianGoString : String := "binary.LittleEndian"

theorem littleEndianGoString_value : littleEndianGoString = "binary.LittleEndian" := rfl

def nativeEndianString : String := "NativeEndian"

theorem nativeEndianString_value : nativeEndianString = "NativeEndian" := rfl

def nativeEndianGoString : String := "binary.NativeEndian"

theorem nativeEndianGoString_value : nativeEndianGoString = "binary.NativeEndian" := rfl

end GoLean.Encoding.Binary
