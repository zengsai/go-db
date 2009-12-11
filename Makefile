include $(GOROOT)/src/Make.$(GOARCH)

TARG=db
GOFILES=db.go classic.go

include $(GOROOT)/src/Make.pkg
