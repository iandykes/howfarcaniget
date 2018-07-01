var map;
var markers = [];

function initMap() {
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

$("#btnSubmit").on('click', function () {
    deleteMarkers();
    // TODO: Get start location and hours selected
    $.getJSON("/distances", function(result){
        showResult(result);
    });
});