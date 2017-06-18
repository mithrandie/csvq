$(document).ready(function(){
  $('.button-collapse').sideNav({
      menuWidth: 280,
      draggable: true
    }
  );

  var colOpen = function (el) {
    var icon = el.find('.collapsible-header > i:eq(0)');

    icon.text('arrow_drop_up');
    icon.addClass('low-profile')
  };

  var colClose = function (el) {
    var icon = el.find('.collapsible-header > i:eq(0)');

    icon.text('arrow_drop_down');
    icon.removeClass('low-profile')
  };

  $('.collapsible').collapsible(
    {
      onOpen: colOpen,
      onClose: colClose
    },
    {
      onOpen: colOpen,
      onClose: colClose
    }
  );
});
