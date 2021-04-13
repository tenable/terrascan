// Replace all json-objects elements to be a JSON tree
let jsonElements = document.getElementsByClassName("json-object")
for (let val of jsonElements) {
  let element = val;
  if (element.innerText.length < 1) {
    continue
  }

  let data = JSON.parse(element.innerText);
  element.innerText = ""

  jsonTree.create(data, element);
}

// Replace all time-object elements to be in the 'DD/MM/YYYY hh:mm:ss A' format of moment.js
let timeElements = document.getElementsByClassName("time-object")
for (let val of timeElements) {
  let element = val;
  let elapsedTimeUntilNow = Date.now() - new Date(element.innerText)
  if (elapsedTimeUntilNow / 1000 < 120) {
    // In case elapsed less than 2 minutes, show "A few seconds ago"
    element.innerText = moment(element.innerText).fromNow()
  }
  else {
    element.innerText = moment(element.innerText).format('DD/MM/YYYY hh:mm:ss A')
  }
}

// Change the colors of the review status
let statusElements = document.getElementsByClassName("review-status")
for (let val of statusElements) {
  let element = val;
  switch (element.innerText) {
    case "Allowed":
      element.classList.add("allowed")
      break
    case "Rejected":
      element.classList.add("rejected")
      break
    default:
      element.classList.add("warn")
      break
  }
}