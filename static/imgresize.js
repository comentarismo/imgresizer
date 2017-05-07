
function imgresizer(image,slug,width,height,quality,cb){
    if(image == "") {
        $(image).remove();
    }else {
        //http post
        // var host = "http://img.comentarismo.com/r";
        var host = "http://localhost:3666/r";

        var request = $.ajax({
            url: host + '/img/',
            type: 'post',
            data: {
                url: image,
                width : width,
                height : height,
                quality : quality
            },
            mimeType: "text/plain; charset=x-user-defined"
        });

        request.done(function (binaryData) {
            if(binaryData && binaryData !== "" ){
                console.log("imgresizer DONE OK");
                var base64Data = base64Encode(binaryData);
                $(slug).attr("src", "data:image/jpeg;base64," + base64Data);
            }else {
                console.log(binaryData);
            }
            return cb();
        });

        request.fail(function (e) {
            console.log(e);
            //setTimeout(function(){
            //    window.location.reload();
            //},5000);
        });
    }
}

function base64Encode(str) {
    var CHARS = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/";
    var out = "", i = 0, len = str.length, c1, c2, c3;
    while (i < len) {
        c1 = str.charCodeAt(i++) & 0xff;
        if (i == len) {
            out += CHARS.charAt(c1 >> 2);
            out += CHARS.charAt((c1 & 0x3) << 4);
            out += "==";
            break;
        }
        c2 = str.charCodeAt(i++);
        if (i == len) {
            out += CHARS.charAt(c1 >> 2);
            out += CHARS.charAt(((c1 & 0x3)<< 4) | ((c2 & 0xF0) >> 4));
            out += CHARS.charAt((c2 & 0xF) << 2);
            out += "=";
            break;
        }
        c3 = str.charCodeAt(i++);
        out += CHARS.charAt(c1 >> 2);
        out += CHARS.charAt(((c1 & 0x3) << 4) | ((c2 & 0xF0) >> 4));
        out += CHARS.charAt(((c2 & 0xF) << 2) | ((c3 & 0xC0) >> 6));
        out += CHARS.charAt(c3 & 0x3F);
    }
    return out;
}