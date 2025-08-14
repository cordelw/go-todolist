// SHOW UNFINISHED / FINISHED list

function showComplete() {
  document.getElementById("unfinished-tasks").style.display = "none";
  document.getElementById("finished-tasks").style.display = "block";
  document.getElementById("show-complete-button").style.display = "none";
  document.getElementById("show-back-button").style.display = "inline";
}

// TASK EDIT / COMPLETE FORMS
function enableTaskEdit(id) {
  // Swap label for input
  document.getElementById("item-label-" + id).hidden = true;
  document.getElementById("item-input-" + id).hidden = false;

  // Swap buttons
  document.getElementById("item-edit-button-" + id).hidden = true;
  document.getElementById("item-complete-button-" + id).hidden = true;
  document.getElementById("item-complete-button-" + id).type = "button";
  document.getElementById("item-unedit-button-" + id).hidden = false;
  document.getElementById("item-update-button-" + id).hidden = false;
  document.getElementById("item-update-button-" + id).type = "submit";
}

function disableTaskEdit(id) {
  // Swap label for input
  document.getElementById("item-label-" + id).hidden = false;
  document.getElementById("item-input-" + id).hidden = true;

  // Swap buttons
  document.getElementById("item-edit-button-" + id).hidden = false;
  document.getElementById("item-complete-button-" + id).hidden = false;
  document.getElementById("item-complete-button-" + id).type = "submit";
  document.getElementById("item-unedit-button-" + id).hidden = true;
  document.getElementById("item-update-button-" + id).hidden = true;
  document.getElementById("item-update-button-" + id).type = "button";
}
