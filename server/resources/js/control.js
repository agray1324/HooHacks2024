function handleLoad() {
    page = document.getElementById("page");
    urlText = document.getElementById("urlText").value;
    searchText = document.getElementById("searchText").value;
    if(urlText != "" && searchText != ""){
        page.innerHTML = "<img src=\"img/pixil-gif-drawing.gif\" alt=\"Milk Loading\"  width=\"250\" />";
    }
}