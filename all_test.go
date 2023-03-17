package ecs_cpp_style

import (
	"fmt"
	"testing"
)

type DummySystemA struct {
	Entities []Entity
}

func (dsa *DummySystemA) AddEntity(entity Entity) {
	dsa.Entities = append(dsa.Entities, entity)
}

func (dsa *DummySystemA) RemoveEntity(entity Entity) {
	pos := -1
	for i, ent := range dsa.Entities {
		if ent == entity {
			pos = i
			break
		}
	}
	if pos >= 0 {
		dsa.Entities[pos] = dsa.Entities[len(dsa.Entities)-1]
		dsa.Entities = dsa.Entities[:len(dsa.Entities)-1]
	}
}

type DummySystemB struct {
	Entities []Entity
}

func (dsa *DummySystemB) AddEntity(entity Entity) {
	dsa.Entities = append(dsa.Entities, entity)
}

func (dsa *DummySystemB) RemoveEntity(entity Entity) {
	pos := -1
	for i, ent := range dsa.Entities {
		if ent == entity {
			pos = i
			break
		}
	}
	if pos >= 0 {
		dsa.Entities[pos] = dsa.Entities[len(dsa.Entities)-1]
		dsa.Entities = dsa.Entities[:len(dsa.Entities)-1]
	}
}

func TestDumbSignature(t *testing.T) {
	system1 := DumbSignature{false, true, false, true, false}

	entity1 := DumbSignature{false, true, false, true, false}
	entity2 := DumbSignature{false, true, true, true, false}
	entity3 := DumbSignature{false, false, false, true, false}

	r1, _ := system1.Contains(entity1)
	r2, _ := system1.Contains(entity2)
	r3, _ := system1.Contains(entity3)

	if !r1 || !r2 {
		t.Fatalf("Should both be true, got: %v and %v", r1, r2)
	}

	if r3 {
		t.Fatalf("Should both be false, got: %v", r3)
	}
}

func TestQueue(t *testing.T) {

	q := make(Queue[Entity], 0)

	for i := 0; i < 5; i++ {
		q.enqueue(Entity(i))
	}

	if len(q) != 5 {
		t.Fatalf("Wanted len(q) == 5 | Got %v", len(q))
	}

	a, aOk := q.dequeue()
	b, bOk := q.dequeue()

	if a != 0 || !aOk {
		t.Fatalf("Wanted a == 0 and ok = true | Got %v and %v", a, aOk)
	}
	if b != 1 || !bOk {
		t.Fatalf("Wanted b == 1 and ok = true | Got %v and %v", b, bOk)
	}
}

func TestEntityManager(t *testing.T) {

	entManager := NewEntityManager(3)

	e0 := entManager.CreateEntity()
	e1 := entManager.CreateEntity()

	wantOne := 0
	wantTwo := 1

	if Entity(wantOne) != e0 {
		t.Fatalf("Wanted %v | got %v", wantOne, e0)
	}

	if Entity(wantTwo) != e1 {
		t.Fatalf("Wanted %v | got %v", wantTwo, e1)
	}

	entManager.DestroyEntity(1)

	e2 := entManager.CreateEntity()
	e3 := entManager.CreateEntity()

	if e2 != 2 {
		t.Fatalf("Wanted 2 | got %v", e2)
	}
	if e3 != 1 {
		t.Fatalf("Wanted 1 | got %v", e3)
	}

	sig := make(DumbSignature, MAX_COMPONENTS)
	sig.Set(0, true)

	entManager.SetSignature(e3, sig)
	_sig := entManager.GetSignature(e3)
	if _sig[0] != true {
		t.Fatal("Wanted sig true | got false")
	}

	entManager.DestroyEntity(e3)
	_sig = entManager.GetSignature(e3)
	if _sig[0] != false {
		t.Fatal("Wanted resetted sig false | got true")
	}
}

func TestComponentArray(t *testing.T) {
	ca := NewComponentArray[string](5)

	ca.Insert(Entity(0), "A")
	ca.Insert(Entity(1), "B")
	ca.Insert(Entity(2), "C")
	ca.Insert(Entity(3), "D")

	e, _ := ca.Get(Entity(1))
	v := *e

	if v != "B" {
		t.Fatalf(`Expected e = "B" | Got v = %v `, v)
	}

	ca.EntityDestroyed(Entity(1)) // ADC ; A0, D3, C2
	ca.EntityDestroyed(Entity(3)) // AC; A0, C2
	ca.Insert(Entity(4), "E")     // ACE, A0, C2, E4

	if ca.Components[0] != "A" || ca.Components[1] != "C" || ca.Components[2] != "E" {
		t.Fatalf("Got wrong components, expected [A, C, E] | got %v", ca.Components[:ca.size])
	}
}

