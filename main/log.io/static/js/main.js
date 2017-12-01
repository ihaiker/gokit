String.prototype.replaceAll = function (s1, s2) {
    return this.replace(new RegExp(s1, "gm"), s2);
}

$(function () {
    var grepMessage = "";
    var $messages = $("#messages");
    var appendMessage = function (msg) {
        var theDiv = $messages[0];
        var doScroll = theDiv.scrollTop == theDiv.scrollHeight - theDiv.clientHeight;
        $(msg).appendTo($messages);
        if (doScroll) {
            theDiv.scrollTop = theDiv.scrollHeight - theDiv.clientHeight;
        }
    };
    var w = new WebSocket("ws://" + window.location.host + "/ws" + Remote + FID);
    w.onopen = function () {
        appendMessage("the logs open!");
    };

    w.onclose = function () {
        appendMessage("<div><center><h3>Disconnected</h3></center></div>");
    };
    w.onmessage = function (msg) {
        if (grepMessage == "") {
            appendMessage("<div>" + msg.data + "</div>");
        } else {
            appendMessage("<div>" + msg.data.replaceAll(grepMessage, "<red>" + grepMessage + "</red>") + "</div>");
        }
    };

    $("#grepBtn").click(function () {
        var messageTxt = $("#grepTxt");
        grepMessage = messageTxt.val().toString();
        w.send(grepMessage);
    });
    $("#cleanBtn").click(function () {
        $messages.html("");
    });

    var ws = windowsSize();
    $messages.height(ws.height - 65);
});

function windowsSize() {
    var winWidth, winHeight;
    // 获取窗口宽度
    if (window.innerWidth) {
        winWidth = window.innerWidth;
    }
    else if ((document.body) && (document.body.clientWidth)) {
        winWidth = document.body.clientWidth;
    }

    // 获取窗口高度
    if (window.innerHeight) {
        winHeight = window.innerHeight;
    }
    else if ((document.body) && (document.body.clientHeight)) {
        winHeight = document.body.clientHeight;
    }
    //通过深入 Document 内部对 body 进行检测，获取窗口大小
    if (document.documentElement && document.documentElement.clientHeight && document.documentElement.clientWidth) {
        winHeight = document.documentElement.clientHeight;
        winWidth = document.documentElement.clientWidth;
    }
    return {"width": winWidth, "height": winHeight};
}