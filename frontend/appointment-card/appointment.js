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
const visitFiles = document.getElementById("visitFiles")

const fileInput = document.getElementById("fileInput")
const uploadBtn = document.getElementById("uploadBtn")
const dropZone = document.getElementById("dropZone")


// ---------- Редактирование приема: элементы ----------

let currentVisit = null

const editVisitBtn = document.getElementById("editVisitBtn")
const editVisitModal = document.getElementById("editVisitModal")
const saveVisitBtn = document.getElementById("saveVisitBtn")
const closeVisitModalBtn = document.getElementById("closeVisitModalBtn")

const editVisitDate = document.getElementById("editVisitDate")
const editVisitAnamnesis = document.getElementById("editVisitAnamnesis")
const editVisitDisease = document.getElementById("editVisitDisease")
const editVisitTreatment = document.getElementById("editVisitTreatment")
const editVisitRecommendations = document.getElementById("editVisitRecommendations")


// ---------- Навигация ----------

backBtn.addEventListener("click", () => {
    console.log(`Переход к питомцу с id=${petId} и владельцем id=${ownerId}`)
    window.location.href = `http://localhost:8080/pets?id=${petId}&ownerId=${ownerId}`
})


// ---------- Вспомогательные функции для даты ----------

// "20260331" -> Date
function parseVisitDate(raw) {
    const year = parseInt(raw.slice(0, 4), 10)
    const month = parseInt(raw.slice(4, 6), 10)
    const day = parseInt(raw.slice(6, 8), 10)
    return new Date(year, month - 1, day)
}

// "20260331" -> "2026-03-31" (для <input type="date">)
function toInputDate(raw) {
    return `${raw.slice(0, 4)}-${raw.slice(4, 6)}-${raw.slice(6, 8)}`
}

// "2026-03-31" -> "20260331" (формат, в котором дата хранится на сервере)
function fromInputDate(value) {
    return value.replaceAll("-", "")
}


