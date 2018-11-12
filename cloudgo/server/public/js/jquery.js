let username = /^[a-zA-Z][0-9a-zA-Z\_]{5,17}$/;
let email = /^[a-zA-Z0-9_\-]+@(([a-zA-Z0-9_\-])+\.)+[a-zA-Z]{2,4}$/;
let number = /^[1-9]\d{7}$/;
let phone = /^[1-9]\d{10}$/;
let data = {};
$('#clear').click(function(){
    $('input').val('');
    $('#tip').html('');
    $('p').html('');
});

$('input').focus(function(){
    switch(this.name){
        case 'username':
            $('#in_user').css('color', '#808080');
            $('#in_user').html('用户名为6~18位英文字母、数字或下划线，且必须以英文字母开头');
            break;
        case 'number':
            $('#in_num').css('color', '#808080');
            $('#in_num').html('学号为8位数字，且不能以0开头');
            break;
        case 'phone':
            $('#in_phone').css('color', '#808080');
            $('#in_phone').html('电话为11位数字，且不能以0开头');
            break;
        case 'email':
            $('#in_email').css('color', '#808080');
            $('#in_email').html('邮箱应满足通用格式');
            break;
    }
})

$('input').blur(function(){
        switch(this.name){
            case 'username':
                if(username.test($(this).val()) === false){
                    $('#in_user').css('color', 'red');
                    $('#in_user').html('用户名格式错误！');
                }
                else{
                    $('#in_user').css('color', '#808080');
                    $('#in_user').html('√');
                }
                break;
            case 'number':
                if(number.test($(this).val()) === false){
                    $('#in_num').css('color', 'red');
                    $('#in_num').html('学号格式错误！');
                }
                else{
                    $('#in_num').css('color', '#808080');
                    $('#in_num').html('√');
                }
                break;
            case 'phone':
                if(phone.test($(this).val()) === false){
                    $('#in_phone').css('color', 'red');
                    $('#in_phone').html('电话号码格式错误！');
                }
                else{
                    $('#in_phone').css('color', '#808080');
                    $('#in_phone').html('√');
                }
                break;
            case 'email':
                if(email.test($(this).val()) === false){
                    $('#in_email').css('color', 'red');
                    $('#in_email').html('邮箱格式错误！');
                }
                else{
                    $('#in_email').css('color', '#808080');
                    $('#in_email').html('√');
                }
                break;
        }
});

$('#submit').click(loadMess);

function loadMess(){
    $('#tip').html('');
    if(checkValid()) {
        data.username = $('input[name=\'username\']').val();
        data.number = $('input[name=\'number\']').val();
        data.phone = $('input[name=\'phone\']').val();
        data.email = $('input[name=\'email\']').val();
        let xmlhttp = $.post('', data)
            .done(function(data, status, xhr){
                window.location.href = '?username=' + $('input[name=\'username\']').val();
            })
            .fail(function(xhr, status){
                $('#tip').html(xhr.responseText);
            })
    }
}

function checkValid(){
    let state = 1;
    if(username.test($('input[name=\'username\']').val()) === false ||
        number.test($('input[name=\'number\']').val()) === false ||
        phone.test($('input[name=\'phone\']').val()) === false ||
        email.test($('input[name=\'email\']').val()) === false)
        state = 0;
    return state;
}