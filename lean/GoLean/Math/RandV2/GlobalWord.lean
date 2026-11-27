import GoLean.Math.RandV2.Word

namespace GoLean.Math.RandV2

def globalInt64 (state : RandState) := randInt64 (some state)

theorem globalInt64_delegates (state : RandState) :
    globalInt64 state = randInt64 (some state) := rfl

def globalUint32 (state : RandState) := randUint32 (some state)

theorem globalUint32_delegates (state : RandState) :
    globalUint32 state = randUint32 (some state) := rfl

def globalUint64 (state : RandState) := randUint64 (some state)

theorem globalUint64_delegates (state : RandState) :
    globalUint64 state = randUint64 (some state) := rfl

def globalInt32 (state : RandState) := randInt32 (some state)

theorem globalInt32_delegates (state : RandState) :
    globalInt32 state = randInt32 (some state) := rfl

def globalInt (width : Nat) (state : RandState) := randInt width (some state)

theorem globalInt_delegates (width : Nat) (state : RandState) :
    globalInt width state = randInt width (some state) := rfl

def globalUint (width : Nat) (state : RandState) := randUint width (some state)

theorem globalUint_delegates (width : Nat) (state : RandState) :
    globalUint width state = randUint width (some state) := rfl

def globalFloat64 (state : RandState) := randFloat64 (some state)

theorem globalFloat64_delegates (state : RandState) :
    globalFloat64 state = randFloat64 (some state) := rfl

def globalFloat32 (state : RandState) := randFloat32 (some state)

theorem globalFloat32_delegates (state : RandState) :
    globalFloat32 state = randFloat32 (some state) := rfl

end GoLean.Math.RandV2