// ---------- Питомец ----------

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

        petInfoBlock.textContent = ""
        petAllergies.textContent = ""

        addPetField("Имя", pet.name)
        addPetField("Тип", pet.type)
        addPetField("Возраст", pet.age)
        addPetField("Вес", pet.weight + " кг")

        if (pet.allergies && pet.allergies.length > 0) {
            const allergiesArr = pet.allergies
                .split(",")
                .map(a => a.trim())
                .filter(a => a !== "")

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


// ---------- Прием ----------

function renderVisit(visit) {
    visitDate.textContent = "Прием от: " + parseVisitDate(visit.date).toLocaleDateString("ru-RU", {
        day: "2-digit",
        month: "long",
        year: "numeric"
    })

    visitAnamnesis.textContent = visit.anamnesis || "Нет данных"
    visitDisease.textContent = visit.disease || "Не установлено"
    visitTreatment.textContent = visit.treatment || "Не назначено"
    visitRecommendations.textContent = visit.recommendations || "Нет рекомендаций"
}

function renderFiles(files) {
    visitFiles.textContent = ""

    if (!files || files.length === 0) {
        const empty = document.createElement("span")
        empty.classList.add("text-content", "text-muted")
        empty.textContent = "Нет прикрепленных файлов"
        visitFiles.appendChild(empty)
        return
    }

    files.forEach(file => {
        const btnGroup = document.createElement("div")
        btnGroup.className = "btn-group me-2 mb-2"
        btnGroup.setAttribute("role", "group")

        // Кнопка скачивания
        const downloadBtn = document.createElement("button")
        downloadBtn.className = "btn btn-outline-primary btn-sm file-btn"
        downloadBtn.textContent = "📄 " + file.name
        downloadBtn.title = "Скачать файл"
        downloadBtn.addEventListener("click", () => downloadFile(file.id, file.name, downloadBtn))

        // Кнопка удаления
        const deleteBtn = document.createElement("button")
        deleteBtn.className = "btn btn-outline-danger btn-sm"
        deleteBtn.textContent = "✖"
        deleteBtn.title = "Удалить файл"
        deleteBtn.addEventListener("click", () => deleteDocument(file.id, file.name))

        btnGroup.appendChild(downloadBtn)
        btnGroup.appendChild(deleteBtn)
        visitFiles.appendChild(btnGroup)
    })
}

async function loadVisit() {
    try {
        const response = await fetch("http://localhost:8080/api/appointments?id=" + visitId)
        if (!response.ok) {
            throw new Error("Не удалось получить данные приема")
        }

        currentVisit = await response.json()
        renderVisit(currentVisit)

        const fileResponse = await fetch(`http://localhost:8080/api/documents?appointmentId=${visitId}`)
        const files = await fileResponse.json()
        renderFiles(files)

    } catch (error) {
        console.error("Ошибка загрузки приема", error)
    }
}


// ---------- Редактирование приема ----------

function openEditVisitModal() {
    if (!currentVisit) return

    editVisitDate.value = toInputDate(currentVisit.date)
    editVisitAnamnesis.value = currentVisit.anamnesis || ""
    editVisitDisease.value = currentVisit.disease || ""
    editVisitTreatment.value = currentVisit.treatment || ""
    editVisitRecommendations.value = currentVisit.recommendations || ""

    editVisitModal.classList.remove("hidden")
}

function closeEditVisitModal() {
    editVisitModal.classList.add("hidden")
}

async function updateVisit() {
    if (!editVisitDate.value) {
        alert("Укажите дату приема")
        return
    }

    try {
        saveVisitBtn.disabled = true
        saveVisitBtn.textContent = "Сохранение..."

        // Берем весь загруженный объект (чтобы не потерять pet_id и другие поля)
        // и подменяем только то, что редактируется в форме
        const payload = {
            ...currentVisit,
            id: parseInt(visitId, 10),
            date: fromInputDate(editVisitDate.value),
            anamnesis: editVisitAnamnesis.value.trim(),
            disease: editVisitDisease.value.trim(),
            treatment: editVisitTreatment.value.trim(),
            recommendations: editVisitRecommendations.value.trim()
        }

        const response = await fetch("http://localhost:8080/api/appointments", {
            method: "PUT",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(payload)
        })

        if (!response.ok) {
            throw new Error("Ошибка при обновлении приема (статус: " + response.status + ")")
        }

        closeEditVisitModal()
        await loadVisit()

    } catch (error) {
        console.error("Ошибка обновления приема:", error)
        alert("Произошла ошибка: " + error.message)
    } finally {
        saveVisitBtn.disabled = false
        saveVisitBtn.textContent = "Сохранить"
    }
}

editVisitBtn.addEventListener("click", openEditVisitModal)
closeVisitModalBtn.addEventListener("click", closeEditVisitModal)
saveVisitBtn.addEventListener("click", updateVisit)

// Закрытие по клику на затемненный фон
editVisitModal.addEventListener("click", (e) => {
    if (e.target === editVisitModal) closeEditVisitModal()
})

// Закрытие по Escape
document.addEventListener("keydown", (e) => {
    if (e.key === "Escape" && !editVisitModal.classList.contains("hidden")) {
        closeEditVisitModal()
    }
})


// ---------- Документы ----------

// Скачивание файла: сервер отдает JSON, где content — строка Base64
async function downloadFile(fileId, fileName, buttonElement) {
    const originalText = buttonElement.textContent
    buttonElement.textContent = "⏳ Загрузка..."
    buttonElement.disabled = true

    try {
        const response = await fetch(`http://localhost:8080/api/documents?id=${fileId}`)
        if (!response.ok) throw new Error("Ошибка при скачивании файла")

        const doc = await response.json()

        // Декодируем Base64 обратно в байты
        const byteCharacters = atob(doc.content)
        const byteNumbers = new Array(byteCharacters.length)
        for (let i = 0; i < byteCharacters.length; i++) {
            byteNumbers[i] = byteCharacters.charCodeAt(i)
        }
        const byteArray = new Uint8Array(byteNumbers)

        const blob = new Blob([byteArray], { type: "application/octet-stream" })
        const url = window.URL.createObjectURL(blob)

        const a = document.createElement("a")
        a.style.display = "none"
        a.href = url
        a.download = doc.name || fileName
        document.body.appendChild(a)
        a.click()

        window.URL.revokeObjectURL(url)
        document.body.removeChild(a)
    } catch (error) {
        console.error("Ошибка скачивания:", error)
        alert("Не удалось скачать файл: " + fileName)
    } finally {
        buttonElement.textContent = originalText
        buttonElement.disabled = false
    }
}

// Читает файл и возвращает строку Base64 (без префикса data:...)
function getBase64(file) {
    return new Promise((resolve, reject) => {
        const reader = new FileReader()
        reader.readAsDataURL(file)
        reader.onload = () => resolve(reader.result.split(",")[1])
        reader.onerror = error => reject(error)
    })
}

uploadBtn.addEventListener("click", async () => {
    const file = fileInput.files[0]

    if (!file) {
        return
    }

    try {
        uploadBtn.disabled = true
        uploadBtn.textContent = "Отправка..."
        fileInput.disabled = true

        const base64Content = await getBase64(file)

        const payload = {
            pet_id: parseInt(petId, 10),
            appointment_id: parseInt(visitId, 10),
            date: new Date().toISOString(),
            name: file.name,
            content: base64Content
        }

        const response = await fetch("http://localhost:8080/api/documents", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(payload)
        })

        if (!response.ok) {
            throw new Error("Ошибка сервера при загрузке документа")
        }

        fileInput.value = ""
        await loadVisit()

    } catch (error) {
        console.error("Ошибка при отправке файла:", error)
        alert("Не удалось загрузить файл")
    } finally {
        uploadBtn.disabled = false
        uploadBtn.textContent = "Загрузить"
        fileInput.disabled = false
    }
})

async function deleteDocument(fileId, fileName) {
    const isConfirmed = confirm(`Вы уверены, что хотите безвозвратно удалить файл "${fileName}"?`)
    if (!isConfirmed) return

    try {
        const response = await fetch(`http://localhost:8080/api/documents?id=${fileId}`, {
            method: "DELETE"
        })

        if (!response.ok) {
            throw new Error(`Ошибка сервера при удалении документа (статус: ${response.status})`)
        }

        await loadVisit()

    } catch (error) {
        console.error("Ошибка при удалении файла:", error)
        alert("Произошла ошибка при удалении файла. Попробуйте еще раз.")
    }
}


// ---------- Drag & Drop ----------

function preventDefaults(e) {
    e.preventDefault()
    e.stopPropagation()
}

;["dragenter", "dragover", "dragleave", "drop"].forEach(eventName => {
    dropZone.addEventListener(eventName, preventDefaults, false)
    document.body.addEventListener(eventName, preventDefaults, false)
})

;["dragenter", "dragover"].forEach(eventName => {
    dropZone.addEventListener(eventName, () => dropZone.classList.add("dragover"), false)
})

;["dragleave", "drop"].forEach(eventName => {
    dropZone.addEventListener(eventName, () => dropZone.classList.remove("dragover"), false)
})

dropZone.addEventListener("drop", handleDrop, false)

function handleDrop(e) {
    const files = e.dataTransfer.files

    if (files && files.length > 0) {
        fileInput.files = files
    }
}


// ---------- Старт ----------

if (petId) loadPet()
if (visitId) loadVisit()
