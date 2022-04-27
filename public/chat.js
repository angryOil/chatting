$(function () {
    if (!window.EventSource) {
        alert("No EventSource!")
        return
    }
    let $chatLog = $('#chat-log')
    let $chatMsg = $('#chat-msg')

    let isBlank = function (string) {
        return string == null || string.trim() === "";
    };
    let username;
    while (isBlank(username)) {
        username = prompt("name?")
        if (!isBlank(username)) {
            $('#user-name').html('<b>' + username + '</b>')
        }
    }
    $('#input-form').on('submit', function (e) {
        $.post('/messages', {
            msg: $chatMsg.val(),
            name: username
        });
        $chatMsg.val("");
        $chatMsg.focus();
        return false;
    });
    let addMessage = function (data) {
        let text = "";
        if (!isBlank(data.name)) {
            text += '<strong>' + data.name + ':</strong>';
        }
        text += data.msg;
        $chatLog.prepend('<div><span>' + text + '</span></div>');
    };
    addMessage({
        msg: 'hello',
        name: 'aaa'
    });
    addMessage({
        msg: 'hello2',
        name: 'bbb'
    })

    const es = new EventSource('/stream');
    es.onopen = function (e) {
        $.post('/users',{
            name:username
        })
    };
    es.onmessage = function (e) {
        let msg = JSON.parse(e.data);
        addMessage(msg);
    };

    window.onbeforeunload = function () {
        $.ajax({
            url: "/users?username=" + username,
            type: "delete"
        });
        es.close();
    };
});