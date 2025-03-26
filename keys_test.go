package art

import (
	"math"
	"testing"

	"golang.org/x/text/collate"
	"golang.org/x/text/language"
)

func TestAlphaKeys(t *testing.T) {
	var aok AlphabeticalOrderKey[string]

	tmp, _ := aok.Transform("hello")
	res := aok.Restore(tmp)

	if res != "hello" {
		t.Fatalf("expected 'hello', got %q", res)
	}
}

func TestCollateKeys(t *testing.T) {
	cok := CollationOrderKey[string]{
		c:   collate.New(language.Und),
		buf: &collate.Buffer{},
	}

	tmp, _ := cok.Transform("hello")
	res := cok.Restore(tmp)

	if res != "hello" {
		t.Fatalf("expected 'hello', got %q", res)
	}
}

func TestUnsignedKeysUint8(t *testing.T) {
	var ubk UnsignedBinaryKey[uint8]

	tmp, _ := ubk.Transform(1)
	res := ubk.Restore(tmp)

	if res != 1 {
		t.Fatalf("expected 1, got %d", res)
	}
}

func TestUnsignedKeysUint16(t *testing.T) {
	var ubk UnsignedBinaryKey[uint16]

	tmp, _ := ubk.Transform(1)
	res := ubk.Restore(tmp)

	if res != 1 {
		t.Fatalf("expected 1, got %d", res)
	}
}

func TestUnsignedKeysUint32(t *testing.T) {
	var ubk UnsignedBinaryKey[uint32]

	tmp, _ := ubk.Transform(1)
	res := ubk.Restore(tmp)

	if res != 1 {
		t.Fatalf("expected 1, got %d", res)
	}
}

func TestUnsignedKeysUint64(t *testing.T) {
	var ubk UnsignedBinaryKey[uint64]

	tmp, _ := ubk.Transform(1)
	res := ubk.Restore(tmp)

	if res != 1 {
		t.Fatalf("expected 1, got %d", res)
	}
}

func TestUnsignedKeysUint(t *testing.T) {
	var ubk UnsignedBinaryKey[uint]

	tmp, _ := ubk.Transform(1)
	res := ubk.Restore(tmp)

	if res != 1 {
		t.Fatalf("expected 1, got %d", res)
	}
}

func TestSignedKeysInt8(t *testing.T) {
	var sbk SignedBinaryKey[int8]

	tmp, _ := sbk.Transform(-1)
	res := sbk.Restore(tmp)

	if res != -1 {
		t.Fatalf("expected -1, got %d", res)
	}
}

func TestSignedKeysInt16(t *testing.T) {
	var sbk SignedBinaryKey[int16]

	tmp, _ := sbk.Transform(-1)
	res := sbk.Restore(tmp)

	if res != -1 {
		t.Fatalf("expected -1, got %d", res)
	}
}

func TestSignedKeysInt32(t *testing.T) {
	var sbk SignedBinaryKey[int32]

	tmp, _ := sbk.Transform(-1)
	res := sbk.Restore(tmp)

	if res != -1 {
		t.Fatalf("expected -1, got %d", res)
	}
}

func TestSignedKeysInt64(t *testing.T) {
	var sbk SignedBinaryKey[int64]

	tmp, _ := sbk.Transform(-1)
	res := sbk.Restore(tmp)

	if res != -1 {
		t.Fatalf("expected -1, got %d", res)
	}
}

func TestSignedKeysInt(t *testing.T) {
	var sbk SignedBinaryKey[int]

	tmp, _ := sbk.Transform(-1)
	res := sbk.Restore(tmp)

	if res != -1 {
		t.Fatalf("expected -1, got %d", res)
	}
}

func TestFloatKeysFloat32(t *testing.T) {
	var sbk FloatBinaryKey[float32]

	tmp, _ := sbk.Transform(-1.5)
	res := sbk.Restore(tmp)

	if res != -1.5 {
		t.Fatalf("expected -1.5, got %f", res)
	}
}

func TestFloatKeysFloat64(t *testing.T) {
	var sbk FloatBinaryKey[float64]

	tmp, _ := sbk.Transform(-1.5)
	res := sbk.Restore(tmp)

	if res != -1.5 {
		t.Fatalf("expected -1.5, got %f", res)
	}
}

func TestFloatKeysFloatNaNFloat32(t *testing.T) {
	var sbk FloatBinaryKey[float32]

	tmp, _ := sbk.Transform(float32(math.NaN()))
	res := sbk.Restore(tmp)

	if !math.IsNaN(float64(res)) {
		t.Fatalf("expected NaN, got %f", res)
	}
}

func TestFloatKeysFloatNaNFloat64(t *testing.T) {
	var sbk FloatBinaryKey[float64]

	tmp, _ := sbk.Transform(math.NaN())
	res := sbk.Restore(tmp)

	if !math.IsNaN(res) {
		t.Fatalf("expected NaN, got %f", res)
	}
}

func TestFloatKeysFloatInfFloat32(t *testing.T) {
	var sbk FloatBinaryKey[float32]

	tmp, _ := sbk.Transform(float32(math.Inf(-1)))
	res := sbk.Restore(tmp)

	if !math.IsInf(float64(res), -1) {
		t.Fatalf("expected -Inf, got %f", res)
	}

	tmp, _ = sbk.Transform(float32(math.Inf(1)))
	res = sbk.Restore(tmp)

	if !math.IsInf(float64(res), 1) {
		t.Fatalf("expected Inf, got %f", res)
	}
}

func TestFloatKeysFloatInfFloat64(t *testing.T) {
	var sbk FloatBinaryKey[float64]

	tmp, _ := sbk.Transform(math.Inf(-1))
	res := sbk.Restore(tmp)

	if !math.IsInf(res, -1) {
		t.Fatalf("expected -Inf, got %f", res)
	}

	tmp, _ = sbk.Transform(math.Inf(1))
	res = sbk.Restore(tmp)

	if !math.IsInf(res, 1) {
		t.Fatalf("expected Inf, got %f", res)
	}
}

func TestFloatKeysFloatZeroFloat32(t *testing.T) {
	var sbk FloatBinaryKey[float32]

	tmp, _ := sbk.Transform(0)
	res := sbk.Restore(tmp)

	if res != 0 {
		t.Fatalf("expected 0, got %f", res)
	}
}

func TestFloatKeysFloatZeroFloat64(t *testing.T) {
	var sbk FloatBinaryKey[float64]

	tmp, _ := sbk.Transform(0)
	res := sbk.Restore(tmp)

	if res != 0 {
		t.Fatalf("expected 0, got %f", res)
	}
}
