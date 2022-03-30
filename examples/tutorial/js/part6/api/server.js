const express = require('express');
const path = require('path');
const cors = require("cors");
const app = express(),
bodyParser = require("body-parser");
port = 3080;

app.use(bodyParser.json());
app.use(express.static(path.join(__dirname, '../my-app/build')));

app.use(cors());

const users = [
	{
		'first_name': 'Lee',
		'last_name' : 'Earth'
	}
]

app.get('/api/users', (req, res) => {
  console.log('api/users called!')
  res.json(users);
});

app.listen(port, '0.0.0.0', () => {
  console.log(`Server listening on the port::${port}`);
});
