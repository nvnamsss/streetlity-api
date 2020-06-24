package middleware_test

import (
	"errors"
	"regexp"
	"testing"
)

//sumOfRunes return the sum of runes by int in a string
func sumOfRunes(s string) (sum int) {
	for _, r := range s {
		sum += int(r)
	}

	return
}

func compareVersion(currentVersion, minVersion, maxVersion string) (status bool, err error) {
	regFormat, _ := regexp.Compile("\\d*?[.]\\d*")
	reg, _ := regexp.Compile("[.]")
	if !regFormat.Match([]byte(currentVersion)) {
		return false, errors.New("currentVersion format is invalid")
	}

	if !regFormat.Match([]byte(minVersion)) {
		return false, errors.New("currentVersion format is invalid")
	}

	if !regFormat.Match([]byte(maxVersion)) {
		return false, errors.New("currentVersion format is invalid")
	}

	current := reg.Split(currentVersion, -1)
	min := reg.Split(minVersion, -1)
	max := reg.Split(maxVersion, -1)

	if len(current) < 3 {
		return false, errors.New("currentVersion format is invalid")
	}

	if len(min) < 3 {
		return false, errors.New("minVersion format is invalid")
	}

	if len(max) < 3 {
		return false, errors.New("maxVersion format is invalid")
	}
	minValue := sumOfRunes(min[0])*100 + sumOfRunes(min[1])*10 + sumOfRunes(min[2])
	currentValue := sumOfRunes(current[0])*100 + sumOfRunes(current[1])*10 + sumOfRunes(current[2])
	maxValue := sumOfRunes(max[0])*100 + sumOfRunes(max[1])*10 + sumOfRunes(max[2])

	if currentValue >= minValue && currentValue <= maxValue {
		return true, nil
	}

	return false, errors.New("This version is not supported")
}

func TestCompareVersion(t *testing.T) {
	current := "1.0.0"
	min := "1.0.0"
	max := "1.0.0"

	status, err := compareVersion(current, min, max)
	expected := true

	if err != nil {
		t.Errorf("Compare version failed, expected %v but got error %v", expected, err.Error())
	} else {
		t.Logf("Compare version success, expected %v got %v", expected, status)
	}

	//test 2
	current = "1.1.0"
	min = "1.0.0"
	max = "1.0.0"

	status, err = compareVersion(current, min, max)
	expected = false

	if status != expected {
		t.Errorf("Compare version failed, expected %v but got %v", expected, status)
	} else {
		t.Logf("Compare version success, expected %v got %v", expected, status)
	}

	//test 3
	current = "1.1.0"
	min = "1.0.0"
	max = "2.0.0"

	status, err = compareVersion(current, min, max)
	expected = true

	if status != expected {
		t.Errorf("Compare version failed, expected %v but got %v", expected, status)
	} else {
		t.Logf("Compare version success, expected %v got %v", expected, status)
	}

	current = "1.1.0fdf"
	min = "1.0.0"
	max = "2.0.0"

	status, err = compareVersion(current, min, max)
	expected = false

	if status != expected {
		t.Errorf("Compare version failed, expected %v but got %v", expected, status)
	} else {
		t.Logf("Compare version success, expected %v got %v", expected, status)
	}
}
