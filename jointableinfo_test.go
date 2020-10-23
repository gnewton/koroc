package main

import (
	"errors"
	"testing"
)

func TestJoinTableInfo_makeKey(t *testing.T) {
	dialect := new(DialectSqlite3)
	personTable, pid, pname, pHasCar, err := personTableFull(dialect)
	if err != nil {
		t.Fatal(err)
	}

	// Person
	person, err := personTable.Record()
	if err != nil {
		t.Fatal(err)
	}
	if err := person.Add(pid, uint32(42)); err != nil {
		t.Fatal(err)
	}
	if err := person.Add(pname, "Bill"); err != nil {
		t.Fatal(err)
	}
	if err := person.Add(pHasCar, true); err != nil {
		t.Fatal(err)
	}

	// Car
	carTable, carId, manufacturer, model, year, err := carTableFull(dialect)
	if err != nil {
		t.Fatal(err)
	}
	if carTable == nil || carId == nil || manufacturer == nil || model == nil || year == nil {
		t.Fatal("Car is broken")
	}
	car, err := carRecord1(carTable, carId, manufacturer, model, year)

	jt, err := NewJoinTable2(personTable, carTable, manufacturer, model, year)
	if err != nil {
		t.Fatal(err)
	}

	key, err := jt.joinTableInfo.makeKey(car)
	if err != nil {
		t.Logf("%+v", manufacturer)
		t.Logf("%+v", manufacturer.typ)
		t.Logf("**0 %T", car.values[0])
		t.Logf("**1 %T", car.values[1])
		t.Logf("**2 %T", car.values[2])
		t.Logf("**3 %T", car.values[3])
		t.Logf("%+v", *person)
		t.Logf("%+v", *car)
		t.Fatal(err)
	}
	if key != "0:Ford|1:Escort|2:1988" {
		t.Fatal(errors.New("Bad key string:" + key))
	}
	t.Log(key)

}
