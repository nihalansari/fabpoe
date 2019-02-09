const express = require('express')
const app = express()
var bodyParser = require('body-parser')
//const ExampleNetwork = require('./query');
	
var ExampleNetwork = require('./queryAllDocs');
var ExampleAdmin = require('./enrollAdmin');
var ExampleUser = require('./registerUser');
//Attach the middleware
app.use( bodyParser.json() );

//create the user=admin who will create other users
//should this endpoint be called each time a new user registers?? 
// Answer) NO. users from same organization will have only 1 admin.
app.get('/enrolladmin', function(req, res) {
		console.log("In endpoint enrolladmin");
        var data = req.body.data;
		var adm = ExampleAdmin.fabcarenroll(function(data){
			//callback here
			console.log("inside callback enrolladmin");
			console.log(data);
			res.end(data);
			
		}); 
	});

//register a new user idenitified by "userID"
//Note: The register function must accept a parameter synonymous to UserId		
//      Because many users will use this server. Whenever a user registers,
//      a new set of priv/pub/cert are generated, indentified by userId.
//      whenever any invoke or query operation is performed on the chaincode, the 
//      client program will need to use the corresponding hfc-key-store credentials
//      as per the logged in user
app.get('/registeruser', function(req, res) {
		console.log("In endpoint query");
        var data = req.body.data;
		var usr = ExampleUser.fabcarreg(function(data){
			//callback here
			console.log("inside callback registerUser");
		});
		console.log("returned from registerUser?");
		console.log(usr);
		res.end();
		});
		
//query All assets
app.get('/queryalldocs', function(req, res) {
		console.log("In endpoint queryAlldocs");
        var data = req.body.data;
		ExampleNetwork.fabcarqueryall(function(data){
			//callback here
			res.end(data);
			//res.end();
		});
		
		});


		
				
var server = app.listen(3000, function() {
    console.log("Listening on port %s...", server.address().port);
});

