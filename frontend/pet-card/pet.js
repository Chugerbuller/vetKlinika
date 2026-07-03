const petCard = document.getElementById("petCard");
const allergyList = document.getElementById("allergyList");
const visitsList = document.getElementById("visitsList");

const modal = document.getElementById("visitModal");
const closeModal = document.getElementById("closeVisitModal");
const addVisitBtn = document.getElementById("addVisitBtn");

// Новые элементы для редактирования и возврата
const backToOwnerBtn = document.getElementById("backToOwnerBtn");
const editAllergyBtn = document.getElementById("editAllergyBtn");
const saveAllergyBtn = document.getElementById("saveAllergyBtn");
const cancelAllergyBtn = document.getElementById("cancelAllergyBtn");
const allergyEditContainer = document.getElementById("allergyEditContainer");
const allergyInput = document.getElementById("allergyInput");

const params = new URLSearchParams(window.location.search);
const petId = params.get("id");
const ownerId = params.get("ownerId");

// Глобальная переменная для хранения данных текущего питомца
let currentPet = null;

closeModal.addEventListener("click", () => {
    modal.classList.add("hidden");
});

addVisitBtn.addEventListener("click", () => {
    window.location.href = `/visit/create?petId=${petId}&ownerId=${ownerId}`;
});

// НОВОЕ: Возврат на страницу владельца
backToOwnerBtn.addEventListener("click", () => {
    if (ownerId != 0) {
        // Предполагается, что в объекте pet есть поле ownerId
        window.location.href = `/person?id=${ownerId}`;
    } else {
        alert("ID владельца не найден");
    }
});

function addField(label, value) {
    const field = document.createElement("div");
    field.classList.add("pet-field");

    const text = document.createElement("strong");
    text.textContent = label + ": ";

    const valueText = document.createElement("span");
    valueText.textContent = value;

    field.appendChild(text);
    field.appendChild(valueText);

    petCard.appendChild(field);
}

// НОВОЕ: Отдельная функция для отрисовки списка аллергий
function renderAllergies(allergies) {
    allergyList.innerHTML = ""; // Очищаем список перед отрисовкой

    if (allergies && allergies.length > 0) {
        allergiesArr = allergies.split(',').map(a => a.trim()).filter(a => a !== ""); // Преобразуем строку в массив
        allergiesArr.forEach(allergy => {
            const item = document.createElement("div");
            item.textContent = "• " + allergy; // Добавил маркер для красоты
            allergyList.appendChild(item);
        });
    } else {
        allergyList.innerHTML = "<span style='color: #888;'>Нет аллергий</span>";
    }
}

function renderPet(pet) {
    petCard.innerHTML = ""; // Очищаем карточку
    addField("Имя", pet.name);
    addField("Тип", pet.type);
    addField("Порода", pet.breed || pet.type); // У вас было pet.type два раза, исправил на breed
    addField("Возраст", pet.age);
    addField("Вес", pet.weight);

    renderAllergies(pet.allergies);
}

// --- ЛОГИКА РЕДАКТИРОВАНИЯ АЛЛЕРГИЙ ---

// Открыть режим редактирования
editAllergyBtn.addEventListener("click", () => {
    allergyList.classList.add("hidden");
    editAllergyBtn.classList.add("hidden");
    allergyEditContainer.classList.remove("hidden");

    // Заполняем текстарею текущими аллергиями (каждая с новой строки)
    allergyInput.value = (currentPet.allergies || []).join('\n');
});

// Отмена редактирования
cancelAllergyBtn.addEventListener("click", () => {
    allergyList.classList.remove("hidden");
    editAllergyBtn.classList.remove("hidden");
    allergyEditContainer.classList.add("hidden");
});

// Сохранение (PUT запрос)
saveAllergyBtn.addEventListener("click", async () => {
    // Получаем текст, разбиваем по переносу строки, удаляем пробелы и пустые строки
    const newAllergies = allergyInput.value
        .split('\n')
        .map(a => a.trim())
        .filter(a => a !== "")
        .join(', ');

    // Создаем копию питомца с новыми аллергиями
    const updatedPet = { ...currentPet, allergies: newAllergies };

    try {
        saveAllergyBtn.disabled = true;
        saveAllergyBtn.textContent = "Сохранение...";

        // Отправляем PUT запрос для обновления питомца
        const response = await fetch(`/api/pets`, { // <-- Если ваш API требует /api/pets?id=... измените URL
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(updatedPet)
        });

        if (response.ok) {
            currentPet = updatedPet; // Обновляем локальный объект
            renderAllergies(currentPet.allergies); // Перерисовываем список

            // Возвращаем UI в исходное состояние
            allergyList.classList.remove("hidden");
            editAllergyBtn.classList.remove("hidden");
            allergyEditContainer.classList.add("hidden");
        } else {
            console.log("Не удалось сохранить данные.");
        }
    } catch (error) {
        console.error("Ошибка при сохранении аллергий:", error);
    } finally {
        saveAllergyBtn.disabled = false;
        saveAllergyBtn.textContent = "Сохранить";
    }
});

// --------------------------------------

function createVisitCard(visit) {
    const card = document.createElement("div");
    card.classList.add("visit-card");

    const diseaseTitle = document.createElement("div");
    diseaseTitle.classList.add("visit-disease");
    diseaseTitle.textContent = visit.disease || "Неизвестное заболевание";

    const rawDate = visit.date;
    const formattedStr = `${rawDate.slice(0, 4)}-${rawDate.slice(4, 6)}-${rawDate.slice(6, 8)}`;
    const formattedDate = new Date(formattedStr).toLocaleDateString("ru-RU", {
        day: '2-digit',
        month: 'long',
        year: 'numeric'
    });

    const dateText = document.createElement("div");
    dateText.classList.add("visit-date");
    dateText.textContent = formattedDate;

    card.appendChild(diseaseTitle);
    card.appendChild(dateText);

    card.addEventListener("click", () => {
        window.location.href = `/visit?id=${visit.id}&petId=${petId}&ownerId=${ownerId}`;
    });

    return card;
}

function renderVisits(visits) {
    visitsList.replaceChildren();
    visits.forEach(v => {
        const card = createVisitCard(v);
        visitsList.appendChild(card);
    });
}

async function loadPet() {
    const response = await fetch("/api/pets?id=" + petId);
    currentPet = await response.json(); // Сохраняем в глобальную переменную
    renderPet(currentPet);
}

async function loadVisits() {
    const response = await fetch("/api/appointments?petId=" + petId);
    const visits = await response.json();
    renderVisits(visits);
}

loadPet();
loadVisits();