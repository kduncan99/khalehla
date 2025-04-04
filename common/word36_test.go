// khalehla Project
// Copyright © 2023-2025 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package common

import (
	"testing"
)

func checkEquals(t *testing.T, expected uint64, actual uint64, msg string) {
	if expected != actual {
		t.Errorf("%s expected:%012o, actual:%012o", msg, expected, actual)
	}
}

func TestGetPartialWord(t *testing.T) {
	checkEquals(t, GetW(0_111_123321_456654), 0_123321_456654, "GetW() failed")
	checkEquals(t, GetH1(0_101023_123123), 0_101023, "GetH1() failed")
	checkEquals(t, GetH1(0_777777_123123), 0_777777, "GetH1() failed")
	checkEquals(t, GetH2(0_321364_101023), 0_101023, "GetH1() failed")
	checkEquals(t, GetH2(0_123123_777777), 0_777777, "GetH1() failed")
	checkEquals(t, GetQ1(0_112233_445566), 0112, "GetQ1() failed")
	checkEquals(t, GetQ2(0_112233_445566), 0233, "GetQ2() failed")
	checkEquals(t, GetQ3(0_112233_445566), 0445, "GetQ3() failed")
	checkEquals(t, GetQ4(0_112233_445566), 0566, "GetQ4() failed")
	checkEquals(t, GetS1(0_112233_445566), 011, "GetS1() failed")
	checkEquals(t, GetS2(0_112233_445566), 022, "GetS2() failed")
	checkEquals(t, GetS3(0_112233_445566), 033, "GetS3() failed")
	checkEquals(t, GetS4(0_112233_445566), 044, "GetS4() failed")
	checkEquals(t, GetS5(0_112233_445566), 055, "GetS5() failed")
	checkEquals(t, GetS6(0_112233_445566), 066, "GetS6() failed")
	checkEquals(t, GetT1(0_112233_445566), 01122, "GetT1() failed")
	checkEquals(t, GetT2(0_112233_445566), 03344, "GetT2() failed")
	checkEquals(t, GetT3(0_112233_445566), 05566, "GetT3() failed")
	checkEquals(t, GetXH1(0_101023_123123), 0_101023, "GetXH1() failed")
	checkEquals(t, GetXH1(0_401023_123123), 0_777777_401023, "GetXH1() failed")
	checkEquals(t, GetXH2(0_321364_371023), 0_371023, "GetXH2() failed")
	checkEquals(t, GetXH2(0_321364_400000), 0_777777_400000, "GetXH2() failed")
	checkEquals(t, GetXT1(0_112233_445566), 01122, "GetXT1() failed")
	checkEquals(t, GetXT1(0_412233_445566), 0777777_774122, "GetXT1() failed")
	checkEquals(t, GetXT2(0_112233_445566), 03344, "GetXT2() failed")
	checkEquals(t, GetXT2(0_112273_445566), 0_777777_777344, "GetXT2() failed")
	checkEquals(t, GetXT3(0_112233_443777), 03777, "GetXT3() failed")
	checkEquals(t, GetXT3(0_112233_444777), 0_777777_774777, "GetXT3() failed")
}

func TestSetPartialWord(t *testing.T) {
	checkEquals(t, SetH1(0, 0765432_331455), 0_331455_000000, "SetH1() failed")
	checkEquals(t, SetH1(0_777777_777777, 0331455), 0_331455_777777, "SetH1() failed")
	checkEquals(t, SetH2(0, 0765432_331455), 0_000000_331455, "SetH2() failed")
	checkEquals(t, SetH2(0_777777_777777, 0331455), 0_777777_331455, "SetH2() failed")
	// TODO more
}

