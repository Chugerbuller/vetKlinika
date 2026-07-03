const form = document.getElementById("createPersonForm")
const cancelBtn = document.getElementById("cancelBtn")

form.addEventListener("submit", async (event) => {

    event.preventDefault()

    const person = {

        surname: document.getElementById("surname").value,
        name: document.getElementById("name").value,
        patronymic: document.getElementById("patronymic").value,
        phone: document.getElementById("phone").value,
        address: document.getElementById("address").value

    }

    const response = await fetch(
        "http://localhost:8080/api/persons",
        {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify(person)
        }
    )

    if(response.ok){
        id = (await response.json()).id
        window.location.href = "/person?id=" + id
    }else{
        alert("Ошибка создания клиента")
    }

})


cancelBtn.addEventListener("click", () => {

    window.location.href = "/"

})