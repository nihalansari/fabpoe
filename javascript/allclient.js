/*
 * SPDX-License-Identifier: Apache-2.0
 */

'use strict';
const FabricCAServices = require('fabric-ca-client');

const { FileSystemWallet, Gateway, X509WalletMixin } = require('fabric-network');
const fs = require('fs');
const path = require('path');

const ccpPath = path.resolve(__dirname, '..', '..', 'basic-network', 'connection.json');
const ccpJSON = fs.readFileSync(ccpPath, 'utf8');
const ccp = JSON.parse(ccpJSON);


module.exports = { 

invoke  : async function(request, onok, onerr) {
    try {

        // Create a new file system based wallet for managing identities.
        const walletPath = path.join(process.cwd(), 'wallet');
        const wallet = new FileSystemWallet(walletPath);
        console.log(`Wallet path: ${walletPath}`);

        // Check to see if we've already enrolled the user.
        const userExists = await wallet.exists('user1');
        if (!userExists) {
            console.log('An identity for the user "user1" does not exist in the wallet');
            console.log('Run the registerUser.js application before retrying');
            return;
        }

	// Get request parameters send by Frontend
	const requestBody = JSON.parse(request);

        // Create a new gateway for connecting to our peer node.
        const gateway = new Gateway();
        await gateway.connect(ccp, { wallet, identity: 'user1', discovery: { enabled: false } });

        // Get the network (channel) our contract is deployed to.
        const network = await gateway.getNetwork(requestBody.chainId);

        // Get the contract from the network.
        const contract = network.getContract(requestBody.chaincodeId);

        // Submit the transaction sent by Frontend App
	//Note: Based on the length of arguments present in the list a customized call to submitTransaction will be made
	switch(requestBody.args.length){
		case 0:
	        await contract.submitTransaction(requestBody.fcn); 
		break;

		case 1:
	        await contract.submitTransaction(requestBody.fcn,requestBody.args[0]); 
		break;
	
		case 2:
	        await contract.submitTransaction(requestBody.fcn,requestBody.args[0], requestBody.args[1]); 
		break;

		case 3:
	        await contract.submitTransaction(requestBody.fcn,requestBody.args[0], requestBody.args[1],requestBody.args[2]); 
		break;

		case 4:
	        await contract.submitTransaction(requestBody.fcn,requestBody.args[0], requestBody.args[1],requestBody.args[2], requestBody.args[3]); 
		break;

		case 4:
	        await contract.submitTransaction(requestBody.fcn,requestBody.args[0], requestBody.args[1],requestBody.args[2], requestBody.args[3], requestBody.args[4]); 

	}
        console.log('Transaction has been submitted:' );

        // Disconnect from the gateway.
        await gateway.disconnect();

	//callback
	onok('Transaction was successful: txn->' + requestBody);
	return;

    } catch (error) {
        console.error(`Failed to submit transaction: ${error}`);
	onerr(error.message.toString());
	return;
    }
},	//END function invoke

query : async function(request, onok, onerr) {
    try {

        // Create a new file system based wallet for managing identities.
        const walletPath = path.join(process.cwd(), 'wallet');
        const wallet = new FileSystemWallet(walletPath);
        console.log(`Wallet path: ${walletPath}`);

        // Check to see if we've already enrolled the user.
        const userExists = await wallet.exists('user1');
        if (!userExists) {
            console.log('An identity for the user "user1" does not exist in the wallet');
            console.log('Run the registerUser.js application before retrying');
            return;
        }

	// Get request parameters send by Frontend
	const requestBody = JSON.parse(request);

        // Create a new gateway for connecting to our peer node.
        const gateway = new Gateway();
        await gateway.connect(ccp, { wallet, identity: 'user1', discovery: { enabled: false } });

        // Get the network (channel) our contract is deployed to.
        const network = await gateway.getNetwork(requestBody.chainId);

        // Get the contract from the network.
        const contract = network.getContract(requestBody.chaincodeId);

        // Evaluate the specified transaction.
	//Note: Based on the length of arguments present in the list a customized call to submitTransaction will be made
	var result=0;
	switch(requestBody.args.length){
		case 0:
		result = await contract.evaluateTransaction(requestBody.fcn);
		break;

		case 1:
	        result = await contract.evaluateTransaction(requestBody.fcn,requestBody.args[0]);
		break;
	
		case 2:
	        result = await contract.evaluateTransaction(requestBody.fcn,requestBody.args[0],requestBody.args[1]);
		break;

		case 3:
	        result = await contract.evaluateTransaction(requestBody.fcn,requestBody.args[0],requestBody.args[1],requestBody.args[2]);
		break;

		case 4:
	        result = await contract.evaluateTransaction(requestBody.fcn,requestBody.args[0],requestBody.args[1],requestBody.args[2],requestBody.args[3]);
		break;


	}
        console.log(`Transaction has been evaluated, result is: ${result.toString()}`);
	onok(result.toString());
	return;

    } catch (error) {
        console.error(`Failed to evaluate transaction: ${error}`);
	onerr(error.message.toString());
	return;
    }
},// END function query

registeruser : async function (req,onok,onerr) {
    try {

        // Create a new file system based wallet for managing identities.
        const walletPath = path.join(process.cwd(), 'wallet');
        const wallet = new FileSystemWallet(walletPath);
        console.log(`Wallet path: ${walletPath}`);

        // Check to see if we've already enrolled the user.
        const userExists = await wallet.exists('user1');
        if (userExists) {
            console.log('An identity for the user "user1" already exists in the wallet');
            return;s
        }

        // Check to see if we've already enrolled the admin user.
        const adminExists = await wallet.exists('admin');
        if (!adminExists) {
            console.log('An identity for the admin user "admin" does not exist in the wallet');
            console.log('Run the enrollAdmin.js application before retrying');
            return;
        }

        // Create a new gateway for connecting to our peer node.
        const gateway = new Gateway();
        await gateway.connect(ccp, { wallet, identity: 'admin', discovery: { enabled: false } });

        // Get the CA client object from the gateway for interacting with the CA.
        const ca = gateway.getClient().getCertificateAuthority();
        const adminIdentity = gateway.getCurrentIdentity();

        // Register the user, enroll the user, and import the new identity into the wallet.
        const secret = await ca.register({ affiliation: 'org1.department1', enrollmentID: 'user1', role: 'client' }, adminIdentity);
        const enrollment = await ca.enroll({ enrollmentID: 'user1', enrollmentSecret: secret });
        const userIdentity = X509WalletMixin.createIdentity('Org1MSP', enrollment.certificate, enrollment.key.toBytes());
        wallet.import('user1', userIdentity);
        console.log('Successfully registered and enrolled admin user "user1" and imported it into the wallet');
	onok('Successfully registered and enrolled admin user "user1" and imported it into the wallet');
	return;

    } catch (error) {
        console.error(`Failed to register user "user1": ${error}`);
	onerr(error.message.toString());
        return;
    }
},// END function registeruser


enrollAdmin: async function(onok,err) {



    try {

        // Create a new CA client for interacting with the CA.
        const caURL = ccp.certificateAuthorities['ca.example.com'].url;
        const ca = new FabricCAServices(caURL);

        // Create a new file system based wallet for managing identities.
        const walletPath = path.join(process.cwd(), 'wallet');
        const wallet = new FileSystemWallet(walletPath);
        console.log(`Wallet path: ${walletPath}`);

        // Check to see if we've already enrolled the admin user.
        const adminExists = await wallet.exists('admin');
        if (adminExists) {
            console.log('An identity for the admin user "admin" already exists in the wallet');
	    onok('An identity for the admin user "admin" already exists in the wallet');
            return;
        }

        // Enroll the admin user, and import the new identity into the wallet.
        const enrollment = await ca.enroll({ enrollmentID: 'admin', enrollmentSecret: 'adminpw' });
        const identity = X509WalletMixin.createIdentity('Org1MSP', enrollment.certificate, enrollment.key.toBytes());
        wallet.import('admin', identity);
        console.log('Successfully enrolled admin user "admin" and imported it into the wallet');
	success('Successfully enrolled admin user "admin" and imported it into the wallet');
	onok('Successfully enrolled admin user "admin" and imported it into the wallet');
	return;

    } catch (error) {
        console.error(`Failed to enroll admin user "admin": ${error}`);
	onerr(error.message.toString());
	return;
        //process.exit(1);
    } //end try block
} //end function

} //end exports
