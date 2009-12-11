include $(GOROOT)/src/Make.$(GOARCH)

TARG=db
GOFILES=db.go classic.go util.go result.go

include $(GOROOT)/src/Make.pkg
