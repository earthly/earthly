async function getUsers() {

  const response = await fetch('http://0.0.0.0:3080/api/users');
  return await response.json();
}

function component() {
  const element = document.createElement('div');
  getUsers()
    .then( users => {
      element.innerHTML = `hello world <b>${users[0].first_name} ${users[0].last_name}</b>`
    })

	return element;
}

document.body.appendChild(component());
