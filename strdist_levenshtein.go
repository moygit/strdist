package strdist

import "github.com/mozillazg/go-unidecode"

// Calculate unweighted Levenshtein distance between two strings.
// This function incurs a performance penalty to convert both strings to ASCII (0-127) with `unidecode`.
func DefaultLevenshtein(str1, str2 string) (DistanceType, DistanceType) {
	return LevenshteinSafe(str1, str2, &unitArray, &unitArray, &unit2dArray)
}

// Calculate weighted Levenshtein distance between two strings.
// This function incurs a performance penalty to convert both strings to ASCII (0-127) with `unidecode`.
func LevenshteinSafe(str1, str2 string,
	insertCosts *[alphabet_size]DistanceType,
	deleteCosts *[alphabet_size]DistanceType,
	substituteCosts *[alphabet_size][alphabet_size]DistanceType) (DistanceType, DistanceType) {
	str1 = unidecode.Unidecode(str1)
	str2 = unidecode.Unidecode(str2)
	return LevenshteinUnsafe(str1, str2, insertCosts, deleteCosts, substituteCosts)
}

// Calculate weighted Levenshtein distance between two strings.
// This function assumes both strings have been converted to only have values in [0,...,127].
// If this is false then the function will `panic`.
func LevenshteinUnsafe(str1, str2 string,
	insertCosts *[alphabet_size]DistanceType,
	deleteCosts *[alphabet_size]DistanceType,
	substituteCosts *[alphabet_size][alphabet_size]DistanceType) (DistanceType, DistanceType) {
	return levenshteinCore([]byte(str1), []byte(str2), insertCosts, deleteCosts, substituteCosts)
}

// Calculate weighted Levenshtein distance between two byte-arrays.
// This function assumes both underlying strings have been converted to only have values in [0,...,127].
// If this is false then the function will `panic`.
func levenshteinCore(str1, str2 []byte,
	insertCosts *[alphabet_size]DistanceType,
	deleteCosts *[alphabet_size]DistanceType,
	substituteCosts *[alphabet_size][alphabet_size]DistanceType) (distance, maxWeightedDistance DistanceType) {
	len1 := len(str1)
	len2 := len(str2)
	d := makeFast2dDistanceSlice(len1+1, len2+1)

	// d[0][0] is 0
	for i := 1; i <= len1; i++ {
		ch := str1[i-1]
		d[i][0] = d[i-1][0] + deleteCosts[ch]
	}
	for j := 1; j <= len2; j++ {
		ch := str2[j-1]
		d[0][j] = d[0][j-1] + insertCosts[ch]
	}

	for i := 1; i <= len1; i++ {
		ch1 := str1[i-1]
		for j := 1; j <= len2; j++ {
			ch2 := str2[j-1]
			if ch1 == ch2 {
				d[i][j] = d[i-1][j-1]
			} else {
				d[i][j] = min3(d[i-1][j]+deleteCosts[ch1], d[i][j-1]+insertCosts[ch2], d[i-1][j-1]+substituteCosts[ch1][ch2])
			}
		}
	}

	distance = d[len1][len2]
	maxWeightedDistance = max2(d[len1][0], d[0][len2])
	return
}
