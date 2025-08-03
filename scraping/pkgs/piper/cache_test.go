package piper

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"gorm.io/gorm"
)

func TestPiperCache(t *testing.T) {
	sqliteDbPath := "tmp/test_cache.db"

	c, err := NewCacheDb(sqliteDbPath)
	if err != nil {
		t.Fatal(err)
	}

	if err = c.Validate(); err != nil && err != ErrIncorrectSchema {
		t.Fatal(err)
	} else if err == ErrIncorrectSchema {
		if err = c.Setup(); err != nil {
			t.Fatal(err)
		}
	}

	type s struct {
		Foo int
		Bar string
	}

	if err := c.Set("foo", []byte("bar"), time.Second*3); err != nil {
		t.Fatal(err)
	}

	dat, err := json.Marshal(s{Foo: 69, Bar: "foo"})
	if err != nil {
		t.Fatal(err)
	}

	if err := c.Set("bar", dat, time.Second*3); err != nil {
		t.Fatal(err)
	}

	foo, err := c.Get("foo")
	if err != nil {
		t.Fatal(err)
	}

	bar, err := c.Get("bar")
	if err != nil {
		t.Fatal(err)
	}

	var rs s

	if err := json.Unmarshal(bar, &rs); err != nil {
		t.Fatal(err)
	}

	fmt.Println(string(foo))
	fmt.Println(rs)

	time.Sleep(time.Second * 3)

	foo, err = c.Get("foo")
	if err == nil || err != gorm.ErrRecordNotFound {
		t.Errorf("Want %v, get %v", gorm.ErrRecordNotFound, err)
	}

	bar, err = c.Get("bar")
	if err == nil || err != gorm.ErrRecordNotFound {
		t.Errorf("Want %v, get %v", gorm.ErrRecordNotFound, err)
	}
}