func TestCountBits(t *testing.T) {
	type parameterSet struct {
		input          uint64
		expectedResult uint64
	}

	parameterSets := []parameterSet{
		{0, 0},
		{1, 1},
		{2, 1},
		{0_707070_707070, 18},
		{0_070707_070707, 18},
		{0_525252_525252, 18},
		{0_252525_252525, 18},
		{0_000000_777777, 18},
		{0_777777_777777, 36},
	}

	for _, parameterSet := range parameterSets {
		result := CountBits(parameterSet.input)
		if result != parameterSet.expectedResult {
			t.Errorf("CountBits(%012o) == %012o, should be %012o", parameterSet.input, result, parameterSet.expectedResult)
		}
	}
}

func TestEliminateNegativeZero(t *testing.T) {
	type parameterSet struct {
		input          uint64
		expectedResult uint64
	}

	parameterSets := []parameterSet{
		{0, 0},
		{1, 1},
		{2, 2},
		{0_707070_707070, 0_707070_707070},
		{0_777777_777776, 0_777777_777776},
		{0_777777_777777, 0},
	}

	for _, parameterSet := range parameterSets {
		result := EliminateNegativeZero(parameterSet.input)
		if result != parameterSet.expectedResult {
			t.Errorf("EliminateNegativeZero(%012o) == %012o, should be %012o", parameterSet.input, result, parameterSet.expectedResult)
		}
	}
}

func TestIsNegative(t *testing.T) {
	type parameterSet struct {
		input          uint64
		expectedResult bool
	}

	parameterSets := []parameterSet{
		{0, false},
		{1, false},
		{0_377777_777777, false},
		{0_400000_000000, true},
		{0_777777_777777, true},
	}

	for _, parameterSet := range parameterSets {
		result := IsNegative(parameterSet.input)
		if result != parameterSet.expectedResult {
			t.Errorf("IsNegative(%012o) == %v, should be %v", parameterSet.input, result, parameterSet.expectedResult)
		}
	}
}

func TestIsPositive(t *testing.T) {
	type parameterSet struct {
		input          uint64
		expectedResult bool
	}

	parameterSets := []parameterSet{
		{0, true},
		{1, true},
		{0_377777_777777, true},
		{0_400000_000000, false},
		{0_777777_777777, false},
	}

	for _, parameterSet := range parameterSets {
		result := IsPositive(parameterSet.input)
		if result != parameterSet.expectedResult {
			t.Errorf("IsPositive(%012o) == %v, should be %v", parameterSet.input, result, parameterSet.expectedResult)
		}
	}
}

func TestIsZero(t *testing.T) {
	type parameterSet struct {
		input          uint64
		expectedResult bool
	}

	parameterSets := []parameterSet{
		{0, true},
		{0_777777_777777, true},
		{0_707070_707071, false},
		{1, false},
		{0_377777_777777, false},
		{0_400000_000000, false},
	}

	for _, parameterSet := range parameterSets {
		result := IsZero(parameterSet.input)
		if result != parameterSet.expectedResult {
			t.Errorf("IsZero(%012o) == %v, should be %v", parameterSet.input, result, parameterSet.expectedResult)
		}
	}
}

