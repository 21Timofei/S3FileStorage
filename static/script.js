function toggleSubmitButton() {
    const fileInput = document.getElementById("fileInput");
    const submitButton = document.getElementById("submitButton");

    if (fileInput.value) {
        submitButton.style.display = "inline-block";
    } else {
        submitButton.style.display = "none";
    }
}