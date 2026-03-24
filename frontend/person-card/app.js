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
    document.getElementById("animalAnamnesis").value = ""

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

function createAnimalCard(animal){

    const card = document.createElement("div")
    card.classList.add("animal-card")

    const name = document.createElement("div")
    name.textContent = animal.name

    card.appendChild(name)

    card.addEventListener("click", () => {

        window.location.href =
            "/pets?id=" + animal.id

    })

    return card
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

async function createAnimal(){

    let typeValue = typeSelect.value

    if(typeValue === "other"){

        const newType = {
            type: newTypeInput.value
        }

        const response = await fetch("/api/types",{
            method:"POST",
            headers:{
                "Content-Type":"application/json"
            },
            body:JSON.stringify(newType)
        })

        const createdType = await response.json()

        typeValue = createdType.type

    }

    const pet = {

        name: nameInput.value,
        weight: parseFloat(document.getElementById("animalWeight").value),
        age: parseInt(document.getElementById("animalAge").value),
        age_write_date: new Date().toISOString(),
        type: typeValue,
        anamnesis: document.getElementById("animalAnamnesis").value

    }

    await fetch("/api/pets?ownerId=" + personId,{
        method:"POST",
        headers:{
            "Content-Type":"application/json"
        },
        body:JSON.stringify(pet)
    })

    modal.classList.add("hidden")

    resetAnimalForm()

    loadAnimals()

}

saveBtn.addEventListener("click", createAnimal)

loadPerson()
loadAnimals()

loadTypes()