/*
    //  Arithmetic -----------------------------------------------------------------------------------------------------------------

    @Test
    public void addPosPos() {
        Word36 w1 = new Word36(25);
        Word36 w2 = new Word36(1027);
        Word36.AdditionResult ar = w1.add(w2);
        assertEquals(1052, ar._result.getTwosComplement());
        assertFalse(ar._flags._carry);
        assertFalse(ar._flags._overflow);
    }

    @Test
    public void addPosPosOverflow() {
        Word36 w1 = new Word36(0_377777_777777L);
        Word36 w2 = new Word36(1);
        Word36.AdditionResult ar = w1.add(w2);
        assertFalse(ar._flags._carry);
        assertTrue(ar._flags._overflow);
    }

    @Test
    public void addPosNegResultPos() {
        Word36 w1 = new Word36(Word36.getOnesComplement(1234));
        Word36 w2 = new Word36(Word36.getOnesComplement(-234));
        Word36.AdditionResult ar = w1.add(w2);
        assertEquals(1000, ar._result.getTwosComplement());
        assertTrue(ar._flags._carry);
        assertFalse(ar._flags._overflow);
    }

    @Test
    public void addPosNegResultNeg() {
        Word36 w1 = new Word36(Word36.getOnesComplement(234));
        Word36 w2 = new Word36(Word36.getOnesComplement(-1234));
        Word36.AdditionResult ar = w1.add(w2);
        assertEquals(-1000, ar._result.getTwosComplement());
        assertTrue(ar._flags._carry);
        assertFalse(ar._flags._overflow);
    }

    @Test
    public void addNegNeg() {
        Word36 w1 = new Word36(Word36.getOnesComplement(-1992));
        Word36 w2 = new Word36(Word36.getOnesComplement(-2933));
        Word36.AdditionResult ar = w1.add(w2);
        assertEquals(-1992-2933, ar._result.getTwosComplement());
        assertTrue(ar._flags._carry);
        assertFalse(ar._flags._overflow);
    }

    @Test
    public void addNegNegOverflow() {
        Word36 w1 = new Word36(0_400000_000000L);
        Word36 w2 = new Word36(0_777777_777776L);
        Word36.AdditionResult ar = w1.add(w2);
        assertTrue(ar._flags._carry);
        assertTrue(ar._flags._overflow);
    }

    @Test
    public void addPosZPosZ() {
        Word36 w1 = Word36.W36_POSITIVE_ZERO;
        Word36 w2 = Word36.W36_POSITIVE_ZERO;
        Word36.AdditionResult ar = w1.add(w2);
        assertEquals(Word36.POSITIVE_ZERO, ar._result.getW());
        assertFalse(ar._flags._carry);
        assertFalse(ar._flags._overflow);
    }

    @Test
    public void addPosZNegZ() {
        Word36 w1 = Word36.W36_POSITIVE_ZERO;
        Word36 w2 = Word36.W36_NEGATIVE_ZERO;
        Word36.AdditionResult ar = w1.add(w2);
        assertEquals(Word36.POSITIVE_ZERO, ar._result.getW());
        assertTrue(ar._flags._carry);
        assertFalse(ar._flags._overflow);
    }

    @Test
    public void addNegZPosZ() {
        Word36 w1 = Word36.W36_NEGATIVE_ZERO;
        Word36 w2 = Word36.W36_POSITIVE_ZERO;
        Word36.AdditionResult ar = w1.add(w2);
        assertEquals(Word36.POSITIVE_ZERO, ar._result.getW());
        assertTrue(ar._flags._carry);
        assertFalse(ar._flags._overflow);
    }

    @Test
    public void addNegZNegZ() {
        Word36 w1 = Word36.W36_NEGATIVE_ZERO;
        Word36 w2 = Word36.W36_NEGATIVE_ZERO;
        Word36.AdditionResult ar = w1.add(w2);
        assertEquals(Word36.NEGATIVE_ZERO, ar._result.getW());
        assertTrue(ar._flags._carry);
        assertFalse(ar._flags._overflow);
    }

    @Test
    public void addInverses() {
        Word36 w1 = new Word36(Word36.getOnesComplement(19883));
        Word36 w2 = new Word36(Word36.getOnesComplement(-19883));
        Word36.AdditionResult ar = w1.add(w2);
        assertEquals(Word36.POSITIVE_ZERO, ar._result.getW());
        assertTrue(ar._flags._carry);
        assertFalse(ar._flags._overflow);
    }

    @Test
    public void multiply_1() {
        long factor1 = 0_003234_715364L;
        long factor2 = 0_073654_717623L;
        BigInteger expProduct = BigInteger.valueOf(factor1).multiply(BigInteger.valueOf(factor2));

        Word36 w36Factor1 = new Word36(Word36.getOnesComplement(factor1));
        Word36 w36Factor2 = new Word36(Word36.getOnesComplement(factor2));
        DoubleWord36 product = w36Factor1.multiply(w36Factor2);
        BigInteger biProduct = product.getTwosComplement();

        assertEquals(expProduct, biProduct);
    }

    @Test
    public void multiply_2() {
        long factor1 = -29937;
        long factor2 = 0_073654_717623L;
        BigInteger expProduct = BigInteger.valueOf(factor1).multiply(BigInteger.valueOf(factor2));

        Word36 w36Factor1 = new Word36(Word36.getOnesComplement(factor1));
        Word36 w36Factor2 = new Word36(Word36.getOnesComplement(factor2));
        DoubleWord36 product = w36Factor1.multiply(w36Factor2);
        BigInteger biProduct = product.getTwosComplement();

        assertEquals(expProduct, biProduct);
    }

    @Test
    public void negate_PositiveOne() {
        Word36 word36 = Word36.W36_POSITIVE_ONE;
        Word36 expectedResult = word36.negate();
        assertEquals(Word36.NEGATIVE_ONE, expectedResult.getW());
    }

    @Test
    public void negate_PositiveZero() {
        Word36 word36 = Word36.W36_POSITIVE_ZERO;
        Word36 expectedResult = word36.negate();
        assertEquals(Word36.NEGATIVE_ZERO, expectedResult.getW());
    }

    @Test
    public void negate_NegativeOne() {
        Word36 word36 = Word36.W36_NEGATIVE_ONE;
        Word36 expectedResult = word36.negate();
        assertEquals(Word36.POSITIVE_ONE, expectedResult.getW());
    }

    @Test
    public void negate_NegativeZero() {
        Word36 word36 = Word36.W36_NEGATIVE_ZERO;
        Word36 expectedResult = word36.negate();
        assertEquals(Word36.POSITIVE_ZERO, expectedResult.getW());
    }


    //  Shifts ---------------------------------------------------------------------------------------------------------------------

    //TODO a few more leftShiftAlgebraic tests

    @Test
    public void leftShiftAlgebraic() {
        //  sign bit always remains unchanged...
        long parameter = 0_3123_4537_0123L;
        long expected =  0_2247_1276_0246L;
        Word36 word36 = new Word36(parameter);
        Word36 expectedResult = word36.leftShiftAlgebraic(1);
        assertEquals(expected, expectedResult.getW());
    }

    @Test
    public void leftShiftCircular_by0() {
        long parameter = 0_111222_333444L;
        long expected = 0_111222_333444L;
        Word36 word36 = new Word36(parameter);
        Word36 expectedResult = word36.leftShiftCircular(0);
        assertEquals(expected, expectedResult.getW());
    }

    @Test
    public void leftShiftCircular_by3() {
        long parameter = 0_111222_333444L;
        long expected = 0_112223_334441L;
        Word36 word36 = new Word36(parameter);
        Word36 expectedResult = word36.leftShiftCircular(3);
        assertEquals(expected, expectedResult.getW());
    }

    @Test
    public void leftShiftCircular_by36() {
        long parameter = 0_111222_333444L;
        long expected = 0_111222_333444L;
        Word36 word36 = new Word36(parameter);
        Word36 expectedResult = word36.leftShiftCircular(36);
        assertEquals(expected, expectedResult.getW());
    }

    @Test
    public void leftShiftCircular_byNeg() {
        long parameter = 0_111222_333444L;
        long expected = 0_441112_223334L;
        Word36 word36 = new Word36(parameter);
        Word36 expectedResult = word36.leftShiftCircular(-6);
        assertEquals(expected, expectedResult.getW());
    }

    @Test
    public void leftShiftLogical_by3() {
        long parameter = 0_111222_333444L;
        long expected = 0_112223_334440L;
        Word36 word36 = new Word36(parameter);
        Word36 expectedResult = word36.leftShiftLogical(3);
        assertEquals(expected, expectedResult.getW());
    }

    @Test
    public void leftShiftLogical_by36() {
        long parameter = 0_111222_333444L;
        long expected = 0;
        Word36 word36 = new Word36(parameter);
        Word36 expectedResult = word36.leftShiftLogical(36);
        assertEquals(expected, expectedResult.getW());
    }

    @Test
    public void leftShiftLogical_negCount() {
        long parameter = 0_111222_333444L;
        long expected = 0_001112_223334L;
        Word36 word36 = new Word36(parameter);
        Word36 expectedResult = word36.leftShiftLogical(-6);
        assertEquals(expected, expectedResult.getW());
    }

    @Test
    public void leftShiftLogical_zeroCount() {
        long parameter = 0_111222_333444L;
        long expected = 0_111222_333444L;
        Word36 word36 = new Word36(parameter);
        Word36 expectedResult = word36.leftShiftLogical(0);
        assertEquals(expected, expectedResult.getW());
    }

    @Test
    public void rightShiftAlgebraic_negCount() {
        long parameter = 033225L;
        long expResult = 0332250L;
        Word36 word36 = new Word36(parameter);
        Word36 expectedResult = word36.rightShiftAlgebraic(-3);
        assertEquals(expResult, expectedResult.getW());
    }

    @Test
    public void rightShiftAlgebraic_neg_3Count() {
        long parameter = 0_400000_112233L;
        long expResult = 0_740000_011223L;
        Word36 word36 = new Word36(parameter);
        Word36 expectedResult = word36.rightShiftAlgebraic(3);
        assertEquals(expResult, expectedResult.getW());
    }

    @Test
    public void rightShiftAlgebraic_neg_34Count() {
        long parameter = 0_421456_321456L;
        long expResult = 0_777777_777742L;
        Word36 word36 = new Word36(parameter);
        Word36 expectedResult = word36.rightShiftAlgebraic(30);
        assertEquals(expResult, expectedResult.getW());
    }

    @Test
    public void rightShiftAlgebraic_neg_minus18Count() {
        long parameter = 0_423232_123123L;
        long expResult = 0_523123_000000L;
        Word36 word36 = new Word36(parameter);
        Word36 expectedResult = word36.rightShiftAlgebraic(-18);
        assertEquals(expResult, expectedResult.getW());
    }

    @Test
    public void rightShiftAlgebraic_neg_36Count() {
        long parameter = 0_421456_321456L;
        long expResult = 0_777777_777777L;
        Word36 word36 = new Word36(parameter);
        Word36 expectedResult = word36.rightShiftAlgebraic(36);
        assertEquals(expResult, expectedResult.getW());
    }

    @Test
    public void rightShiftAlgebraic_pos_3Count() {
        long parameter = 033225L;
        long expResult = parameter >> 3;
        Word36 word36 = new Word36(parameter);
        Word36 expectedResult = word36.rightShiftAlgebraic(3);
        assertEquals(expResult, expectedResult.getW());
    }

    @Test
    public void rightShiftAlgebraic_pos_34Count() {
        long parameter = 0_321456_321456L;
        long expResult = parameter >> 34;
        Word36 word36 = new Word36(parameter);
        Word36 expectedResult = word36.rightShiftAlgebraic(34);
        assertEquals(expResult, expectedResult.getW());
    }

    @Test
    public void rightShiftAlgebraic_pos_35Count() {
        long parameter = 0_321456_321456L;
        long expResult = parameter >> 35;
        Word36 word36 = new Word36(parameter);
        Word36 expectedResult = word36.rightShiftAlgebraic(35);
        assertEquals(expResult, expectedResult.getW());
    }

    @Test
    public void rightShiftAlgebraic_pos_36Count() {
        long parameter = 0_321456_321456L;
        long expResult = 0;
        Word36 word36 = new Word36(parameter);
        Word36 expectedResult = word36.rightShiftAlgebraic(36);
        assertEquals(expResult, expectedResult.getW());
    }

    @Test
    public void rightShiftAlgebraic_zeroCount() {
        long parameter = 033225L;
        long expResult = 033225L;
        Word36 word36 = new Word36(parameter);
        Word36 expectedResult = word36.rightShiftAlgebraic(0);
        assertEquals(expResult, expectedResult.getW());
    }

    @Test
    public void rightShiftCircular_1() {
        Word36 word36 = new Word36(0_112233_445566L);
        Word36 expectedResult = word36.rightShiftCircular(6);
        assertEquals(0_661122_334455L, expectedResult.getW());
    }

    @Test
    public void rightShiftCircular_2() {
        Word36 word36 = new Word36(0_112200_334400L);
        Word36 expectedResult = word36.rightShiftCircular(3);
        assertEquals(0_011220_033440L, expectedResult.getW());
    }

    @Test
    public void rightShiftLogical() {
        Word36 word36 = new Word36(0_112233_445566L);
        Word36 expectedResult = word36.rightShiftLogical(9);
        assertEquals(0_000112_233445L, expectedResult.getW());
    }


    //  Logic tests ----------------------------------------------------------------------------------------------------------------

    @Test
    public void and() {
        Word36 op1 = new Word36(0_776655_221100L);
        Word36 op2 = new Word36(0_765432_543210L);
        Word36 exp = new Word36(0_764410_001000L);
        Word36 expectedResult = op1.logicalAnd(op2);
        assertEquals(exp, expectedResult);
    }

    @Test
    public void not() {
        Word36 op1 = new Word36(0_776655_221100L);
        Word36 exp = new Word36(0_001122_556677L);
        Word36 expectedResult = op1.logicalNot();
        assertEquals(exp, expectedResult);
    }

    @Test
    public void or() {
        Word36 op1 = new Word36(0_776655_221100L);
        Word36 op2 = new Word36(0_765432_543210L);
        Word36 exp = new Word36(0_777677_763310L);
        Word36 expectedResult = op1.logicalOr(op2);
        assertEquals(exp, expectedResult);
    }

    @Test
    public void xor() {
        Word36 op1 = new Word36(0_776655_221100L);
        Word36 op2 = new Word36(0_765432_543210L);
        Word36 exp = new Word36(0_013267_762310L);
        Word36 expectedResult = op1.logicalXor(op2);
        assertEquals(exp, expectedResult);
    }


    //  Display --------------------------------------------------------------------------------------------------------------------

    @Test
    public void toASCII() {
        long word = 0_101_102_103_104L;
        assertEquals("ABCD", Word36.toStringFromASCII(word));
    }

    @Test
    public void toFieldata() {
        long word = 0_05_06_07_10_11_12L;
        assertEquals(" ABCDE", Word36.toStringFromFieldata(word));
    }

    @Test
    public void toOctal() {
        long word = 0_05_06_07_10_11_12L;
        assertEquals("050607101112", Word36.toOctal(word));
    }


    //  Misc -----------------------------------------------------------------------------------------------------------------------

    @Test
    public void stringToWordASCII() {
        var w = Word36.stringToWordASCII("Help");
        assertEquals(0_110_145_154_160L, w);
    }

    @Test
    public void stringToWordASCII_over() {
        var w = Word36.stringToWordASCII("HelpSlop");
        assertEquals(0_110_145_154_160L, w);
    }

    @Test
    public void stringToWordASCII_partial() {
        var w = Word36.stringToWordASCII("01");
        assertEquals(0_060_061_040_040L, w);
    }

    @Test
    public void stringToWordFieldata() {
        var w = Word36.stringToWordFieldata("Abc@23");
        assertEquals(0_060710_006263L, w);
    }

    @Test
    public void stringToWordFieldata_over() {
        var w = Word36.stringToWordFieldata("A B C@D E F");
        assertEquals(0_060507_051000L, w);
    }

    @Test
    public void stringToWordFieldata_partial() {
        var w = Word36.stringToWordFieldata("1234");
        assertEquals(0_616263_640505L, w);
    }


    //  Sign-extension tests -------------------------------------------------------------------------------------------------------

    @Test
    public void getSignExtended12_positive() {
        assertEquals(03765, Word36.getSignExtended12(03765));
    }

    @Test
    public void getSignExtended12_negative() {
        assertEquals(0_777777_774765L, Word36.getSignExtended12(04765));
    }

    @Test
    public void getSignExtended18_positive() {
        assertEquals(0_376500, Word36.getSignExtended18(0_376500));
    }

    @Test
    public void getSignExtended18_negative() {
        assertEquals(0_777777_400001L, Word36.getSignExtended18(0_400001));
    }

    @Test
    public void getSignExtended24_positive() {
        assertEquals(0_000037_776500L, Word36.getSignExtended24(0_000037_776500L));
    }

    @Test
    public void getSignExtended24_negative() {
        assertEquals(0_777767_776500L, Word36.getSignExtended24(0_000067_776500L));
    }
}
*/

