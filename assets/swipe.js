// Swipe and click navigation for GoHeadlines vertical carousel (mobile) and 3-column (desktop)
(function() {
  const sections = document.querySelectorAll('.carousel-section');
  const dots = document.querySelectorAll('.swipe-dot');
  let current = 0;
  let isMobile = window.matchMedia('(max-width: 900px)').matches;

  function showSection(idx) {
    sections.forEach((sec, i) => {
      if (isMobile) {
        sec.style.display = i === idx ? 'block' : 'none';
      } else {
        sec.style.display = 'block'; // All visible on desktop
      }
    });
    dots.forEach((dot, i) => {
      dot.classList.toggle('active', i === idx);
    });
    current = idx;
  }

  function nextSection() {
    showSection((current + 1) % sections.length);
  }
  function prevSection() {
    showSection((current - 1 + sections.length) % sections.length);
  }

  // Touch events for swipe (mobile)
  let startY = null;
  document.addEventListener('touchstart', function(e) {
    if (!isMobile) return;
    startY = e.touches[0].clientY;
  });
  document.addEventListener('touchend', function(e) {
    if (!isMobile || startY === null) return;
    let endY = e.changedTouches[0].clientY;
    if (endY - startY > 50) prevSection(); // swipe down
    else if (startY - endY > 50) nextSection(); // swipe up
    startY = null;
  });

  // Dot click navigation
  dots.forEach((dot, i) => {
    dot.addEventListener('click', () => showSection(i));
  });

  // Keyboard navigation (optional)
  document.addEventListener('keydown', function(e) {
    if (!isMobile) return;
    if (e.key === 'ArrowUp') prevSection();
    if (e.key === 'ArrowDown') nextSection();
  });

  // Responsive: update on resize
  window.addEventListener('resize', function() {
    isMobile = window.matchMedia('(max-width: 900px)').matches;
    showSection(isMobile ? current : 0);
    // Show all sections on desktop
    if (!isMobile) {
      sections.forEach(sec => sec.style.display = 'block');
    }
  });

  // Init
  showSection(0);
})();
