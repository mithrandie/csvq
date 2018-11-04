$(document).ready(function(){
  $('.button-collapse').sideNav({
      menuWidth: 280,
      edge: 'left',
      closeOnClick: false,
      draggable: true
    }
  );

  $('.collapsible').collapsible(
    {
      onOpen: function (el) {
        var icon = el.find('.collapsible-header > i:eq(0)');

        icon.text('arrow_drop_up');
        icon.addClass('low-profile')
      },
      onClose: function (el) {
        var icon = el.find('.collapsible-header > i:eq(0)');

        icon.text('arrow_drop_down');
        icon.removeClass('low-profile')
      }
    }
  );

  var current = window.location.pathname;
  $('#slide-out a[href="' + current + '"]').addClass('current-page');

  $('.collapsible-header').each(function (i) {
    if ($(this).parent().find('a[href="' + current + '"]').length) {
      $(this).addClass('active');
      $(this).parent().find('.collapsible-body').css('display', 'block');
      var icon = $(this).find('i:eq(0)');
      icon.text('arrow_drop_up');
      icon.addClass('low-profile')
    }
  });
});
