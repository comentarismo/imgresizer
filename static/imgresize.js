function imgresizer(image, slug, width, height, quality, cb) {
    if (image == "") {
        $(image).remove();
    } else {
        //http post
        var host = "http://img.comentarismo.com/r";

        if (document.location.hostname.indexOf("localhost") !== -1) {
            host = "http://localhost:3666";
        }
        var request = $.ajax({
//                url: host + '/img/',
            url: host + '/meme/',
            type: 'post',
            data: {
                url: image,
                width: width,
                height: height,
                quality: quality
            },
            mimeType: "text/plain; charset=x-user-defined"
        });

        request.done(function (binaryData) {
            if (binaryData && binaryData !== "") {
                console.log("imgresizer DONE OK");
                var base64Data = base64Encode(binaryData);
                $(slug).attr("src", "data:image/jpeg;base64," + base64Data);
            } else {
                console.log(binaryData);
            }
            return cb();
        });

        request.fail(function (e) {
            console.log(e);
            // setTimeout(function(){
            //    window.location.reload();
            // },5000);
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
            out += CHARS.charAt(((c1 & 0x3) << 4) | ((c2 & 0xF0) >> 4));
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

function r(f) {
    /in/.test(document.readyState) ? setTimeout('r(' + f + ')', 9) : f()
}
r(function () {

    var curImages = new Array();

    $('textarea').liveUrl({
        loadStart: function () {
            $('.liveurl-loader').show();
        },
        loadEnd: function () {
            $('.liveurl-loader').hide();
        },
        success: function (data) {
            var output = $('.liveurl');
            output.find('.title').text(data.title);
            output.find('.description').text(data.description);
            output.find('.url').text(data.url);
            output.find('.image').empty();

            output.find('.close').one('click', function () {
                var liveUrl = $(this).parent();
                liveUrl.hide('fast');
                liveUrl.find('.video').html('').hide();
                liveUrl.find('.image').html('');
                liveUrl.find('.controls .prev').addClass('inactive');
                liveUrl.find('.controls .next').addClass('inactive');
                liveUrl.find('.thumbnail').hide();
                liveUrl.find('.image').hide();

                $('textarea').trigger('clear');
                curImages = new Array();
            });

            output.show('fast');

            if (data.video != null) {
                var ratioW = data.video.width / 350;
                data.video.width = 350;
                data.video.height = data.video.height / ratioW;

                var video =
                    '<object width="' + data.video.width + '" height="' + data.video.height + '">' +
                    '<param name="movie"' +
                    'value="' + data.video.file + '"></param>' +
                    '<param name="allowScriptAccess" value="always"></param>' +
                    '<embed src="' + data.video.file + '"' +
                    'type="application/x-shockwave-flash"' +
                    'allowscriptaccess="always"' +
                    'width="' + data.video.width + '" height="' + data.video.height + '"></embed>' +
                    '</object>';
                output.find('.video').html(video).show();


            }
        },
        addImage: function (image) {
            var output = $('.liveurl');
            var jqImage = $(image);
            jqImage.attr('alt', 'Preview');

            if ((image.width / image.height) > 7
                || (image.height / image.width) > 4) {
                // we dont want extra large images...
                return false;
            }

            var width = "388";
            var height = "195";
            var quality = "50";
            var url = jqImage.attr('src');

            console.log("INIT imgresizer ", url, "#img-test", width, height, quality);
            imgresizer(url, "#img-test", width, height, quality, function () {
                console.log("END imgresizer")


                if (curImages.length == 1) {
                    // first image...
                    output.find('.thumbnail .current').text('1');
                    output.find('.thumbnail').show();
                    output.find('.image').show();
                    jqImage.addClass('active');
                }

                if (curImages.length == 2) {
                    output.find('.controls .next').removeClass('inactive');
                }

                output.find('.thumbnail .max').text(curImages.length);

            });
        }
    });


    $('.liveurl ').on('click', '.controls .button', function () {
        var self = $(this);
        var liveUrl = $(this).parents('.liveurl');
        var content = liveUrl.find('.image');
        var images = $('img', content);
        var activeImage = $('img.active', content);

        if (self.hasClass('next'))
            var elem = activeImage.next("img");
        else var elem = activeImage.prev("img");

        if (elem.length > 0) {
            activeImage.removeClass('active');
            elem.addClass('active');
            liveUrl.find('.thumbnail .current').text(elem.index() + 1);

            if (elem.index() + 1 == images.length || elem.index() + 1 == 1) {
                self.addClass('inactive');
            }
        }

        if (self.hasClass('next'))
            var other = elem.prev("img");
        else var other = elem.next("img");

        if (other.length > 0) {
            if (self.hasClass('next'))
                self.prev().removeClass('inactive');
            else   self.next().removeClass('inactive');
        } else {
            if (self.hasClass('next'))
                self.prev().addClass('inactive');
            else   self.next().addClass('inactive');
        }


        var imgs = document.getElementsByTagName('img')
        for (var i = 0, j = imgs.length; i < j; i++) {
            imgs[i].onerror = function (e) {
                this.parentNode.removeChild(this);
            }
        }

    });


    $('textarea').focus();
    $("textarea").keyup();


});