func Test_AddDouble(t *testing.T) {
	result := AddDouble(DoublePositiveZero, DoublePositiveZero)
	if CompareDouble(result, DoublePositiveZero) != 0 {
		t.Errorf("Error expected result to be 0:0, but it was %012o:%012o", result[0], result[1])
	}

	result = AddDouble(DoubleNegativeZero, DoubleNegativeZero)
	if CompareDouble(result, DoubleNegativeZero) != 0 {
		t.Errorf("Error expected result to be 777777777777:777777777777, but it was %012o:%012o", result[0], result[1])
	}

	result = AddDouble(DoubleNegativeZero, DoublePositiveZero)
	if CompareDouble(result, DoublePositiveZero) != 0 {
		t.Errorf("Error expected result to be 0:0, but it was %012o:%012o", result[0], result[1])
	}

	addend1 := []uint64{0, 0_210335_732001}
	addend2 := []uint64{0, 0_104772_100001}
	expected := []uint64{0, addend1[1] + addend2[1]}
	result = AddDouble(addend1, addend2)
	if CompareDouble(result, expected) != 0 {
		t.Errorf("Error expected result to be %012o:%012o, but it was %012o:%012o",
			expected[0], expected[1], result[0], result[1])
	}

	addend1 = []uint64{0_777777_777777, 0_777777_777777}
	addend2 = []uint64{0, 0_104772_100001}
	expected = addend2
	result = AddDouble(addend1, addend2)
	if CompareDouble(result, expected) != 0 {
		t.Errorf("Error expected result to be %012o:%012o, but it was %012o:%012o",
			expected[0], expected[1], result[0], result[1])
	}

	addend1 = []uint64{0_777777_777777, 0_777777_777774}
	addend2 = []uint64{0, 0_104772_100017}
	expected = []uint64{0, 0_104772_100014}
	result = AddDouble(addend1, addend2)
	if CompareDouble(result, expected) != 0 {
		t.Errorf("Error expected result to be %012o:%012o, but it was %012o:%012o",
			expected[0], expected[1], result[0], result[1])
	}

	addend1 = []uint64{0_000000_000000, 0_777777_777777}
	addend2 = []uint64{0, 0_104772_100017}
	expected = []uint64{01, 0_104772_100016}
	result = AddDouble(addend1, addend2)
	if CompareDouble(result, expected) != 0 {
		t.Errorf("Error expected result to be %012o:%012o, but it was %012o:%012o",
			expected[0], expected[1], result[0], result[1])
	}

	addend1 = []uint64{0_000000_543210, 0_210056_523004}
	addend2 = []uint64{0_777777_347677, 0_777735_667775}
	expected = []uint64{0_000000_113110, 0_210014_413002}
	result = AddDouble(addend1, addend2)
	if CompareDouble(result, expected) != 0 {
		t.Errorf("Error expected result to be %012o:%012o, but it was %012o:%012o",
			expected[0], expected[1], result[0], result[1])
	}
}

