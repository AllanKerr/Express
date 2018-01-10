 /***************************************************************************************
 *    Fun with Java string hashing in Go
 *    Author: Manni Wood
 *    Date: Dec. 24, 2017
 *    Code version: 1.0
 *    Availability: https://www.manniwood.com/2016_03_20/fun_with_java_string_hashing.html
 *
 ***************************************************************************************/

package kube

// Hash function to create hash codes for strings using the Java hashing algorithm
func HashString(s string) int {
	h := 0
	for i := 0; i < len(s); i++ {
		h = 31*h + int(s[i])
	}
	return h
}