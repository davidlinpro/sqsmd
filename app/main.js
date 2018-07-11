function begin() {
    // get initial stats via regular XHR
    doXHR();
    // continue to get stats via websockets
    doWS();
}


function getColor(n) {
    // these color limits can be configurable
    var i = parseInt(n);
    if ((i > 1) && (i < 5)) {
        return "orange";
    }
    if (i >= 5) {
        return "red";
    }
    return "green";                        
}

function displayStats(tbody, data) {
    var tbodydata = "";
    var stats = JSON.parse(data);
    for (var queue_name in stats) {
        let thisQ = stats[queue_name];
        tbodydata += "<tr align='center'>";
        tbodydata += "<td><b>" + thisQ.name + "</b></td>";
        tbodydata += "<td bgcolor='" + getColor(thisQ.delayed) + "'>" + thisQ.delayed + " (max:" + thisQ.delayed_max + ")</td>";
        tbodydata += "<td bgcolor='" + getColor(thisQ.messages) + "'>" + thisQ.messages + " (max:" + thisQ.messages_max + ")</td>";
        tbodydata += "<td bgcolor='" + getColor(thisQ.not_visible) + "'>" + thisQ.not_visible + " (max:" + thisQ.not_visible_max + ")</td>";
        tbodydata += "</tr>";
    }
    tbody.innerHTML = tbodydata;
}


function doXHR() {
    var tbody = document.getElementById("statsbody");
    var request = new XMLHttpRequest();
    request.onreadystatechange = function() {
        if(request.readyState === 4) {
            if(request.status === 200) { 
                displayStats(tbody, request.responseText);
            } else {
                console.log("xhr error: invalid response data");
            } 
        }
    }
    request.open('GET', '/stats');
    request.send();    
}


function doWS() {
    var tbody = document.getElementById("statsbody");
    var conn = new WebSocket("ws://" + window.location.host + "/ws_stats");
    conn.onopen = function() {
        console.log("ws connected.");
        conn.send("give me the stats!");
    }
    conn.onclose = function(e) {
        console.log("ws closed: " + e.code + " " + e.reason);
    }
    conn.onerror = function(e) {
        console.log("ws error: " + e.code + " " + e.reason);
    }
    conn.onmessage = function(e) {
        console.log("update received.");
        if (typeof e.data === "string") {
            displayStats(tbody, e.data);
            // console.log(e.data);
        } else {
            console.log("ws error: invalid response data");
        }
    }
}
