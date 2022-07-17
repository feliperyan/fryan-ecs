package ecs_cpp_style

import (
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
	if e != "B" {
		t.Fatalf(`Expected e = "one" | Got e = %v"`, e)
	}

	ca.EntityDestroyed(Entity(1))
	ca.EntityDestroyed(Entity(3))
	ca.Insert(Entity(4), "E")

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

	cm := NewComponentManager()

	RegisterComponent[DummyComponentA](cm)
	RegisterComponent[DummyComponentB](cm)

	tp := GetComponentType[DummyComponentB](cm)

	if tp != 1 {
		t.Fatalf("Wanted 1, got %v", tp)
	}

	entManager := NewEntityManager(3)
	entCompA := entManager.CreateEntity()
	AddComponent[DummyComponentA](cm, entCompA, compA)

	sameA, _ := GetComponent[DummyComponentA](cm, entCompA)

	if sameA.SomeVal != "component A" {
		t.Fatalf("wanted component A | got %v", sameA.SomeVal)
	}

	RemoveComponent[DummyComponentA](cm, entCompA)
	sameA, ok := GetComponent[DummyComponentA](cm, entCompA)
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

	cm := NewComponentManager()
	entManager := NewEntityManager(3)

	ent := entManager.CreateEntity()

	RegisterComponent[DummyComponentA](cm)
	RegisterComponent[DummyComponentB](cm)

	AddComponent[DummyComponentA](cm, ent, compA)
	AddComponent[DummyComponentB](cm, ent, compB)

	EntityDestroyed(cm, ent)

	_, okA := GetComponent[DummyComponentA](cm, ent)
	_, okB := GetComponent[DummyComponentB](cm, ent)
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
	cm := NewComponentManager()
	sm := NewSystemManager()

	ent1 := em.CreateEntity()

	RegisterComponent[Transform](cm)
	RegisterComponent[Vec2](cm)

	sys := DummySystemA{Entities: make([]Entity, 0)}
	sm.RegisterSystem(&sys)
	sysSig := make(DumbSignature, MAX_COMPONENTS)
	sysSig.Set(int(GetComponentType[Transform](cm)), true)
	sysSig.Set(int(GetComponentType[Vec2](cm)), true)

	sm.SetSignature(&sys, sysSig)

	transf := Transform{
		position: Vec2{X: 10, Y: 10},
		scale:    Vec2{X: 1, Y: 1},
	}
	vec := Vec2{
		X: 1,
		Y: 1,
	}

	AddComponent[Transform](cm, ent1, transf)
	AddComponent[Vec2](cm, ent1, vec)

	sig := em.GetSignature(ent1)
	sig.Set(int(GetComponentType[Transform](cm)), true)
	sig.Set(int(GetComponentType[Vec2](cm)), true)

	em.SetSignature(ent1, sig)
	sm.EntitySignatureChanged(ent1, sig)

}
