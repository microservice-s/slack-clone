"use strict";

const express = require('express'),
    { Wit } = require('node-wit'), 
    bodyParser = require('body-parser'),
    mongo = require('mongodb'),
    MongoClient = mongo.MongoClient,
   // ObjectID = require('mongodb').ObjectID,
    co = require('co'),
    assert = require('assert'),
    moment = require('moment');

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

const dbAddr = process.env.DBADDR;
if (!dbAddr) {
	console.error("please set DBADDR to where your database is");
	process.exit(1);
}

// set up the mongo connection
co(function*() {
  // Connection URL
  var url = `mongodb://${dbAddr}/production`;
  console.log("Connected correctly to server");
  // Use connect method to connect to the Server
  app.locals.db = yield MongoClient.connect(url);
  
}).catch(function(err) {
  console.log(err.stack);
});

// who is in the XYZ channel?
function handleUsers (req, res, data) {
    if (!data.entities.channel) {
        res.send("I'm sorry dave, I'm afraid I can't do that. You didn't provide a channel for me to query.");
        return;
    }
     co(function*() {
        const db = req.app.locals.db;

        // get the channel the user is asking for
        let witChan = data.entities.channel[0].value;
        let channel = yield db.collection('channels').find({name: new RegExp(witChan, "i")}, {_id: 0, members: 1}).limit(1).toArray();
        if (channel.length == 0) {
            res.send("I'm sorry dave, I'm afraid I can't do that.");
            return;
        } else if (!channel[0].members || channel[0].members.length == 0) {
            res.send("I'm sorry dave, there are no members in that channel.");
            return;
        }
        let memberNames = [];
        // get the users first name
        for(var i = 0; i < channel[0].members.length; i++) {
            console.log(channel[0].members[i])
            let member = yield db.collection('users').find({_id: new mongo.ObjectID(channel[0].members[i])}).limit(1).toArray();
            if (member.length > 0) {
                memberNames.push(member[0].firstname);
            }
        }

        // respond to the user 
        if (memberNames.length == 1) {
            res.send(`${memberNames[0]} is in the ${witChan} channel`);
        } else {
            res.send(`${memberNames.join()} are in the ${witChan} channel`);
        }
     }).catch(function(err) {
        res.status(500).send();
        console.log(err.stack);
     });
}

function handleUsersNegation (req, res, data) {
    if (!data.entities.channel) {
        res.send("I'm sorry dave, I'm afraid I can't do that. You didn't provide a channel for me to query.");
        return;
    }

    co(function*() {
        const db = req.app.locals.db;

        // get the channel the user is asking for
        let witChan = data.entities.channel[0].value;
        let channel = yield db.collection('channels').find({name: new RegExp(witChan, "i")}, { members: 1}).limit(1).toArray();
        if (channel.length == 0) {
            res.send("I'm sorry dave, I'm afraid I can't do that.");
            return;
        } else if (!channel[0].members || channel[0].members.length == 0) {
            res.send("I'm sorry dave, there are no members in that channel.");
            return;
        }

        let noPosts = [];
        // query the database for each member, and add it to the list if there are no results
        for (var i = 0; i < channel[0].members.length; i++) {
            let user = yield db.collection('messages').find({$and : [{creatorid: new mongo.ObjectID(channel[0].members[i])}, {channelid : new mongo.ObjectID(channel[0]._id)}]}).toArray();
            if (user.length == 0) {
                noPosts.push(channel[0].members[i]);
            }
        }

        let memberNames = [];
        for (var j = 0; j < noPosts.length; j++) {
            console.log(noPosts[j])
            let user = yield db.collection('users').find({_id: new mongo.ObjectID(noPosts[j])}, {_id: 0, firstname: 1}).limit(1).toArray();
            if (user.length != 0) {
                memberNames.push(user[0].firstname);
            }   
        }

        // respond to the user 
        if (memberNames.length == 0) {
            res.send(`No users haven't posted to the ${witChan} channel`);
        } else if (memberNames.length == 1) {
            res.send(`Only ${memberNames[0]} hasn't posted to the ${witChan} channel`);
        } else {
            res.send(`${memberNames.join()} haven't posted to the ${witChan} channel`);
        }
        
     }).catch(function(err) {
        res.status(500).send();
        console.log(err.stack);
     });
}

function handleHigestPoster (req, res, data) {
    if (!data.entities.channel) {
        res.send("I'm sorry dave, I'm afraid I can't do that. You didn't provide a channel for me to query.");
        return;
    }
    
    co(function*() {

        const db = req.app.locals.db;

        // get the channel the user is asking for
        let witChan = data.entities.channel[0].value;
        let channel = yield db.collection('channels').find({name: new RegExp(witChan, "i")}).limit(1).toArray();
        if (channel.length == 0) {
            res.send("I'm sorry dave, I'm afraid I can't do that.");
            return;
        }
        let cID = new mongo.ObjectID(channel[0]._id);
        // Get the collection
        let col = db.collection('messages');
        let count = 0;
        
        let ag = yield col.aggregate([
            {$match: {channelid: cID}},
            {$group: {_id : "$creatorid", count: {$sum: 1}}},
            {$sort : {count : -1}}
        ]).toArray();
        if (ag.length == 0) {
            res.send("No users have posted to that channel");
            return;
        }
        console.log(ag);
        let hPosters = [];
        hPosters.push(ag[0]);
        for (let i = 1; i < ag.length; i++) {
            if (ag[i].count == hPosters[i-1].count) {
                hPosters.push(ag[i]);
            } else {
                break;
            }
        }
        let nameArray = [];
        // go through and query the db for the names of the users
        for(var i = 0; i < hPosters.length; i++) {
            let user =  yield db.collection('users').find({_id: new mongo.ObjectID(hPosters[i]._id)}).toArray();
            nameArray.push(user[0].firstname);
        }
        // respond to the user
        if (nameArray.length > 1) {
            res.send(`${nameArray.join()} have made the most posts to the ${witChan} channel`)
        } else {
            res.send(`${nameArray[0]} has made the most posts to the ${witChan} channel`)
        }

     }).catch(function(err) {
        res.status(500).send();
        console.log(err.stack);
    });
}

