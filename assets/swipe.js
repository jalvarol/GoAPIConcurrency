// Swipe and click navigation for GoHeadlines vertical carousel (mobile) and 3-column (desktop)
document.addEventListener('DOMContentLoaded', function () {
  const sections = document.querySelectorAll('.carousel-section');
  const leftArrow = document.querySelector('.carousel-arrow.left');
  const rightArrow = document.querySelector('.carousel-arrow.right');
  const dots = document.querySelectorAll('.swipe-dot');
  let current = 0;
  let isMobile = window.matchMedia('(max-width: 900px)').matches;

  function showSection(idx) {
    sections.forEach((sec, i) => {
      sec.style.display = i === idx ? 'block' : 'none';
      sec.classList.toggle('active', i === idx);
      if (dots[i]) dots[i].classList.toggle('active', i === idx);
    });
    current = idx;
  }

  function nextSection() {
    showSection((current + 1) % sections.length);
  }
  function prevSection() {
    showSection((current - 1 + sections.length) % sections.length);
  }

  // Remove arrow click and keyboard navigation for mobile and desktop
  if (leftArrow && rightArrow) {
    leftArrow.style.display = 'none';
    rightArrow.style.display = 'none';
  }

  // --- Swipe for all devices (mobile and desktop) ---
  let startX = null;
  let isTouching = false;
  const container = document.querySelector('.swipe-container');
  if (container) {
    container.addEventListener('touchstart', function (e) {
      if (e.touches.length === 1) {
        startX = e.touches[0].clientX;
        isTouching = true;
      }
    });
    container.addEventListener('touchmove', function (e) {
      if (isTouching) e.preventDefault();
    }, { passive: false });
    container.addEventListener('touchend', function (e) {
      if (!isTouching || startX === null) return;
      let endX = e.changedTouches[0].clientX;
      let diff = endX - startX;
      if (diff > 50) prevSection();
      else if (diff < -50) nextSection();
      startX = null;
      isTouching = false;
    });
    // Mouse drag for desktop swipe
    let mouseDown = false;
    let mouseStartX = null;
    container.addEventListener('mousedown', function (e) {
      mouseDown = true;
      mouseStartX = e.clientX;
    });
    container.addEventListener('mousemove', function (e) {
      if (!mouseDown) return;
      e.preventDefault();
    });
    container.addEventListener('mouseup', function (e) {
      if (!mouseDown || mouseStartX === null) return;
      let mouseEndX = e.clientX;
      let diff = mouseEndX - mouseStartX;
      if (diff > 50) prevSection();
      else if (diff < -50) nextSection();
      mouseDown = false;
      mouseStartX = null;
    });
    container.addEventListener('mouseleave', function () {
      mouseDown = false;
      mouseStartX = null;
    });
  }

  // Dots navigation (optional)
  if (dots.length) {
    dots.forEach((dot, idx) => {
      dot.addEventListener('click', () => showSection(idx));
    });
  }

  // Init
  showSection(0);
});
