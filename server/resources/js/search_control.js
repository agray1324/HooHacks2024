window.onload = () => {
    document.getElementById("website").innerHTML="Currently Searching:<br>" + _website +"<br><br>Searching For:<br>" + _searchText
    delimiter = "\\,\\"
    _urls = _urls.split(delimiter)
    _titles = _titles.split(delimiter)
    _data = _data.split(delimiter)
    list = ""
    console.log(_urls.length)
    for (let i = 0; i < _urls.length; i++){
        /*if(list != ""){
            list += "<br>"
        }*/
        list += "<div class=\"listrow row\" onclick=goto(\"" + _urls[i] + "\")><pre><h1 id=\"title_" + String(i) + "\">" + _titles[i] + "</h1><pre><p1 class=\"col-12\" id=\"data_" + String(i) + "\">" + _data[i] + "</p1></div>"
    } 
    document.getElementById("list").innerHTML = list
}

function goto(url){
    window.open(url,'_blank');
}