func TestComponentManager(t *testing.T) {

	type DummyComponentA struct {
		SomeVal string
	}

	type DummyComponentB struct {
		SomeOtherVal string
	}

	compA := DummyComponentA{SomeVal: "component A"}

	cm := newComponentManager()

	registerComponent[DummyComponentA](cm)
	registerComponent[DummyComponentB](cm)

	tp := getComponentType[DummyComponentB](cm)

	if tp != 1 {
		t.Fatalf("Wanted 1, got %v", tp)
	}

	entManager := NewEntityManager(3)
	entCompA := entManager.CreateEntity()
	addComponent[DummyComponentA](cm, entCompA, compA)

	sameA, _ := getComponent[DummyComponentA](cm, entCompA)

	if sameA.SomeVal != "component A" {
		t.Fatalf("wanted component A | got %v", sameA.SomeVal)
	}

	removeComponent[DummyComponentA](cm, entCompA)
	sameA, ok := getComponent[DummyComponentA](cm, entCompA)
	if ok == nil {
		t.Fatalf("expected error indicating component had already been removed for Entity %v", entCompA)
	}
}

func TestComponentManagerEntityDestroy(t *testing.T) {
	type DummyComponentA struct {
		SomeVal string
	}
	type DummyComponentB struct {
		SomeOtherVal string
	}

	compA := DummyComponentA{SomeVal: "component A"}
	compB := DummyComponentB{SomeOtherVal: "component B"}

	cm := newComponentManager()
	entManager := NewEntityManager(3)

	ent := entManager.CreateEntity()

	registerComponent[DummyComponentA](cm)
	registerComponent[DummyComponentB](cm)

	addComponent[DummyComponentA](cm, ent, compA)
	addComponent[DummyComponentB](cm, ent, compB)

	entityDestroyed(cm, ent)

	_, okA := getComponent[DummyComponentA](cm, ent)
	_, okB := getComponent[DummyComponentB](cm, ent)
	if okA == nil || okB == nil {
		t.Fatalf("expected error indicating component had already been removed for Entity %v", ent)
	}
}

func TestNewSystemManager(t *testing.T) {

	sm := NewSystemManager()

	sys := &DummySystemA{}
	sm.RegisterSystem(sys)
}

func TestComprehensive(t *testing.T) {
	em := NewEntityManager(1000)
	cm := newComponentManager()
	sm := NewSystemManager()

	ent1 := em.CreateEntity()

	registerComponent[Transform](cm)
	registerComponent[Vec2](cm)

	sys := DummySystemA{Entities: make([]Entity, 0)}
	sm.RegisterSystem(&sys)
	sysSig := make(DumbSignature, MAX_COMPONENTS)
	sysSig.Set(int(getComponentType[Transform](cm)), true)
	sysSig.Set(int(getComponentType[Vec2](cm)), true)

	sm.SetSignature(&sys, sysSig)

	transf := Transform{
		position: Vec2{X: 10, Y: 10},
		scale:    Vec2{X: 1, Y: 1},
	}
	vec := Vec2{
		X: 1,
		Y: 1,
	}

	addComponent(cm, ent1, transf)
	addComponent(cm, ent1, vec)

	sig := em.GetSignature(ent1)
	sig.Set(int(getComponentType[Transform](cm)), true)
	sig.Set(int(getComponentType[Vec2](cm)), true)

	em.SetSignature(ent1, sig)
	sm.EntitySignatureChanged(ent1, sig)
}

func TestNewCoordinator(t *testing.T) {

	coord := NewCoordinator(100)

	// Register all components
	RegisterNewComponentType[Transform](coord)
	RegisterNewComponentType[Vec2](coord)

	sys1 := DummySystemA{Entities: make([]Entity, 0)}
	sys2 := DummySystemB{Entities: make([]Entity, 0)}
	RegisterNewSystem(coord, &sys1)
	RegisterNewSystem(coord, &sys2)

	sys1Sig := make(DumbSignature, MAX_COMPONENTS)
	sys1Sig.Set(int(GetComponentType[Transform](coord)), true)
	sys1Sig.Set(int(GetComponentType[Vec2](coord)), true)
	SetSystemSignature(coord, &sys1, sys1Sig)

	sys2Sig := make(DumbSignature, MAX_COMPONENTS)
	sys2Sig.Set(int(GetComponentType[Vec2](coord)), true)
	SetSystemSignature(coord, &sys2, sys2Sig)

	ent1 := CreateNewEntity(coord)
	ent2 := CreateNewEntity(coord)

	transf := Transform{
		position: Vec2{X: 10, Y: 10},
		scale:    Vec2{X: 1, Y: 1},
	}
	vec := Vec2{
		X: 1,
		Y: 1,
	}
	AddNewComponent(coord, ent1, transf)
	AddNewComponent(coord, ent1, vec)

	AddNewComponent(coord, ent2, vec)

	for ent := range sys1.Entities {
		trans := GetExistingComponent[Transform](coord, Entity(ent))
		vec := GetExistingComponent[Vec2](coord, Entity(ent))

		fmt.Println("trans ", trans, " vec ", vec)
	}

	transArray := coord.cm.componentArrays["Transform"].(*ComponentArray[Transform])
	for ent := range sys1.Entities {
		idx := transArray.EntityToIndex[Entity(ent)]

		result := transArray.Components[idx]
		fmt.Println("trans ", result)
	}

	fmt.Println(coord)
}
