window.onload = () => {
    document.getElementById("website").innerHTML="Currently Searching:<br>" + _website
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
        list += "<div class=\"listrow row\"><pre><h1 id=\"title_" + String(i) + "\">" + _titles[i] + "  <a href=\"" + _urls[i] + "\"><i class=\"fa-solid fa-right-to-bracket\"></i></a></h1><pre><p1 class=\"col-12\" id=\"data_" + String(i) + "\">" + _data[i] + "</p1></div>"
    } 
    document.getElementById("list").innerHTML = list
}