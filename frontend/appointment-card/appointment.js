const params = new URLSearchParams(window.location.search)
const petId = params.get("petId")
const visitId = params.get("id")
const ownerId = params.get("ownerId")
const backBtn = document.getElementById("backBtn")
const petInfoBlock = document.getElementById("petInfoBlock")
const petAllergies = document.getElementById("petAllergies")
const visitDate = document.getElementById("visitDate")
const visitAnamnesis = document.getElementById("visitAnamnesis")
const visitDisease = document.getElementById("visitDisease")
const visitTreatment = document.getElementById("visitTreatment")
const visitRecommendations = document.getElementById("visitRecommendations")
const visitFiles = document.getElementById("visitFiles") // НОВАЯ ПЕРЕМЕННАЯ

const fileInput = document.getElementById("fileInput")
const uploadBtn = document.getElementById("uploadBtn")
const dropZone = document.getElementById("dropZone");

backBtn.addEventListener("click", () => {
    console.log(`Переход к питомцу с id=${petId} и владельцем id=${ownerId}`)
    window.location.href = `http://localhost:8080/pets?id=${petId}&ownerId=${ownerId}`
})

function addPetField(label, value) {
    const field = document.createElement("div")
    field.classList.add("pet-field")
    const strongLabel = document.createElement("strong")
    strongLabel.textContent = label + ": "

    const spanValue = document.createElement("span")
    spanValue.textContent = value || "—"

    field.appendChild(strongLabel)
    field.appendChild(spanValue)

    petInfoBlock.appendChild(field)
}

async function loadPet() {
    try {
        const response = await fetch("http://localhost:8080/api/pets?id=" + petId)
        const pet = await response.json()
        addPetField("Имя", pet.name)
        addPetField("Тип", pet.type)
        addPetField("Возраст", pet.age)
        addPetField("Вес", pet.weight + " кг")

        if (pet.allergies && pet.allergies.length > 0) {
            let allergiesArr = pet.allergies.split(',').map(a => a.trim()).filter(a => a !== "");
            allergiesArr.forEach(allergy => {
                const item = document.createElement("div")
                item.classList.add("allergy-item")
                item.textContent = "• " + allergy
                petAllergies.appendChild(item)
            })
        } else {
            petAllergies.textContent = "Аллергии не выявлены"
        }
    } catch (error) {
        console.error("Ошибка загрузки питомца", error)
    }
}

async function loadVisit() {
    try {
        const response = await fetch("http://localhost:8080/api/appointments?id=" + visitId)
        const visit = await response.json()
        const rawDate = visit.date; // "20260331"

        // Превращаем "20260331" в "2026-03-31"
        const formattedStr = `${rawDate.slice(0, 4)}-${rawDate.slice(4, 6)}-${rawDate.slice(6, 8)}`;

        visitDate.textContent = "Прием от: " + new Date(formattedStr).toLocaleDateString("ru-RU", {
            day: '2-digit',
            month: 'long',
            year: 'numeric'
        });

        visitAnamnesis.textContent = visit.anamnesis || "Нет данных"
        visitDisease.textContent = visit.disease || "Не установлено"
        visitTreatment.textContent = visit.treatment || "Не назначено"
        visitRecommendations.textContent = visit.recommendations || "Нет рекомендаций"

        const fileResponse = await fetch(`http://localhost:8080/api/documents?appointmentId=${visitId}`)
        const files = await fileResponse.json();

        // ОТОБРАЖЕНИЕ ФАЙЛОВ
        // Предполагается, что сервер возвращает массив файлов: visit.files = [{id: 1, name: "uzi.pdf"}, ...]
        if (files || files.length === 0) {
            visitFiles.innerHTML = ""; // Очищаем текст "Нет файлов"
            files.forEach(file => {
                const btnGroup = document.createElement("div");
                btnGroup.className = "btn-group me-2 mb-2";
                btnGroup.setAttribute("role", "group");

                // Кнопка скачивания
                const downloadBtn = document.createElement("button");
                downloadBtn.className = "btn btn-outline-primary btn-sm file-btn";
                downloadBtn.innerHTML = `📄 ${file.name}`;
                downloadBtn.title = "Скачать файл";
                downloadBtn.addEventListener("click", () => downloadFile(file.id, file.name, downloadBtn));

                // Кнопка удаления (Красный крестик)
                const deleteBtn = document.createElement("button");
                deleteBtn.className = "btn btn-outline-danger btn-sm";
                deleteBtn.innerHTML = `✖`;
                deleteBtn.title = "Удалить файл";
                deleteBtn.addEventListener("click", () => deleteDocument(file.id, file.name));

                // Добавляем кнопки в группу, а группу в контейнер
                btnGroup.appendChild(downloadBtn);
                btnGroup.appendChild(deleteBtn);
                visitFiles.appendChild(btnGroup);
            });
        }
    } catch (error) {
        console.error("Ошибка загрузки приема", error)
    }
}

