package utils

import (
	"testing"
)

func TestValidPhoneNumber(t *testing.T) {
	testsSucceed := []bool{ValidPhoneNumber("0748326171"), ValidPhoneNumber("0748-326-171"),
		ValidPhoneNumber("+40748326171"), ValidPhoneNumber("+40748326171"),
		ValidPhoneNumber("+40748-326-171"), ValidPhoneNumber("+40748 326 171")}
	testsFail := []bool{ValidPhoneNumber("748326171"), ValidPhoneNumber("0040748326171")}

	for i, test := range testsSucceed {
		if test != true {
			t.Fatalf("Test %d should have succeded, but failed", i+1)
		}
	}
	for i, test := range testsFail {
		if test != false {
			t.Fatalf("Test %d should have failed, but succeded", i+1)
		}
	}
}

func TestValidatePassword(t *testing.T) {
	testsSucceed := []bool{
		ValidatePassword("!BarOsAnu1"),
		ValidatePassword("-Megabigb0ss"),
		ValidatePassword("bAROSANU!1")}
	testsFail := []bool{
		ValidatePassword("BarOsAnu"),
		ValidatePassword("megabigb0ss"),
		ValidatePassword("hermosa"),
		ValidatePassword("BAROSANu1! ")}

	tryTests(t, &testsSucceed, &testsFail)
}

func tryTests(t *testing.T, succeed *[]bool, fail *[]bool) {
	for i, test := range *succeed {
		if test != true {
			t.Fatalf("Test %d should have succeded, but failed", i+1)
		}
	}
	for i, test := range *fail {
		if test != false {
			t.Fatalf("Test %d should have failed, but succeded", i+1)
		}
	}
}
