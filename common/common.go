package common

const DefaultOutFile = "mesa_trace_raw.csv"
const TotalFlag = 200
const FlagListFile = "flagList.csv"
const DefaultCallToCompare = "GLX_SWAP_BUFFERS"
const Threshold = 100

const GLX_SWAP_BUFFERS_FLAG = 30
//const defineFmt = `#define ([A-Z_]+) ([0-9]+)`

// Package common contains the common interfaces used by other module.

type GetFlagName func (cgoType int) string