func Test_AddSimple_2(t *testing.T) {
	value1 := uint64(0_300000_000000)
	value2 := uint64(0_077777_777777)
	expected := uint64(0_377777_777777)
	result := AddSimple(value1, value2)
	if result != expected {
		t.Errorf("Error expected %12o, got %12o", expected, result)
	}
}

func Test_AddSimple_1(t *testing.T) {
	value1 := uint64(0_777777_777722)
	value2 := uint64(0_000000_000055)
	expected := uint64(0)
	result := AddSimple(value1, value2)
	if result != expected {
		t.Errorf("Error expected %12o, got %12o", expected, result)
	}
}

func Test_GetOnesComplement_1(t *testing.T) {
	value := uint64(100234)
	expected := uint64(100234)
	result := GetOnesComplement(value)
	if result != expected {
		t.Errorf("Error expected %12o, got %12o", expected, result)
	}
}

func Test_GetOnesComplement_2(t *testing.T) {
	// -17dec is -021oct
	value := uint64(0xFFFFFFFF_FFFFFFEF)
	expected := uint64(0_777777_777756)
	result := GetOnesComplement(value)
	if result != expected {
		t.Errorf("Error expected %12o, got %12o", expected, result)
	}
}

func Test_GetSignExtended12_1(t *testing.T) {
	value := uint64(0_776644_011111)
	expected := uint64(0_000000_001111)
	result := GetSignExtended12(value)
	if result != expected {
		t.Errorf("Error expected %12o, got %12o", expected, result)
	}
}

