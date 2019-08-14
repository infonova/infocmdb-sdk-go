package main

// This library is used for the communication with the infoCMDB
//
// The API provided by the infoCMDB has different versions and therefore this library is split up v1/v2.
//
// v1 - custom HTTP API
//
// Is the legacy version which is based on configured sql queries provided via a custom http api
//
// v2 - Restful API
//
// This is the first re-engineering or the api to access core models and services to have native access.
// This api properly handles all permission checks and access to native functions.