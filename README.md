# Securely and durably back up your private key using QR codes printed on paper

[![GoDoc](https://godoc.org/github.com/mtraver/qrbak?status.svg)](https://godoc.org/github.com/mtraver/qrbak)
[![Go Report Card](https://goreportcard.com/badge/github.com/mtraver/qrbak)](https://goreportcard.com/report/github.com/mtraver/qrbak)

You know that your private key is very important. You know that you should back it up. You know
that when you do so you should encrypt it. You even know that the most reliable backup medium is
ink on paper, and because a key is a fairly small amount of data it's feasible to back it up that
way.

But it would be so annoying to manually transcribe hundreds or thousands of characters in the
event that you do need to recover your key from the backup. If only there were some kind of
machine-readable represen...oh wait there is, it's QR codes, the answer is QR codes.

qrbak takes your key's ID as input and produces a PDF of QR codes as output. Simple as that.