// НОВАЯ ФУНКЦИЯ ДЛЯ СКАЧИВАНИЯ ФАЙЛОВ В ФОРМАТЕ BLOB
async function downloadFile(fileId, fileName, buttonElement) {
    const originalText = buttonElement.innerHTML;
    buttonElement.innerHTML = `⏳ Загрузка...`;
    buttonElement.disabled = true;

    try {
        // Укажите правильный URL вашего API для скачивания файлов
        const response = await fetch(`http://localhost:8080/api/documents?id=${fileId}`);

        if (!response.ok) throw new Error("Ошибка при скачивании файла со стороны сервера");

        // Получаем файл как Blob
        const blob = await response.blob();

        // Создаем временную ссылку на скачивание
        const url = window.URL.createObjectURL(blob);
        const a = document.createElement("a");
        a.style.display = "none";
        a.href = url;
        a.download = fileName; // Имя, с которым файл сохранится у пользователя
        document.body.appendChild(a);

        // Эмулируем клик
        a.click();

        // Подчищаем память и DOM
        window.URL.revokeObjectURL(url);
        document.body.removeChild(a);
    } catch (error) {
        console.error("Ошибка скачивания:", error);
    } finally {
        // Возвращаем кнопке исходный вид
        buttonElement.innerHTML = originalText;
        buttonElement.disabled = false;
    }
}
// СКАЧИВАНИЕ ФАЙЛА (предполагаем, что эндпоинт отдает сырые байты - Blob)
async function downloadFile(fileId, fileName, buttonElement) {
    const originalText = buttonElement.innerHTML;
    buttonElement.innerHTML = `⏳ Загрузка...`;
    buttonElement.disabled = true;

    try {
        // 1. Делаем GET запрос на получение документа (возвращает JSON)
        const response = await fetch(`http://localhost:8080/api/documents?id=${fileId}`);
        if (!response.ok) throw new Error("Ошибка при скачивании файла");

        const doc = await response.json();

        // 2. Декодируем Base64 строку (doc.content) обратно в бинарные данные
        // atob расшифровывает Base64 в строку бинарных символов
        const byteCharacters = atob(doc.content);

        // Создаем массив чисел (байт)
        const byteNumbers = new Array(byteCharacters.length);
        for (let i = 0; i < byteCharacters.length; i++) {
            byteNumbers[i] = byteCharacters.charCodeAt(i);
        }

        // Оборачиваем в типизированный массив
        const byteArray = new Uint8Array(byteNumbers);

        // 3. Создаем Blob. Тип 'application/octet-stream' заставит браузер скачать файл
        const blob = new Blob([byteArray], { type: 'application/octet-stream' });
        const url = window.URL.createObjectURL(blob);

        // 4. Создаем невидимую ссылку и эмулируем клик
        const a = document.createElement("a");
        a.style.display = "none";
        a.href = url;
        // Используем имя файла из ответа сервера (или переданное в функцию, если на сервере пусто)
        a.download = doc.name || fileName;
        document.body.appendChild(a);

        a.click();

        // 5. Очищаем память
        window.URL.revokeObjectURL(url);
        document.body.removeChild(a);
    } catch (error) {
        console.error("Ошибка скачивания:", error);
        alert("Не удалось скачать файл: " + fileName);
    } finally {
        // Возвращаем кнопку в исходное состояние
        buttonElement.innerHTML = originalText;
        buttonElement.disabled = false;
    }
}

