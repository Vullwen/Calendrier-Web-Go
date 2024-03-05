//fontion pour l'affichage de la date et de l'heure dans le header
function getCurrentHour() {
    const currentHour = document.querySelector('.currentHour h1');
    const date = new Date();
    const hour = date.getHours();
    let minutes = date.getMinutes();
    if (minutes < 10) {
        minutes = "0" + minutes;
    }
    currentHour.innerHTML = `${hour}:${minutes}`;
}

function getCurrentDate() {
    const currentDate = document.querySelector('.currentDate h2');
    const date = new Date();
    const day = date.getDate();
    const month = date.toLocaleString('fr-FR', { month: 'long' });
    const year = date.getFullYear();
    currentDate.innerHTML = `${day} ${month} ${year}`;
}

setInterval(getCurrentHour, 1000);
setInterval(getCurrentDate, 1000);
setInterval(updateNotification, 1000*60*10); // 10 minutes
/***************************************************************** */
//si il y a une date dans l'url, on la récupère
const urlParams = new URLSearchParams(window.location.search);
const date = urlParams.get('date');
if (date) {
    tmp_date = new Date(date);
}else{
    tmp_date = new Date();
}
let currentWeekStartDate = getStartOfWeek(tmp_date);
//fonction pour afficher la semaine actuelle
function showCurrentWeek() {
    updateDateDisplay(currentWeekStartDate);
    updateList();
}
//fonction pour afficher la semaine suivante
function showNextWeek() {
    console.log("next");
    currentWeekStartDate.setDate(currentWeekStartDate.getDate() + 7);
    updateDateDisplayAndURL(currentWeekStartDate);
    updateList();
}
//fonction pour afficher la semaine précédente
function showPreviousWeek() {
    console.log("previous");
    currentWeekStartDate.setDate(currentWeekStartDate.getDate() - 7);
    updateDateDisplayAndURL(currentWeekStartDate);
    updateList();
}
function updateDateDisplayAndURL(startDate) {
    updateDateDisplay(startDate);

    const formattedDate = startDate.toISOString().slice(0, 19).replace("T", " ");

    const urlParams = new URLSearchParams(window.location.search);
    urlParams.set('date', formattedDate);

    history.replaceState(null, null, "?" + urlParams.toString());
    window.location.href = window.location.pathname + "?" + urlParams.toString();
}


//fonction pour mettre à jour l'affichage de la date dans la semaine actuelle
function updateDateDisplay(startDate) {
    const firstDayOfWeek = document.querySelector('.mainToolsLayout .date h3');
    const day = startDate.getDate();
    const month = startDate.toLocaleString('fr-FR', { month: 'long' });
    const year = startDate.getFullYear();
    const formattedDay = (day < 10) ? '0' + day : day;
    firstDayOfWeek.innerHTML = `${formattedDay} ${month} ${year}`;
}

//fonction pour récupérer le premier jour de la semaine en fonction de la date
function getStartOfWeek(date) {
    const dayOfWeek = date.getDay();
    console.log(dayOfWeek);
    let diff = date.getDate() - dayOfWeek + (dayOfWeek === 0 ? -6 : 1);
    return new Date(date.setDate(diff));
}

document.addEventListener('DOMContentLoaded', showCurrentWeek);

//fonction pour le format d'affichage du planning
function selectItem() {
    const selectedElements = document.querySelectorAll('.format .button.selected');
    selectedElements.forEach(selectedElement => {
        selectedElement.classList.remove('selected');
        selectedElement.classList.add('notSelected');
        const img = selectedElement.querySelector('img');
        const src = img.getAttribute('src');
        const newSrc = src.replace('selected', 'notSelected');
        img.setAttribute('src', newSrc);
    });

    this.classList.remove('notSelected');
    this.classList.add('selected');
    const img = this.querySelector('img');
    const src = img.getAttribute('src');
    const newSrc = src.replace('notSelected', 'selected');
    img.setAttribute('src', newSrc);
}
//on ajoute un event listener sur chaque bouton pour échanger le format
const elements = document.querySelectorAll('.format .button');
elements.forEach(element => {
    element.addEventListener('click', selectItem);
});

