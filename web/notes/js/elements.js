function addKeyword(value = "")
{
    var keywordArea; 

    if (mode == "read")
    {

    }
    else
    {
        keywordArea =`
        <div class="input-group mb-3">
            <div class="input-group-prepend">
                <span class="input-group-text ">Keyword</span>
            </div>
            <input spellcheck="true" id="keyword-element" type="text" class="form-control" value="` + value + `">
        </div>`;
    }

    $("#meta-elements").append(keywordArea);
    refreshBin();
}

function addDescription(value = "")
{
    var descArea;

    if (mode == "read")
    {
        descArea = `<h5>` + processKeywords(value) + `</h5>`;
    }
    else
    {
        descArea = `
        <div class="input-group mb-3">
            <div class="input-group-prepend">
                <span class="input-group-text">Description</span>
            </div>
            <input spellcheck="true" id="description-element" type="text" class="form-control" value="` + value + `">
        </div>`;
    }

    $("#meta-elements").append(descArea);
    refreshBin();
}

function addHeader(value = "")
{
    var headerArea;

    if (mode == "read")
    {
        headerArea = `<h3>` + value + `</h3>`;
    }
    else
    {
        headerArea = `
        <div class="input-group mb-1 bin-group">
            <div class="input-group-prepend">
                <span class="input-group-text">Header</span>
            </div>
            <input spellcheck="true" name="header" style="border-top:none;" type="text" class="form-control element" value="` + value + `">
            <div class="input-group-append">
                <button class="input-group-text fa fa-chevron-up up" type="button"></button>
                <button class="input-group-text fa fa-chevron-down down" type="button"></button>
                <button class="input-group-text fa fa-trash bin" type="button"></button>
            </div>
        </div>`;
    }

    $("#note-elements").append(headerArea);
    refreshBin();
}

function addTextPlain(value = "")
{
    var textArea;

    if (mode == "read")
    {
        textArea = `<p style="white-space: pre-wrap;">` + processKeywords(value) + `</p>`;
    }
    else
    {
        textArea = `
        <div class="input-group mb-1 bin-group">
            <textarea spellcheck="true" name="text-plain" style="height:300px;" class='form-control element mb-1 element-group'>` + value + `</textarea>
            <div class="input-group-append">
                <button class="input-group-text fa fa-chevron-up up" type="button"></button>
                <button class="input-group-text fa fa-chevron-down down" type="button"></button>
                <button class="input-group-text fa fa-trash bin" type="button"></button>
            </div>
        </div>`;
    }

    $("#note-elements").append(textArea);
    refreshBin();
}

function addLink(link = "", label = "")
{
    var linkArea;

    if (mode == "read")
    {
        if (label == "")
            linkArea = `<a target="_blank" href="` + link + `">` + link + ` <i class="fa fa-link"></i></a>`;
        else
            linkArea = `<a target="_blank" class="" href="` + link + `">` + label + ` <i class="fa fa-link"></i></a>`;
    }
    else
    {
        linkArea = `
        <div class="input-group mb-1 bin-group">
            <div class="input-group-prepend">
                <span class="input-group-text" id="">URL</span>
            </div>
            <input spellcheck="true" name="label" type="text" class="form-control col-3" placeholder="Label" value="` + label + `">
            <input name="link" type="text" class="form-control element" placeholder="URL" value="` + link + `">
            <div class="input-group-append">
                <button class="input-group-text fa fa-chevron-up up" type="button"></button>
                <button class="input-group-text fa fa-chevron-down down" type="button"></button>
                <button class="input-group-text fa fa-trash bin" type="button"></button>
            </div>
        </div>`;
    }

    $("#note-elements").append(linkArea);
    refreshBin();
}