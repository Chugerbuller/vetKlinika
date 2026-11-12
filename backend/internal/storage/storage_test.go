package storage

import (
	"bytes"
	"database/sql"
	"errors"
	"path/filepath"
	"strconv"
	"testing"
	"time"
)

const dateLayout = "20060102"

// newTestStorage создаёт чистую БД во временной папке на каждый тест.
func newTestStorage(t *testing.T) *Storage {
	t.Helper()
	st, err := New(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	t.Cleanup(func() { _ = st.Close() })
	return st
}

func mustAddPerson(t *testing.T, st *Storage, p Person) int {
	t.Helper()
	id, err := st.AddPerson(p)
	if err != nil {
		t.Fatalf("AddPerson() error = %v", err)
	}
	return id
}

func mustCreateType(t *testing.T, st *Storage, name string) int {
	t.Helper()
	id, err := st.CreateType(Types{Type: name})
	if err != nil {
		t.Fatalf("CreateType(%q) error = %v", name, err)
	}
	return id
}

func mustAddPet(t *testing.T, st *Storage, pet Pet, ownerID int) int {
	t.Helper()
	id, err := st.AddPet(pet, ownerID)
	if err != nil {
		t.Fatalf("AddPet() error = %v", err)
	}
	return id
}

func TestNew_EmptyDatabase(t *testing.T) {
	st := newTestStorage(t)

	persons, err := st.GetPersons()
	if err != nil {
		t.Fatalf("GetPersons() error = %v", err)
	}
	if len(*persons) != 0 {
		t.Errorf("GetPersons() = %d persons, want 0", len(*persons))
	}

	types, err := st.GetTypes()
	if err != nil {
		t.Fatalf("GetTypes() error = %v", err)
	}
	if len(*types) != 0 {
		t.Errorf("GetTypes() = %d types, want 0", len(*types))
	}
}

func TestPersons_CRUD(t *testing.T) {
	st := newTestStorage(t)

	id := mustAddPerson(t, st, Person{
		Name: "Иван", Surname: "Петров", Patronymic: "Сергеевич",
		Phone: "9001234567", Address: "Пермь",
	})
	if id == 0 {
		t.Fatal("AddPerson() returned id 0")
	}

	got, err := st.GetPersonById(id)
	if err != nil {
		t.Fatalf("GetPersonById() error = %v", err)
	}
	// AddPerson приводит ФИО к нижнему регистру.
	if got.Name != "иван" || got.Surname != "петров" || got.Patronymic != "сергеевич" {
		t.Errorf("GetPersonById() = %+v, want lowercased name", got)
	}
	if got.Phone != "9001234567" || got.Address != "Пермь" {
		t.Errorf("GetPersonById() = %+v, phone/address changed", got)
	}

	found, err := st.GetPersonsByName("петров")
	if err != nil {
		t.Fatalf("GetPersonsByName() error = %v", err)
	}
	if len(*found) != 1 || (*found)[0].Id != id {
		t.Errorf("GetPersonsByName(\"петров\") = %+v, want person %d", *found, id)
	}

	got.Phone = "9007654321"
	if err := st.UpdatePerson(*got); err != nil {
		t.Fatalf("UpdatePerson() error = %v", err)
	}
	updated, err := st.GetPersonById(id)
	if err != nil {
		t.Fatalf("GetPersonById() after update error = %v", err)
	}
	if updated.Phone != "9007654321" {
		t.Errorf("phone = %q, want %q", updated.Phone, "9007654321")
	}

	if err := st.DeletePerson(id); err != nil {
		t.Fatalf("DeletePerson() error = %v", err)
	}
	if _, err := st.GetPersonById(id); !errors.Is(err, sql.ErrNoRows) {
		t.Errorf("GetPersonById() after delete error = %v, want sql.ErrNoRows", err)
	}
}

func TestPersons_SearchByNameNoMatch(t *testing.T) {
	st := newTestStorage(t)
	mustAddPerson(t, st, Person{Name: "Анна", Surname: "Иванова"})

	found, err := st.GetPersonsByName("сидоров")
	if err != nil {
		t.Fatalf("GetPersonsByName() error = %v", err)
	}
	if len(*found) != 0 {
		t.Errorf("GetPersonsByName() = %+v, want empty", *found)
	}
}

func TestTypes(t *testing.T) {
	st := newTestStorage(t)

	catID := mustCreateType(t, st, "кошка")
	mustCreateType(t, st, "собака")

	byName, err := st.GetTypeByName("кошка")
	if err != nil {
		t.Fatalf("GetTypeByName() error = %v", err)
	}
	if byName.Id != catID {
		t.Errorf("GetTypeByName().Id = %d, want %d", byName.Id, catID)
	}

	byID, err := st.GetTypeById(catID)
	if err != nil {
		t.Fatalf("GetTypeById() error = %v", err)
	}
	if byID.Type != "кошка" {
		t.Errorf("GetTypeById().Type = %q, want %q", byID.Type, "кошка")
	}

	all, err := st.GetTypes()
	if err != nil {
		t.Fatalf("GetTypes() error = %v", err)
	}
	if len(*all) != 2 {
		t.Errorf("GetTypes() = %d types, want 2", len(*all))
	}

	// Колонка type уникальна.
	if _, err := st.CreateType(Types{Type: "кошка"}); err == nil {
		t.Error("CreateType() duplicate error = nil, want UNIQUE constraint error")
	}

	if _, err := st.GetTypeByName("попугай"); !errors.Is(err, sql.ErrNoRows) {
		t.Errorf("GetTypeByName() unknown type error = %v, want sql.ErrNoRows", err)
	}
}

func TestPets_CreateAndRead(t *testing.T) {
	st := newTestStorage(t)
	typeID := mustCreateType(t, st, "кошка")
	ownerID := mustAddPerson(t, st, Person{Name: "Иван", Surname: "Петров"})

	today := time.Now().Format(dateLayout)
	petID := mustAddPet(t, st, Pet{
		Name: "Мурка", Weight: 4.2, Age: 3, AgeWriteDate: today,
		Type: strconv.Itoa(typeID), Allergies: "нет",
	}, ownerID)

	pet, err := st.GetPetById(petID)
	if err != nil {
		t.Fatalf("GetPetById() error = %v", err)
	}
	if pet.Name != "Мурка" || pet.Weight != 4.2 || pet.Age != 3 || pet.Allergies != "нет" {
		t.Errorf("GetPetById() = %+v", pet)
	}
	// Тип возвращается названием из справочника, а не id.
	if pet.Type != "кошка" {
		t.Errorf("GetPetById().Type = %q, want %q", pet.Type, "кошка")
	}

	all, err := st.GetPets()
	if err != nil {
		t.Fatalf("GetPets() error = %v", err)
	}
	if len(*all) != 1 {
		t.Errorf("GetPets() = %d pets, want 1", len(*all))
	}

	byOwner, err := st.GetPetsByOwner(ownerID)
	if err != nil {
		t.Fatalf("GetPetsByOwner() error = %v", err)
	}
	if len(*byOwner) != 1 || (*byOwner)[0].Name != "Мурка" {
		t.Errorf("GetPetsByOwner() = %+v, want [Мурка]", *byOwner)
	}

	byType, err := st.GetPetsByType("кошка")
	if err != nil {
		t.Fatalf("GetPetsByType() error = %v", err)
	}
	if len(*byType) != 1 {
		t.Errorf("GetPetsByType() = %d pets, want 1", len(*byType))
	}
}

// Возраст хранится на дату записи и пересчитывается при чтении.
// Все способы чтения питомца должны отдавать одинаковый возраст.
func TestPets_AgeIsRecalculatedEverywhere(t *testing.T) {
	st := newTestStorage(t)
	typeID := mustCreateType(t, st, "собака")
	ownerID := mustAddPerson(t, st, Person{Name: "Анна"})

	twoYearsAgo := time.Now().AddDate(-2, 0, 0).Format(dateLayout)
	petID := mustAddPet(t, st, Pet{
		Name: "Шарик", Age: 1, AgeWriteDate: twoYearsAgo, Type: strconv.Itoa(typeID),
	}, ownerID)

	const want = 3 // 1 год на дату записи + 2 года прошло

	pet, err := st.GetPetById(petID)
	if err != nil {
		t.Fatalf("GetPetById() error = %v", err)
	}
	if pet.Age != want {
		t.Errorf("GetPetById().Age = %d, want %d", pet.Age, want)
	}

	byOwner, err := st.GetPetsByOwner(ownerID)
	if err != nil {
		t.Fatalf("GetPetsByOwner() error = %v", err)
	}
	if len(*byOwner) != 1 {
		t.Fatalf("GetPetsByOwner() = %d pets, want 1", len(*byOwner))
	}
	if got := (*byOwner)[0].Age; got != want {
		t.Errorf("GetPetsByOwner()[0].Age = %d, want %d (возраст не пересчитан)", got, want)
	}
}

func TestPets_UpdateAndDelete(t *testing.T) {
	st := newTestStorage(t)
	typeID := mustCreateType(t, st, "кошка")
	ownerID := mustAddPerson(t, st, Person{Name: "Иван"})
	today := time.Now().Format(dateLayout)
	petID := mustAddPet(t, st, Pet{Name: "Мурка", AgeWriteDate: today, Type: strconv.Itoa(typeID)}, ownerID)

	// UpdatePet принимает тип по названию и создаёт его, если такого нет.
	err := st.UpdatePet(Pet{Id: petID, Name: "Мурка", Weight: 5, AgeWriteDate: today, Type: "британская кошка"})
	if err != nil {
		t.Fatalf("UpdatePet() error = %v", err)
	}
	pet, err := st.GetPetById(petID)
	if err != nil {
		t.Fatalf("GetPetById() error = %v", err)
	}
	if pet.Weight != 5 || pet.Type != "британская кошка" {
		t.Errorf("after update pet = %+v", pet)
	}
	if _, err := st.GetTypeByName("британская кошка"); err != nil {
		t.Errorf("UpdatePet() did not create new type: %v", err)
	}

	if err := st.DeletePet(petID); err != nil {
		t.Fatalf("DeletePet() error = %v", err)
	}
	if _, err := st.GetPetById(petID); !errors.Is(err, sql.ErrNoRows) {
		t.Errorf("GetPetById() after delete error = %v, want sql.ErrNoRows", err)
	}
}

func TestAppointments(t *testing.T) {
	st := newTestStorage(t)
	typeID := mustCreateType(t, st, "кошка")
	ownerID := mustAddPerson(t, st, Person{Name: "Иван"})
	petID := mustAddPet(t, st, Pet{Name: "Мурка", AgeWriteDate: time.Now().Format(dateLayout), Type: strconv.Itoa(typeID)}, ownerID)

	appt := &Appointment{
		PetId: petID, Date: "20260110", Anamnesis: "не ест",
		Disease: "гастрит", Treatment: "диета", Recommendations: "повторный приём через неделю",
	}
	id, err := st.CreateAppointment(appt)
	if err != nil {
		t.Fatalf("CreateAppointment() error = %v", err)
	}

	got, err := st.GetAppointmentById(id)
	if err != nil {
		t.Fatalf("GetAppointmentById() error = %v", err)
	}
	if got.PetId != petID || got.Disease != "гастрит" || got.Date != "20260110" {
		t.Errorf("GetAppointmentById() = %+v", got)
	}

	if _, err := st.CreateAppointment(&Appointment{PetId: petID, Date: "20260117"}); err != nil {
		t.Fatalf("CreateAppointment() second error = %v", err)
	}
	list, err := st.GetAppointmentsByPetId(petID)
	if err != nil {
		t.Fatalf("GetAppointmentsByPetId() error = %v", err)
	}
	if len(*list) != 2 {
		t.Errorf("GetAppointmentsByPetId() = %d, want 2", len(*list))
	}

	got.Treatment = "антибиотики"
	if err := st.UpdateAppointment(got); err != nil {
		t.Fatalf("UpdateAppointment() error = %v", err)
	}
	updated, _ := st.GetAppointmentById(id)
	if updated.Treatment != "антибиотики" {
		t.Errorf("Treatment = %q, want %q", updated.Treatment, "антибиотики")
	}

	err = st.UpdateAppointment(&Appointment{Id: 9999, PetId: petID})
	if !errors.Is(err, sql.ErrNoRows) {
		t.Errorf("UpdateAppointment() missing id error = %v, want sql.ErrNoRows", err)
	}
}

func TestDocuments(t *testing.T) {
	st := newTestStorage(t)
	typeID := mustCreateType(t, st, "кошка")
	ownerID := mustAddPerson(t, st, Person{Name: "Иван"})
	petID := mustAddPet(t, st, Pet{Name: "Мурка", AgeWriteDate: time.Now().Format(dateLayout), Type: strconv.Itoa(typeID)}, ownerID)
	apptID, err := st.CreateAppointment(&Appointment{PetId: petID, Date: "20260110"})
	if err != nil {
		t.Fatalf("CreateAppointment() error = %v", err)
	}

	content := []byte("%PDF-1.4 fake content")
	docID, err := st.CreateDocument(Document{
		PetId: petID, AppointmentId: apptID, Date: "20260110", Name: "анализ.pdf", Content: content,
	})
	if err != nil {
		t.Fatalf("CreateDocument() error = %v", err)
	}

	doc, err := st.GetDocumentById(docID)
	if err != nil {
		t.Fatalf("GetDocumentById() error = %v", err)
	}
	if doc.Name != "анализ.pdf" || doc.PetId != petID || doc.AppointmentId != apptID {
		t.Errorf("GetDocumentById() = %+v", doc)
	}
	if !bytes.Equal(doc.Content, content) {
		t.Errorf("content = %q, want %q", doc.Content, content)
	}

	byAppt, err := st.GetDocumentsByAppointment(apptID)
	if err != nil {
		t.Fatalf("GetDocumentsByAppointment() error = %v", err)
	}
	if len(byAppt) != 1 || byAppt[0].Id != docID {
		t.Errorf("GetDocumentsByAppointment() = %+v", byAppt)
	}

	byPet, err := st.GetDocumentByPet(petID)
	if err != nil {
		t.Fatalf("GetDocumentByPet(%d) error = %v", petID, err)
	}
	if byPet.Id != docID {
		t.Errorf("GetDocumentByPet().Id = %d, want %d", byPet.Id, docID)
	}

	if err := st.DeleteDocument(docID); err != nil {
		t.Fatalf("DeleteDocument() error = %v", err)
	}
	if _, err := st.GetDocumentById(docID); !errors.Is(err, sql.ErrNoRows) {
		t.Errorf("GetDocumentById() after delete error = %v, want sql.ErrNoRows", err)
	}
}
