package main

import (
	"git.appteam.infonova.cloud/AppTeam/infocmdbGoLib"
	"log"
	"net/url"
)

func main() {
	i := infocmdbGoLib.NewInfoCmdbGoLib()
	err := i.Login("http://infocmdb.local", "admin", "admin")
	if err != nil {
		log.Fatal(err)
	}

	params := url.Values{
		"argv1": {"1"},
	}

	ret, err := i.WC.CallGet("query", "int_getListOfCiIdsOfCiType", params)
	if err != nil {
		log.Panic("Error: ", err)

	}
	log.Println("Get: ", ret)

	ret, err = i.WC.CallPost("query", "int_getListOfCiIdsOfCiType", params)
	if err != nil {
		log.Panic("Error: ", err)
	}
	log.Println("Post: ", ret)

	r, err := i.GetListOfCiIdsOfCiType(1)
	if err != nil {
		log.Panic("WS Error", err)
	}
	log.Println(r)
	rCi, err := i.GetCi(1)
	if err != nil {
		log.Panic("WS Error", err)
	}
	log.Println(rCi)

	rCiA, err := i.GetCiAttributes(1)
	if err != nil {
		log.Panic("WS Error", err)
	}
	log.Println(rCiA)
}
