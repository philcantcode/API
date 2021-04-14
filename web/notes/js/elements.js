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
            <input id="keyword-element" type="text" class="form-control" value="` + value + `">
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
            <input id="description-element" type="text" class="form-control" value="` + value + `">
        </div>`;
    }

    $("#meta-elements").append(descArea);
    refreshBin();
}

function addTitle(value = "")
{
    var titleArea;

    if (mode == "read")
    {
        titleArea = `<h3>` + value + `</h3>`;
    }
    else
    {
        titleArea = `
        <div class="input-group mb-1 bin-group">
            <div class="input-group-prepend">
                <span class="input-group-text">Title</span>
            </div>
            <input name="title" style="border-top:none;" type="text" class="form-control element" value="` + value + `">
            <div class="input-group-append">
                <button class="input-group-text fa fa-trash bin" type="button"></button>
            </div>
        </div>`;
    }

    $("#note-elements").append(titleArea);
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
            <textarea name="text-plain" style="height:300px;" class='form-control element mb-1 element-group'>` + value + `</textarea>
            <div class="input-group-append">
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
            <input name="label" type="text" class="form-control col-3" placeholder="Label" value="` + label + `">
            <input name="link" type="text" class="form-control element" placeholder="URL" value="` + link + `">
            <div class="input-group-append">
                <button class="input-group-text fa fa-trash bin" type="button"></button>
            </div>
        </div>`;
    }

    $("#note-elements").append(linkArea);
    refreshBin();
}