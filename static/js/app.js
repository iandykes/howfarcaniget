var map;
var markers = [];

function initMap() {
    // TODO: Geolocation to centre on current location
    map = new google.maps.Map(document.getElementById('map'), {
        center: {lat: 52.697, lng: -1.644},
        zoom: 6
    });
}
function addMarker(location, label) {
    var marker = new google.maps.Marker({
      position: location,
      label: label,
      map: map
    });

    markers.push(marker);
  }

function showResult(result) {
    if (result.statusCode != 200) {
        // TODO: Feedback on errors
        console.log("Error", result);
        return;
    }
    // TODO: Style each marker specifically: Origin, Each durationGroup, plus no results ones
    addMarker(result.origin, "#");
    for (i=0; i < result.points.length; i++) {
        if (result.points[i].durationGroup) {
            addMarker(result.points[i].destination, result.points[i].durationGroup.toString());
        }
        else {
            addMarker(result.points[i].destination, "!");
        }
    }
}

function deleteMarkers() {
    for (var i = 0; i < markers.length; i++) {
        markers[i].setMap(map);
    }

    markers = [];
}

function hideTitleCard() {
   $("#titleCard").fadeOut();
   $("#searchAgainBox").fadeIn();
   $("#versionBox").fadeOut();
}

function showTitleCard() {
    $("#titleCard").fadeIn();
    $("#searchAgainBox").fadeOut();
    $("#versionBox").fadeIn();
}

// TODO: Sort out click vs tap
$("#btnSubmit").on('click', function () {
    deleteMarkers();
    hideTitleCard();

    var searchValue = $("#txtSearch").val();

    // TODO: Get start location and hours selected
    $.getJSON("/api/distances?s="+searchValue, function(result){
        showResult(result);
    });
});

$("#btnSearchAgain").click(function(){
    showTitleCard();
});

$("#btnCloseMainCard").click(function() {
    hideTitleCard();
})