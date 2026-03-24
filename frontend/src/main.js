const list = document.getElementById("clientsList")
const searchInput = document.getElementById("searchInput")
const addButton = document.getElementById("addClientBtn")


addButton.addEventListener("click", () => {

    window.location.href =
        "http://localhost:8080/person/create"

})


function createClientCard(client){

    const card = document.createElement("div")
    card.classList.add("client-card")

    const name = document.createElement("div")
    name.classList.add("client-name")
    name.textContent =
        client.surname + " " +
        client.name + " " +
        (client.patronymic || "")

    const meta = document.createElement("div")
    meta.classList.add("client-meta")
    meta.textContent =
        "📞 " + client.phone +
        " • " +
        "📍 " + client.address

    card.appendChild(name)
    card.appendChild(meta)

    card.addEventListener("click", () => {

        window.location.href =
            "http://localhost:8080/person?id=" + client.id

    })

    return card
}


function renderClients(clients){

    list.textContent = ""

    clients.forEach(client => {

        const card = createClientCard(client)

        list.appendChild(card)

    })

}


async function loadClients(){

    const response =
        await fetch("http://localhost:8080/api/persons")

    const data = await response.json()

    renderClients(data)

}


async function searchClients(name){

    const response =
        await fetch(
            "http://localhost:8080/api/persons?name=" + name
        )

    const data = await response.json()

    renderClients(data)

}


searchInput.addEventListener("input", (e) => {

    const value = e.target.value.trim()

    if(value === ""){
        loadClients()
    }else{
        searchClients(value)
    }

})


loadClients()