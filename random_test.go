package helper

import "testing"

func TestRandomOption_GenerateRandomString(t *testing.T) {
	opt := &RandomOption{
		LowerCase:   true,
		UpperCase:   true,
		Numeric:     true,
		SpecialChar: true,
	}
	randomString, err := opt.GenerateRandomString(10)

	if err != nil {
		t.Errorf("got error : %s", err.Error())
	} else {
		if len(randomString) != 10 {
			t.Errorf("got error : %s", "length missmatch")
		}
	}
}
