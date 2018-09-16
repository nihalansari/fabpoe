package main

import (
		pb "github.com/hyperledger/fabric/protos/peer"
)

//=================================================================================================
//================================================================================= RETURN HANDLING
// Return handling: for return, we either return "shim.Success (payload []byte) with HttpRetCode=200"
// or "shim.Error(doc string) with HttpRetCode=500". However, we want to set our own status codes to
// map into HTTP return codes. A few utility functions:
/*
func Success(rc int32, doc string, payload []byte) pb.Response {
	return pb.Response{
		Status:  rc,
		Message: doc,
		Payload: payload,
	}
}

func Error(rc int32, doc string) pb.Response {
	logger.Errorf("Error %d = %s", rc, doc)
	return pb.Response{
		Status:  rc,
		Message: doc,
	}
}
*/
func routeShimError(message string) pb.Response{
	
		return pb.Response{
		Status:  -1,
		Message: message,
		Payload: nil,
	}
}

func routeShimSuccess(payload []byte) pb.Response {

	
		return pb.Response{
		Status:  0,
		Message: "OK",
		Payload: payload,
	}
}