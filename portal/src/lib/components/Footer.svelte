<footer>

  <!-- Upper: white -->
  <div class="bg-white border-t border-[var(--border)] px-4 py-3">
    <div class="grid grid-cols-1 sm:grid-cols-3 gap-8">

      <!-- Security notices -->
      <div class="space-y-2">
        <div class="flex items-center gap-2">
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="#1f2978" stroke-width="2" class="h-5 w-5 shrink-0">
            <path stroke-linejoin="miter" stroke-linecap="square" d="M12 3L22.5 21H1.5L12 3Z" />
            <line x1="12" y1="9" x2="12" y2="14.5" stroke-linecap="square" />
            <rect x="11.25" y="16.5" width="1.5" height="1.5" fill="#1f2978" stroke="none" />
          </svg>
          <h2 class="text-base font-semibold text-[var(--foreground)]">Sicherheitshinweise</h2>
        </div>
        <p class="text-sm text-black-600 leading-relaxed">
          Achtung! Betrüger geben sich am Telefon als Bankmitarbeitende aus.
        </p>
        <div class="space-y-1">
          <a href="#" class="block text-sm font-semibold text-[#1f2978] underline hover:opacity-75">Link zu den aktuellen Sicherheitshinweisen</a>
              <a href="#" class="block text-sm font-semibold text-[#1f2978] underline hover:opacity-75">Link zu Sicherheit im Überblick</a>
        </div>
      </div>

      <!-- Online banking -->
      <div class="space-y-2">
        <div class="flex items-center gap-2">
          <img src="/laptop_mobile_icon.png" alt="" class="h-7 w-auto shrink-0" />
          <h2 class="text-base font-semibold text-[var(--foreground)]">Online-Banking Zugang</h2>
        </div>
        <p class="text-sm text-black-600 leading-relaxed">
          Hier können Sie Ihren persönlichen Zugang zum Online-Banking beantragen.
        </p>
        <a href="#" class="block text-sm font-semibold text-[#1f2978] underline hover:opacity-75">Zugang zum Online-Banking beantragen</a>
      </div>

      <!-- Security procedures -->
      <div class="space-y-2">
        <div class="flex items-center gap-2">
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="#1f2978" stroke-width="2" class="h-5 w-5 shrink-0">
            <!-- shackle: rounded arch -->
            <path stroke-linecap="round" stroke-linejoin="round" d="M8 11V7a4 4 0 018 0v4" />
            <!-- body: sharp rectangle -->
            <rect x="5" y="11" width="14" height="11" stroke-linejoin="miter" />
            <!-- dot in center -->
            <rect x="10.5" y="15" width="3" height="3" fill="#1f2978" stroke="none" />
          </svg>
          <h2 class="text-base font-semibold text-[var(--foreground)]">Unsere Sicherheitsverfahren</h2>
        </div>
        <p class="text-sm text-black-600 leading-relaxed">
          Alles Wissenswerte rund um Ihren Login.
        </p>
        <a href="#" class="block text-sm font-semibold text-[#1f2978] underline hover:opacity-75">Link zu den Sicherheitsverfahren</a>
      </div>

    </div>
  </div>

  <!-- Lower: dark blue -->
  <div class="bg-[#1E2A78] px-4 py-5">
    <div class="space-y-3">

      <div class="flex flex-wrap items-center gap-x-4 gap-y-2">
        <a href="#" class="text-sm font-bold text-white hover:text-white/70 transition-colors">English Version</a>
        <a href="#" class="text-sm font-bold text-white hover:text-white/70 transition-colors">Hilfe</a>
        <a href="#" class="text-sm font-bold text-white hover:text-white/70 transition-colors">Demo-Konto</a>
        <a href="#" class="text-sm font-bold text-white hover:text-white/70 transition-colors">Impressum</a>
        <a href="#" class="text-sm font-bold text-white hover:text-white/70 transition-colors">Rechtliche Hinweise</a>
        <a href="#" class="text-sm font-bold text-white hover:text-white/70 transition-colors">Datenschutz</a>
        <a href="#" class="text-sm font-bold text-white hover:text-white/70 transition-colors">Cookie-Einstellungen</a>
      </div>
      <div>
        <button class="w-full sm:w-auto bg-[#0D43B3] hover:bg-[#0a389a] px-5 py-2 text-sm font-semibold text-white transition-colors">
          Vertrag widerrufen
        </button>
      </div>

      <p class="text-sm text-white">
        © {new Date().getFullYear()} Dagual Bank AG..
      </p>

    </div>
  </div>

</footer>
