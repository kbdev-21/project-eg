package auth

import (
	"backend/src/db"
	"backend/src/shared"
)

var adjs = []string{
	"Lazy",
	"Angry",
	"Silly",
	"Sleepy",
	"Chubby",
	"Happy",
	"Fluffy",
	"Fancy",
	"Lucky",
	"Grumpy",
}

var avtCodes = []db.UserAvtCode{
	db.UserAvtCodeBUNNY,
	db.UserAvtCodeKITTEN,
	db.UserAvtCodeGRIZZLE,
	db.UserAvtCodeHAMSTER,
	db.UserAvtCodeMONKEY,
}

func getNameFromAdjAndAvtCode(adj string, avtCode db.UserAvtCode) string {
	return adj + shared.CapitalizeString(string(avtCode))
}

func genRandomIdentity() (db.UserAvtCode, string) {
	randAvtCode := avtCodes[shared.RandomInt(0, len(avtCodes)-1)]
	randAdj := adjs[shared.RandomInt(0, len(adjs)-1)]
	randName := getNameFromAdjAndAvtCode(randAdj, randAvtCode)
	return randAvtCode, randName
}
