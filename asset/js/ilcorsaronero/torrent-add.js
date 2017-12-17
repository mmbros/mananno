function torrentAdd( hash, href, paused ) {

  var request = {
    jsonrpc: "2.0",
    method: "torrent-add",
    params: {
      hash: hash,
      href: href,
      paused: false,
    },
    id: 1,
  }

  var xhr = new XMLHttpRequest();

  xhr.onreadystatechange = function() {
    if (xhr.readyState == 4) {
      if (xhr.status == 200) {
        var json = JSON.parse(xhr.responseText);

        if (json.hasOwnProperty("error")) {
          console.error("err 1: " + json.error.message);
		  console.log(json);
          alert(json.error.message);
        }
        else {
          console.log("torrent-added");
		  console.log(json);
          alert("torrent-added");
        }
      } else {
          console.log(xhr.status + ":" + xhr.statusText);
		  console.log(json);
          alert(xhr.status + ":" + xhr.statusText);
      }
    }
  };

  xhr.ontimeout = function () {
    console.error("The request timed out.");
    alert("The request timed out.");
  };

  var jsonRequest = JSON.stringify(request);
  console.log(jsonRequest);

  xhr.open("POST", "/jsonrpc");
  xhr.setRequestHeader("Content-Type","application/json; charset=utf-8");
  xhr.timeout = 500; // msecs
  xhr.send(jsonRequest);
  return false;
}


function doJsonReq( request ) {
  console.log(request)

  var xhr = new XMLHttpRequest();

  xhr.onreadystatechange = function() {
    if (xhr.readyState == 4) {
      if (xhr.status == 200) {
        var json = JSON.parse(xhr.responseText);

        if (json.hasOwnProperty("error")) {
          console.error("err 1: " + json.error.message);
		  console.log(json);
          alert(json.error.message);
        }
        else {
		  var msg = request.method + ": OK";
          console.log(msg);
		  console.log(json);
          alert(msg);
        }
      } else {
          console.log(xhr.status + ":" + xhr.statusText);
		  console.log(json);
          alert(xhr.status + ":" + xhr.statusText);
      }
    }
  };

  xhr.ontimeout = function () {
    console.error("The request timed out.");
    alert("The request timed out.");
  };

  var jsonRequest = JSON.stringify(request);
  console.log(jsonRequest);

  xhr.open("POST", "/jsonrpc");
  xhr.setRequestHeader("Content-Type","application/json; charset=utf-8");
  xhr.timeout = 500; // msecs
  xhr.send(jsonRequest);
  return false;
}


function magnetAdd( magnet, paused ) {

  var request = {
    jsonrpc: "2.0",
    method: "magnet-add",
    params: {
      magnet: magnet,
      paused: paused,
    },
    id: 1,
  }
  doJsonReq( request )
}

