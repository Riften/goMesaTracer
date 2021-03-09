package common

// DefaultOutFile ...
const DefaultOutFile = "mesa_trace_raw.csv"

// TotalFlag ...
const TotalFlag = 200

// FlagListFile ...
const FlagListFile = "flagList.csv"

// DefaultCallToCompare ...
const DefaultCallToCompare = "GLX_SWAP_BUFFERS"

// Threshold ...
const Threshold = 150

// GlxSwapBuffersFlag ...
const GlxSwapBuffersFlag = 30

//const defineFmt = `#define ([A-Z_]+) ([0-9]+)`

// Package common contains the common interfaces used by other module.

// GetFlagName ...
type GetFlagName func(cgoType int) string
