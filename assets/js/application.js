import xhook from 'xhook';
import $ from 'jquery';
window.jQuery = $;
window.$ = $;
import 'popper.js';
import 'parsleyjs';
import 'bootstrap';
import 'select2';
import 'jquery-smooth-scroll';
import Siema from 'siema';
import AOS from 'aos';

import '../scss/application.scss'

$(document).ready(function () {

    //make dropdown link navigatable
    $('.navbar .dropdown-toggle').click(function () {
        if (!isMobileDevice())
            window.location = $(this).attr('href');
    });

    if (document.querySelector('#ck-content')) {
        //add csrf protection to ckeditor uploads
        xhook.before(function (request) {
            if (!/^(GET|HEAD|OPTIONS|TRACE)$/i.test(request.method)) {
                request.xhr.setRequestHeader("X-CSRF-TOKEN", window.csrf_token);
            }
        });

        ClassicEditor
            .create(document.querySelector('#ck-content'), {
                language: 'en', //to set different lang include <script src="/public/js/ckeditor/build/translations/{lang}.js"></script> along with core ckeditor script
                ckfinder: {
                    uploadUrl: '/admin/upload'
                },
            })
            .catch(error => {
                console.error(error);
            });
    }

    $('#upload').change(function (e) {
        var formData = new FormData();
        var file = document.getElementById('upload').files[0];
        formData.append('upload', file, file.name);
        var xhr = new XMLHttpRequest();
        xhr.open('POST', '/admin/new_image', true);
        // Set up a handler for when the request finishes.
        xhr.onload = function () {
            if (xhr.status === 200) {
                var img = JSON.parse(this.responseText);
                $('#product-form').append('<input type="hidden" name="imageids" value="' + img.ID + '" id="imageids-' + img.ID + '" />');
                $("#images").append('<div class="img-wrapper" data-id="' + img.ID + '"><img src="' + img.URL + '" class="card-img-top" /><div class="text-center mb-2 mt-auto"><a href="#" class="btn btn-outline-secondary btn-sm remove-btn">Удалить</a></div></div>');
                //append img-wrapper
                setDefaultImage();
            } else {
                alert('Error occurred while uploading the file!');
            }
        };
        // Send the Data.
        xhr.send(formData);
    });

    $('#product-form').on('click', '.img-wrapper', function () {
        $('.img-wrapper.default').removeClass("default");
        $(this).addClass("default");
        setDefaultImage();
    });
    $('#product-form').on('click', '.remove-btn', function () {
        var id = $(this).closest('.img-wrapper').attr('data-id');
        var xhr = new XMLHttpRequest();
        xhr.open('POST', '/admin/images/' + id + '/delete', true);
        // Set up a handler for when the request finishes.
        xhr.onload = function () {
            if (xhr.status === 201) {
                $('.img-wrapper[data-id="' + id + '"]').remove();
                $('#imageids-' + id).remove();
                setDefaultImage();
            } else {
                alert('Error occurred while deleting the file!');
            }
        };
        // Send the Data.
        xhr.send(null);
        return false;
    });

    $('.products-show .image-previews').on('click', 'a', function () {
        $('.image-wrapper a').attr('data-featherlight', $(this).attr('href'));
        $('.image-wrapper img').attr('src', $(this).attr('href'));
        return false;
    })

    $('.table.table-hover tr').on('click', function () {
        var url = $(this).attr('data-url');
        if (url)
            window.location = url;
    });
    $('.table-hover tr .btn').on('click', function (e) {
        e.stopPropagation();
    });

    //product-preview image selector
    $('.product-preview-wrapper .images img').on('click', function (e) {
        e.stopPropagation();
        $(this).closest('.product-preview-wrapper').find('.card-image-top').attr('src', $(this).attr('src'));
    });

    //smooth scroll
    $('a#smooth-scroll').smoothScroll({
        easing: 'linear',
        speed: 400,
        preventDefault: true,
    });

    //animate on scroll
    AOS.init({
        duration: 300,
    });

    //add to cart 
    function onAddToCart(button) {
        var form = $(button.form);
        var url = form.attr('action');
        $.ajax({
            type: "POST",
            url: url,
            data: form.serialize(),
            success: function (data) {
                $(form).find('.fa-shopping-cart').removeClass('fa-shopping-cart').addClass('fa-check');
                $('#cart-total').text(data);
            }
        });
    }
    window.onAddToCart = onAddToCart;

    //submit search on enter press
    $('#search input').keyup(function (e) {
        if (e.which == 10 || e.which == 13) {
            this.form.submit();
        }
    });

});

function setDefaultImage() {
    var definput = $('#product-form #default-image-id');
    var def = $('#product-form .img-wrapper.default');
    var defid = (def.length > 0) ? def.attr("data-id") : 0;
    if (defid == 0) {
        var wrapper = $("#product-form .img-wrapper").first();
        defid = (wrapper.length > 0) ? wrapper.attr("data-id") : 0;
    }
    $("#product-form .img-wrapper[data-id='" + defid + "']").addClass("default");
    definput.val(defid);
}

// Write your Javascript code.
function isMobileDevice() {
    return (typeof window.orientation !== "undefined") || (navigator.userAgent.indexOf('IEMobile') !== -1);
};

$(window).bind('scroll', function () {
    if ($(window).scrollTop() > 110) {
        $('.navbar.public-navbar').addClass('fixed-top');
    } else {
        $('.navbar.public-navbar').removeClass('fixed-top');
    }

    if ($(window).scrollTop() > 300) {
        $('a#smooth-scroll').removeClass('hidden');
    } else {
        $('a#smooth-scroll').addClass('hidden');
    }
});