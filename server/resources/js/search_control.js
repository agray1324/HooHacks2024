window.onload = () => {
    document.getElementById("website").innerHTML="Currently Searching: " + _website
    _urls = _urls.split("\\,\\")
    _titles = _titles.split("\\,\\")
    _data = _data.split("\\,\\")
    list = ""
    console.log(_urls.length)
    for (let i = 0; i < _urls.length; i++){
        list += "<div class=\"listrow row\"><pre><h1 id=\"title_" + String(i) + "\">" + _titles[i] + "  <a href=\"" + _urls[i] + "\"><i class=\"fa-solid fa-right-to-bracket\"></i></a></h1><pre><p1 class=\"col-12\" id=\"data_" + String(i) + "\">" + _data[i] + "</p1></div><br>"
    } 
    document.getElementById("list").innerHTML = list
}