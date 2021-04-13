function del()
{
    $.ajax({
      url: "/notes/delete",
      type: "post",
      data: { 
        id: note.ID
      },
      success: function(response) {
        var r = JSON.parse(response);

        switch (r.Type)
        {
            case "Error":
                $("#note-error").text(r.Message);
                $("#error-alert").show();
                break;
            case "Success":
                $("#note-success").text(r.Message);
                $("#success-alert").show();
                break;
        }
      },
      error: function(xhr) {
        $("#note-error").text("Server Issue");
        $("#error-alert").show();
        console.log("ERROR: Server Issue");
      }
  });
}

// Handles uploading the note (new note only)
function save()
{
    var postLocation = "/notes/update"; // Existing note

    // New note
    if (note.ID == 0)
        postLocation = "/notes/create";

    var contents = {
        "ID": note.ID,
        "Keyword": $("#keyword-element").val().toLowerCase(),
        "Desc": $("#description-element").val(),
        "Elements": []
    };

    // Loop over all the .elements and add to JSON object
    $(".element").each(function()
    {
        var row = {
            "Key": $(this).attr('name'), 
            "Value": $(this).val(),
            "Meta": []
        };

        // Get the element meta
        switch (row.Key)
        {
            case "link":
                // Find the closest label to the link & set as meta
                var label = $(this).closest(".input-group").find("[name='label']").val();
                row.Meta[0] = {
                    "Key": "label",
                    "Value": label
                };
                break;
        }
        
        contents.Elements.push(row);
    });

    console.log("Sending: " + JSON.stringify(contents))

    $.ajax({
      url: postLocation,
      type: "post",
      data: { 
        contents: JSON.stringify(contents)
      },
      success: function(response) {
        var r = JSON.parse(response);

        switch (r.Type)
        {
            case "Error":
                $("#note-error").text(r.Message);
                $("#error-alert").show();
                break;
            case "Success":
                $("#note-success").text(r.Message);
                $("#success-alert").show();
                window.location.href = "/notes/k/" + contents.Keyword.toLowerCase();
                break;
        }
      },
      error: function(xhr) {
        $("#note-error").text("Server Issue");
        $("#error-alert").show();
        console.log("ERROR: Server Issue");
      }
  });
}