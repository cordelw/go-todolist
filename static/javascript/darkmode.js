let darkmode = localStorage.getItem("dark");

function enableDarkmode() {
  document.body.classList.add("dark-mode");
  localStorage.setItem("dark", "true");
}

function disableDarkmode() {
  document.body.classList.remove("dark-mode");
  localStorage.setItem("dark", "false");
}

function toggleDarkmode() {
  darkmode = localStorage.getItem("dark");
  darkmode !== "true" ? enableDarkmode() : disableDarkmode();
}

if (darkmode === "true") enableDarkmode();
