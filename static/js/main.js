function checkStatus() {
    $.getJSON('api/config', function(jd) {
        console.log(jd);
        $('#my-ip-input').val(jd.ClientIP);
        $('#current-config').text(`Redirecting to ${jd.Ip}:${jd.Port}`);
        if ($('#proxy-ip-input').val() === "") {
            $('#proxy-ip-input').val(jd.Ip);
            $('#proxy-port-input').val(jd.Port);
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
        }).fail(function(r) {
            alert("Error " + r.status + ": " + r.responseText);
        });
    });

    $('#disable-button').click(checkStatus);

    checkStatus();
})