// Вспомогательная функция: читает файл и возвращает строку Base64 (без префикса data:...)
function getBase64(file) {
    return new Promise((resolve, reject) => {
        const reader = new FileReader();
        reader.readAsDataURL(file);
        reader.onload = () => {
            // reader.result выглядит так: "data:application/pdf;base64,JVBERi0xLjQK..."
            // Нам нужна только часть после запятой
            const base64String = reader.result.split(',')[1];
            resolve(base64String);
        };
        reader.onerror = error => reject(error);
    });
}

// ОБРАБОТЧИК ЗАГРУЗКИ НОВОГО ФАЙЛА
uploadBtn.addEventListener("click", async () => {
    const file = fileInput.files[0];

    if (!file) {
        return;
    }

    try {
        uploadBtn.disabled = true;
        uploadBtn.textContent = "Отправка...";
        fileInput.disabled = true;

        // 1. Получаем содержимое файла в виде строки Base64
        const base64Content = await getBase64(file);

        // 2. Формируем JSON-объект в точности под структуру Go:
        /*
            type Document struct {
                Id            int    `json:"id"`             // пропускаем, сгенерирует БД
                PetId         int    `json:"pet_id"`
                Date          string `json:"date"`
                AppointmentId int    `json:"appointment_id"`
                Name          string `json:"name"`
                Content       []byte `json:"content"`        // Go распарсит Base64 строку в []byte
            }
        */
        const payload = {
            pet_id: parseInt(petId, 10),
            appointment_id: parseInt(visitId, 10),
            date: new Date().toISOString(), // Передаем текущую дату в формате ISO 
            name: file.name,
            content: base64Content
        };

        // 3. Отправляем JSON на сервер
        const response = await fetch("http://localhost:8080/api/documents", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify(payload)
        });

        if (!response.ok) {
            throw new Error("Ошибка сервера при загрузке документа");
        }

        fileInput.value = ""; // Очищаем поле файла

        // Перезагружаем прием, чтобы файл появился в списке
        await loadVisit();

    } catch (error) {
        console.error("Ошибка при отправке файла:", error);
    } finally {
        uploadBtn.disabled = false;
        uploadBtn.textContent = "Загрузить";
        fileInput.disabled = false;
    }
});
['dragenter', 'dragover', 'dragleave', 'drop'].forEach(eventName => {
    dropZone.addEventListener(eventName, preventDefaults, false);
    document.body.addEventListener(eventName, preventDefaults, false);
});

function preventDefaults(e) {
    e.preventDefault();
    e.stopPropagation();
}

// Добавляем и убираем визуальную подсветку зоны при наведении
['dragenter', 'dragover'].forEach(eventName => {
    dropZone.addEventListener(eventName, () => dropZone.classList.add('dragover'), false);
});

['dragleave', 'drop'].forEach(eventName => {
    dropZone.addEventListener(eventName, () => dropZone.classList.remove('dragover'), false);
});

// Обрабатываем само "сбрасывание" файла
dropZone.addEventListener('drop', handleDrop, false);

function handleDrop(e) {
    const dt = e.dataTransfer;
    const files = dt.files;

    if (files && files.length > 0) {
        // Программно помещаем перетащенный файл в наш стандартный input
        fileInput.files = files;

        // Опционально: можно раскомментировать строку ниже, чтобы файл загружался автоматически сразу после перетаскивания без нажатия кнопки "Загрузить"
        // uploadBtn.click(); 
    }
}
async function deleteDocument(fileId, fileName) {
    // Спрашиваем подтверждение перед удалением
    const isConfirmed = confirm(`Вы уверены, что хотите безвозвратно удалить файл "${fileName}"?`);
    if (!isConfirmed) return;

    try {
        // Отправляем DELETE запрос (предполагаем, что ID передается в query, как и при GET)
        const response = await fetch(`http://localhost:8080/api/documents?id=${fileId}`, {
            method: "DELETE"
        });

        if (!response.ok) {
            throw new Error(`Ошибка сервера при удалении документа (статус: ${response.status})`);
        }

        // Если удаление прошло успешно — просто перезагружаем данные приема, 
        // чтобы файл исчез из списка на экране
        await loadVisit();

    } catch (error) {
        console.error("Ошибка при удалении файла:", error);
        alert("Произошла ошибка при удалении файла. Попробуйте еще раз.");
    }
}
if (petId) loadPet()
if (visitId) loadVisit()