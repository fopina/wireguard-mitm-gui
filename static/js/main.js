function checkStatus() {
    $.getJSON('api/config', function(jd) {
        console.log(jd);
        $('#my-ip-input').val(jd.YourIP);
        if (jd.Config == null) {
            $('#current-config').text(`Disabled`);
        } else {
            $('#current-config').text(`Redirecting to ${jd.Config.Ip}:${jd.Config.Port}`);
            if ($('#proxy-ip-input').val() === "") {
                $('#proxy-ip-input').val(jd.Config.Ip);
                $('#proxy-port-input').val(jd.Config.Port);
            }
        }
    }).fail(function(r) {
        alert("Error " + r.status + ": " + r.responseText);
    });
}

$(function(){
    $('#reload-button').click(checkStatus);

    $('#use-ip-button').click(function(){
        $('#proxy-ip-input').val($('#my-ip-input').val());
    });

    $('#save-button').click(function(){
        $.ajax('api/config', {
            data: JSON.stringify({
                'Ip': $('#proxy-ip-input').val(),
                'Port': parseInt($('#proxy-port-input').val())
            }),
            contentType: 'application/json',
            type: 'POST',
        }).done(function(){
            checkStatus()
        }).fail(function(r) {
            alert("Error " + r.status + ": " + r.responseText);
        });
    });

    $('#disable-button').click(function(){
        $.ajax('api/config', {
            data: "null",
            contentType: 'application/json',
            type: 'POST',
        }).done(function(){
            checkStatus()
        }).fail(function(r) {
            alert("Error " + r.status + ": " + r.responseText);
        });
    });

    checkStatus();
})