function handlePostTime (req, res, data) {
    
    // check if they are asking about a specific channel or generally
    co(function*() {
        // get the userid
        let user = JSON.parse(req.get("X-User"));
        let userID = user.id;
        const db = req.app.locals.db;
        // Get the collection
        let col = db.collection('messages');
        let message = '';
        let channelR = '';
        if (data.entities.channel) {
            // find the channel id from the name
            let witChan = data.entities.channel[0].value;
            let channel = yield db.collection('channels').find({name: new RegExp(witChan, "i")}).limit(1).toArray();
            console.log(witChan)
            let cID = new mongo.ObjectID(channel[0]._id);
            message = yield col.find({ $and : [{creatorid: new mongo.ObjectID(userID)}, {channelid: cID}]} ).sort({createdat: -1}).limit(1).toArray();
            channelR = ` to the ${witChan} channel`
        } else {
            message = yield col.find({creatorid: new mongo.ObjectID(userID)}).sort({createdat: -1}).limit(1).toArray();
        }

        if (message.length == 0) {
            res.send("I'm sorry dave, I'm afraid I can't do that.");
        } else {
            let m = moment(message[0].createdat)
            let time = m.format("h:mm a")
            let day = ordinalSuffix(m.format("DD"));
            let month = m.format("MMM");
            res.send(`The last time you posted${channelR} was on ${month} ${day} at ${time}`);
        }
        
        // Get first two documents that match the query
        //var docs = yield col.find({a:1}).limit(2).toArray();
    }).catch(function(err) {
        res.status(500).send();
        console.log(err.stack);
    });

}

function handlePostCount (req, res, data) {
    if (!data.entities.channel) {
        res.send("I'm sorry dave, I'm afraid I can't do that. You didn't provide a channel for me to query.");
        return;
    }
    
    co(function*() {
        // get the userid
        let user = JSON.parse(req.get("X-User"));
        let userID = user.id;
        const db = req.app.locals.db;
        let witChan = data.entities.channel[0].value;
        let channel = yield db.collection('channels').find({name: new RegExp(witChan, "i")}).limit(1).toArray();
        if (channel.length == 0) {
            res.send("I'm sorry dave, I'm afraid I can't do that.");
            return;
        }
        let cID = new mongo.ObjectID(channel[0]._id);
        // Get the collection
        let col = db.collection('messages');
        let count = 0;
        if (data.entities.datetime) {
            // get the beginning and end of the day
            let start = new Date(data.entities.datetime[0].value);
            let end = new Date(start.getTime() + 86400000);
            count = yield col.count({ $and : [{creatorid: new mongo.ObjectID(userID)}, {channelid: cID}, {createdat: {$gte: start, $lt: end}}]});
            let m = moment(data.entities.datetime[0].value);
            let day = ordinalSuffix(m.format("DD"));
            let month = m.format("MMM");
            res.send(`You posted to the ${channel[0].name} ${count} times on ${month} ${day}`);
        } else {
            count = yield col.count({ $and : [{creatorid: new mongo.ObjectID(userID)}, {channelid: cID}]});
            console.log(count);
            res.send(`You have posted to the ${channel[0].name} channel ${count} times`);
        }

     }).catch(function(err) {
        res.status(500).send();
        console.log(err.stack);
    });
}

function ordinalSuffix(i) {
    i = parseInt(i);
    var j = i % 10,
        k = i % 100;
    if (j == 1 && k != 11) {
        return i + "st";
    }
    if (j == 2 && k != 12) {
        return i + "nd";
    }
    if (j == 3 && k != 13) {
        return i + "rd";
    }
    return i + "th";
}



app.post("/v1/bot", (req, res, next) => {
	//use witaiClient.message() to
	//extract meaning from the value in the
	//`q` query string param and respond
	//accordingly
    console.log("handling route")
	let q = req.body;
	console.log(`user is asking ${q}`);
	witaiClient.message(q)
		.then(data => {
			console.log(JSON.stringify(data, undefined, 2));
            if (!data.entities.intent) {
                res.send("I'm sorry dave, I'm afraid I can't do that.");
                return;
            }
			switch (data.entities.intent[0].value) {
				case "get_users":
					handleUsers(req, res, data);
				    break;
                case "get_highest_poster":
					handleHigestPoster(req, res, data);
				    break;
                case "get_post_time":
                    handlePostTime(req, res, data);
				    break;
                case "get_post_count":
					handlePostCount(req, res, data);
				    break;
                case "get_users_negation":
					handleUsersNegation(req, res, data);
				    break;
				default:
					res.send("I'm sorry dave, I'm afraid I can't do that.");
			}
		})
		.catch(next); // passes the error to express (who will report)

});

//error handler
app.use((err, req, res, next) => {
    console.error(err);
    res.status(err.status || 500).send(err.message);
});

app.listen(port, host, () => {
	console.log(`server is listening at http://${host}:${port}`);
});