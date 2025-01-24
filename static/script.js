document.getElementById("addItemForm").addEventListener("submit", function(event) {
    const name = document.getElementById("name").value;
    const category = document.getElementById("category").value;
    const value = document.getElementById("value").value;
    const quantity = document.getElementById("quantity").value;

    let errorMessages = [];
    if (!name || !category || !value || !quantity) {
        errorMessages.push("すべてのフィールドを入力してください");
    }

    if (isNaN(value) || value <= 0) {
        errorMessages.push("価値（ゴールド）は正の数で入力してください");
    }

    if (isNaN(quantity) || quantity <= 0) {
        errorMessages.push("個数は正の数で入力してください");
    }

    const errorMessagesDiv = document.getElementById("errorMessages");
    errorMessagesDiv.innerHTML = errorMessages.join("<br>");

    if (errorMessages.length > 0) {
        event.preventDefault();
    }
});

document.getElementById("cancelButton").addEventListener("click", function() {
    document.getElementById("name").value = '';
    document.getElementById("category").value = '';
    document.getElementById("value").value = '';
    document.getElementById("quantity").value = '';
    document.getElementById("errorMessages").innerHTML = '';
});
