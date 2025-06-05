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

  // Show arrows for desktop (mouse) only
  function updateArrowVisibility() {
    const isDesktop = window.matchMedia('(pointer: fine) and (hover: hover)').matches;
    if (leftArrow && rightArrow) {
      leftArrow.style.display = isDesktop ? 'block' : 'none';
      rightArrow.style.display = isDesktop ? 'block' : 'none';
    }
  }
  updateArrowVisibility();
  window.addEventListener('resize', updateArrowVisibility);

  // Arrow navigation for desktop (mouse)
  if (leftArrow && rightArrow) {
    leftArrow.addEventListener('click', prevSection);
    rightArrow.addEventListener('click', nextSection);
  }

  // --- Swipe for all devices (mobile and trackpad) ---
  let startX = null;
  let isDragging = false;
  const container = document.querySelector('.swipe-container');
  if (container) {
    // Touch events for mobile/trackpad
    container.addEventListener('touchstart', function (e) {
      if (e.touches.length === 1) {
        startX = e.touches[0].clientX;
        isDragging = true;
      }
    });
    container.addEventListener('touchmove', function (e) {
      if (isDragging) e.preventDefault();
    }, { passive: false });
    container.addEventListener('touchend', function (e) {
      if (!isDragging || startX === null) return;
      let endX = e.changedTouches[0].clientX;
      let diff = endX - startX;
      if (diff > 50) prevSection();
      else if (diff < -50) nextSection();
      startX = null;
      isDragging = false;
    });
    // Remove mouse drag for desktop, keep only arrows for mouse
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
