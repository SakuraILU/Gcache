package group

import (
	"testing"
)

func TestGroup1(t *testing.T) {
	// define several kvs
	kvs := map[string]string{
		"Tom":       "cat",
		"Jerry":     "mouse",
		"Tom&Jerry": "friend",
	}

	// define a Getter
	var getter GetFunc = func(key string) ([]byte, error) {
		// return kvs[key], nil
		if v, ok := kvs[key]; ok {
			return []byte(v), nil
		}
		return nil, nil
	}

	// define a Group
	g := NewGroup("Tom&Jerry", 1000, getter)

	// get values from Group
	for k, v := range kvs {
		if val, err := g.Get(k); err != nil {
			t.Errorf("[Error]: %v", err)
		} else if val.ToString() != v {
			t.Errorf("[Error]: %v != %v", val.ToString(), v)
		}
	}
}

// multi groups
func TestGroup2(t *testing.T) {
	// group1
	kvs1 := map[string]string{
		"Tom":       "cat",
		"Jerry":     "mouse",
		"Tom&Jerry": "friend",
	}
	var getter1 GetFunc = func(key string) ([]byte, error) {
		if v, ok := kvs1[key]; ok {
			return []byte(v), nil
		}
		return nil, nil
	}
	g1 := NewGroup("Tom&Jerry", 1000, getter1)

	// group2
	kvs2 := map[string]string{
		"apple":  "fruit",
		"banana": "fruit",
		"orange": "fruit",
	}
	var getter2 GetFunc = func(key string) ([]byte, error) {
		if v, ok := kvs2[key]; ok {
			return []byte(v), nil
		}
		return nil, nil
	}
	g2 := NewGroup("fruit", 1000, getter2)

	// group3
	kvs3 := map[string]string{
		"bag":   "thing",
		"book":  "thing",
		"ship":  "vehicle",
		"plane": "vehicle",
	}
	var getter3 GetFunc = func(key string) ([]byte, error) {
		if v, ok := kvs3[key]; ok {
			return []byte(v), nil
		}
		return nil, nil
	}
	g3 := NewGroup("thing", 1000, getter3)

	// get values from Group
	g1, _ = GetGroup("Tom&Jerry")
	for k, v := range kvs1 {
		if val, err := g1.Get(k); err != nil {
			t.Errorf("[Error]: %v", err)
		} else if val.ToString() != v {
			t.Errorf("[Error]: %v != %v", val.ToString(), v)
		}
	}

	g2, _ = GetGroup("fruit")
	for k, v := range kvs2 {
		if val, err := g2.Get(k); err != nil {
			t.Errorf("[Error]: %v", err)
		} else if val.ToString() != v {
			t.Errorf("[Error]: %v != %v", val.ToString(), v)
		}
	}

	g3, _ = GetGroup("thing")
	for k, v := range kvs3 {
		if val, err := g3.Get(k); err != nil {
			t.Errorf("[Error]: %v", err)
		} else if val.ToString() != v {
			t.Errorf("[Error]: %v != %v", val.ToString(), v)
		}
	}

}
