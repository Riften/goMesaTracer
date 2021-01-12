package common

const DefaultOutFile = "mesa_trace_raw.csv"
const TotalFlag = 100
const FlagListFile = "flagList.csv"
const DefaultCallToCompare = "GLX_SWAP_BUFFERS"

//const defineFmt = `#define ([A-Z_]+) ([0-9]+)`

// Package common contains the common interfaces used by other module.

type GetFlagName func (cgoType int) string
