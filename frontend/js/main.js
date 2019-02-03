if (/Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(navigator.userAgent)) {
    $("#wrap").css("height", $(window).height())
        .css("width", $(window).width())
}

$(document).ready(() => {
    $("[data-toggle='offcanvas']").click(() => {
        $(".row-offcanvas").toggleClass("active")
    })
})
