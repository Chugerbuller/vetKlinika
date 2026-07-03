const params = new URLSearchParams(window.location.search)
const petId = params.get("petId")
const ownerId = params.get("ownerId")
const petSummaryContent = document.getElementById("petSummaryContent")
const form = document.getElementById("createVisitForm")
const cancelBtn = document.getElementById("cancelBtn")
document.getElementById("currentDate").textContent = new Date().toLocaleDateString("ru-RU")
async function loadPetInfo() {
    try {
        const response = await fetch("http://localhost:8080/api/pets?id=" + petId)
        const pet = await response.json()
        petSummaryContent.textContent = "" // Очищаем текст перед заполнением

        // Добавляем имя
        const strongName = document.createElement("strong")
        strongName.textContent = pet.name
        petSummaryContent.appendChild(strongName)

        // Добавляем основную инфу (Тип, Возраст, Вес)
        const textInfo = document.createTextNode(
            ` • ${pet.type} • Возраст: ${pet.age} • Вес: ${pet.weight} кг`
        )
        petSummaryContent.appendChild(textInfo)

        // Добавляем аллергии, если есть
        if (pet.allergies && pet.allergies.length > 0) {
            const br = document.createElement("br")
            petSummaryContent.appendChild(br)

            const allergySpan = document.createElement("span")
            allergySpan.classList.add("text-danger", "mt-1", "d-inline-block")
            allergiesArr = pet.allergies.split(',').map(a => a.trim()).filter(a => a !== ""); // Преобразуем строку в массив
            allergySpan.textContent = "⚠️ Аллергии: " + allergiesArr.join(", ")

            petSummaryContent.appendChild(allergySpan)
        }

    } catch (error) {
        console.error("Ошибка загрузки питомца", error)

        petSummaryContent.textContent = ""
        const errorSpan = document.createElement("span")
        errorSpan.classList.add("text-danger")
        errorSpan.textContent = "Не удалось загрузить данные пациента"

        petSummaryContent.appendChild(errorSpan)
    }
}
form.addEventListener("submit", async (event) => {
    event.preventDefault()
    const visitData = {
        "pet_id": parseInt(petId),
        date: new Date().toISOString(),
        anamnesis: document.getElementById("anamnesis").value,
        disease: document.getElementById("disease").value,
        treatment: document.getElementById("treatment").value,
        recommendations: document.getElementById("recommendations").value
    }

    try {
        const response = await fetch("http://localhost:8080/api/appointments", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify(visitData)
        })

        if (response.ok) {
            window.location.href = `/pets?id=${petId}&ownerId=${ownerId}`;
        } else {
            alert("Ошибка при сохранении приема")
        }
    } catch (error) {
        console.error("Ошибка", error)
        alert("Ошибка сети при сохранении приема")
    }
})
cancelBtn.addEventListener("click", () => {
    window.location.href = `/pets?id=${petId}&ownerId=${ownerId}`
})
if (petId) loadPetInfo()