//fonction qui met à jour la liste des jours en fonction de la date
function updateList() {
    const listItems = document.querySelectorAll('.listDay ol li');
    listItems.forEach((li, index) => {
        if (index > 0) {
            const dayDate = new Date(currentWeekStartDate);
            dayDate.setDate(dayDate.getDate() + index-1);
            const day = dayDate.getDate();
            const month = dayDate.toLocaleString('fr-FR', { month: 'long' });
            const liSelected = li.querySelector('h6');
            liSelected.innerHTML = `${day} ${month}`;
        }
    });
}
//fonction pour du responsive
window.addEventListener('resize', updateLayout);

function updateLayout() {
    const heightLi = document.querySelector('.listDay ol li:nth-child(2)').offsetHeight;
    //console.log(heightLi);

    const heightListEventItems = document.querySelectorAll('.listEvent ol li');
    const heightListDayItems = document.querySelectorAll('.listDay ol li');

    heightListEventItems.forEach((li, index) => {
        li.style.height = 'calc(100% - ' + heightListDayItems[index].offsetHeight + 'px)';
        //console.log(li.style.height);
    });
}

updateLayout();
updateNotification();

function searchEvent(value) {
    
    var resultLink

    var myHeaders = new Headers();
    myHeaders.append("Authorization", "Basic QVBJX1VTRVI6UVRidGozUDNHNGVwbkJLdVVYN1JXYU11Zw==");

    var requestOptions = {
    method: 'GET',
    headers: myHeaders,
    redirect: 'follow'
    };

    fetch("http://dedream.fr/api/event?title="+value+"&description="+value+"&localisation="+value+"&like=true", requestOptions).then(response => response.json()).then(result => {
        document.querySelector('#searchEvent > #eventList').innerHTML = ""

        for (let i = 0; i < result.length; i++) {
            
            resultLink = fetch("http://dedream.fr/api/link?user_id="+document.cookie.split("id=")[1].split(";")[0]+"&event_id="+result[i].id, requestOptions).then(response => {
                if (response.status == 200) {
                    document.querySelector('#searchEvent > #eventList').innerHTML += "<div class='card'><div class='card-header'><h4>"+result[i].title+"</h4></div><div class='card-body'><p class='card-text'>"+result[i].description+"</p><p class='card-text'>Debut : "+result[i].start_date+"</p><p class='card-text'>Fin : "+result[i].end_date+"</p><button class='button' onclick='openEvent("+result[i].id+")'>Voir l'évènement</button></div></div>"
                }
            })
        }
    }).catch(error => console.log('error', error));

}

function updateNotification() {

    var resultLink

    var myHeaders = new Headers();
    myHeaders.append("Authorization", "Basic QVBJX1VTRVI6UVRidGozUDNHNGVwbkJLdVVYN1JXYU11Zw==");

    var requestOptions = {
    method: 'GET',
    headers: myHeaders,
    redirect: 'follow'
    };

    fetch("http://dedream.fr/api/events?start_date="+new Date().toISOString().slice(0, 19).replace("T", " "), requestOptions).then(response => response.json()).then(result => {

        for (let i = 0; i < result.length; i++) {

            console.log(result[i]);
        
            resultLink = fetch("http://dedream.fr/api/link?user_id="+document.cookie.split("id=")[1].split(";")[0]+"&event_id="+result[0].id, requestOptions).then(response => {
                if (response.status == 200) {
                    showPopup(result[i].title, result[i].description, result[i].start_date, result[i].end_date)
                }
            })
        }
    }).catch(error => console.log('error', error));
}

function showPopup(title, description, startDate, endDate) {
    const popupContainer = document.createElement('div');
    popupContainer.className = 'popup-container';

    const popupContent = document.createElement('div');
    popupContent.className = 'popup-content';

    const closeBtn = document.createElement('span');
    closeBtn.className = 'close-btn';
    closeBtn.innerHTML = '&times;';
    closeBtn.addEventListener('click', function () {
        document.body.removeChild(popupContainer);
    });

    const titleElement = document.createElement('h2');
    titleElement.textContent = title;

    const descriptionElement = document.createElement('p');
    descriptionElement.textContent = description;

    const startDateElement = document.createElement('p');
    startDateElement.innerHTML = 'Date de début: <span class="event-date">' + startDate + '</span>';

    const endDateElement = document.createElement('p');
    endDateElement.innerHTML = 'Date de fin: <span class="event-date">' + endDate + '</span>';

    popupContent.appendChild(closeBtn);
    popupContent.appendChild(titleElement);
    popupContent.appendChild(descriptionElement);
    popupContent.appendChild(startDateElement);
    popupContent.appendChild(endDateElement);

    popupContainer.appendChild(popupContent);
    document.body.appendChild(popupContainer);
}