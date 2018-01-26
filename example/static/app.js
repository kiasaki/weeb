document.addEventListener('DOMContentLoaded', function () {
  var navbarBurgerEls = Array.prototype.slice.call(document.querySelectorAll('.navbar-burger'), 0);
  navbarBurgerEls.forEach(function(el) {
    el.addEventListener('click', function() {
      var target = el.dataset.target;
      var targetEl = document.getElementById(target);
      el.classList.toggle('is-active');
      targetEl.classList.toggle('is-active');
    });
  });

  var dropdownEls = Array.prototype.slice.call(document.querySelectorAll('.dropdown'), 0);
  dropdownEls.forEach(function(el) {
    el.querySelector('.dropdown-trigger').addEventListener('click', function() {
      el.classList.toggle('is-active');
    });
  });
});
