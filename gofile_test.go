package gofile

import "testing"

func TestNew(t *testing.T) {

	Register("one", buildOne)
	Register("two", buildTwo)

	emptyConfig := map[string]string{}

	one, err := New("one", emptyConfig)

	if err != nil {
		t.Errorf("Could not create driver 'one'")
	}

	switch v := one.(type) {
	default:
		t.Errorf("'one' not a 'driverOne', is %T instead", v)
	case *driverOne:
	}

	two, err := New("two", emptyConfig)

	if err != nil {
		t.Errorf("Could not create driver 'two'")
	}

	switch v := two.(type) {
	default:
		t.Errorf("'two' not a 'driverTwo', is %T instead", v)
	case *driverTwo:
	}

	_, err = New("three", emptyConfig)

	if err == nil {
		t.Errorf("Expected error when trying to create non-existant driver 'three'")
	}

}

type driverOne struct {
	Writer
	Reader
}

type driverTwo struct {
	Writer
	Reader
}

func buildOne(config map[string]string) (Driver, error) {
	return &driverOne{}, nil
}

func buildTwo(config map[string]string) (Driver, error) {
	return &driverTwo{}, nil
}
