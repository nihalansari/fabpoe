const express = require('express')
const app = express()
var bodyParser = require('body-parser')
//const ExampleNetwork = require('./query');
	
var poeClient = require('./allclient');
var ExampleAdmin = require('./enrollAdmin');
var ExampleUser = require('./registerUser');
const fs = require('fs');



//Attach the middleware
app.use( bodyParser.json() );
app.use(bodyParser.urlencoded({     // to support URL-encoded bodies
  extended: true
})); 


app.get('/api', function(req, res) {
		console.log("In endpoint api");
        
		var rawdata = fs.readFileSync('./Testing/POE.postman_collection.json');
				res.end(rawdata);
			
}); 


//create the user=admin who will create other users
//should this endpoint be called each time a new user registers?? 
// Answer) NO. users from same organization will have only 1 admin.
// Not sure if we need this endpoint or not. an admin may well be created
// at the time the blockchain network is deployed
app.post('/enrolladmin', function(req, res) {
		console.log("In endpoint enrolladmin");
        var data = req.body;
		var adm = ExampleAdmin.fabcarenroll(function(data){
			//callback here
			console.log("inside OK callback enrolladmin");
			console.log(data);
			res.end(data);
			
		},
		
		function(data){ //error
			//callback here
			console.log("inside error callback enrolladmin");
			console.log(data);
			res.end(data);
			
		}
		
		); 
	});

//register a new user idenitified by "userID"
//Note: The register function must accept a parameter synonymous to UserId		
//      Because many users will use this server. Whenever a user registers,
//      a new set of priv/pub/cert are generated, indentified by userId.
//      whenever any invoke or query operation is performed on the chaincode, the 
//      client program will need to use the corresponding hfc-key-store credentials
//      as per the logged in user
app.post('/registeruser', function(req, res) {
		
        var data = JSON.stringify(req.body);
		console.log("In endpoint registeruser payload=" + data);
		var usr = ExampleUser.registernewuser(
			data,
			function(respdata){
				//callback here
				console.log("registerUser: OK callback");
				res.end(respdata);
			},
			function(resperr){
				//callback here
				console.log("registerUser: ERR callback");
				res.end(resperr);
			}
			
		);
		});
		

//query an existing asset by ID
app.post('/query', function(req, res) {
		var data = JSON.stringify(req.body);
		console.log("In endpoint query payload=" + data);
        
		poeClient.query(
			data, //input json from frontend/client
			function(respdata){
				//callback here
				console.log("query:OK callback");
				res.end(respdata);
				//res.end();
			},
			function(resperr){
				//callback here
				console.log("query:ERR callback");
				res.end(resperr);
				//res.end();
			},	
		);
});

//query an existing asset by ID
app.post('/invoke', function(req, res) {
		var data = JSON.stringify(req.body);
		console.log("In endpoint invoke payload=" + data);
        
		poeClient.invoke(
			data, //input json from frontend/client
			function(respdata){
				//callback here
				console.log("invoke:OK callback");
				res.end(respdata);
				//res.end();
			},
			function(resperr){
				//callback here
				console.log("invoke:ERR callback");
				res.end(resperr);
				//res.end();
			},	
		);
});


/* NOTE: other endpoints like block query, blockchain height, system level information
         will go next 																*/
	

//******************** SERVER ******************************************/
//create a server that listen to all the requests for the above endpoints				
var server = app.listen(3000, function() {
    console.log("Listening on port %s...", server.address().port);
});
//******************** SERVER ******************************************/

