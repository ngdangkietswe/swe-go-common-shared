package util

import (
	"fmt"
	"testing"
)

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func (p *Person) Greet() string {
	return fmt.Sprintf("Hello, my name is %s, I am %d years old", p.Name, p.Age)
}

// TestGetField is a function to test GetField function.
func TestGetField(t *testing.T) {
	p := &Person{
		Name: "John Doe",
		Age:  30,
	}

	val, err := GetField(p, "Name")
	if err != nil {
		t.Errorf("GetField failed: %v", err)
	}

	if val != "John Doe" {
		t.Errorf("GetField failed: expected %v but got %v", "John Doe", val)
	}
}

// TestSetField is a function to test SetField function.
func TestSetField(t *testing.T) {
	p := &Person{
		Name: "John Doe",
		Age:  30,
	}

	err := SetField(p, "Name", "Jane Doe")
	if err != nil {
		t.Errorf("SetField failed: %v", err)
	}

	if p.Name != "Jane Doe" {
		t.Errorf("SetField failed: expected %v but got %v", "Jane Doe", p.Name)
	}
}

// TestHasField is a function to test HasField function.
func TestHasField(t *testing.T) {
	p := &Person{
		Name: "John Doe",
		Age:  30,
	}

	if !HasField(p, "Name") {
		t.Errorf("HasField failed: expected %v but got %v", true, false)
	}

	if HasField(p, "Address") {
		t.Errorf("HasField failed: expected %v but got %v", false, true)
	}
}

// TestIsZero is a function to test IsZero function.
func TestIsZero(t *testing.T) {
	// Test zero value
	if !IsZero(Person{}) {
		t.Error("Expected zero value for Person")
	}

	// Test non-zero value
	p := Person{Name: "Alice", Age: 30}
	if IsZero(p) {
		t.Error("Expected non-zero value for Person")
	}
}

// TestIsNil is a function to test IsNil function.
func TestIsNil(t *testing.T) {
	// Test non-nil value
	if IsNil("not nil") {
		t.Error("Expected non-nil value to be false")
	}

	// Test nil pointer
	var ptr *Person
	if !IsNil(ptr) {
		t.Error("Expected nil pointer to be true")
	}

	// Test non-nil pointer
	p := &Person{Name: "Alice"}
	if IsNil(p) {
		t.Error("Expected non-nil pointer to be false")
	}
}

// TestIsEmpty is a function to test IsEmpty function.
func TestIsEmpty(t *testing.T) {
	p := &Person{
		Name: "John Doe",
		Age:  30,
	}

	if IsEmpty(p) {
		t.Errorf("IsEmpty failed: expected %v but got %v", false, true)
	}

	slice := []int{1, 2, 3}
	if IsEmpty(slice) {
		t.Errorf("IsEmpty failed: expected %v but got %v", false, true)
	}

	m := map[string]int{"a": 1, "b": 2}
	if IsEmpty(m) {
		t.Errorf("IsEmpty failed: expected %v but got %v", false, true)
	}
}

// TestIsNilOrEmpty is a function to test IsNilOrEmpty function.
func TestIsNilOrEmpty(t *testing.T) {
	var p *Person

	if !IsNilOrEmpty(p) {
		t.Errorf("IsNilOrEmpty failed: expected %v but got %v", true, false)
	}

	p = &Person{
		Name: "John Doe",
		Age:  30,
	}
	if IsNilOrEmpty(p) {
		t.Errorf("IsNilOrEmpty failed: expected %v but got %v", false, true)
	}

	slice := []int{1, 2, 3}
	if IsNilOrEmpty(slice) {
		t.Errorf("IsNilOrEmpty failed: expected %v but got %v", false, true)
	}

	m := map[string]int{"a": 1, "b": 2}
	if IsNilOrEmpty(m) {
		t.Errorf("IsNilOrEmpty failed: expected %v but got %v", false, true)
	}
}
