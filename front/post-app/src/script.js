let dropbox;

dropbox = document.getElementById("dropzone");
dropbox.addEventListener("dragenter", dragenter, false);
dropbox.addEventListener("dragover", dragover, false);
dropbox.addEventListener("drop", drop, false);

function dragenter(e) {
  e.stopPropagation();
  e.preventDefault();
}

function dragover(e) {
  e.stopPropagation();
  e.preventDefault();
}

function drop(e) {
  e.stopPropagation();
  e.preventDefault();

  const dt = e.dataTransfer;
  const files = dt.files;

  handleFiles(files);
}

function handleFiles(files) {
  for (let i = 0; i < files.length; i++) {
    const file = files[i];

    if (!file.type.startsWith("image/")) {
      continue;
    }

    const img = document.createElement("img");
    img.classList.add("post", "to-upload");
    img.file = file;
    let preview = document.getElementById("preview-list");
    preview.appendChild(img);

    const reader = new FileReader();
    reader.onload = (e) => {
      img.src = e.target.result;
    };
    reader.readAsDataURL(file);
  }
}

async function uploadImage(file, content) {
  try {
    const response = await fetch("/posts", {
      method: "POST",
      headers: {
        "Content-Type": file.type,
      },
      body: content,
    })
    return response.text();
  } catch (err) {
    console.error("Error uploading image: " + err.message);
    throw err;
  }
}


function FileUpload(file) {
  const reader = new FileReader();
  reader.onload = (evt) => {
    uploadImage(file, evt.target.result)
      .then(() => {
        updateFileList();
      })
      .catch(console.error);
  };
  reader.readAsArrayBuffer(file);
}

async function getImageList() {
  try {
    const response = await fetch("/posts", {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
      }
    });
    return response.json();
  } catch (err) {
    console.error("Error listing images: " + err.message);
    throw err;
  }
}

async function getImageElements(gallery) {
  const response = await getImageList();
  let imgElements = [];
  const images = response;
  for (const image of images) {
    const img = document.createElement("img");
    img.classList.add("post");
    img.src = "data:" + image.mime + ";base64," + image.image;
    imgElements.push(img);
  }
  return imgElements;
}

async function updateFileList() {
  let gallery = document.getElementsByClassName("gallery")[0];
  gallery.innerHTML = "";
  const loading = document.createElement("div");
  loading.classList.add("loading");
  gallery.appendChild(loading);
  const imgElements = await getImageElements(gallery);
  gallery.innerHTML = "";
  imgElements.forEach(img => gallery.appendChild(img));
}

function clearPreviewList() {
  let preview = document.getElementById("preview-list");
  preview.innerHTML = "";
}

function sendFiles(submit) {
  const imgs = document.querySelectorAll(".to-upload");
  if (imgs.length == 0) return;

  submit.disabled = true;

  for (let i = 0; i < imgs.length; i++) {
    new FileUpload(imgs[i].file);
    imgs[i].innerHTML = "";
  }

  clearPreviewList();
  updateFileList();

  submit.disabled = false;
}

let submit = document.getElementById("submit");
submit.addEventListener("click", (submit) => sendFiles(submit), false);

document.addEventListener("DOMContentLoaded", function(){
  updateFileList();
});
