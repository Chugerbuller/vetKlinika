const personInfo = document.getElementById("personInfo")
const editPersonBtn = document.getElementById("editPersonBtn")
const animalsList = document.getElementById("animalsList")
const backBtn = document.getElementById("backBtn")
const newTypeBlock = document.getElementById("newTypeBlock")
const params = new URLSearchParams(window.location.search)
const personId = params.get("id")


// ---------- Элементы модалки питомца ----------

const addAnimalBtn = document.getElementById("addAnimalBtn")
const modal = document.getElementById("animalModal")

const closeBtn = document.getElementById("closeAnimalBtn")
const saveBtn = document.getElementById("saveAnimalBtn")

const typeSelect = document.getElementById("animalType")
const newTypeInput = document.getElementById("newTypeInput")

const nameInput = document.getElementById("animalName")

let types = []


// ---------- Элементы модалки клиента ----------

let currentPerson = null

const editPersonModal = document.getElementById("editPersonModal")
const closePersonModalBtn = document.getElementById("closePersonModalBtn")
const savePersonBtn = document.getElementById("savePersonBtn")

const editSurnameInput = document.getElementById("editSurname")
const editNameInput = document.getElementById("editName")
const editPatronymicInput = document.getElementById("editPatronymic")
const editPhoneInput = document.getElementById("editPhone")
const editAddressInput = document.getElementById("editAddress")


// ---------- Навигация ----------

backBtn.addEventListener("click", () => {

    window.location.href = "/"

})


// ---------- Клиент ----------

function createPersonCard(person) {

    personInfo.textContent = ""

    const name = document.createElement("div")
    name.classList.add("person-name")
    name.textContent = [person.surname, person.name, person.patronymic]
        .filter(Boolean)
        .join(" ")

    const phone = document.createElement("div")
    phone.classList.add("person-meta")
    phone.textContent = "📞 " + (person.phone || "не указан")

    const address = document.createElement("div")
    address.classList.add("person-meta")
    address.textContent = "📍 " + (person.address || "не указан")

    personInfo.appendChild(name)
    personInfo.appendChild(phone)
    personInfo.appendChild(address)
}


async function loadPerson() {

    try {
        const response = await fetch("/api/persons?id=" + personId)

        if (!response.ok) {
            throw new Error("Не удалось получить данные о клиенте")
        }

        currentPerson = await response.json()
        createPersonCard(currentPerson)

    } catch (error) {
        console.error(error)
        personInfo.textContent = "Не удалось загрузить данные клиента"
    }
}


function openEditPersonModal() {

    if (!currentPerson) return

    editSurnameInput.value = currentPerson.surname || ""
    editNameInput.value = currentPerson.name || ""
    editPatronymicInput.value = currentPerson.patronymic || ""
    editPhoneInput.value = currentPerson.phone || ""
    editAddressInput.value = currentPerson.address || ""

    editPersonModal.classList.remove("hidden")
}


function closeEditPersonModal() {
    editPersonModal.classList.add("hidden")
}


async function updatePerson() {

    try {
        savePersonBtn.disabled = true
        savePersonBtn.textContent = "Сохранение..."

        const payload = {
            id: parseInt(personId, 10),
            name: editNameInput.value.trim(),
            surname: editSurnameInput.value.trim(),
            patronymic: editPatronymicInput.value.trim(),
            phone: editPhoneInput.value.trim(),
            address: editAddressInput.value.trim()
        }

        const response = await fetch("/api/persons", {
            method: "PUT",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(payload)
        })

        if (!response.ok) {
            throw new Error("Ошибка при обновлении клиента")
        }

        closeEditPersonModal()
        await loadPerson()

    } catch (error) {
        console.error("Ошибка:", error)
        alert("Произошла ошибка: " + error.message)
    } finally {
        savePersonBtn.disabled = false
        savePersonBtn.textContent = "Сохранить"
    }
}


editPersonBtn.addEventListener("click", openEditPersonModal)
closePersonModalBtn.addEventListener("click", closeEditPersonModal)
savePersonBtn.addEventListener("click", updatePerson)


// ---------- Животные ----------

