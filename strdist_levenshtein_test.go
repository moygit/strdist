package strdist

import (
	"math"
	"testing"
)

const float_tolerance = 0.01

var insCosts [alphabet_size]DistanceType
var delCosts [alphabet_size]DistanceType
var subCosts [alphabet_size][alphabet_size]DistanceType

func lev(str1, str2 string) (DistanceType, DistanceType) {
	return LevenshteinUnsafe(str1, str2, &insCosts, &delCosts, &subCosts)
}

func test(t *testing.T, f func(string, string) (DistanceType, DistanceType), name, s1, s2 string, expDist DistanceType) {
	// when:
	dist, _ := f(s1, s2)
	// then:
	if dist != expDist {
		t.Errorf("%s(%s, %s) == %f, expected %f", name, s1, s2, dist, expDist)
	}
}

func TestLevenshteinDefaultCosts(t *testing.T) {
	// given:
	insCosts = *New1dCostArray(1)
	delCosts = *New1dCostArray(1)
	subCosts = *New2dCostArray(1)

	test(t, lev, "lev", "1234", "1234", 0)
	test(t, lev, "lev", "", "1234", 4)
	test(t, lev, "lev", "1234", "", 4)
	test(t, lev, "lev", "", "", 0)
	test(t, lev, "lev", "1234", "12", 2)
	test(t, lev, "lev", "1234", "14", 2)
	test(t, lev, "lev", "1111", "1", 3)
}

func TestLevenshteinInsertCost(t *testing.T) {
	// given:
	insCosts = *New1dCostArray(1)
	delCosts = *New1dCostArray(1)
	subCosts = *New2dCostArray(1)

	// Change insert cost of 'a' to 5
	insCosts[int('a')] = 5
	test(t, lev, "lev", "", "a", 5)
	test(t, lev, "lev", "a", "", 1)
	test(t, lev, "lev", "", "aa", 10)
	test(t, lev, "lev", "a", "aa", 5)
	test(t, lev, "lev", "aa", "a", 1)
	test(t, lev, "lev", "asdf", "asdf", 0)
	test(t, lev, "lev", "xyz", "abc", 3)
	test(t, lev, "lev", "xyz", "axyz", 5)
	test(t, lev, "lev", "x", "ax", 5)
}

func TestLevenshteinDeleteCost(t *testing.T) {
	// given:
	insCosts = *New1dCostArray(1)
	delCosts = *New1dCostArray(1)
	subCosts = *New2dCostArray(1)

	// Change delete cost of 'z' to 7.5
	delCosts[int('z')] = 7.5
	test(t, lev, "lev", "", "z", 1)
	test(t, lev, "lev", "z", "", 7.5)
	test(t, lev, "lev", "xyz", "zzxz", 3)
	test(t, lev, "lev", "zzxzzz", "xyz", 18)
}

func TestLevenshteinSubstituteCost(t *testing.T) {
	// given:
	insCosts = *New1dCostArray(1)
	delCosts = *New1dCostArray(1)
	subCosts = *New2dCostArray(1)

	// Change delete cost of 'z' to 7.5
	subCosts[int('a')][int('z')] = 1.2
	subCosts[int('z')][int('a')] = 0.1
	test(t, lev, "lev", "a", "z", 1.2)
	test(t, lev, "lev", "z", "a", 0.1)
	test(t, lev, "lev", "a", "", 1)
	test(t, lev, "lev", "", "a", 1)
	test(t, lev, "lev", "asdf", "zzzz", 4.2)
	test(t, lev, "lev", "asdf", "zz", 4.0)
	test(t, lev, "lev", "asdf", "zsdf", 1.2)
	test(t, lev, "lev", "zsdf", "asdf", 0.1)
}

func TestLevenshteinUnsafeFailsOnAccent(t *testing.T) {
	// given:
	insCosts = *New1dCostArray(1)
	delCosts = *New1dCostArray(1)
	subCosts = *New2dCostArray(1)
	s1 := "curaçao"
	s2 := ""
	// when we check *safely*
	dist, maxDist := LevenshteinSafe(s1, s2, &insCosts, &delCosts, &subCosts)
	expDist, expMaxDist := 7.0, 7.0
	// then all is well
	if math.Abs(float64(dist)-expDist) > float_tolerance {
		t.Errorf("LevenshteinSafe(%s, %s) == %f, expected %f", s1, s2, dist, expDist)
	}
	if math.Abs(float64(maxDist)-expMaxDist) > float_tolerance {
		t.Errorf("LevenshteinSafe(%s, %s) max-distance == %f, expected %f", s1, s2, maxDist, expMaxDist)
	}
	// but when we check *un*safely then go correctly panics
	defer catchPanicOrElse(t, `Unsafely matching "curaçao" vs "" should have failed but didn't.`)
	test(t, lev, "lev", s1, s2, -1.0)
}

func catchPanicOrElse(t *testing.T, msg string) {
	if r := recover(); r == nil {
		t.Errorf(msg)
	}
}
