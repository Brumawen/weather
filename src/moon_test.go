package main

import (
	"fmt"
	"testing"
	"time"
)

func TestCanGetMoon(t *testing.T) {
	m := Moon{}
	m.ForDate(time.Now())
	fmt.Println(m)
}
