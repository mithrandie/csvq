$(document).ready(function(){
  $('.button-collapse').sideNav({
      menuWidth: 280,
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
});
