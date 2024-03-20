function submitForm() {
    document.getElementById("listeForm").submit();
}
document.addEventListener('DOMContentLoaded', function () {
    var toggleButton = document.getElementById('toggleSortForm');
    var sortForm = document.getElementById('categoryForm');
    var sortSection = document.querySelector('.tri');
});

//Permet de garder en mémoire le select de nombre de cartes à afficher
document.addEventListener('DOMContentLoaded', function () {

    var selectNbCartes = document.getElementById('select_nb_cartes');

    var selectedValue = localStorage.getItem('selectedNbCartes');

    if (selectedValue) {
        selectNbCartes.value = selectedValue;
    }

    selectNbCartes.addEventListener('change', function () {
        localStorage.setItem('selectedNbCartes', selectNbCartes.value);
        submitForm();
    });

    function submitForm() {
        document.getElementById("listeForm").submit();
    }
});