func Test_GetSignExtended12_2(t *testing.T) {
	value := uint64(0_776644_004111)
	expected := uint64(0_777777_774111)
	result := GetSignExtended12(value)
	if result != expected {
		t.Errorf("Error expected %12o, got %12o", expected, result)
	}
}

func Test_GetSignExtended18_1(t *testing.T) {
	value := uint64(0_776644_311111)
	expected := uint64(0_000000_311111)
	result := GetSignExtended18(value)
	if result != expected {
		t.Errorf("Error expected %12o, got %12o", expected, result)
	}
}

func Test_GetSignExtended18_2(t *testing.T) {
	value := uint64(0_000004_404111)
	expected := uint64(0_777777_404111)
	result := GetSignExtended18(value)
	if result != expected {
		t.Errorf("Error expected %12o, got %12o", expected, result)
	}
}

func Test_GetSignExtended24_1(t *testing.T) {
	value := uint64(0_776637_311111)
	expected := uint64(0_000037_311111)
	result := GetSignExtended24(value)
	if result != expected {
		t.Errorf("Error expected %12o, got %12o", expected, result)
	}
}

func Test_GetSignExtended24_2(t *testing.T) {
	value := uint64(0_0066_4440_4111)
	expected := uint64(0_7777_4440_4111)
	result := GetSignExtended24(value)
	if result != expected {
		t.Errorf("Error expected %12o, got %12o", expected, result)
	}
}

func Test_GetTwosComplement_1(t *testing.T) {
	value := uint64(100000)
	expected := uint64(100000)
	result := GetTwosComplement(value)
	if result != expected {
		t.Errorf("Error expected %12o, got %12o", expected, result)
	}
}

func Test_GetTwosComplement_2(t *testing.T) {
	value := uint64(0_777777_777770)
	expected := uint64(0xFFFFFFFF_FFFFFFF9)
	result := GetTwosComplement(uint64(value))
	if result != expected {
		t.Errorf("Error expected %12o, got %12o", expected, result)
	}
}

func Test_GetTwosComplement_3(t *testing.T) {
	value := uint64(0_777777_777777)
	expected := uint64(0)
	result := GetTwosComplement(value)
	if result != expected {
		t.Errorf("Error expected %12o, got %12o", expected, result)
	}
}

func TestNegate(t *testing.T) {
	value := uint64(0377_123456)
	expected := uint64(0_777400_654321)
	result := Negate(value)
	if result != expected {
		t.Errorf("Error expected %12o, got %12o", expected, result)
	}
}
