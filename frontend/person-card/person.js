const personCard = document.getElementById("personCard")
const animalsList = document.getElementById("animalsList")
const backBtn = document.getElementById("backBtn")
const newTypeBlock = document.getElementById("newTypeBlock")
const params = new URLSearchParams(window.location.search)
const personId = params.get("id")


backBtn.addEventListener("click", () => {

    window.location.href = "/"

})


function createPersonCard(person){

    const name = document.createElement("div")
    name.classList.add("person-name")
    name.textContent =
        person.surname + " " +
        person.name + " " +
        (person.patronymic || "")

    const phone = document.createElement("div")
    phone.classList.add("person-meta")
    phone.textContent = "📞 " + person.phone

    const address = document.createElement("div")
    address.classList.add("person-meta")
    address.textContent = "📍 " + person.address

    personCard.appendChild(name)
    personCard.appendChild(phone)
    personCard.appendChild(address)

}


function createAnimalCard(animal){

    const card = document.createElement("div")
    card.classList.add("animal-card")

    const name = document.createElement("div")
    name.textContent = animal.name

    card.appendChild(name)

    return card
}


function renderAnimals(animals){

    animalsList.textContent = ""

    animals.forEach(animal => {

        const card = createAnimalCard(animal)

        animalsList.appendChild(card)

    })

}


async function loadPerson(){

    const response =
        await fetch("http://localhost:8080/api/persons?id=" + personId)

    const person = await response.json()

    createPersonCard(person)

}


async function loadAnimals(){

    const response =
        await fetch(
            "http://localhost:8080/api/pets?ownerId=" + personId
        )

    const animals = await response.json()

    renderAnimals(animals)

}

const addAnimalBtn = document.getElementById("addAnimalBtn")
const modal = document.getElementById("animalModal")

const closeBtn = document.getElementById("closeAnimalBtn")
const saveBtn = document.getElementById("saveAnimalBtn")

const typeSelect = document.getElementById("animalType")
const newTypeInput = document.getElementById("newTypeInput")

const nameInput = document.getElementById("animalName")

let types = []


addAnimalBtn.addEventListener("click", () => {

    modal.classList.remove("hidden")

})


closeBtn.addEventListener("click", () => {

    modal.classList.add("hidden")

    resetAnimalForm()

})

function resetAnimalForm(){

    nameInput.value = ""
    document.getElementById("animalWeight").value = ""
    document.getElementById("animalAge").value = ""

    newTypeInput.value = ""

    typeSelect.selectedIndex = 0

    newTypeBlock.classList.add("hidden")

}
typeSelect.addEventListener("change", () => {

    if(typeSelect.value === "other"){
        newTypeBlock.classList.remove("hidden")
    }else{
        newTypeBlock.classList.add("hidden")
    }

})

async function deleteAnimal(animalId) {
    if (!confirm("Вы уверены, что хотите удалить запись о животном?")) {
        return;
    }

    try {
        const response = await fetch(`/api/pets?id=${animalId}`, {
            method: "DELETE"
        });

        if (!response.ok) {
            throw new Error("Не удалось удалить животное");
        }

        // Обновляем список животных после удаления
        loadAnimals();

    } catch (error) {
        console.error("Ошибка при удалении:", error);
        alert("Произошла ошибка при удалении");
    }
}

// Обновленная функция создания карточки (без innerHTML)
function createAnimalCard(animal) {
    const card = document.createElement("div");
    card.classList.add("animal-card");

    // Контейнер для текста (чтобы кнопка была сбоку)
    const infoWrapper = document.createElement("div");
    infoWrapper.classList.add("animal-info");

    const name = document.createElement("div");
    name.classList.add("animal-name-label");
    name.textContent = animal.name;

    infoWrapper.appendChild(name);

    // Кнопка удаления
    const deleteBtn = document.createElement("button");
    deleteBtn.classList.add("delete-cross");
    deleteBtn.textContent = "\u00d7"; // Специальный символ крестика
    deleteBtn.title = "Удалить";

    // Обработчик удаления
    deleteBtn.addEventListener("click", (e) => {
        e.stopPropagation(); // Чтобы не сработал переход на страницу животного
        deleteAnimal(animal.id);
    });

    card.appendChild(infoWrapper);
    card.appendChild(deleteBtn);

    // Клик по самой карточке для перехода
    card.addEventListener("click", () => {
        window.location.href = "/pets?id=" + animal.id + "&ownerId=" + personId;
    });

    return card;
}
async function loadTypes(){

    typeSelect.textContent = ""

    try{

        const response = await fetch("/api/types")
        const data = await response.json()

        types = data || []

        if(types.length > 0){

            types.forEach(type => {

                const option = document.createElement("option")

                option.value = type.id
                option.textContent = type.type

                typeSelect.appendChild(option)

            })

        }

    }catch(error){

        console.error("Ошибка загрузки типов", error)

    }

    const otherOption = document.createElement("option")
    otherOption.value = "other"
    otherOption.textContent = "Другое"

    typeSelect.appendChild(otherOption)
    if(types.length === 0){
        typeSelect.value = "other"
        newTypeBlock.classList.remove("hidden")
    }
}

async function createAnimal() {
    const submitButton = saveBtn; // Предположим, у вас есть кнопка
    
    try {
        // 1. Блокируем кнопку, чтобы избежать повторных нажатий (Race Condition в UI)
        submitButton.disabled = true;
        submitButton.textContent = "Сохранение...";

        let typeValue = typeSelect.value;

        // 2. Если выбран "Другой", сначала создаем тип
        if (typeValue === "other") {
            const newType = {
                type: newTypeInput.value
            };

            const typeResponse = await fetch("/api/types", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify(newType)
            });

            // ОБЯЗАТЕЛЬНО проверяем, прошел ли запрос успешно
            if (!typeResponse.ok) {
                throw new Error("Не удалось создать новый тип животного");
            }

            const createdType = await typeResponse.json();
            
            // Убедитесь, что поле называется именно .value или .id, как в вашем Go-структе
            typeValue = createdType.value || createdType.id; 
        }

        // 3. Теперь создаем само животное
        const pet = {
            name: nameInput.value,
            weight: parseFloat(document.getElementById("animalWeight").value),
            age: parseInt(document.getElementById("animalAge").value),
            age_write_date: new Date().toISOString(),
            type: parseInt(typeValue, 10), // Используем полученный ID
        };

        const petResponse = await fetch("/api/pets?ownerId=" + personId, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(pet)
        });

        if (!petResponse.ok) {
            throw new Error("Не удалось сохранить данные о животном");
        }

        // 4. Если всё успешно — уходим на страницу карточки или обновляем список
        window.location.href = "/person?id=" + personId;

    } catch (error) {
        // Обработка ошибок
        console.error("Ошибка:", error);
        alert("Произошла ошибка: " + error.message);
    } finally {
        // 5. Разблокируем кнопку в любом случае (если не ушли со страницы)
        submitButton.disabled = false;
        submitButton.textContent = "Создать";
    }
    modal.classList.add("hidden")
    resetAnimalForm()
    loadTypes()
    loadAnimals()
}

saveBtn.addEventListener("click", createAnimal)

loadPerson()
loadAnimals()

loadTypes()