function createAnimalCard(animal) {

    const card = document.createElement("div")
    card.classList.add("animal-card")

    // Контейнер для текста (чтобы кнопка была сбоку)
    const infoWrapper = document.createElement("div")
    infoWrapper.classList.add("animal-info")

    const name = document.createElement("div")
    name.classList.add("animal-name-label")
    name.textContent = animal.name

    infoWrapper.appendChild(name)

    // Кнопка удаления
    const deleteBtn = document.createElement("button")
    deleteBtn.classList.add("delete-cross")
    deleteBtn.textContent = "\u00d7"
    deleteBtn.title = "Удалить"

    deleteBtn.addEventListener("click", (e) => {
        e.stopPropagation() // Чтобы не сработал переход на страницу животного
        deleteAnimal(animal.id)
    })

    card.appendChild(infoWrapper)
    card.appendChild(deleteBtn)

    // Клик по карточке — переход на страницу животного
    card.addEventListener("click", () => {
        window.location.href = "/pets?id=" + animal.id + "&ownerId=" + personId
    })

    return card
}


function renderAnimals(animals) {

    animalsList.textContent = ""

    animals.forEach(animal => {

        const card = createAnimalCard(animal)

        animalsList.appendChild(card)

    })

}


async function loadAnimals() {

    const response = await fetch("/api/pets?ownerId=" + personId)

    const animals = await response.json()

    renderAnimals(animals)

}


async function deleteAnimal(animalId) {

    if (!confirm("Вы уверены, что хотите удалить запись о животном?")) {
        return
    }

    try {
        const response = await fetch(`/api/pets?id=${animalId}`, {
            method: "DELETE"
        })

        if (!response.ok) {
            throw new Error("Не удалось удалить животное")
        }

        // Обновляем список животных после удаления
        loadAnimals()

    } catch (error) {
        console.error("Ошибка при удалении:", error)
        alert("Произошла ошибка при удалении")
    }
}


// ---------- Модалка питомца ----------

addAnimalBtn.addEventListener("click", () => {

    modal.classList.remove("hidden")

})


closeBtn.addEventListener("click", () => {

    modal.classList.add("hidden")

    resetAnimalForm()

})


function resetAnimalForm() {

    nameInput.value = ""
    document.getElementById("animalWeight").value = ""
    document.getElementById("animalAge").value = ""

    newTypeInput.value = ""

    typeSelect.selectedIndex = 0

    newTypeBlock.classList.add("hidden")

}


typeSelect.addEventListener("change", () => {

    if (typeSelect.value === "other") {
        newTypeBlock.classList.remove("hidden")
    } else {
        newTypeBlock.classList.add("hidden")
    }

})


async function loadTypes() {

    typeSelect.textContent = ""

    try {

        const response = await fetch("/api/types")
        const data = await response.json()

        types = data || []

        if (types.length > 0) {

            types.forEach(type => {

                const option = document.createElement("option")

                option.value = type.id
                option.textContent = type.type

                typeSelect.appendChild(option)

            })

        }

    } catch (error) {

        console.error("Ошибка загрузки типов", error)

    }

    const otherOption = document.createElement("option")
    otherOption.value = "other"
    otherOption.textContent = "Другое"

    typeSelect.appendChild(otherOption)

    if (types.length === 0) {
        typeSelect.value = "other"
        newTypeBlock.classList.remove("hidden")
    }
}


async function createAnimal() {

    const submitButton = saveBtn

    try {
        // Блокируем кнопку, чтобы избежать повторных нажатий
        submitButton.disabled = true
        submitButton.textContent = "Сохранение..."

        let typeValue = typeSelect.value

        // Если выбран "Другое", сначала создаем тип
        if (typeValue === "other") {
            const newType = {
                type: newTypeInput.value
            }

            const typeResponse = await fetch("/api/types", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify(newType)
            })

            if (!typeResponse.ok) {
                throw new Error("Не удалось создать новый тип животного")
            }

            const createdType = await typeResponse.json()

            typeValue = createdType.value || createdType.id
        }

        // Создаем само животное
        const pet = {
            name: nameInput.value,
            weight: parseFloat(document.getElementById("animalWeight").value),
            age: parseInt(document.getElementById("animalAge").value),
            age_write_date: new Date().toISOString(),
            type: parseInt(typeValue, 10),
        }

        const petResponse = await fetch("/api/pets?ownerId=" + personId, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(pet)
        })

        if (!petResponse.ok) {
            throw new Error("Не удалось сохранить данные о животном")
        }

        // Успех: закрываем окно и обновляем данные
        modal.classList.add("hidden")
        resetAnimalForm()
        loadTypes()
        loadAnimals()

    } catch (error) {
        console.error("Ошибка:", error)
        alert("Произошла ошибка: " + error.message)
    } finally {
        submitButton.disabled = false
        submitButton.textContent = "Сохранить запись"
    }
}


saveBtn.addEventListener("click", createAnimal)


// ---------- Старт ----------

loadPerson()
loadAnimals()
loadTypes()
