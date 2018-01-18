// Copyright 2017 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

const template = `
<!DOCTYPE html>
<html>
<head>
    <meta name="viewport" content="initial-scale=1.0, user-scalable=no">
    <meta charset="utf-8">
    <title>maproute</title>
    <style>
        #map {
            height: 100%;
        }

        /* Optional: Makes the sample page fill the window. */
        html, body {
            height: 100%;
            margin: 0;
            padding: 0;
        }
    </style>
</head>
<body>
<div id="map"></div>
<script>
    function initMap() {
        var map = new google.maps.Map(document.getElementById('map'), {
            zoom: 3,
            center: {lat: 0, lng: -130},
            mapTypeId: 'terrain'
        });

        var flightPlanCoordinates = [
			GEOS_HERE
        ];
        var flightPath = new google.maps.Polyline({
            path: flightPlanCoordinates,
            geodesic: true,
            strokeColor: '#FF0000',
            strokeOpacity: 1.0,
            strokeWeight: 2
        });

        flightPath.setMap(map);
    }
</script>
<script async defer src="https://maps.googleapis.com/maps/api/js?key=API_KEY_HERE&callback=initMap">
</script>
</body>
</html>
`
