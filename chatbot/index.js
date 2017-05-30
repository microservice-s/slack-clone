"use strict";

const express = require('express');
const { Wit } = require('node-wit');
const bodyParser = require('body-parser');

const app = express();
// set up body parser middleware to parse any text from post requests (queries)
app.use(bodyParser.text());
const port = process.env.PORT || '80';
const host = process.env.HOST || '';
const witaiToken = process.env.WITAITOKEN;

if (!witaiToken) {
	console.error("please set WITAITOKEN to your wit.ai app token");
	process.exit(1);
}

const witaiClient = new Wit({ accessToken: witaiToken });

app.post("/bot", (req, res, next) => {
	//use witaiClient.message() to
	//extract meaning from the value in the
	//`q` query string param and respond
	//accordingly
	let q = req.body;
	console.log(`user is asking ${q}`);
	witaiClient.message(q)
		.then(data => {
			console.log(JSON.stringify(data, undefined, 2));
			switch (data.entities.intent[0].value) {
				case "get_users":
					
				    break;
                case "get_highest_poster":
					
				    break;
                case "get_post_time":
					
				    break;
                case "get_post_count":
					
				    break;
                case "get_users_negation":
					
				    break;
				default:
					res.send("I'm sorry dave, I'm afraid I can't do that.");
			}
		})
		.catch(next); // passes the error to express (who will report)


});

app.listen(port, host, () => {
	console.log(`server is listening at http://${host}:${port}`);
});