package router_test

import (
	"errors"
	"regexp"
	"testing"
)

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

	minValue := min[0]*100 + min[1]*10 + min[2]
	currentValue := current[0]*100 + current[1]*10 + current[2]
	maxValue := max[0]*100 + max[1]*10 + max[2]

	if currentValue >= minValue && currentvalue <= maxValue {
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
