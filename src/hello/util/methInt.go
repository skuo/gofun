package util

import (
	"fmt"
)

// =======================
// method
// pass pointer for function param

type ounce float64

func (o ounce) cup() cup {
	return cup(o * 0.1250)
}

type cup float64

func (c cup) quart() quart {
	return quart(c * 0.25)
}
func (c cup) ounce() ounce {
	return ounce(c * 8.0)
}

type quart float64

func (q quart) gallon() gallon {
	return gallon(q * 0.25)
}
func (q quart) cup() cup {
	return cup(q * 4.0)
}

type gallon float64

func (g gallon) quart() quart {
	return quart(g * 4)
}

// pass by value. Does not change the origin value
func (g gallon) half() {
	g = gallon(g * 0.5)
}

// pass by reference. Original value is changed.
func (g *gallon) double() {
	*g = gallon(*g * 2)
}

func tryMethod() {
	gal := gallon(5)
	fmt.Printf("%.2f gallons = %.2f quarts\n", gal, gal.quart())
	ozs := gal.quart().cup().ounce()
	fmt.Printf("%.2f gallons = %.2f ounces\n", gal, ozs)

	gal.half()
	fmt.Println("half() is passed by value, does not change gal: ", gal)
	gal.double()
	fmt.Println("double() is passed by reference, gal is changed: ", gal)
}

// =============================
// struct as object
type fuel int

const (
	GASOLINE fuel = iota
	BIO
	ELECTRIC
	JET
)

type vehicle struct {
	make  string
	model string
}

type engine struct {
	fuel   fuel
	thrust int
}

func (e *engine) start() {
	fmt.Println("Engine started.")
}

type truck struct {
	// vehicle and engine are embedded structs
	vehicle
	engine
	axels  int
	wheels int
	class  int
}

func (t *truck) drive() {
	fmt.Printf("Truck %s %s, on the go!\n", t.make, t.model)
}

func newTruck(mk, mdl string) *truck {
	return &truck{vehicle: vehicle{mk, mdl}}
}

type plane struct {
	vehicle
	engine
	engineCount int
	fixedWings  bool
	maxAltitude int
}

func newPlane(mk, mdl string) *plane {
	p := &plane{}
	p.make = mk
	p.model = mdl
	return p
}

func (p *plane) fly() {
	fmt.Printf("Aircraft %s %s clear for takeoff!\n", p.make, p.model)
}

func tryStructObj() {
	t := newTruck("Ford", "F750")
	t.axels = 2
	t.wheels = 6
	t.class = 3
	t.start()
	t.drive()

	p := newPlane("HondaJet", "HA-420")
	p.fuel = JET
	p.thrust = 2050
	p.engineCount = 2
	p.fixedWings = true
	p.maxAltitude = 43000
	p.start()
	p.fly()

}

// =============================
// struct method
type volume struct {
	unit string
	qty  float64
}

func (v volume) String() string {
	return fmt.Sprintf("%.2f %s", v.qty, v.unit)
}

func tryStructMethod() {
	v := volume{unit: "quart", qty: 7}
	v.String()
	fmt.Println(v.String())
}

// =============================
// Interface
type food interface {
	eat()
}

type plant interface {
    color()    
}

type veggie string

func (v veggie) eat() {
	fmt.Println("Eating", v)
}
func (v veggie) color() {
    fmt.Println("Multiple colors", v)
}

type meat string

func (m meat) eat() {
	fmt.Println("Eating tasty", m)
}

func eat(f food) {
	switch morsel := f.(type) {
	case veggie:
		if morsel == "okra" {
			fmt.Println("Yuk! not eating ", morsel)
		} else {
			morsel.eat()
		}
	case meat:
		if morsel == "beef" {
			fmt.Println("Yuk! not eating ", morsel)
		} else {
			morsel.eat()
		}
	default:
		fmt.Println("Not eating whatever that is: ", f)
	}
}

func printAnyType(val interface{}) {
	fmt.Println(val)
}

func tryInterface() {
	eat(veggie("carrot"))
	eat(meat("lamb"))
	eat(veggie("okra"))
	eat(meat("beef"))
	
	// check the interface type
	cabbage := veggie("cabbage")
	_, ok := interface{}(cabbage).(plant)
	if ok {
	    fmt.Println("cabbage is a plant")
	} else {
	    fmt.Println("cabbage is not a plant")
	}
	_, ok = interface{}(cabbage).(food)
	if ok {
	    fmt.Println("cabbage is food")
	} else {
	    fmt.Println("cabbage is not food")
	}
    
	var anyType interface{}
	anyType = 77.0
	anyType = "I am a string now"
	fmt.Println(anyType)

	printAnyType("The car is slow")
	m := map[string]string{"ID": "12345", "name": "Kerry"}
	printAnyType(m)
	printAnyType(1253443455)
}

// =============================
func TryMethInt() {
	tryMethod()
	tryStructObj()
	tryStructMethod()
	tryInterface()
}
