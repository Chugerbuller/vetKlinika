const petCard = document.getElementById("petCard")
const anamnesisText = document.getElementById("anamnesisText")
const allergyList = document.getElementById("allergyList")
const visitsList = document.getElementById("visitsList")

const modal = document.getElementById("visitModal")
const closeModal = document.getElementById("closeVisitModal")

const params = new URLSearchParams(window.location.search)
const petId = params.get("id")


closeModal.addEventListener("click", () => {

    modal.classList.add("hidden")

})


function addField(label, value){

    const field = document.createElement("div")
    field.classList.add("pet-field")

    const text = document.createElement("strong")
    text.textContent = label + ": "

    const valueText = document.createElement("span")
    valueText.textContent = value

    field.appendChild(text)
    field.appendChild(valueText)

    petCard.appendChild(field)

}


function renderPet(pet){

    addField("Имя", pet.name)
    addField("Тип", pet.type)
    addField("Порода", pet.breed)
    addField("Возраст", pet.age)
    addField("Вес", pet.weight)

    anamnesisText.textContent = pet.anamnesis

    if(pet.allergies){

        pet.allergies.forEach(allergy => {

            const item = document.createElement("div")
            item.textContent = allergy

            allergyList.appendChild(item)

        })

    }

}


function createVisitCard(visit){

    const card = document.createElement("div")
    card.classList.add("visit-card")

    const text = document.createElement("div")
    text.textContent = visit.date

    card.appendChild(text)

    card.addEventListener("click", () => {

        modal.classList.remove("hidden")

    })

    return card

}


function renderVisits(visits){

    visits.forEach(v => {

        const card = createVisitCard(v)

        visitsList.appendChild(card)    

    })

}


async function loadPet(){

    const response = await fetch("/api/pets?id=" + petId)
    const pet = await response.json()

    renderPet(pet)

}


async function loadVisits(){
    return
    const response = await fetch("/api/visits?petId=" + petId)
    const visits = await response.json()

    renderVisits(visits)

}


loadPet()
loadVisits()